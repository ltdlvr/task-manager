package service

import (
	"context"

	"github.com/ltdlvr/task-manager/internal/core/adapter/db"
	"github.com/ltdlvr/task-manager/internal/core/adapter/repo"
	"github.com/ltdlvr/task-manager/internal/core/adapter/tool"
	"github.com/ltdlvr/task-manager/internal/core/model"
)

type Auth struct {
	usersRepo repo.Users
	dbClient  db.Client
	pswdTool  tool.Password
}

func NewAuth(u repo.Users, d db.Client, pt tool.Password) *Auth {
	return &Auth{
		usersRepo: u,
		dbClient:  d,
		pswdTool:  pt,
	}
}

func (s *Auth) Register(ctx context.Context, u *model.User) error {

}
