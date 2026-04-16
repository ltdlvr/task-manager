package service

import (
	"context"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/adapter/repo"
	"github.com/ltdlvr/task-manager/internal/core/model"
)

type Auth struct {
	usersRepo repo.Users
	dbClient  db.Client
}

func NewAuth(u repo.Users, d db.Client) *Auth {
	return &Auth{
		usersRepo: u,
		dbClient:  d,
	}
}

func (s *Auth) Register(ctx context.Context, u *model.User) error {

}
