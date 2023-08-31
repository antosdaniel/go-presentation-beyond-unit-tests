package api

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
)

type Expense struct {
	ID       string  `json:"id"`
	Amount   float64 `json:"amount"`
	Category string  `json:"category"`
	Date     string  `json:"date"`
	Notes    string  `json:"notes"`
}

func addExpenseHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var expense Expense
		err := json.NewDecoder(r.Body).Decode(&expense)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		insertQuery := `
			INSERT INTO expenses (id, amount, category, date, notes)
			VALUES ($1, $2, $3, $4, $5)
		`
		_, err = db.Exec(insertQuery, expense.ID, expense.Amount, expense.Category, expense.Date, expense.Notes)
		if err != nil {
			slog.With(slog.String("error", err.Error())).Error("could not create expense")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func getExpensesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		rows, err := db.Query("SELECT id, amount, category, date, notes FROM expenses")
		if err != nil {
			slog.With(slog.String("error", err.Error())).Error("could not query expenses")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		expenses := make([]Expense, 0)
		for rows.Next() {
			var expense Expense
			err := rows.Scan(&expense.ID, &expense.Amount, &expense.Category, &expense.Date, &expense.Notes)
			if err != nil {
				slog.With(slog.String("error", err.Error())).Error("could not scan expense")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			expenses = append(expenses, expense)
		}

		json.NewEncoder(w).Encode(expenses)
	}
}
