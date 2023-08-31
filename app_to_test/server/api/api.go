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
			internalError(w, err, "could not create expense")
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
			internalError(w, err, "could not query expenses")
			return
		}
		defer rows.Close()

		expenses := make([]Expense, 0)
		for rows.Next() {
			var expense Expense
			err := rows.Scan(&expense.ID, &expense.Amount, &expense.Category, &expense.Date, &expense.Notes)
			if err != nil {
				internalError(w, err, "could not scan expense")
				return
			}
			expenses = append(expenses, expense)
		}

		err = json.NewEncoder(w).Encode(expenses)
		if err != nil {
			internalError(w, err, "could not encode expenses")
			return
		}
	}
}

func summarizeExpensesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		query := `
			SELECT
			    EXTRACT(YEAR FROM date) AS year, 
			    EXTRACT(MONTH FROM date) AS month, 
			    category, 
			    SUM(amount) AS total_amount
			FROM expenses
			GROUP BY year, month, category
			ORDER BY year, month, category
		`

		rows, err := db.Query(query)
		if err != nil {
			internalError(w, err, "could not query expense summary")
			return
		}
		defer rows.Close()

		type ExpenseSum struct {
			Year        int     `json:"year"`
			Month       int     `json:"month"`
			Category    string  `json:"category"`
			TotalAmount float64 `json:"total_amount"`
		}

		expenseSums := make([]ExpenseSum, 0)
		for rows.Next() {
			var expenseSum ExpenseSum
			err := rows.Scan(&expenseSum.Year, &expenseSum.Month, &expenseSum.Category, &expenseSum.TotalAmount)
			if err != nil {
				internalError(w, err, "could not scan expense summary row")
				return
			}
			expenseSums = append(expenseSums, expenseSum)
		}

		err = json.NewEncoder(w).Encode(expenseSums)
		if err != nil {
			internalError(w, err, "could not encode expense summary")
			return
		}
	}
}

func internalError(w http.ResponseWriter, err error, msg string) {
	slog.With(slog.String("error", err.Error())).Error(msg)
	w.WriteHeader(http.StatusInternalServerError)
}
