package rest

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/ltdlvr/task-manager/internal/core/model"
	"github.com/ltdlvr/task-manager/internal/core/service"
)

type createColumnReq struct {
	BoardID   uint64 `json:"boardId" uri:"boardId"`
	Name      string `json:"name"`
	TargetPos int    `json:"position"`
}

type getColumnsByBoardReq struct {
	BoardID uint64 `uri:"boardId"`
}

type deleteColumnReq struct {
	ID uint64 `uri:"id"`
}

type moveColumnReq struct {
	ID        uint64 `uri:"id"`
	TargetPos int    `json:"targetPos"`
}

type ColumnRes struct {
	ID        uint64    `json:"id"`
	BoardID   uint64    `json:"boardId"`
	Name      string    `json:"name"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"createdAt"`
}

type Columns struct {
	columnService *service.Columns
}

func NewColumns(c *service.Columns) *Columns {
	return &Columns{
		columnService: c,
	}
}

func (h *Columns) Create(c fiber.Ctx) error {
	var req createColumnReq
	if err := c.Bind().All(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	col := model.Column{
		BoardID:  req.BoardID,
		Name:     req.Name,
		Position: req.TargetPos,
	}

	if err := h.columnService.Create(c.Context(), &col); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(ColumnRes{
		ID:        col.ID,
		BoardID:   col.BoardID,
		Name:      col.Name,
		Position:  col.Position,
		CreatedAt: col.CreatedAt,
	})
}

func (h *Columns) GetAllByBoard(c fiber.Ctx) error {
	var uri getColumnsByBoardReq
	if err := c.Bind().URI(&uri); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid board id",
		})
	}

	boardID := uri.BoardID

	columns, err := h.columnService.GetAllByBoard(c.Context(), boardID)
	if err != nil {
		return err
	}

	res := make([]ColumnRes, len(columns))

	for i, col := range columns {
		res[i] = ColumnRes{
			ID:        col.ID,
			BoardID:   col.BoardID,
			Name:      col.Name,
			Position:  col.Position,
			CreatedAt: col.CreatedAt,
		}
	}
	return c.Status(200).JSON(res)
}

func (h *Columns) DeleteByID(c fiber.Ctx) error {
	var uri deleteColumnReq
	if err := c.Bind().URI(&uri); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid column id",
		})
	}

	if err := h.columnService.DeleteByID(c.Context(), uri.ID); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.SendStatus(204)
}

func (h *Columns) MoveColumn(c fiber.Ctx) error {
	var req moveColumnReq
	if err := c.Bind().All(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid request",
		})
	}

	if err := h.columnService.MoveColumn(c.Context(), req.ID, req.TargetPos); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.SendStatus(204)
}
