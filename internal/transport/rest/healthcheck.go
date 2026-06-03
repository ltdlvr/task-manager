package rest

import (
	"github.com/gofiber/fiber/v3"
)

type Healthcheck struct{}

func NewHealthcheck() *Healthcheck {
	return &Healthcheck{}
}

func (h *Healthcheck) Check(c fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status": "ok",
	})
}
