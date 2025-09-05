package test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func CreateDbTestContainer(ctx context.Context, t *testing.T) (*pgxpool.Pool, func()) {
	t.Helper()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:14-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get mapped port: %v", err)
	}

	dsn := "postgres://testuser:testpassword@" + host + ":" + port.Port() + "/testdb?sslmode=disable"
	// ... run migrations

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Failed to create connection pool: %v", err)
	}

	cleanup := func() {
		pool.Close()
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("Failed to terminate container: %v", err)
		}
	}

	return pool, cleanup
}
