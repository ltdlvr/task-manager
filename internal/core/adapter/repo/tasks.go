package repo

import (
	"context"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/model"
)

type Tasks interface {
	Create(ctx context.Context, client db.DB, t *model.Task) error

	GetByID(ctx context.Context, client db.DB, id uint64) (*model.Task, error)

	GetByIDForUpdate(ctx context.Context, client db.DB, id uint64) (*model.Task, error)

	GetAllByColumn(ctx context.Context, client db.DB, columnID uint64) ([]*model.Task, error)

	GetOtherByColumn(ctx context.Context, client db.DB, columnID, taskID uint64) ([]*model.Task, error)

	DeleteByID(ctx context.Context, client db.DB, id uint64) error

	UpdateColumnAndPosition(ctx context.Context, client db.DB, columnID, taskID uint64, newPos int) error
}
