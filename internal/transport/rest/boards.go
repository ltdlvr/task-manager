package rest

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/ltdlvr/task-manager/internal/core/model"
	"github.com/ltdlvr/task-manager/internal/core/service"
)

type createBoardReq struct {
	Name string `json:"name"`
}

type byIdBoardReq struct {
	ID uint64 `uri:"id"`
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
	var body createBoardReq
	if err := c.Bind().Body(&body); err != nil {
		return fiber.ErrBadRequest
	}

	b := model.Board{
		Name: body.Name,
	}

	if err := h.boardsService.Create(c.Context(), &b); err != nil {
		return err
	}

	return c.Status(201).JSON(boardRes{
		ID:        b.ID,
		Name:      b.Name,
		CreatedAt: b.CreatedAt,
	})
}

func (h *Boards) GetByID(c fiber.Ctx) error {
	var uri byIdBoardReq

	if err := c.Bind().URI(&uri); err != nil {
		return fiber.ErrBadRequest
	}

	board, err := h.boardsService.GetByID(c.Context(), uri.ID)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(boardRes{
		ID:        board.ID,
		Name:      board.Name,
		CreatedAt: board.CreatedAt,
	})
}

func (h *Boards) DeleteByID(c fiber.Ctx) error {
	var uri byIdBoardReq

	if err := c.Bind().URI(&uri); err != nil {
		return fiber.ErrBadRequest
	}

	if err := h.boardsService.DeleteById(c.Context(), uri.ID); err != nil {
		return err
	}
	return c.SendStatus(204)
}
