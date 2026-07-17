package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/adapter/repo"
	"github.com/ltdlvr/task-manager/internal/core/model"
)

type Boards struct {
	boardsRepo       repo.Boards
	boardMembersRepo repo.BoardMembers
	dbClient         db.Client
}

// TODO - решить с названиями для db, а то в репо всегда client, в сервисе где db.Client, там dbClient, а где db.DB, там db
func NewBoards(b repo.Boards, bm repo.BoardMembers, d db.Client) *Boards {
	return &Boards{
		boardsRepo:       b,
		boardMembersRepo: bm,
		dbClient:         d,
	}
}

func (s *Boards) Create(ctx context.Context, userID uint64, b *model.Board) error {
	return s.dbClient.ExecTrx(ctx, func(trx db.Trx) error {
		if err := s.boardsRepo.Create(ctx, trx, b); err != nil {
			return fmt.Errorf("create board: %w", err)
		}

		bm := &model.BoardMember{
			BoardID: b.ID,
			UserID:  userID,
			Role:    model.BoardRoleOwner,
		}

		if err := s.boardMembersRepo.Create(ctx, trx, bm); err != nil {
			return fmt.Errorf("add board member: %w", err)
		}

		return nil
	})
}

func (s *Boards) GetByID(ctx context.Context, userID, boardID uint64) (*model.Board, error) {
	if _, err := s.requireMember(ctx, boardID, userID); err != nil {
		return nil, err
	}

	board, err := s.boardsRepo.GetByID(ctx, s.dbClient, boardID)
	if err != nil {
		return nil, fmt.Errorf("get board: %w", err)
	}

	return board, nil
}

func (s *Boards) DeleteById(ctx context.Context, userID, boardID uint64) error {
	role, err := s.requireMember(ctx, boardID, userID)
	if err != nil {
		return err
	}

	if role != model.BoardRoleOwner {
		return db.ErrForbidden
	}

	if err := s.boardsRepo.DeleteByID(ctx, s.dbClient, boardID); err != nil {
		return fmt.Errorf("delete board: %w", err)
	}

	return nil
}

func (s *Boards) requireMember(ctx context.Context, boardID, userID uint64) (model.BoardRole, error) {
	if _, err := s.boardsRepo.GetByID(ctx, s.dbClient, boardID); err != nil {
		return "", err
	}

	role, err := s.boardMembersRepo.GetRole(ctx, s.dbClient, boardID, userID)
	if err != nil {
		if errors.Is(err, db.ErrEntityNotFound) {
			return "", db.ErrEntityNotFound
		}
		return "", fmt.Errorf("get role: %w", err)
	}

	return role, nil
}
