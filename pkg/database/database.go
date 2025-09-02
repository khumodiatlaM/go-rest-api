package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-rest-api/pkg/logger"
	"time"
)

// Connect establishes a connection to the PostgreSQL database
func Connect(cfg DbConfig, logger logger.Logger) (*pgxpool.Pool, error) {
	logger.Info("Connecting to the database...")
	connConfig, err := pgxpool.ParseConfig(cfg.GetDBConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database connection string: %w", err)
	}

	// Adjust connection pool settings
	connConfig.MaxConns = 5
	connConfig.MinConns = 1
	connConfig.MaxConnIdleTime = 2 * time.Minute // 2 minutes

	dbpool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection pool: %w", err)
	}

	// Ping the database to ensure the connection is established
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := dbpool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established")
	return dbpool, nil
}

// GetDBConnectionString constructs the database connection string
func (cfg *DbConfig) GetDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
}
