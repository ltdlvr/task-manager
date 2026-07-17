package service

import (
	"context"
	"database/sql"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
)

type transactionDBStub struct{}

func (*transactionDBStub) QueryContext(context.Context, string, ...any) (*sql.Rows, error) {
	panic("unexpected QueryContext call")
}

func (*transactionDBStub) QueryRowContext(context.Context, string, ...any) *sql.Row {
	panic("unexpected QueryRowContext call")
}

func (*transactionDBStub) ExecContext(context.Context, string, ...any) (sql.Result, error) {
	panic("unexpected ExecContext call")
}

type transactionClientStub struct {
	trx db.Trx
}

func newTransactionClientStub() *transactionClientStub {
	return &transactionClientStub{trx: &transactionDBStub{}}
}

func (*transactionClientStub) QueryContext(context.Context, string, ...any) (*sql.Rows, error) {
	panic("unexpected QueryContext call")
}

func (*transactionClientStub) QueryRowContext(context.Context, string, ...any) *sql.Row {
	panic("unexpected QueryRowContext call")
}

func (*transactionClientStub) ExecContext(context.Context, string, ...any) (sql.Result, error) {
	panic("unexpected ExecContext call")
}

func (c *transactionClientStub) ExecTrx(ctx context.Context, fn func(db.Trx) error) error {
	return fn(c.trx)
}

func assertClientsUseTransaction(t interface {
	Helper()
	Fatalf(string, ...any)
}, want db.Trx, clients ...db.DB) {
	t.Helper()
	for i, client := range clients {
		if client != want {
			t.Fatalf("repo client[%d] = %T %p, want transaction %T %p", i, client, client, want, want)
		}
	}
}
