package repo

import (
	"context"
	"fmt"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/model"
)

type Columns struct{}

func NewColumns() *Columns {
	return &Columns{}
}

func (r *Columns) Create(ctx context.Context, client db.DB, c *model.Column) error {
	row := client.QueryRowContext(
		ctx,
		"INSERT INTO columns (board_id, name, position) VALUES ($1, $2, $3) RETURNING id, created_at",
		c.BoardID, c.Name, c.Position,
	)
	if err := row.Scan(&c.ID, &c.CreatedAt); err != nil {
		return db.MapError(err)
	}
	return nil
}

func (r *Columns) GetByID(ctx context.Context, client db.DB, id uint64) (*model.Column, error) {
	row := client.QueryRowContext(
		ctx, "SELECT board_id, name, position, created_at FROM columns WHERE id = $1", id,
	)
	c := &model.Column{
		ID: id,
	}

	if err := row.Scan(&c.BoardID, &c.Name, &c.Position, &c.CreatedAt); err != nil {
		return nil, db.MapError(err)
	}

	return c, nil
}

func (r *Columns) LockByID(ctx context.Context, client db.DB, id uint64) error {
	row := client.QueryRowContext(ctx, "SELECT id FROM columns WHERE id = $1 FOR UPDATE", id)
	var lockedID uint64
	return db.MapError(row.Scan(&lockedID))
}

func (r *Columns) GetAllByBoard(ctx context.Context, client db.DB, boardID uint64) ([]*model.Column, error) {
	rows, err := client.QueryContext(ctx, `
		SELECT id, name, position, created_at
		FROM columns
		WHERE board_id = $1
		ORDER BY position ASC, id ASC
	`, boardID)
	if err != nil {
		return nil, db.MapError(err)
	}
	defer rows.Close()

	var columns []*model.Column
	for rows.Next() {
		c := &model.Column{BoardID: boardID}
		if err := rows.Scan(&c.ID, &c.Name, &c.Position, &c.CreatedAt); err != nil {
			return nil, db.MapError(err)
		}
		columns = append(columns, c)
	}
	if err := rows.Err(); err != nil {
		return nil, db.MapError(err)
	}
	return columns, nil
}

func (r *Columns) GetOtherByBoard(ctx context.Context, client db.DB, boardID uint64, columnID uint64) ([]*model.Column, error) {
	rows, err := client.QueryContext(ctx, `
		SELECT id, name, position, created_at
		FROM columns
		WHERE board_id = $1 AND id != $2
		ORDER BY position ASC, id ASC
	`, boardID, columnID)
	if err != nil {
		return nil, db.MapError(err)
	}
	defer rows.Close()

	columns := make([]*model.Column, 0)
	for rows.Next() {
		c := &model.Column{BoardID: boardID}
		if err := rows.Scan(&c.ID, &c.Name, &c.Position, &c.CreatedAt); err != nil {
			return nil, db.MapError(err)
		}
		columns = append(columns, c)
	}
	if err := rows.Err(); err != nil {
		return nil, db.MapError(err)
	}
	return columns, nil
}

// NOTE - можно переделать с учетом RowsAffected чтобы отлавливать удаление несуществующего
func (r *Columns) DeleteByID(ctx context.Context, client db.DB, id uint64) error {
	_, err := client.ExecContext(
		ctx, "DELETE FROM columns WHERE id = $1", id,
	)
	return db.MapError(err)
}

func (r *Columns) UpdatePosition(ctx context.Context, client db.DB, newPos int, id uint64) error {
	res, err := client.ExecContext(ctx, `
		UPDATE columns
		SET position = $1
		WHERE id = $2
	`, newPos, id)
	if err != nil {
		return db.MapError(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if rows == 0 {
		return db.ErrEntityNotFound
	}
	return nil
}

// func (r *Columns) GetNeighbours(ctx context.Context, client db.Client, c *model.Column) (prev, next *model.Column, err error) {
// 	prev = &model.Column{BoardID: c.BoardID}
// 	row := client.QueryRowContext(ctx, `
// 		SELECT id, name, position, created_at
// 		FROM columns
// 		WHERE board_id = $1 and position < $2
// 		ORDER BY position DESC
// 		LIMIT 1
// 	`, c.BoardID, c.Position,
// 	)
// 	if err = row.Scan(&prev.ID, &prev.Name, &prev.Position, &prev.CreatedAt); err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			prev = nil
// 		} else {
// 			return nil, nil, db.MapError(err)
// 		}
// 	}
// 	next = &model.Column{BoardID: c.BoardID}
// 	row = client.QueryRowContext(ctx, `
// 		SELECT id, name, position, created_at
// 		FROM columns
// 		WHERE board_id = $1 and position > $2
// 		ORDER BY position ASC
// 		LIMIT 1
// 	`, c.BoardID, c.Position,
// 	)
// 	if err = row.Scan(&next.ID, &next.Name, &next.Position, &next.CreatedAt); err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			next = nil
// 		} else {
// 			return nil, nil, db.MapError(err)
// 		}
// 	}
// 	return prev, next, nil
// }
