package rest

import (
	"github.com/gofiber/fiber/v3"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/adapter/repo"
	"github.com/ltdlvr/task-manager/internal/core/model"
)

type Auth struct {
	users    repo.Users
	dbClient db.Client
}

func NewAuth(users repo.Users, dbClient db.Client) *Auth {
	return &Auth{
		users:    users,
		dbClient: dbClient,
	}
}

func (h *Auth) Register(c fiber.Ctx) error {
	var user model.User
	err := c.Bind().Body(&user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	if err := h.users.Save(c.Context(), h.dbClient, &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": user.ID,
	})
}
