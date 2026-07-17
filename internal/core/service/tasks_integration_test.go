//go:build integration

package service_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ltdlvr/task-manager/internal/core/service"
	"github.com/ltdlvr/task-manager/internal/infra/db/pg"
	"github.com/ltdlvr/task-manager/internal/infra/repo"
)

func TestMoveTaskIntegration(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx := context.Background()
	client, err := pg.NewClient(databaseURL)
	if err != nil {
		t.Fatalf("connect to test database: %v", err)
	}
	defer client.Close()

	suffix := time.Now().UnixNano()
	var userID, boardID, sourceColumnID, targetColumnID, taskID uint64
	if err := client.QueryRowContext(ctx,
		"INSERT INTO users (name, password) VALUES ($1, $2) RETURNING id",
		fmt.Sprintf("move-task-test-%d", suffix), "test-password",
	).Scan(&userID); err != nil {
		t.Fatalf("create user: %v", err)
	}
	defer client.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)

	if err := client.QueryRowContext(ctx,
		"INSERT INTO boards (name) VALUES ($1) RETURNING id",
		fmt.Sprintf("Move task test %d", suffix),
	).Scan(&boardID); err != nil {
		t.Fatalf("create board: %v", err)
	}
	defer client.ExecContext(ctx, "DELETE FROM boards WHERE id = $1", boardID)

	if _, err := client.ExecContext(ctx,
		"INSERT INTO board_members (board_id, user_id, role) VALUES ($1, $2, 'member')",
		boardID, userID,
	); err != nil {
		t.Fatalf("create board member: %v", err)
	}
	if err := client.QueryRowContext(ctx,
		"INSERT INTO columns (board_id, name, position) VALUES ($1, 'Source', 1024) RETURNING id",
		boardID,
	).Scan(&sourceColumnID); err != nil {
		t.Fatalf("create source column: %v", err)
	}
	if err := client.QueryRowContext(ctx,
		"INSERT INTO columns (board_id, name, position) VALUES ($1, 'Target', 2048) RETURNING id",
		boardID,
	).Scan(&targetColumnID); err != nil {
		t.Fatalf("create target column: %v", err)
	}
	if err := client.QueryRowContext(ctx,
		"INSERT INTO tasks (column_id, title, position) VALUES ($1, 'Moving', 1024) RETURNING id",
		sourceColumnID,
	).Scan(&taskID); err != nil {
		t.Fatalf("create moving task: %v", err)
	}
	if _, err := client.ExecContext(ctx, `
		INSERT INTO tasks (column_id, title, position)
		VALUES ($1, 'First', 1024), ($1, 'Second', 3072)
	`, targetColumnID); err != nil {
		t.Fatalf("create neighbouring tasks: %v", err)
	}

	svc := service.NewTasks(repo.NewTasks(), repo.NewColumns(), repo.NewBoardMembers(), client)
	if err := svc.MoveTask(ctx, userID, taskID, targetColumnID, 1); err != nil {
		t.Fatalf("MoveTask() error = %v", err)
	}

	var actualColumnID uint64
	var actualPosition int
	if err := client.QueryRowContext(ctx,
		"SELECT column_id, position FROM tasks WHERE id = $1", taskID,
	).Scan(&actualColumnID, &actualPosition); err != nil {
		t.Fatalf("load moved task: %v", err)
	}
	if actualColumnID != targetColumnID || actualPosition != 2048 {
		t.Fatalf(
			"moved task = column %d, position %d; want column %d, position 2048",
			actualColumnID, actualPosition, targetColumnID,
		)
	}
}
