package test_repos

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/antosdaniel/go-presentation-beyond-unit-tests/app_to_test/server/api"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This would usually live next to repo it is testing

func TestExpenseRepo_Add(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	db := StartDB(t, ctx).DB
	expenseRepo := api.NewExpenseRepo(db)

	t.Run("successfully adds expense", func(t *testing.T) {
		expense := api.Expense{
			ID:       "c811c5d4-c38a-4f61-932d-d656c203b5f6",
			Amount:   123_50,
			Category: "food",
			Date:     time.Date(2020, 9, 5, 0, 0, 0, 0, time.UTC),
			Notes:    "some notes",
		}

		err := expenseRepo.Add(expense)

		require.NoError(t, err, "could not add expense")
		result := getAllExpenses(t, db).FindByID(expense.ID)
		if assert.NotNil(t, result, "expense not added") {
			assert.Equal(t, expense, *result, "added expense is different")
		}
	})

	t.Run("fails on invalid ID", func(t *testing.T) {
		expense := api.Expense{
			ID:       "invalid-id",
			Amount:   123_50,
			Category: "food",
			Date:     time.Date(2020, 9, 5, 0, 0, 0, 0, time.UTC),
			Notes:    "some notes",
		}

		err := expenseRepo.Add(expense)

		if assert.Error(t, err) {
			assert.Equal(t, `ERROR: invalid input syntax for type uuid: "invalid-id" (SQLSTATE 22P02)`, err.Error())
		}
	})
}

// getAllExpenses Test helper for getting all expenses from the database.
func getAllExpenses(t *testing.T, db *sql.DB) api.Expenses {
	t.Helper()

	row, err := db.Query(`
		SELECT id, amount, category, date, notes 
		FROM expenses
	`)
	if err != nil {
		t.Fatalf("could not query expenses: %v", err)
	}
	defer row.Close()

	expenses := make([]api.Expense, 0)
	for row.Next() {
		var expense api.Expense
		err := row.Scan(&expense.ID, &expense.Amount, &expense.Category, &expense.Date, &expense.Notes)
		if err != nil {
			t.Fatalf("could not scan expense: %v", err)
		}
		expenses = append(expenses, expense)
	}

	return expenses
}
