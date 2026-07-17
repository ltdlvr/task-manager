package rest

import (
	"bytes"
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/model"
	"github.com/ltdlvr/task-manager/internal/core/service"
)

type createColumnsRepoStub struct {
	created *model.Column
}

func (r *createColumnsRepoStub) Create(_ context.Context, _ db.DB, column *model.Column) error {
	copy := *column
	r.created = &copy
	return nil
}

func (*createColumnsRepoStub) GetByID(context.Context, db.DB, uint64) (*model.Column, error) {
	panic("unexpected GetByID call")
}

func (*createColumnsRepoStub) LockByID(context.Context, db.DB, uint64) error {
	panic("unexpected LockByID call")
}

func (*createColumnsRepoStub) GetAllByBoard(context.Context, db.DB, uint64) ([]*model.Column, error) {
	return nil, nil
}

func (*createColumnsRepoStub) GetOtherByBoard(context.Context, db.DB, uint64, uint64) ([]*model.Column, error) {
	panic("unexpected GetOtherByBoard call")
}

func (*createColumnsRepoStub) DeleteByID(context.Context, db.DB, uint64) error {
	panic("unexpected DeleteByID call")
}

func (*createColumnsRepoStub) UpdatePosition(context.Context, db.DB, int, uint64) error {
	panic("unexpected UpdatePosition call")
}

type createTasksRepoStub struct {
	created *model.Task
}

func (r *createTasksRepoStub) Create(_ context.Context, _ db.DB, task *model.Task) error {
	copy := *task
	r.created = &copy
	return nil
}

func (*createTasksRepoStub) GetByID(context.Context, db.DB, uint64) (*model.Task, error) {
	panic("unexpected GetByID call")
}

func (*createTasksRepoStub) GetByIDForUpdate(context.Context, db.DB, uint64) (*model.Task, error) {
	panic("unexpected GetByIDForUpdate call")
}

func (*createTasksRepoStub) GetAllByColumn(context.Context, db.DB, uint64) ([]*model.Task, error) {
	return nil, nil
}

func (*createTasksRepoStub) GetOtherByColumn(context.Context, db.DB, uint64, uint64) ([]*model.Task, error) {
	panic("unexpected GetOtherByColumn call")
}

func (*createTasksRepoStub) DeleteByID(context.Context, db.DB, uint64) error {
	panic("unexpected DeleteByID call")
}

func (*createTasksRepoStub) UpdateColumnAndPosition(context.Context, db.DB, uint64, uint64, int) error {
	panic("unexpected UpdateColumnAndPosition call")
}

type createTaskColumnsRepoStub struct {
	requestedColumnID uint64
}

func (*createTaskColumnsRepoStub) Create(context.Context, db.DB, *model.Column) error {
	panic("unexpected Create call")
}

func (r *createTaskColumnsRepoStub) GetByID(_ context.Context, _ db.DB, id uint64) (*model.Column, error) {
	r.requestedColumnID = id
	return &model.Column{ID: id, BoardID: 77}, nil
}

func (*createTaskColumnsRepoStub) LockByID(context.Context, db.DB, uint64) error {
	return nil
}

func (*createTaskColumnsRepoStub) GetAllByBoard(context.Context, db.DB, uint64) ([]*model.Column, error) {
	panic("unexpected GetAllByBoard call")
}

func (*createTaskColumnsRepoStub) GetOtherByBoard(context.Context, db.DB, uint64, uint64) ([]*model.Column, error) {
	panic("unexpected GetOtherByBoard call")
}

func (*createTaskColumnsRepoStub) DeleteByID(context.Context, db.DB, uint64) error {
	panic("unexpected DeleteByID call")
}

func (*createTaskColumnsRepoStub) UpdatePosition(context.Context, db.DB, int, uint64) error {
	panic("unexpected UpdatePosition call")
}

type createMembersRepoStub struct {
	requestedBoardID uint64
}

type createBoardsRepoStub struct{}

func (*createBoardsRepoStub) Create(context.Context, db.DB, *model.Board) error {
	panic("unexpected Create call")
}

func (*createBoardsRepoStub) DeleteByID(context.Context, db.DB, uint64) error {
	panic("unexpected DeleteByID call")
}

func (*createBoardsRepoStub) GetByID(context.Context, db.DB, uint64) (*model.Board, error) {
	panic("unexpected GetByID call")
}

func (*createBoardsRepoStub) LockByID(context.Context, db.DB, uint64) error {
	return nil
}

