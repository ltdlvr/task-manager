package repo

import (
	"context"
	"time"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/model"
)

type Users struct{}

func NewUsers() *Users {
	return &Users{}
}

func (r *Users) Save(ctx context.Context, client db.Client, u *model.User) error {
	row := client.QueryRowContext(
		ctx,
		"INSERT INTO users (name, password) VALUES ($1, $2) RETURNING id, created_at",
		u.Name, u.Password,
	)

	var id uint64
	var createdAt time.Time
	if err := row.Scan(&id, &createdAt); err != nil {
		return db.MapError(err)
	}

	u.ID = id
	u.CreatedAt = createdAt
	return nil
}

func (r *Users) DeleteByID(ctx context.Context, client db.Client, id uint64) error {
	_, err := client.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	return db.MapError(err)
}
