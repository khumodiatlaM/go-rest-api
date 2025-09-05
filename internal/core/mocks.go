package core

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (r *MockUserRepository) CreateUser(ctx context.Context, user *User) error {
	args := r.Called(ctx, user)
	return args.Error(0)
}

func (r *MockUserRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
	args := r.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}
