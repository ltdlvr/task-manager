package service

import (
	"context"
	"fmt"

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
	hash, err := s.pswdTool.Hash(u.Password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	u.Password = hash
	if err := s.usersRepo.Save(ctx, s.dbClient, u); err != nil {
		return fmt.Errorf("save user: %w", err)
	}
	u.Password = ""
	return nil
}

// or maybe just name? Idk why but it sounds like a bad idea
func (s *Auth) LogIn(ctx context.Context, u *model.User) error {
	dbUser, err := s.usersRepo.GetByName(ctx, s.dbClient, u.Name)
	if err != nil {
		return fmt.Errorf("get user data: %w", err)
	}
	if res := s.pswdTool.Verify(u.Password, dbUser.Password); res != true {
		return fmt.Errorf("Invalid Password")
	}
	return nil
}
