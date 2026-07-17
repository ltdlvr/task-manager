package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/model"
)

type authUsersRepoStub struct {
	user *model.User
	err  error
}

func (r *authUsersRepoStub) Save(ctx context.Context, client db.DB, u *model.User) error {
	panic("unexpected Save call")
}

func (r *authUsersRepoStub) DeleteByID(ctx context.Context, client db.DB, id uint64) error {
	panic("unexpected DeleteByID call")
}

func (r *authUsersRepoStub) GetByName(ctx context.Context, client db.DB, name string) (*model.User, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.user, nil
}

type authPasswordToolStub struct {
	verifyPassword string
	verifyHash     string
}

func (t *authPasswordToolStub) Hash(password string) (string, error) {
	panic("unexpected Hash call")
}

func (t *authPasswordToolStub) Verify(password string, hash string) bool {
	t.verifyPassword = password
	t.verifyHash = hash
	return false
}

type authTokenToolStub struct{}

func (t *authTokenToolStub) Generate(userID uint64) (string, error) {
	panic("unexpected Generate call")
}

func (t *authTokenToolStub) Parse(token string) (uint64, error) {
	panic("unexpected Parse call")
}

func TestAuthLogInUsesDummyHashForMissingUser(t *testing.T) {
	passwordTool := &authPasswordToolStub{}
	svc := NewAuth(
		&authUsersRepoStub{err: db.ErrEntityNotFound},
		nil,
		passwordTool,
		&authTokenToolStub{},
	)

	_, err := svc.LogIn(context.Background(), &model.User{
		Name:     "unknown",
		Password: "secret",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("LogIn() error = %v, want %v", err, ErrInvalidCredentials)
	}
	if passwordTool.verifyPassword != "secret" {
		t.Fatalf("Verify() password = %q, want %q", passwordTool.verifyPassword, "secret")
	}
	if passwordTool.verifyHash != dummyHash {
		t.Fatalf("Verify() hash = %q, want dummy hash", passwordTool.verifyHash)
	}
}

func TestAuthLogInReturnsInvalidCredentialsForWrongPassword(t *testing.T) {
	passwordTool := &authPasswordToolStub{}
	svc := NewAuth(
		&authUsersRepoStub{user: &model.User{ID: 1, Password: "stored-hash"}},
		nil,
		passwordTool,
		&authTokenToolStub{},
	)

	_, err := svc.LogIn(context.Background(), &model.User{
		Name:     "known",
		Password: "wrong",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("LogIn() error = %v, want %v", err, ErrInvalidCredentials)
	}
	if passwordTool.verifyHash != "stored-hash" {
		t.Fatalf("Verify() hash = %q, want stored hash", passwordTool.verifyHash)
	}
}
