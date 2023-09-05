package test_repos

import (
	"context"
	"database/sql"
	"fmt"
	"path"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DBContainer struct {
	DB  *sql.DB
	DSN string
}

// You will most likely have multiple setups for different test scenarios, depending on your stack: SQL database, redis,
// kafka etc. It's a good idea to share this setup across all tests that need them.

func StartDB(t *testing.T, ctx context.Context) DBContainer {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		// Using `require.NoError()`, or `t.Fatalf()`, keeps test setup cleaner - no ugly error returning.
		t.Fatalf("could not get current file")
	}
	compose, err := tc.NewDockerCompose(path.Join(path.Dir(filename), "docker-compose-db-only.yaml"))
	require.NoError(t, err, "docker compose setup")

	t.Cleanup(func() {
		compose.Down(ctx, tc.RemoveOrphans(true), tc.RemoveImagesLocal)
	})

	// We want to wait for the migrations to finish. In our case, this means exiting of the migrate container.
	err = compose.WaitForService("migrate", wait.ForExit()).Up(ctx)
	require.NoError(t, err, "docker compose up")

	dbContainer, err := compose.ServiceContainer(ctx, "db")
	require.NoError(t, err, "docker compose db container")

	// Port is randomly assigned by docker. We need to get it.
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
