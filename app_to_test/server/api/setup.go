package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/antosdaniel/go-presentation-beyond-unit-tests/app_to_test/server/bank"
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

	expenseRepo := NewExpenseRepo(db)

	bankAPI := bank.New(os.Getenv("BANK_API_URL"))

	mux := http.NewServeMux()
	mux.HandleFunc("/expenses/add", addExpenseHandler(expenseRepo))
	mux.HandleFunc("/expenses/all", getExpensesHandler(expenseRepo))
	mux.HandleFunc("/expenses/summarize", summarizeExpensesHandler(expenseRepo))
	mux.HandleFunc("/expenses/sync", syncFromBankHandler(expenseRepo, bankAPI))

	return &Setup{
		DB:         db,
		APIAddress: apiAddress,
		APIMux:     mux,
	}, nil
}
