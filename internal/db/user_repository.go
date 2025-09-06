package db

import (
	"context"
	"go-rest-api/internal/core"
	"go-rest-api/pkg/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db     *pgxpool.Pool
	logger logger.CustomLogger
}

func NewUserRepository(db *pgxpool.Pool, logger logger.CustomLogger) core.UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (u *UserRepository) CreateUser(ctx context.Context, user *core.User) (*core.User, error) {
	const query = `INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)`

	_, err := u.db.Exec(ctx, query,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
	)

	if err != nil {
		u.logger.Error(err, user.ID)
		return nil, err
	}
	// ... retrieve the created user
	user, err = u.GetUserByID(ctx, user.ID.String())
	if err != nil {
		u.logger.Error("failed to get user by id", err, user.ID)
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) GetUserByID(ctx context.Context, id string) (*core.User, error) {
	const query = `SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1`

	user := &User{}
	err := u.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		u.logger.Error("failed to get user by id", err, id)
		return nil, err
	}

	if user.ID == uuid.Nil {
		u.logger.Info("user not found", id)
		return nil, nil
	}

	// Map to core.User
	return &core.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
