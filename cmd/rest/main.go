package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/ltdlvr/task-manager/internal/config"
	"github.com/ltdlvr/task-manager/internal/core/service"
	"github.com/ltdlvr/task-manager/internal/handler/rest"
	"github.com/ltdlvr/task-manager/internal/infra/db/pg"
	"github.com/ltdlvr/task-manager/internal/infra/repo"
	"github.com/ltdlvr/task-manager/internal/tool"
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

	// Repositories
	usersRepo := repo.NewUsers()
	boardsRepo := repo.NewBoards()
	columnsRepo := repo.NewColumns()

	// Services
	authService := service.NewAuth(usersRepo, dbClient, pswdTool)
	boardsService := service.NewBoards(boardsRepo, dbClient)
	columnsService := service.NewColumns(columnsRepo, boardsRepo, dbClient)

	// Handlers
	authHandler := rest.NewAuth(authService)
	hcHandler := rest.NewHealthcheck()
	boardsHandler := rest.NewBoards(boardsService)
	columnsHandler := rest.NewColumns(columnsService)

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

	// Boards
	v1.Post("/boards", boardsHandler.Create)
	v1.Get("/boards/:id", boardsHandler.GetByID)
	v1.Delete("/boards/:id", boardsHandler.DeleteByID)

	// Columns
	v1.Post("/boards/:boardId/columns", columnsHandler.Create)
	v1.Get("/boards/:boardId/columns", columnsHandler.GetAllByBoard) // NOTE - можно объединить get в колонках и бордах, когда фронт появится (чтобы не вызывать 2 раздельных запроса)
	v1.Delete("/columns/:id", columnsHandler.DeleteByID)
	v1.Patch("/columns/:id/move", columnsHandler.MoveColumn)

	app.Listen(fmt.Sprintf("%s:%s", conf.ServerHost(), conf.ServerPort()))
}
