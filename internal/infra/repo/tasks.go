package repo

import (
	"context"
	"fmt"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/model"
)

type Tasks struct{}

func NewTasks() *Tasks {
	return &Tasks{}
}

func (r *Tasks) Create(ctx context.Context, client db.DB, t *model.Task) error {
	row := client.QueryRowContext(
		ctx,
		"INSERT INTO tasks (column_id, title, description, position) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
		t.ColumnID, t.Title, t.Description, t.Position,
	)
	if err := row.Scan(&t.ID, &t.CreatedAt); err != nil {
		return db.MapError(err)
	}
	return nil
}

func (r *Tasks) GetByID(ctx context.Context, client db.DB, id uint64) (*model.Task, error) {
	return r.getByID(ctx, client, id, false)
}

func (r *Tasks) GetByIDForUpdate(ctx context.Context, client db.DB, id uint64) (*model.Task, error) {
	return r.getByID(ctx, client, id, true)
}

func (r *Tasks) getByID(ctx context.Context, client db.DB, id uint64, forUpdate bool) (*model.Task, error) {
	query := "SELECT column_id, title, description, position, created_at FROM tasks WHERE id = $1"
	if forUpdate {
		query += " FOR UPDATE"
	}
	row := client.QueryRowContext(
		ctx, query, id,
	)
	c := &model.Task{
		ID: id,
	}

	if err := row.Scan(&c.ColumnID, &c.Title, &c.Description, &c.Position, &c.CreatedAt); err != nil {
		return nil, db.MapError(err)
	}

	return c, nil
}

func (r *Tasks) GetAllByColumn(ctx context.Context, client db.DB, columnID uint64) ([]*model.Task, error) {
	rows, err := client.QueryContext(ctx, `
		SELECT id, title, "description", position, created_at
		FROM tasks
		WHERE column_id = $1
		ORDER BY position ASC, id ASC
	`, columnID)
	if err != nil {
		return nil, db.MapError(err)
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		c := &model.Task{ColumnID: columnID}
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Position, &c.CreatedAt); err != nil {
			return nil, db.MapError(err)
		}
		tasks = append(tasks, c)
	}
	if err := rows.Err(); err != nil {
		return nil, db.MapError(err)
	}
	return tasks, nil
}

func (r *Tasks) GetOtherByColumn(ctx context.Context, client db.DB, columnID, taskID uint64) ([]*model.Task, error) {
	rows, err := client.QueryContext(ctx, `
		SELECT id, title, "description", position, created_at
		FROM tasks
		WHERE column_id = $1 AND id != $2
		ORDER BY position ASC, id ASC
	`, columnID, taskID)
	if err != nil {
		return nil, db.MapError(err)
	}
	defer rows.Close()

	tasks := make([]*model.Task, 0)
	for rows.Next() {
		c := &model.Task{ColumnID: columnID}
		if err := rows.Scan(&c.ID, &c.Title, &c.Description, &c.Position, &c.CreatedAt); err != nil {
			return nil, db.MapError(err)
		}
		tasks = append(tasks, c)
	}
	if err := rows.Err(); err != nil {
		return nil, db.MapError(err)
	}
	return tasks, nil
}

// NOTE - можно переделать с учетом RowsAffected чтобы отлавливать удаление несуществующего
func (r *Tasks) DeleteByID(ctx context.Context, client db.DB, id uint64) error {
	_, err := client.ExecContext(
		ctx, "DELETE FROM tasks WHERE id = $1", id,
	)
	return db.MapError(err)
}

func (r *Tasks) UpdateColumnAndPosition(ctx context.Context, client db.DB, columnID, taskID uint64, newPos int) error {
	res, err := client.ExecContext(ctx, `
		UPDATE tasks
		SET column_id = $1, position = $2
		WHERE id = $3
	`, columnID, newPos, taskID)
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
