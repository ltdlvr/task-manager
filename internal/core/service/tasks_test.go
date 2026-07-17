package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/model"
)

type tasksRepoStub struct {
	task            *model.Task
	otherTasks      []*model.Task
	created         *model.Task
	updatedColumnID uint64
	updatedTaskID   uint64
	updatedPosition int
	updateCalled    bool
	clients         []db.DB
}

func (r *tasksRepoStub) Create(_ context.Context, client db.DB, task *model.Task) error {
	r.clients = append(r.clients, client)
	copy := *task
	r.created = &copy
	return nil
}

func (r *tasksRepoStub) GetByID(context.Context, db.DB, uint64) (*model.Task, error) {
	return r.task, nil
}

func (r *tasksRepoStub) GetByIDForUpdate(_ context.Context, client db.DB, _ uint64) (*model.Task, error) {
	r.clients = append(r.clients, client)
	return r.task, nil
}

func (r *tasksRepoStub) GetAllByColumn(_ context.Context, client db.DB, _ uint64) ([]*model.Task, error) {
	r.clients = append(r.clients, client)
	return r.otherTasks, nil
}

func (r *tasksRepoStub) GetOtherByColumn(_ context.Context, client db.DB, _, _ uint64) ([]*model.Task, error) {
	r.clients = append(r.clients, client)
	return r.otherTasks, nil
}

func (r *tasksRepoStub) DeleteByID(context.Context, db.DB, uint64) error {
	panic("unexpected DeleteByID call")
}

func (r *tasksRepoStub) UpdateColumnAndPosition(_ context.Context, client db.DB, columnID, taskID uint64, position int) error {
	r.clients = append(r.clients, client)
	r.updatedColumnID = columnID
	r.updatedTaskID = taskID
	r.updatedPosition = position
	r.updateCalled = true
	return nil
}

type taskColumnsRepoStub struct {
	column  *model.Column
	columns map[uint64]*model.Column
	clients []db.DB
	locked  []uint64
}

func (r *taskColumnsRepoStub) Create(context.Context, db.DB, *model.Column) error {
	panic("unexpected Create call")
}

func (r *taskColumnsRepoStub) GetByID(_ context.Context, client db.DB, id uint64) (*model.Column, error) {
	r.clients = append(r.clients, client)
	if r.columns != nil {
		return r.columns[id], nil
	}
	return r.column, nil
}

func (r *taskColumnsRepoStub) LockByID(_ context.Context, client db.DB, id uint64) error {
	r.clients = append(r.clients, client)
	r.locked = append(r.locked, id)
	return nil
}

func (r *taskColumnsRepoStub) GetAllByBoard(context.Context, db.DB, uint64) ([]*model.Column, error) {
	panic("unexpected GetAllByBoard call")
}

func (r *taskColumnsRepoStub) GetOtherByBoard(context.Context, db.DB, uint64, uint64) ([]*model.Column, error) {
	panic("unexpected GetOtherByBoard call")
}

func (r *taskColumnsRepoStub) DeleteByID(context.Context, db.DB, uint64) error {
	panic("unexpected DeleteByID call")
}

func (r *taskColumnsRepoStub) UpdatePosition(context.Context, db.DB, int, uint64) error {
	panic("unexpected UpdatePosition call")
}

func TestTasksCreateUsesTransaction(t *testing.T) {
	tasksRepo := &tasksRepoStub{otherTasks: []*model.Task{
		{ID: 8, ColumnID: 5, Position: 1024},
		{ID: 9, ColumnID: 5, Position: 3072},
	}}
	columnsRepo := &taskColumnsRepoStub{columns: map[uint64]*model.Column{
		5: {ID: 5, BoardID: 11},
	}}
	membersRepo := &boardMembersRepoStub{}
	dbClient := newTransactionClientStub()
	svc := NewTasks(tasksRepo, columnsRepo, membersRepo, dbClient)
	task := &model.Task{ColumnID: 5, Title: "New", Position: 1}

	if err := svc.Create(context.Background(), 42, task); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if tasksRepo.created == nil || tasksRepo.created.Position != 2048 {
		t.Fatalf("created task = %#v, want position 2048", tasksRepo.created)
	}
	assertClientsUseTransaction(t, dbClient.trx, tasksRepo.clients...)
	assertClientsUseTransaction(t, dbClient.trx, columnsRepo.clients...)
	assertClientsUseTransaction(t, dbClient.trx, membersRepo.clients...)
}

