package repo

import (
	"context"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/model"
)

type Columns interface {
	Create(ctx context.Context, client db.DB, c *model.Column) error

	GetByID(ctx context.Context, client db.DB, id uint64) (*model.Column, error)

	GetAllByBoard(ctx context.Context, client db.DB, boardId uint64) ([]*model.Column, error)

	GetOtherByBoard(ctx context.Context, client db.DB, boardID uint64, columnID uint64) ([]*model.Column, error)

	DeleteByID(ctx context.Context, client db.DB, id uint64) error

	UpdatePosition(ctx context.Context, client db.DB, newPos int, id uint64) error

	//GetNeighbours(ctx context.Context, client db.Client, c *model.Column) (prev, next *model.Column, err error)
}
