package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"

	"github.com/ltdlvr/task-manager/internal/config"
	"github.com/ltdlvr/task-manager/internal/core/service"
	"github.com/ltdlvr/task-manager/internal/infra/db/pg"
	"github.com/ltdlvr/task-manager/internal/infra/repo"
	"github.com/ltdlvr/task-manager/internal/tool"
	"github.com/ltdlvr/task-manager/internal/transport/rest"
)

const requestBodyLimit = 64 * 1024

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
	tasksRepo := repo.NewTasks()
	boardMembersRepo := repo.NewBoardMembers()

	// Services
	authService := service.NewAuth(usersRepo, dbClient, pswdTool, tokenTool)
	boardsService := service.NewBoards(boardsRepo, boardMembersRepo, dbClient)
	columnsService := service.NewColumns(columnsRepo, boardsRepo, boardMembersRepo, dbClient)
	tasksService := service.NewTasks(tasksRepo, columnsRepo, boardMembersRepo, dbClient)

	// Handlers
	authHandler := rest.NewAuth(authService)
	hcHandler := rest.NewHealthcheck()
	boardsHandler := rest.NewBoards(boardsService)
	columnsHandler := rest.NewColumns(columnsService)
	tasksHandler := rest.NewTasks(tasksService)
	authMiddleware := rest.AuthMiddleware(tokenTool)

	// Init app
	app := fiber.New(fiber.Config{
		BodyLimit: requestBodyLimit,
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
	authLimiter := limiter.New(limiter.Config{
		Max:        5,
		Expiration: time.Minute,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP() + ":" + c.Path()
		},
	})
	v1.Post("/register", authLimiter, authHandler.Register)
	v1.Post("/login", authLimiter, authHandler.LogIn)

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

	// Tasks
	private.Post("/columns/:columnId/tasks", tasksHandler.Create)
	private.Get("/columns/:columnId/tasks", tasksHandler.GetAllByColumn)
	private.Delete("/tasks/:id", tasksHandler.DeleteByID)
	private.Patch("/tasks/:id/move", tasksHandler.MoveTask)

	app.Listen(fmt.Sprintf("%s:%s", conf.ServerHost(), conf.ServerPort()))
}
