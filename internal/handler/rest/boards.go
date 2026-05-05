package rest

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/ltdlvr/task-manager/internal/core/model"
	"github.com/ltdlvr/task-manager/internal/core/service"
)

type boardReq struct {
	ID   uint64 `json:"id" uri:"id"`
	Name string `json:"name"`
}

type boardRes struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type Boards struct {
	boardsService *service.Boards
}

func NewBoards(b *service.Boards) *Boards {
	return &Boards{
		boardsService: b,
	}
}

func (h *Boards) Create(c fiber.Ctx) error {
	var body boardReq
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	b := model.Board{
		Name: body.Name,
	}

	if err := h.boardsService.Create(c.Context(), &b); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(201).JSON(boardRes{
		ID:        b.ID,
		Name:      b.Name,
		CreatedAt: b.CreatedAt,
	})
}

func (h *Boards) GetByID(c fiber.Ctx) error {
	var uri boardReq

	if err := c.Bind().URI(&uri); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	board, err := h.boardsService.GetByID(c.Context(), uri.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	return c.Status(200).JSON(boardRes{
		ID:        board.ID,
		Name:      board.Name,
		CreatedAt: board.CreatedAt,
	})
}

func (h *Boards) DeleteByID(c fiber.Ctx) error {
	var uri boardReq

	if err := c.Bind().URI(&uri); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	if err := h.boardsService.DeleteById(c.Context(), uri.ID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	return c.SendStatus(204)
}
