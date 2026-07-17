package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/adapter/repo"
	"github.com/ltdlvr/task-manager/internal/core/model"
	"github.com/ltdlvr/task-manager/internal/tool"
)

type Tasks struct {
	tasksRepo        repo.Tasks
	columnsRepo      repo.Columns
	boardMembersRepo repo.BoardMembers
	dbClient         db.Client
}

func NewTasks(t repo.Tasks, c repo.Columns, bm repo.BoardMembers, d db.Client) *Tasks {
	return &Tasks{
		tasksRepo:        t,
		columnsRepo:      c,
		boardMembersRepo: bm,
		dbClient:         d,
	}
}

func (s *Tasks) Create(ctx context.Context, userID uint64, t *model.Task) error {
	return s.dbClient.ExecTrx(ctx, func(trx db.Trx) error {
		col, err := s.columnsRepo.GetByID(ctx, trx, t.ColumnID)
		if err != nil {
			return fmt.Errorf("get column: %w", err)
		}

		if err := s.requireBoardMember(ctx, trx, col.BoardID, userID); err != nil {
			return err
		}
		if err := s.columnsRepo.LockByID(ctx, trx, t.ColumnID); err != nil {
			return fmt.Errorf("lock column: %w", err)
		}

		targetPos := t.Position
		newPos, err := s.calculatePosition(ctx, trx, t.ColumnID, targetPos, func() ([]*model.Task, error) {
			tasks, err := s.tasksRepo.GetAllByColumn(ctx, trx, t.ColumnID)
			if err != nil {
				return nil, fmt.Errorf("get tasks by column: %w", err)
			}
			return tasks, nil
		})
		if err != nil {
			return fmt.Errorf("calculate task position: %w", err)
		}
		t.Position = newPos

		if err := s.tasksRepo.Create(ctx, trx, t); err != nil {
			return fmt.Errorf("create task: %w", err)
		}
		return nil
	})
}

func (s *Tasks) GetAllByColumn(ctx context.Context, columnID, userID uint64) ([]*model.Task, error) {
	col, err := s.columnsRepo.GetByID(ctx, s.dbClient, columnID)
	if err != nil {
		return nil, fmt.Errorf("get column: %w", err)
	}

	if err := s.requireBoardMember(ctx, s.dbClient, col.BoardID, userID); err != nil {
		return nil, err
	}

	tasks, err := s.tasksRepo.GetAllByColumn(ctx, s.dbClient, columnID)
	if err != nil {
		return nil, fmt.Errorf("get tasks by column id: %w", err)
	}
	return tasks, nil
}

func (s *Tasks) DeleteByID(ctx context.Context, userID, taskID uint64) error {
	task, err := s.tasksRepo.GetByID(ctx, s.dbClient, taskID)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}

	col, err := s.columnsRepo.GetByID(ctx, s.dbClient, task.ColumnID)
	if err != nil {
		return fmt.Errorf("get column: %w", err)
	}

	if err := s.requireBoardMember(ctx, s.dbClient, col.BoardID, userID); err != nil {
		return err
	}

	if err := s.tasksRepo.DeleteByID(ctx, s.dbClient, taskID); err != nil {
		return fmt.Errorf("delete task by id: %w", err)
	}
	return nil
}

