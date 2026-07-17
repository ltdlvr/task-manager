//go:build integration

package service_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ltdlvr/task-manager/internal/core/model"
	"github.com/ltdlvr/task-manager/internal/core/service"
	"github.com/ltdlvr/task-manager/internal/infra/db/pg"
	"github.com/ltdlvr/task-manager/internal/infra/repo"
)

type orderingFixture struct {
	ctx     context.Context
	client  *pg.Client
	userID  uint64
	boardID uint64
}

func newOrderingFixture(t *testing.T) *orderingFixture {
	t.Helper()
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx := context.Background()
	client, err := pg.NewClient(databaseURL)
	if err != nil {
		t.Fatalf("connect to test database: %v", err)
	}

	suffix := time.Now().UnixNano()
	fixture := &orderingFixture{ctx: ctx, client: client}
	if err := client.QueryRowContext(ctx,
		"INSERT INTO users (name, password) VALUES ($1, $2) RETURNING id",
		fmt.Sprintf("ordering-test-%d", suffix), "test-password",
	).Scan(&fixture.userID); err != nil {
		client.Close()
		t.Fatalf("create user: %v", err)
	}
	if err := client.QueryRowContext(ctx,
		"INSERT INTO boards (name) VALUES ($1) RETURNING id",
		fmt.Sprintf("Ordering test %d", suffix),
	).Scan(&fixture.boardID); err != nil {
		client.ExecContext(ctx, "DELETE FROM users WHERE id = $1", fixture.userID)
		client.Close()
		t.Fatalf("create board: %v", err)
	}
	if _, err := client.ExecContext(ctx,
		"INSERT INTO board_members (board_id, user_id, role) VALUES ($1, $2, 'member')",
		fixture.boardID, fixture.userID,
	); err != nil {
		client.ExecContext(ctx, "DELETE FROM boards WHERE id = $1", fixture.boardID)
		client.ExecContext(ctx, "DELETE FROM users WHERE id = $1", fixture.userID)
		client.Close()
		t.Fatalf("create board member: %v", err)
	}

	t.Cleanup(func() {
		client.ExecContext(ctx, "DELETE FROM boards WHERE id = $1", fixture.boardID)
		client.ExecContext(ctx, "DELETE FROM users WHERE id = $1", fixture.userID)
		client.Close()
	})
	return fixture
}

func runConcurrently(t *testing.T, count int, fn func(int) error) {
	t.Helper()
	start := make(chan struct{})
	errs := make(chan error, count)
	var wg sync.WaitGroup
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			<-start
			errs <- fn(index)
		}(i)
	}
	close(start)
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent operation: %v", err)
		}
	}
}

func assertDistinctPositions(t *testing.T, fixture *orderingFixture, table, parentColumn string, parentID uint64, want int) {
	t.Helper()
	query := fmt.Sprintf(
		"SELECT COUNT(*), COUNT(DISTINCT position) FROM %s WHERE %s = $1",
		table, parentColumn,
	)
	var count, distinct int
	if err := fixture.client.QueryRowContext(fixture.ctx, query, parentID).Scan(&count, &distinct); err != nil {
		t.Fatalf("count positions: %v", err)
	}
	if count != want || distinct != want {
		t.Fatalf("positions: count = %d, distinct = %d; want %d distinct rows", count, distinct, want)
	}
}

func TestConcurrentColumnCreateUsesDistinctPositions(t *testing.T) {
	fixture := newOrderingFixture(t)
	svc := service.NewColumns(repo.NewColumns(), repo.NewBoards(), repo.NewBoardMembers(), fixture.client)
	const workers = 8

	runConcurrently(t, workers, func(index int) error {
		return svc.Create(fixture.ctx, fixture.userID, &model.Column{
			BoardID:  fixture.boardID,
			Name:     fmt.Sprintf("Column %d", index),
			Position: 0,
		})
	})

	assertDistinctPositions(t, fixture, "columns", "board_id", fixture.boardID, workers)
}

func TestConcurrentTaskCreateUsesDistinctPositions(t *testing.T) {
	fixture := newOrderingFixture(t)
	var columnID uint64
	if err := fixture.client.QueryRowContext(fixture.ctx,
		"INSERT INTO columns (board_id, name, position) VALUES ($1, 'Tasks', 1024) RETURNING id",
		fixture.boardID,
	).Scan(&columnID); err != nil {
		t.Fatalf("create column: %v", err)
	}
	svc := service.NewTasks(repo.NewTasks(), repo.NewColumns(), repo.NewBoardMembers(), fixture.client)
	const workers = 8

	runConcurrently(t, workers, func(index int) error {
		return svc.Create(fixture.ctx, fixture.userID, &model.Task{
			ColumnID: columnID,
			Title:    fmt.Sprintf("Task %d", index),
			Position: 0,
		})
	})

	assertDistinctPositions(t, fixture, "tasks", "column_id", columnID, workers)
}

func TestConcurrentTaskMoveUsesDistinctPositions(t *testing.T) {
	fixture := newOrderingFixture(t)
	var columnID uint64
	if err := fixture.client.QueryRowContext(fixture.ctx,
		"INSERT INTO columns (board_id, name, position) VALUES ($1, 'Moves', 1024) RETURNING id",
		fixture.boardID,
	).Scan(&columnID); err != nil {
		t.Fatalf("create column: %v", err)
	}

	rows, err := fixture.client.QueryContext(fixture.ctx, `
		INSERT INTO tasks (column_id, title, position)
		VALUES ($1, 'Anchor', 1024), ($1, 'First move', 2048), ($1, 'Second move', 3072)
		RETURNING id
	`, columnID)
	if err != nil {
		t.Fatalf("create tasks: %v", err)
	}
	defer rows.Close()
	var taskIDs []uint64
	for rows.Next() {
		var id uint64
		if err := rows.Scan(&id); err != nil {
			t.Fatalf("scan task id: %v", err)
		}
		taskIDs = append(taskIDs, id)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("read task ids: %v", err)
	}
	if len(taskIDs) != 3 {
		t.Fatalf("created task IDs = %v, want 3 tasks", taskIDs)
	}

	svc := service.NewTasks(repo.NewTasks(), repo.NewColumns(), repo.NewBoardMembers(), fixture.client)
	runConcurrently(t, 2, func(index int) error {
		return svc.MoveTask(fixture.ctx, fixture.userID, taskIDs[index+1], columnID, 0)
	})

	assertDistinctPositions(t, fixture, "tasks", "column_id", columnID, 3)
}
