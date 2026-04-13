package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ltdlvr/task-manager/internal/models"
)

type UserRepository struct { //NOTE - можно добавить интерфейс для тестов
	Pool *pgxpool.Pool
}

// TODO - добавить валидацию
func (r *UserRepository) CreateUser(ctx context.Context, user models.User) (int, error) {
	var id int
	row := r.Pool.QueryRow(ctx, "INSERT INTO users (name, password) VALUES ($1, $2) RETURNING id", user.Name, user.Password)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
