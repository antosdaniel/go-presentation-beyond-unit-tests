package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/antosdaniel/go-presentation-beyond-unit-tests/app_to_test/server/bank"
)

// This project is not about unit tests, so we're not doing interfaces :)

func addExpenseHandler(expenseRepo *ExpenseRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var expense Expense
		err := json.NewDecoder(r.Body).Decode(&expense)
		if err != nil {
			slog.With(slog.String("error", err.Error())).Warn("invalid request body")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "invalid request body"}`))
			return
		}

		if expense.Category == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "category is required"}`))
			return
		}

		err = expenseRepo.Add(expense)
		if err != nil {
			internalError(w, err, "could not create expense")
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func getExpensesHandler(expenseRepo *ExpenseRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		expenses, err := expenseRepo.All()
		if err != nil {
			internalError(w, err, "could not get expenses")
			return
		}

		err = json.NewEncoder(w).Encode(expenses)
		if err != nil {
			internalError(w, err, "could not encode expenses")
			return
		}
	}
}

func summarizeExpensesHandler(expenseRepo *ExpenseRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		expenseSums, err := expenseRepo.Summarize()
		if err != nil {
			internalError(w, err, "could not get expense summary")
			return
		}

		err = json.NewEncoder(w).Encode(expenseSums)
		if err != nil {
			internalError(w, err, "could not encode expense summary")
			return
		}
	}
}

func syncFromBankHandler(expenseRepo *ExpenseRepo, bankAPI *bank.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bankTransactions, err := bankAPI.GetTransactions()
		if err != nil {
			internalError(w, err, "could not retrieve bank transactions")
			return
		}

		synced := 0
		for _, transaction := range bankTransactions {
			if transaction.Category == "" {
				continue
			}
			expense := Expense{
				ID:       transaction.ID,
				Amount:   transaction.Amount,
				Category: transaction.Category,
				Date:     transaction.CreatedAt,
				Notes:    "Synced from Bank API",
			}

			err := expenseRepo.Add(expense)
			if expenseAlreadyExists(err) {
				continue
			}
			if err != nil {
				internalError(w, err, "could not create expense from bank transaction")
				return
			}
			synced++
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"synced": %d}`, synced)))
	}
}

func internalError(w http.ResponseWriter, err error, msg string) {
	slog.With(slog.String("error", err.Error())).Error(msg)
	w.WriteHeader(http.StatusInternalServerError)
}

func expenseAlreadyExists(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), `duplicate key value violates unique constraint "expenses_pkey"`)
}
