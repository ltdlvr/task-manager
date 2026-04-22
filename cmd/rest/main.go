package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"

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

	// Services
	authService := service.NewAuth(usersRepo, dbClient, pswdTool)

	// Handlers
	authHandler := rest.NewAuth(authService)
	hcHandler := rest.NewHealthcheck()

	// Init app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return tool.MapHttpError(c, err)
		},
	})

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Assign handlers
	v1.Get("/healthcheck", hcHandler.Check)
	v1.Post("/register", authHandler.Register)
	v1.Post("/login", authHandler.LogIn)

	app.Listen(fmt.Sprintf("%s:%s", conf.ServerHost(), conf.ServerPort()))
}