func (s *Tasks) MoveTask(ctx context.Context, userID, taskID, targetColumnID uint64, targetPos int) error {
	return s.dbClient.ExecTrx(ctx, func(trx db.Trx) error {
		task, err := s.tasksRepo.GetByIDForUpdate(ctx, trx, taskID)
		if err != nil {
			return fmt.Errorf("get task for update: %w", err)
		}

		sourceColumn, err := s.columnsRepo.GetByID(ctx, trx, task.ColumnID)
		if err != nil {
			return fmt.Errorf("get source column: %w", err)
		}
		if err := s.requireBoardMember(ctx, trx, sourceColumn.BoardID, userID); err != nil {
			return err
		}

		targetColumn, err := s.columnsRepo.GetByID(ctx, trx, targetColumnID)
		if err != nil {
			return fmt.Errorf("get target column: %w", err)
		}
		if sourceColumn.BoardID != targetColumn.BoardID {
			return db.ErrForbidden
		}

		if err := s.lockColumns(ctx, trx, sourceColumn.ID, targetColumn.ID); err != nil {
			return err
		}

		newPos, err := s.calculatePosition(ctx, trx, targetColumnID, targetPos, func() ([]*model.Task, error) {
			tasks, err := s.tasksRepo.GetOtherByColumn(ctx, trx, targetColumnID, taskID)
			if err != nil {
				return nil, fmt.Errorf("get other tasks by column: %w", err)
			}
			return tasks, nil
		})
		if err != nil {
			return fmt.Errorf("calculate task position: %w", err)
		}

		if err := s.tasksRepo.UpdateColumnAndPosition(ctx, trx, targetColumnID, task.ID, newPos); err != nil {
			return fmt.Errorf("update column and position: %w", err)
		}
		return nil
	})
}

func (s *Tasks) lockColumns(ctx context.Context, client db.DB, firstID, secondID uint64) error {
	if firstID > secondID {
		firstID, secondID = secondID, firstID
	}
	if err := s.columnsRepo.LockByID(ctx, client, firstID); err != nil {
		return fmt.Errorf("lock column %d: %w", firstID, err)
	}
	if firstID != secondID {
		if err := s.columnsRepo.LockByID(ctx, client, secondID); err != nil {
			return fmt.Errorf("lock column %d: %w", secondID, err)
		}
	}
	return nil
}

func (s *Tasks) calculatePosition(ctx context.Context, client db.DB, columnID uint64, targetPos int, loadTasks func() ([]*model.Task, error)) (int, error) {
	tasks, err := loadTasks()
	if err != nil {
		return 0, err
	}
	prev, next, err := getNeighbours(tasks, targetPos)
	if err != nil {
		return 0, fmt.Errorf("get neighbours: %w", err)
	}
	newPos, err := tool.CalculateNewPosition(prev, next)
	if err == nil {
		return newPos, nil
	}

	if err := s.rebalanceColumn(ctx, client, columnID); err != nil {
		return 0, fmt.Errorf("rebalance column: %w", err)
	}

	tasks, err = loadTasks()
	if err != nil {
		return 0, fmt.Errorf("load tasks after rebalance: %w", err)
	}
	prev, next, err = getNeighbours(tasks, targetPos)
	if err != nil {
		return 0, fmt.Errorf("get neighbours after rebalance: %w", err)
	}
	newPos, err = tool.CalculateNewPosition(prev, next)
	if err != nil {
		return 0, fmt.Errorf("calculate position after rebalance: %w", err)
	}
	return newPos, nil
}

func (s *Tasks) rebalanceColumn(ctx context.Context, client db.DB, columnID uint64) error {
	tasks, err := s.tasksRepo.GetAllByColumn(ctx, client, columnID)
	if err != nil {
		return fmt.Errorf("get all tasks: %w", err)
	}

	for i, task := range tasks {
		newPos := tool.Startpos + (i+1)*tool.DefaultStep
		if err := s.tasksRepo.UpdateColumnAndPosition(ctx, client, columnID, task.ID, newPos); err != nil {
			return fmt.Errorf("update task %d: %w", task.ID, err)
		}
	}
	return nil
}

func (s *Tasks) requireBoardMember(ctx context.Context, client db.DB, boardID, userID uint64) error {
	_, err := s.boardMembersRepo.GetRole(ctx, client, boardID, userID)
	if err != nil {
		if errors.Is(err, db.ErrEntityNotFound) {
			return db.ErrForbidden
		}
		return fmt.Errorf("get board role: %w", err)
	}

	return nil
}
