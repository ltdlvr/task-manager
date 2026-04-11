package handlers

import (
	"github.com/gofiber/fiber/v3"

	"github.com/ltdlvr/task-manager/internal/models"
	"github.com/ltdlvr/task-manager/internal/repository"
)

type UserHandler struct {
	Repo *repository.UserRepository
}

func HealthCheck(c fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status": "pong",
	})
}

func (h *UserHandler) Register(c fiber.Ctx) error {
	var user models.User
	err := c.Bind().Body(&user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}
	id, err := h.Repo.CreateUser(c.Context(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create user",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id,
	})
}
