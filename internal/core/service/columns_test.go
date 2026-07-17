package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/model"
)

type boardsRepoStub struct {
	board   *model.Board
	err     error
	clients []db.DB
}

func (r *boardsRepoStub) Create(ctx context.Context, client db.DB, b *model.Board) error {
	panic("unexpected Create call")
}

func (r *boardsRepoStub) DeleteByID(ctx context.Context, client db.DB, id uint64) error {
	panic("unexpected DeleteByID call")
}

func (r *boardsRepoStub) GetByID(ctx context.Context, client db.DB, id uint64) (*model.Board, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.board, nil
}

func (r *boardsRepoStub) LockByID(_ context.Context, client db.DB, _ uint64) error {
	r.clients = append(r.clients, client)
	return nil
}

type columnsRepoStub struct {
	column          *model.Column
	columns         []*model.Column
	otherColumns    []*model.Column
	created         *model.Column
	updatedPosition int
	updateCalled    bool
	clients         []db.DB
}

func (r *columnsRepoStub) Create(ctx context.Context, client db.DB, c *model.Column) error {
	r.clients = append(r.clients, client)
	copy := *c
	r.created = &copy
	return nil
}

func (r *columnsRepoStub) GetByID(ctx context.Context, client db.DB, id uint64) (*model.Column, error) {
	r.clients = append(r.clients, client)
	return r.column, nil
}

func (r *columnsRepoStub) GetAllByBoard(ctx context.Context, client db.DB, boardID uint64) ([]*model.Column, error) {
	r.clients = append(r.clients, client)
	return r.columns, nil
}

func (r *columnsRepoStub) LockByID(_ context.Context, client db.DB, _ uint64) error {
	r.clients = append(r.clients, client)
	return nil
}

func (r *columnsRepoStub) GetOtherByBoard(ctx context.Context, client db.DB, boardID uint64, columnID uint64) ([]*model.Column, error) {
	r.clients = append(r.clients, client)
	return r.otherColumns, nil
}

func (r *columnsRepoStub) DeleteByID(ctx context.Context, client db.DB, id uint64) error {
	panic("unexpected DeleteByID call")
}

func (r *columnsRepoStub) UpdatePosition(ctx context.Context, client db.DB, newPos int, id uint64) error {
	r.clients = append(r.clients, client)
	r.updatedPosition = newPos
	r.updateCalled = true
	return nil
}

type boardMembersRepoStub struct {
	role    model.BoardRole
	err     error
	clients []db.DB
}

func (r *boardMembersRepoStub) Create(ctx context.Context, client db.DB, m *model.BoardMember) error {
	panic("unexpected Create call")
}

func (r *boardMembersRepoStub) GetRole(ctx context.Context, client db.DB, boardID, userID uint64) (model.BoardRole, error) {
	r.clients = append(r.clients, client)
	if r.err != nil {
		return "", r.err
	}
	if r.role == "" {
		return model.BoardRoleMember, nil
	}
	return r.role, nil
}

func TestColumnsCreateCalculatesPositionFromTargetIndex(t *testing.T) {
	repo := &columnsRepoStub{
		columns: []*model.Column{
			{ID: 1, BoardID: 10, Position: 1024},
			{ID: 2, BoardID: 10, Position: 2048},
		},
	}
	boardMembersRepo := &boardMembersRepoStub{}
	boardsRepo := &boardsRepoStub{}
	dbClient := newTransactionClientStub()
	svc := NewColumns(repo, boardsRepo, boardMembersRepo, dbClient)

	col := &model.Column{
		BoardID:  10,
		Name:     "Doing",
		Position: 1,
	}

	if err := svc.Create(context.Background(), 42, col); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if col.Position != 1536 {
		t.Fatalf("col.Position = %d, want 1536", col.Position)
	}
	if repo.created == nil || repo.created.Position != 1536 {
		t.Fatalf("created position = %v, want 1536", repo.created)
	}
	assertClientsUseTransaction(t, dbClient.trx, repo.clients...)
	assertClientsUseTransaction(t, dbClient.trx, boardsRepo.clients...)
	assertClientsUseTransaction(t, dbClient.trx, boardMembersRepo.clients...)
}

func TestColumnsCreateRejectsOutOfRangeTargetIndex(t *testing.T) {
	repo := &columnsRepoStub{
		columns: []*model.Column{
			{ID: 1, BoardID: 10, Position: 1024},
		},
	}
	boardMembersRepo := &boardMembersRepoStub{}
	svc := NewColumns(repo, &boardsRepoStub{}, boardMembersRepo, newTransactionClientStub())

	err := svc.Create(context.Background(), 42, &model.Column{
		BoardID:  10,
		Name:     "Doing",
		Position: 2,
	})
	if err == nil {
		t.Fatal("Create() error = nil, want out-of-range error")
	}
	if !strings.Contains(err.Error(), "out of range") {
		t.Fatalf("Create() error = %v, want out-of-range error", err)
	}
	if repo.created != nil {
		t.Fatalf("created = %v, want nil", repo.created)
	}
}

func TestColumnsCreateRejectsNonMember(t *testing.T) {
	repo := &columnsRepoStub{}
	boardMembersRepo := &boardMembersRepoStub{err: db.ErrEntityNotFound}
	svc := NewColumns(repo, &boardsRepoStub{}, boardMembersRepo, newTransactionClientStub())

	err := svc.Create(context.Background(), 42, &model.Column{
		BoardID:  10,
		Name:     "Doing",
		Position: 0,
	})
	if !errors.Is(err, db.ErrForbidden) {
		t.Fatalf("Create() error = %v, want forbidden", err)
	}
	if repo.created != nil {
		t.Fatalf("created = %v, want nil", repo.created)
	}
}

func TestColumnsMoveColumnUsesTransaction(t *testing.T) {
	columnsRepo := &columnsRepoStub{
		column: &model.Column{ID: 7, BoardID: 10, Position: 1024},
		otherColumns: []*model.Column{
			{ID: 8, BoardID: 10, Position: 1024},
			{ID: 9, BoardID: 10, Position: 3072},
		},
	}
	boardsRepo := &boardsRepoStub{}
	membersRepo := &boardMembersRepoStub{}
	dbClient := newTransactionClientStub()
	svc := NewColumns(columnsRepo, boardsRepo, membersRepo, dbClient)

	if err := svc.MoveColumn(context.Background(), 42, 7, 1); err != nil {
		t.Fatalf("MoveColumn() error = %v", err)
	}
	if !columnsRepo.updateCalled || columnsRepo.updatedPosition != 2048 {
		t.Fatalf("updated position = %d, called = %v; want 2048", columnsRepo.updatedPosition, columnsRepo.updateCalled)
	}
	assertClientsUseTransaction(t, dbClient.trx, columnsRepo.clients...)
	assertClientsUseTransaction(t, dbClient.trx, boardsRepo.clients...)
	assertClientsUseTransaction(t, dbClient.trx, membersRepo.clients...)
}

func TestColumnsGetAllByBoardRejectsMissingBoard(t *testing.T) {
	columnsRepo := &columnsRepoStub{}
	boardsRepo := &boardsRepoStub{err: db.ErrEntityNotFound}
	boardMembersRepo := &boardMembersRepoStub{}
	svc := NewColumns(columnsRepo, boardsRepo, boardMembersRepo, nil)

	columns, err := svc.GetAllByBoard(context.Background(), 10, 42)
	if !errors.Is(err, db.ErrEntityNotFound) {
		t.Fatalf("GetAllByBoard() error = %v, want entity not found", err)
	}
	if columns != nil {
		t.Fatalf("columns = %v, want nil", columns)
	}
}