func TestTasksMoveTaskLoadsNeighboursAndUpdatesPosition(t *testing.T) {
	tasksRepo := &tasksRepoStub{
		task: &model.Task{ID: 7, ColumnID: 3},
		otherTasks: []*model.Task{
			{ID: 8, ColumnID: 5, Position: 1024},
			{ID: 9, ColumnID: 5, Position: 3072},
		},
	}
	columnsRepo := &taskColumnsRepoStub{columns: map[uint64]*model.Column{
		3: {ID: 3, BoardID: 11},
		5: {ID: 5, BoardID: 11},
	}}
	membersRepo := &boardMembersRepoStub{}
	dbClient := newTransactionClientStub()
	svc := NewTasks(tasksRepo, columnsRepo, membersRepo, dbClient)

	if err := svc.MoveTask(context.Background(), 42, 7, 5, 1); err != nil {
		t.Fatalf("MoveTask() error = %v", err)
	}

	if tasksRepo.updatedColumnID != 5 || tasksRepo.updatedTaskID != 7 || tasksRepo.updatedPosition != 2048 {
		t.Fatalf(
			"UpdateColumnAndPosition() = column %d, task %d, position %d; want column 5, task 7, position 2048",
			tasksRepo.updatedColumnID,
			tasksRepo.updatedTaskID,
			tasksRepo.updatedPosition,
		)
	}
	assertClientsUseTransaction(t, dbClient.trx, tasksRepo.clients...)
	assertClientsUseTransaction(t, dbClient.trx, columnsRepo.clients...)
	assertClientsUseTransaction(t, dbClient.trx, membersRepo.clients...)
}

func TestTasksMoveTaskRejectsNonMemberOfSourceBoard(t *testing.T) {
	tasksRepo := &tasksRepoStub{task: &model.Task{ID: 7, ColumnID: 3}}
	columnsRepo := &taskColumnsRepoStub{columns: map[uint64]*model.Column{
		3: {ID: 3, BoardID: 11},
	}}
	membersRepo := &boardMembersRepoStub{err: db.ErrEntityNotFound}
	svc := NewTasks(tasksRepo, columnsRepo, membersRepo, newTransactionClientStub())

	err := svc.MoveTask(context.Background(), 42, 7, 5, 0)
	if !errors.Is(err, db.ErrForbidden) {
		t.Fatalf("MoveTask() error = %v, want forbidden", err)
	}
	if tasksRepo.updateCalled {
		t.Fatal("UpdateColumnAndPosition() called for a non-member")
	}
}

func TestTasksMoveTaskRejectsCrossBoardMove(t *testing.T) {
	tasksRepo := &tasksRepoStub{task: &model.Task{ID: 7, ColumnID: 3}}
	columnsRepo := &taskColumnsRepoStub{columns: map[uint64]*model.Column{
		3: {ID: 3, BoardID: 11},
		5: {ID: 5, BoardID: 12},
	}}
	svc := NewTasks(tasksRepo, columnsRepo, &boardMembersRepoStub{}, newTransactionClientStub())

	err := svc.MoveTask(context.Background(), 42, 7, 5, 0)
	if !errors.Is(err, db.ErrForbidden) {
		t.Fatalf("MoveTask() error = %v, want forbidden", err)
	}
	if tasksRepo.updateCalled {
		t.Fatal("UpdateColumnAndPosition() called for a cross-board move")
	}
}

func TestTasksLockColumnsUsesAscendingIDOrder(t *testing.T) {
	columnsRepo := &taskColumnsRepoStub{}
	svc := NewTasks(nil, columnsRepo, nil, nil)
	trx := &transactionDBStub{}

	if err := svc.lockColumns(context.Background(), trx, 5, 3); err != nil {
		t.Fatalf("lockColumns() error = %v", err)
	}
	if len(columnsRepo.locked) != 2 || columnsRepo.locked[0] != 3 || columnsRepo.locked[1] != 5 {
		t.Fatalf("locked columns = %v, want [3 5]", columnsRepo.locked)
	}
}
