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

type Columns struct {
	columnsRepo      repo.Columns
	boardsRepo       repo.Boards
	boardMembersRepo repo.BoardMembers
	dbClient         db.Client
}

func NewColumns(c repo.Columns, b repo.Boards, bm repo.BoardMembers, d db.Client) *Columns {
	return &Columns{
		columnsRepo:      c,
		boardsRepo:       b,
		boardMembersRepo: bm,
		dbClient:         d,
	}
}

func (s *Columns) Create(ctx context.Context, userID uint64, c *model.Column) error {
	return s.dbClient.ExecTrx(ctx, func(trx db.Trx) error {
		if err := s.requireBoardMember(ctx, trx, c.BoardID, userID); err != nil {
			return err
		}
		if err := s.boardsRepo.LockByID(ctx, trx, c.BoardID); err != nil {
			return fmt.Errorf("lock board: %w", err)
		}

		targetPos := c.Position
		newPos, err := s.calculatePosition(ctx, trx, c.BoardID, targetPos, func() ([]*model.Column, error) {
			columns, err := s.columnsRepo.GetAllByBoard(ctx, trx, c.BoardID)
			if err != nil {
				return nil, fmt.Errorf("get columns by board: %w", err)
			}
			return columns, nil
		})
		if err != nil {
			return fmt.Errorf("calculate column position: %w", err)
		}
		c.Position = newPos

		if err := s.columnsRepo.Create(ctx, trx, c); err != nil {
			return fmt.Errorf("create column: %w", err)
		}
		return nil
	})
}

func (s *Columns) GetAllByBoard(ctx context.Context, boardID, userID uint64) ([]*model.Column, error) {
	if err := s.requireBoardMember(ctx, s.dbClient, boardID, userID); err != nil {
		return nil, err
	}

	//check if exists
	if _, err := s.boardsRepo.GetByID(ctx, s.dbClient, boardID); err != nil {
		return nil, fmt.Errorf("get board: %w", err)
	}

	columns, err := s.columnsRepo.GetAllByBoard(ctx, s.dbClient, boardID)
	if err != nil {
		return nil, fmt.Errorf("get columns by board id: %w", err)
	}
	return columns, nil
}

func (s *Columns) DeleteByID(ctx context.Context, userID, columnID uint64) error {
	col, err := s.columnsRepo.GetByID(ctx, s.dbClient, columnID)
	if err != nil {
		return fmt.Errorf("get column: %w", err)
	}

	if err := s.requireBoardMember(ctx, s.dbClient, col.BoardID, userID); err != nil {
		return err
	}

	if err := s.columnsRepo.DeleteByID(ctx, s.dbClient, columnID); err != nil {
		return fmt.Errorf("delete column by id: %w", err)
	}
	return nil
}

func (s *Columns) MoveColumn(ctx context.Context, userID, columnID uint64, targetPos int) error {
	return s.dbClient.ExecTrx(ctx, func(trx db.Trx) error {
		col, err := s.columnsRepo.GetByID(ctx, trx, columnID)
		if err != nil {
			return fmt.Errorf("get column: %w", err)
		}

		if err := s.requireBoardMember(ctx, trx, col.BoardID, userID); err != nil {
			return err
		}
		if err := s.boardsRepo.LockByID(ctx, trx, col.BoardID); err != nil {
			return fmt.Errorf("lock board: %w", err)
		}

		newPos, err := s.calculatePosition(ctx, trx, col.BoardID, targetPos, func() ([]*model.Column, error) {
			columns, err := s.columnsRepo.GetOtherByBoard(ctx, trx, col.BoardID, columnID)
			if err != nil {
				return nil, fmt.Errorf("get other columns by board: %w", err)
			}
			return columns, nil
		})
		if err != nil {
			return fmt.Errorf("calculate column position: %w", err)
		}

		if err := s.columnsRepo.UpdatePosition(ctx, trx, newPos, col.ID); err != nil {
			return fmt.Errorf("update position: %w", err)
		}

		return nil
	})
}

func (s *Columns) calculatePosition(ctx context.Context, client db.DB, boardID uint64, targetPos int, loadColumns func() ([]*model.Column, error)) (int, error) {
	columns, err := loadColumns()
	if err != nil {
		return 0, err
	}
	prev, next, err := getNeighbours(columns, targetPos)
	if err != nil {
		return 0, fmt.Errorf("get neighbours: %w", err)
	}
	newPos, err := tool.CalculateNewPosition(prev, next)
	if err == nil {
		return newPos, nil
	}

	if err := s.rebalanceBoard(ctx, client, boardID); err != nil {
		return 0, fmt.Errorf("rebalance board: %w", err)
	}

	columns, err = loadColumns()
	if err != nil {
		return 0, fmt.Errorf("load columns after rebalance: %w", err)
	}
	prev, next, err = getNeighbours(columns, targetPos)
	if err != nil {
		return 0, fmt.Errorf("get neighbours after rebalance: %w", err)
	}
	newPos, err = tool.CalculateNewPosition(prev, next)
	if err != nil {
		return 0, fmt.Errorf("calculate position after rebalance: %w", err)
	}
	return newPos, nil
}

func (s *Columns) rebalanceBoard(ctx context.Context, client db.DB, boardID uint64) error {
	columns, err := s.columnsRepo.GetAllByBoard(ctx, client, boardID)
	if err != nil {
		return fmt.Errorf("get all columns: %w", err)
	}

	for i, col := range columns {
		newPos := tool.Startpos + (i+1)*tool.DefaultStep
		if err := s.columnsRepo.UpdatePosition(ctx, client, newPos, col.ID); err != nil {
			return fmt.Errorf("update column %d: %w", col.ID, err)
		}
	}
	return nil
}

func (s *Columns) requireBoardMember(ctx context.Context, client db.DB, boardID, userID uint64) error {
	_, err := s.boardMembersRepo.GetRole(ctx, client, boardID, userID)
	if err != nil {
		if errors.Is(err, db.ErrEntityNotFound) {
			return db.ErrForbidden
		}
		return fmt.Errorf("get board role: %w", err)
	}

	return nil
}
