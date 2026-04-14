package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/ltdlvr/task-manager/internal/db"
	"github.com/ltdlvr/task-manager/internal/handlers"
	"github.com/ltdlvr/task-manager/internal/repository"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load() //я думал закинуть это в отдельный файл, но в описании функции написано, что ее в идеале в мейне вызывать
	if err != nil {
		log.Warn("No .env file found")
	}

	pool, err := db.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	} else {
		log.Info("Database initialized")
	}
	repo := &repository.UserRepository{Pool: pool}
	userHandler := &handlers.UserHandler{Repo: repo}

	app := fiber.New()

	//TODO - добить эндпоинты регистрации

	app.Get("/ping", handlers.HealthCheck)
	app.Post("/register", userHandler.Register)
	app.Delete("/delete_account", userHandler.DeleteUser)
	app.Listen(":3000")
}