type createTransactionStub struct{}

func (*createTransactionStub) QueryContext(context.Context, string, ...any) (*sql.Rows, error) {
	panic("unexpected QueryContext call")
}

func (*createTransactionStub) QueryRowContext(context.Context, string, ...any) *sql.Row {
	panic("unexpected QueryRowContext call")
}

func (*createTransactionStub) ExecContext(context.Context, string, ...any) (sql.Result, error) {
	panic("unexpected ExecContext call")
}

type createDBClientStub struct {
	trx db.Trx
}

func newCreateDBClientStub() *createDBClientStub {
	return &createDBClientStub{trx: &createTransactionStub{}}
}

func (*createDBClientStub) QueryContext(context.Context, string, ...any) (*sql.Rows, error) {
	panic("unexpected QueryContext call")
}

func (*createDBClientStub) QueryRowContext(context.Context, string, ...any) *sql.Row {
	panic("unexpected QueryRowContext call")
}

func (*createDBClientStub) ExecContext(context.Context, string, ...any) (sql.Result, error) {
	panic("unexpected ExecContext call")
}

func (c *createDBClientStub) ExecTrx(ctx context.Context, fn func(db.Trx) error) error {
	return fn(c.trx)
}

func (*createMembersRepoStub) Create(context.Context, db.DB, *model.BoardMember) error {
	panic("unexpected Create call")
}

func (r *createMembersRepoStub) GetRole(_ context.Context, _ db.DB, boardID, _ uint64) (model.BoardRole, error) {
	r.requestedBoardID = boardID
	return model.BoardRoleMember, nil
}

func testCreateApp(path string, handler fiber.Handler) *fiber.App {
	app := fiber.New()
	app.Post(path, func(c fiber.Ctx) error {
		c.Locals(userIDLocalKey, uint64(42))
		return c.Next()
	}, handler)
	return app
}

func performJSONRequest(t *testing.T, app *fiber.App, path, body string) *http.Response {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	return res
}

func TestCreateColumnIgnoresBoardIDFromBody(t *testing.T) {
	columnsRepo := &createColumnsRepoStub{}
	membersRepo := &createMembersRepoStub{}
	handler := NewColumns(service.NewColumns(
		columnsRepo,
		&createBoardsRepoStub{},
		membersRepo,
		newCreateDBClientStub(),
	))
	app := testCreateApp("/boards/:boardId/columns", handler.Create)

	res := performJSONRequest(t, app, "/boards/10/columns", `{"boardId":999,"name":"Doing","targetPos":0}`)
	defer res.Body.Close()

	if res.StatusCode != fiber.StatusCreated {
		t.Fatalf("status = %d, want %d", res.StatusCode, fiber.StatusCreated)
	}
	if membersRepo.requestedBoardID != 10 {
		t.Fatalf("authorized board ID = %d, want URI board ID 10", membersRepo.requestedBoardID)
	}
	if columnsRepo.created == nil || columnsRepo.created.BoardID != 10 {
		t.Fatalf("created column = %#v, want board ID 10", columnsRepo.created)
	}
}

func TestCreateTaskIgnoresColumnIDFromBody(t *testing.T) {
	tasksRepo := &createTasksRepoStub{}
	columnsRepo := &createTaskColumnsRepoStub{}
	membersRepo := &createMembersRepoStub{}
	handler := NewTasks(service.NewTasks(tasksRepo, columnsRepo, membersRepo, newCreateDBClientStub()))
	app := testCreateApp("/columns/:columnId/tasks", handler.Create)

	res := performJSONRequest(t, app, "/columns/10/tasks", `{"columnId":999,"title":"Fix binding","description":"test","targetPos":0}`)
	defer res.Body.Close()

	if res.StatusCode != fiber.StatusCreated {
		t.Fatalf("status = %d, want %d", res.StatusCode, fiber.StatusCreated)
	}
	if columnsRepo.requestedColumnID != 10 {
		t.Fatalf("loaded column ID = %d, want URI column ID 10", columnsRepo.requestedColumnID)
	}
	if membersRepo.requestedBoardID != 77 {
		t.Fatalf("authorized board ID = %d, want board ID 77 from URI column", membersRepo.requestedBoardID)
	}
	if tasksRepo.created == nil || tasksRepo.created.ColumnID != 10 {
		t.Fatalf("created task = %#v, want column ID 10", tasksRepo.created)
	}
}
