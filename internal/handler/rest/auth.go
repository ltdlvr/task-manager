package rest

import (
	"github.com/gofiber/fiber/v3"
)

type registerReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

func (h *Auth) Register(c fiber.Ctx) error {
	var body registerReq

	if err := c.Bind().Body(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

}
