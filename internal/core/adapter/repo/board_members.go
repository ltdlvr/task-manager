package repo

import (
	"context"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/model"
)

type BoardMembers interface {
	Create(ctx context.Context, client db.DB, m *model.BoardMember) error

	GetRole(ctx context.Context, client db.DB, boardID, userID uint64) (model.BoardRole, error)
}
