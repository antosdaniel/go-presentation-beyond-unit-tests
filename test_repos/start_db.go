package test_repos

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DBContainer struct {
	DB  *sql.DB
	DSN string
}

func StartDB(t *testing.T, ctx context.Context) DBContainer {
	t.Helper()

	compose, err := tc.NewDockerCompose("./../docker-compose-db-only.yaml")
	require.NoError(t, err, "docker compose setup")

	t.Cleanup(func() {
		compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal)
	})

	err = compose.WaitForService("migrate", wait.ForExit()).Up(ctx)
	require.NoError(t, err, "docker compose up")

	dbContainer, err := compose.ServiceContainer(ctx, "db")
	require.NoError(t, err, "docker compose db container")

	dbPort, err := dbContainer.MappedPort(ctx, "5432")
	require.NoError(t, err, "docker compose db port")

	dsn := fmt.Sprintf("postgres://postgres:secret123@localhost:%s/expense_tracker?sslmode=disable", dbPort.Port())
	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err, "pgx open")
	t.Cleanup(func() {
		db.Close()
	})

	return DBContainer{
		DB:  db,
		DSN: dsn,
	}
}
