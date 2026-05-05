package repo

import (
	"context"
	"time"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/model"
)

type Boards struct{}

func NewBoards() *Boards {
	return &Boards{}
}

func (r *Boards) Create(ctx context.Context, client db.Client, b *model.Board) error {
	row := client.QueryRowContext(ctx, "INSERT INTO boards (name) VALUES ($1) RETURNING id, created_at", b.Name)

	var id uint64
	var createdAt time.Time

	if err := row.Scan(&id, &createdAt); err != nil {
		return db.MapError(err)
	}

	b.ID = id
	b.CreatedAt = createdAt
	return nil
}

func (r *Boards) DeleteByID(ctx context.Context, client db.Client, id uint64) error {
	_, err := client.ExecContext(ctx, "DELETE FROM boards where id = $1", id)
	return db.MapError(err)
}

func (r *Boards) GetByID(ctx context.Context, client db.Client, id uint64) (*model.Board, error) {
	row := client.QueryRowContext(ctx, "SELECT name, created_at  FROM boards where id = $1", id)

	var name string
	var createdAt time.Time

	if err := row.Scan(&name, &createdAt); err != nil {
		return nil, db.MapError(err)
	}

	b := &model.Board{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
	}

	return b, nil
}
