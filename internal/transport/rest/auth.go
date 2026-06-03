package rest

import (
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/ltdlvr/task-manager/internal/core/model"
	"github.com/ltdlvr/task-manager/internal/core/service"
)

type registerReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type registerRes struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type Auth struct {
	authService *service.Auth
}

func NewAuth(s *service.Auth) *Auth {
	return &Auth{
		authService: s,
	}
}

func (h *Auth) Register(c fiber.Ctx) error {
	var body registerReq
	if err := c.Bind().Body(&body); err != nil {
		return err
	}

	u := model.User{
		Name:     body.Name,
		Password: body.Password,
	}
	if err := h.authService.Register(c.Context(), &u); err != nil {
		return err
	}

	return c.Status(201).JSON(registerRes{
		ID:        u.ID,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
	})
}

func (h *Auth) LogIn(c fiber.Ctx) error {
	var body registerReq
	if err := c.Bind().Body(&body); err != nil {
		return fiber.ErrBadRequest
	}
	u := model.User{
		Name:     body.Name,
		Password: body.Password,
	}
	token, err := h.authService.LogIn(c.Context(), &u)
	if err != nil {
		return fiber.ErrUnauthorized
	}
	return c.Status(200).JSON(fiber.Map{
		"token": token,
	})
}
