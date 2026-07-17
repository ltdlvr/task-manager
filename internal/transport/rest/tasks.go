package rest

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/ltdlvr/task-manager/internal/core/model"
	"github.com/ltdlvr/task-manager/internal/core/service"
)

type createTaskReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TargetPos   int    `json:"targetPos"`
}

type createTaskURI struct {
	ColumnID uint64 `uri:"columnId"`
}

type getTasksByColumnReq struct {
	ColumnID uint64 `uri:"columnId"`
}

type deleteTaskReq struct {
	ID uint64 `uri:"id"`
}

type moveTaskReq struct {
	TargetColumnID uint64 `json:"targetColumnId"`
	TargetPos      int    `json:"targetPos"`
}

type moveTaskURI struct {
	ID uint64 `uri:"id"`
}

type TaskRes struct {
	ID          uint64    `json:"id"`
	ColumnID    uint64    `json:"columnId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Position    int       `json:"position"`
	CreatedAt   time.Time `json:"createdAt"`
}

type Tasks struct {
	taskService *service.Tasks
}

func NewTasks(t *service.Tasks) *Tasks {
	return &Tasks{
		taskService: t,
	}
}

func (h *Tasks) Create(c fiber.Ctx) error {
	var uri createTaskURI
	if err := c.Bind().URI(&uri); err != nil {
		return fiber.ErrBadRequest
	}

	var req createTaskReq
	if err := c.Bind().Body(&req); err != nil {
		return fiber.ErrBadRequest
	}

	task := model.Task{
		ColumnID:    uri.ColumnID,
		Title:       req.Title,
		Description: req.Description,
		Position:    req.TargetPos,
	}

	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	if err := h.taskService.Create(c.Context(), userID, &task); err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(TaskRes{
		ID:          task.ID,
		ColumnID:    task.ColumnID,
		Title:       task.Title,
		Description: task.Description,
		Position:    task.Position,
		CreatedAt:   task.CreatedAt,
	})
}

func (h *Tasks) GetAllByColumn(c fiber.Ctx) error {
	var uri getTasksByColumnReq
	if err := c.Bind().URI(&uri); err != nil {
		return fiber.ErrBadRequest
	}

	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	tasks, err := h.taskService.GetAllByColumn(c.Context(), uri.ColumnID, userID)
	if err != nil {
		return err
	}

	res := make([]TaskRes, len(tasks))
	for i, t := range tasks {
		res[i] = TaskRes{
			ID:          t.ID,
			ColumnID:    t.ColumnID,
			Title:       t.Title,
			Description: t.Description,
			Position:    t.Position,
			CreatedAt:   t.CreatedAt,
		}
	}
	return c.Status(200).JSON(res)
}

func (h *Tasks) DeleteByID(c fiber.Ctx) error {
	var uri deleteTaskReq
	if err := c.Bind().URI(&uri); err != nil {
		return fiber.ErrBadRequest
	}

	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	if err := h.taskService.DeleteByID(c.Context(), userID, uri.ID); err != nil {
		return err
	}

	return c.SendStatus(204)
}

func (h *Tasks) MoveTask(c fiber.Ctx) error {
	var uri moveTaskURI
	if err := c.Bind().URI(&uri); err != nil {
		return fiber.ErrBadRequest
	}

	var req moveTaskReq
	if err := c.Bind().Body(&req); err != nil {
		return fiber.ErrBadRequest
	}

	userID, err := getUserID(c)
	if err != nil {
		return err
	}

	if err := h.taskService.MoveTask(c.Context(), userID, uri.ID, req.TargetColumnID, req.TargetPos); err != nil {
		return err
	}

	return c.SendStatus(204)
}
