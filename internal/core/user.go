package core

import (
	"context"
	"go-rest-api/pkg/logger"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id string) (*User, error)
}

type UserService struct {
	repo   UserRepository
	logger logger.Logger
}

func NewUserService(repo UserRepository, logger logger.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}
