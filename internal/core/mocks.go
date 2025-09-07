package core

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// ---------------------------------
// MockUserRepository
// ---------------------------------

type MockUserRepository struct {
	mock.Mock
}

func (r *MockUserRepository) CreateUser(ctx context.Context, user *User) (*User, error) {
	args := r.Called(ctx, user)
	return args.Get(0).(*User), args.Error(1)
}

func (r *MockUserRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
	args := r.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (r *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	args := r.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

// ---------------------------------
// MockUserService
// ---------------------------------

type MockUserService struct {
	mock.Mock
}

func (s *MockUserService) CreateUser(ctx context.Context, user *User) error {
	args := s.Called(ctx, user)
	return args.Error(0)
}

func (s *MockUserService) GetUserByID(ctx context.Context, id string) (*User, error) {
	args := s.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (s *MockUserService) LoginUserRequest(ctx context.Context, email, password string) (*User, error) {
	args := s.Called(ctx, email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

// ---------------------------------
// MockUserEventService
// ---------------------------------

type MockUserEventService struct {
	mock.Mock
}

func (s *MockUserEventService) PublishUserCreatedEvent(ctx context.Context, user *User) error {
	args := s.Called(ctx, user)
	return args.Error(0)
}
