package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/ltdlvr/task-manager/internal/config"
	"github.com/ltdlvr/task-manager/internal/core/service"
	"github.com/ltdlvr/task-manager/internal/infra/db/pg"
	"github.com/ltdlvr/task-manager/internal/infra/repo"
	"github.com/ltdlvr/task-manager/internal/tool"
	"github.com/ltdlvr/task-manager/internal/transport/rest"
)

func main() {
	// Init config
	conf := config.NewRest()

	// Init db
	dbClient, err := pg.NewClient(conf.DatabaseURL())
	if err != nil {
		log.Fatalf("Failed to create database client: %v", err)
	}

	// Init deps
	// Misc
	pswdTool := tool.NewPassword()
	tokenTool := tool.NewToken(conf.JWTSecret())

	// Repositories
	usersRepo := repo.NewUsers()
	boardsRepo := repo.NewBoards()
	columnsRepo := repo.NewColumns()

	// Services
	authService := service.NewAuth(usersRepo, dbClient, pswdTool, tokenTool)
	boardsService := service.NewBoards(boardsRepo, dbClient)
	columnsService := service.NewColumns(columnsRepo, boardsRepo, dbClient)

	// Handlers
	authHandler := rest.NewAuth(authService)
	hcHandler := rest.NewHealthcheck()
	boardsHandler := rest.NewBoards(boardsService)
	columnsHandler := rest.NewColumns(columnsService)
	authMiddleware := rest.AuthMiddleware(tokenTool)

	// Init app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return tool.MapHttpError(c, err)
		},
	})
	app.Use(logger.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Assign handlers
	// Health
	v1.Get("/healthcheck", hcHandler.Check)

	// Auth
	v1.Post("/register", authHandler.Register)
	v1.Post("/login", authHandler.LogIn)

	// Jwt middleware
	private := v1.Group("", authMiddleware)

	// Boards
	private.Post("/boards", boardsHandler.Create)
	private.Get("/boards/:id", boardsHandler.GetByID)
	private.Delete("/boards/:id", boardsHandler.DeleteByID)

	// Columns
	private.Post("/boards/:boardId/columns", columnsHandler.Create)
	private.Get("/boards/:boardId/columns", columnsHandler.GetAllByBoard) // NOTE - можно объединить get в колонках и бордах, когда (если) фронт появится (чтобы не вызывать 2 раздельных запроса)
	private.Delete("/columns/:id", columnsHandler.DeleteByID)
	private.Patch("/columns/:id/move", columnsHandler.MoveColumn)

	app.Listen(fmt.Sprintf("%s:%s", conf.ServerHost(), conf.ServerPort()))
}
