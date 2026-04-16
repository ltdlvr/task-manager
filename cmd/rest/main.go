package main

import (
	"fmt"

	"github.com/gofiber/fiber/v3"

	"github.com/ltdlvr/task-manager/internal/config"
	"github.com/ltdlvr/task-manager/internal/handler/rest"
	"github.com/ltdlvr/task-manager/internal/infra/db/pg"
	"github.com/ltdlvr/task-manager/internal/infra/repo"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Init config
	conf := config.NewRest()

	// Init db
	dbClient, err := pg.NewClient(conf.DatabaseURL())
	if err != nil {
		log.Fatalf("Failed to create database client: %v", err)
	}

	// InitDeps

	// Repositories
	usersRepo := repo.NewUsers()

	// Handlers
	authHandler := rest.NewAuth()
	hcHandler := rest.NewHealthcheck()

	// Init server
	app := fiber.New()

	app.Get("/healthcheck", hcHandler.Check)
	app.Post("/register", authHandler.Register)
	app.Listen(fmt.Sprintf("%s:%s", conf.ServerHost(), conf.ServerPort()))
}
