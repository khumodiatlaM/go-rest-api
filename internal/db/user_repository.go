package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-rest-api/internal/core"
	"go-rest-api/pkg/logger"
)

type UserRepository struct {
	db     *pgxpool.Pool
	logger logger.Logger
}

func NewUserRepository(db *pgxpool.Pool, logger logger.Logger) core.UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (u *UserRepository) CreateUser(ctx context.Context, user *core.User) error {
	const query = `INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3)`

	_, err := u.db.Exec(ctx, query,
		user.ID,
		user.Username,
		user.Email,
		user.Password,
	)

	if err != nil {
		u.logger.Error(err, user.ID)
		return err
	}
	return nil
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

	// Map to core.User
	return &core.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
