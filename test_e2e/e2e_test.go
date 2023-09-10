//go:build e2e_tests

package test_e2e

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

const expenseToSyncID = "677df0c4-d829-42eb-a0c9-29d5b0a2bbe4"

func TestE2E(t *testing.T) {
	ctx := context.Background()

	// At first, spin up mocked bank API, and get its address.
	bankAPIAddress := mockBankAPI(t)
	address := startApp(t, ctx, bankAPIAddress)

	t.Run("app is starting properly", func(t *testing.T) {
		response, err := http.Get(fmt.Sprintf("%s/expenses/all", address))

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode, "status code")
	})
	t.Run("sync expenses", func(t *testing.T) {
		response, err := http.Get(fmt.Sprintf("%s/expenses/sync", address))

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode, "status code")

		// After sync, our app should now store expense from mocked API.
		response, err = http.Get(fmt.Sprintf("%s/expenses/all", address))
		require.NoError(t, err)
		responseBody, err := io.ReadAll(response.Body)
		require.NoError(t, err)
		assert.Contains(t, string(responseBody), expenseToSyncID)
	})
}

func startApp(t *testing.T, ctx context.Context, bankAPIAddress string) (address string) {
	t.Helper()

	compose, err := tc.NewDockerComposeWith(
		tc.WithStackFiles("./docker-compose-for-e2e.yaml"),
		// Giving unique name to each docker compose stack allows us to run tests in parallel.
		tc.StackIdentifier(uuid.New().String()),
	)
	require.NoError(t, err, "docker compose setup")

	t.Cleanup(func() {
		// When test fail, printing logs is usually helpful :)
		if t.Failed() {
			reader, _ := getServerContainer(t, ctx, compose).Logs(ctx)
			bytes, _ := io.ReadAll(reader)
			fmt.Println(`\nLogs from "server" container:\n`, string(bytes))
		}
		assert.NoError(t, compose.Down(ctx, tc.RemoveOrphans(true), tc.RemoveImagesLocal))
	})

	err = compose.
		WithEnv(map[string]string{
			"BANK_API_URL": bankAPIAddress,
		}).
		WaitForService("server", wait.ForLog("running...")).
		Up(ctx)
	require.NoError(t, err, "docker compose up")

	// Port is randomly assigned by docker. We need to get it.
	apiPort, err := getServerContainer(t, ctx, compose).MappedPort(ctx, "8000")
	require.NoError(t, err, "docker compose server port")

	return fmt.Sprintf("http://localhost:%s", apiPort.Port())
}

func getServerContainer(t *testing.T, ctx context.Context, compose tc.ComposeStack) testcontainers.Container {
	t.Helper()

	serverContainer, err := compose.ServiceContainer(ctx, "server")
	require.NoError(t, err, "docker compose server container")

	return serverContainer
}

func mockBankAPI(t *testing.T) (address string) {
	t.Helper()

	// We are getting random available port.
	// Using 0.0.0.0 address if preferable over localhost, as the former will usually not work in CI.
	// Oh, and Desktop Docker doesn't support IPv6 yet, so it's better to specify "tcp4" network.
	listener, err := net.Listen("tcp4", "0.0.0.0:0")
	require.NoError(t, err, "could not start listener")
	addr, err := net.ResolveTCPAddr(listener.Addr().Network(), listener.Addr().String())
	require.NoError(t, err, "could not resolve tcp addr")

	// Setup mocked API.
	mux := http.NewServeMux()
	mux.Handle("/get-transactions", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf(`
			[
			  {
				"id": "%s",
				"amount": 500.00,
				"category": "food",
				"created_at": "2020-01-01T00:00:00Z"
			  }
			]`,
			expenseToSyncID)))
	}))

	server := httptest.NewUnstartedServer(mux)
	server.Listener = listener
	server.Start()
	t.Cleanup(func() {
		server.Close()
	})

	// "host.docker.internal" resolves to host network. Thanks to this we can access our HTTP mocks from container.
	return fmt.Sprintf("http://host.docker.internal:%d", addr.Port)
}
