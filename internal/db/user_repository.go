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
		u.logger.Error(err)
		return nil, err
	}
	createdUser := &User{}
	// Fetch the created user
	getUserQuery := `SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1`
	err = u.db.QueryRow(ctx, getUserQuery, user.ID).Scan(
		&createdUser.ID,
		&createdUser.Username,
		&createdUser.Email,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)

	if err != nil {
		u.logger.Error("failed to fetch created user", err, user.ID)
		return nil, err
	}

	// Map to core.User
	*user = core.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
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
