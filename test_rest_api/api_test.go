package test_rest_api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/antosdaniel/go-presentation-beyond-unit-tests/app_to_test/server/api"
	"github.com/antosdaniel/go-presentation-beyond-unit-tests/test_repos"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	server := startServer(t, ctx)

	t.Run("summarize expenses", func(t *testing.T) {
		// Notice that if we moved this test after add expense test, this test would fail.
		// Usually it's better to reset the database before each test.
		response, responseBody := call(t, server, http.MethodGet, "/expenses/summarize", "")

		assert.Equal(t, http.StatusOK, response.StatusCode, "status code")
		expected := getExpectedResponse(t)
		assert.JSONEq(t, expected, responseBody, "response body")
	})
	t.Run("add expense successfully", func(t *testing.T) {
		response, _ := call(t, server, http.MethodPost, "/expenses/add", getRequest(t))

		assert.Equal(t, http.StatusCreated, response.StatusCode, "status code")
		// We should check if expense was added to the database here.
	})
	t.Run("add expense fails", func(t *testing.T) {
		response, responseBody := call(t, server, http.MethodPost, "/expenses/add", getRequest(t))

		assert.Equal(t, http.StatusBadRequest, response.StatusCode, "status code")
		expected := getExpectedResponse(t)
		assert.JSONEq(t, expected, responseBody, "response body")
	})
}

func startServer(t *testing.T, ctx context.Context) *httptest.Server {
	t.Helper()

	dbContainer := test_repos.StartDB(t, ctx)
	require.NoError(t, os.Setenv("DB_URL", dbContainer.DSN), "set DB_URL env var")

	setup, err := api.NewSetup()
	require.NoError(t, err, "new setup")

	server := httptest.NewServer(setup.APIMux)
	t.Cleanup(func() {
		server.Close()
	})

	return server
}

func call(t *testing.T, server *httptest.Server, method, path, body string) (*http.Response, string) {
	t.Helper()

	var b io.Reader
	if body != "" {
		b = bytes.NewBuffer([]byte(body))
	}
	request, err := http.NewRequest(method, server.URL+path, b)
	require.NoError(t, err, "new request")

	response, err := server.Client().Do(request)
	require.NoError(t, err, "do request")

	responseBody, err := io.ReadAll(response.Body)
	require.NoError(t, err, "read response body")

	return response, string(responseBody)

}

func getRequest(t *testing.T) string {
	t.Helper()

	path := fmt.Sprintf("./testdata/%s/request.json", t.Name())
	file, err := os.ReadFile(path)
	require.NoError(t, err, "read file")

	return string(file)
}

func getExpectedResponse(t *testing.T) string {
	t.Helper()

	path := fmt.Sprintf("./testdata/%s/response.json", t.Name())
	file, err := os.ReadFile(path)
	require.NoError(t, err, "read file")

	return string(file)
}
