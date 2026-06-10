package repo

import (
	"context"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/model"
)

type Users interface {
	Save(ctx context.Context, client db.DB, u *model.User) error

	DeleteByID(ctx context.Context, client db.DB, id uint64) error

	GetByName(ctx context.Context, client db.DB, name string) (*model.User, error)
}
