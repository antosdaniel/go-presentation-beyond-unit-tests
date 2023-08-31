package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
)

type Setup struct {
	DB         *sql.DB
	APIAddress string
	APIMux     *http.ServeMux
}

func NewSetup() (*Setup, error) {
	db, err := sql.Open("pgx", os.Getenv("DB_URL"))
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	apiAddress := fmt.Sprintf(":%s", port)

	mux := http.NewServeMux()
	mux.HandleFunc("/expenses/all", getExpensesHandler(db))
	mux.HandleFunc("/expenses/add", addExpenseHandler(db))

	return &Setup{
		DB:         db,
		APIAddress: apiAddress,
		APIMux:     mux,
	}, nil
}
