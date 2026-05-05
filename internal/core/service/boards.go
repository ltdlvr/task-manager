package service

import (
	"context"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/adapter/repo"
	"github.com/ltdlvr/task-manager/internal/core/model"
)

type Boards struct {
	boardsRepo repo.Boards
	dbClient   db.Client
}

func NewBoards(b repo.Boards, d db.Client) *Boards {
	return &Boards{
		boardsRepo: b,
		dbClient:   d,
	}
}

func (s *Boards) Create(ctx context.Context, b *model.Board) error {
	return s.boardsRepo.Create(ctx, s.dbClient, b)
}

func (s *Boards) GetByID(ctx context.Context, id uint64) (*model.Board, error) {
	return s.boardsRepo.GetByID(ctx, s.dbClient, id)
}

func (s *Boards) DeleteById(ctx context.Context, id uint64) error {
	return s.boardsRepo.DeleteByID(ctx, s.dbClient, id)
}
