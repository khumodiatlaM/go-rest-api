package core

import (
	"context"
	"go-rest-api/pkg/logger"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type UserService struct {
	repo   UserRepository
	logger logger.CustomLogger
}

func NewUserService(repo UserRepository, logger logger.CustomLogger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *User) (*User, error) {
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		s.logger.Error("failed to hash password: ", err)
		return nil, err
	}
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Password = hashedPassword
	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error("failed to get user by email for authentication: ", err)
		return "", err
	}

	if user == nil {
		s.logger.Error("user not found with email: ", email)
		return "", nil
	}

	if err = VerifyPassword(user.Password, password); err != nil {
		s.logger.Error("password verification failed: ", err)
		return "", err
	}

	token, err := GenerateAuthToken(user.ID, "jwt-secret")
	if err != nil {
		s.logger.Error("failed to generate auth token: ", err)
		return "", err
	}

	return token, nil
}
