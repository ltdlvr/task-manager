package repo

import (
	"context"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/model"
)

type Boards interface {
	Create(ctx context.Context, client db.DB, b *model.Board) error

	DeleteByID(ctx context.Context, client db.DB, id uint64) error

	GetByID(ctx context.Context, client db.DB, id uint64) (*model.Board, error)

	LockByID(ctx context.Context, client db.DB, id uint64) error
}
