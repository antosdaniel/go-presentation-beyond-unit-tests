package api

import "database/sql"

type ExpenseRepo struct {
	db *sql.DB
}

func NewExpenseRepo(db *sql.DB) *ExpenseRepo {
	return &ExpenseRepo{db: db}
}

func (r *ExpenseRepo) Add(expense Expense) error {
	insertQuery := `
		INSERT INTO expenses (id, amount, category, date, notes)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(insertQuery, expense.ID, expense.Amount, expense.Category, expense.Date, expense.Notes)
	if err != nil {
		return err
	}
	return nil
}

func (r *ExpenseRepo) All() ([]Expense, error) {
	rows, err := r.db.Query(`
		SELECT id, amount, category, date, notes 
		FROM expenses
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := make([]Expense, 0)
	for rows.Next() {
		var expense Expense
		err := rows.Scan(&expense.ID, &expense.Amount, &expense.Category, &expense.Date, &expense.Notes)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}

func (r *ExpenseRepo) Summarize() ([]ExpenseSummary, error) {
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

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenseSums := make([]ExpenseSummary, 0)
	for rows.Next() {
		var expenseSum ExpenseSummary
		err := rows.Scan(&expenseSum.Year, &expenseSum.Month, &expenseSum.Category, &expenseSum.TotalAmount)
		if err != nil {
			return nil, err
		}
		expenseSums = append(expenseSums, expenseSum)
	}
	return expenseSums, nil
}
