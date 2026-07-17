package rest

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/ltdlvr/task-manager/internal/core/model"
	"github.com/ltdlvr/task-manager/internal/core/service"
)

type createColumnReq struct {
	Name      string `json:"name"`
	TargetPos int    `json:"targetPos"`
}

type createColumnURI struct {
	BoardID uint64 `uri:"boardId"`
}

type getColumnsByBoardReq struct {
	BoardID uint64 `uri:"boardId"`
}

type deleteColumnReq struct {
	ID uint64 `uri:"id"`
}

type moveColumnReq struct {
	TargetPos int `json:"targetPos"`
}

type moveColumnURI struct {
	ID uint64 `uri:"id"`
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
	var uri createColumnURI
	if err := c.Bind().URI(&uri); err != nil {
		return fiber.ErrBadRequest
	}

	var req createColumnReq
	if err := c.Bind().Body(&req); err != nil {
		return fiber.ErrBadRequest
	}

	col := model.Column{
		BoardID:  uri.BoardID,
		Name:     req.Name,
		Position: req.TargetPos,
	}

	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	if err := h.columnService.Create(c.Context(), userID, &col); err != nil {
		return err
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
		return fiber.ErrBadRequest
	}

	boardID := uri.BoardID

	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	columns, err := h.columnService.GetAllByBoard(c.Context(), boardID, userID)
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
		return fiber.ErrBadRequest
	}

	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	if err := h.columnService.DeleteByID(c.Context(), userID, uri.ID); err != nil {
		return err
	}

	return c.SendStatus(204)
}

func (h *Columns) MoveColumn(c fiber.Ctx) error {
	var uri moveColumnURI
	if err := c.Bind().URI(&uri); err != nil {
		return fiber.ErrBadRequest
	}

	var req moveColumnReq
	if err := c.Bind().Body(&req); err != nil {
		return fiber.ErrBadRequest
	}

	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	if err := h.columnService.MoveColumn(c.Context(), userID, uri.ID, req.TargetPos); err != nil {
		return err
	}

	return c.SendStatus(204)
}
