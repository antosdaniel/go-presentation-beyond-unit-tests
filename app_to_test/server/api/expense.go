package api

import "time"

type Expenses []Expense

type Expense struct {
	ID       string    `json:"id"`
	Amount   float64   `json:"amount"`
	Category string    `json:"category"`
	Date     time.Time `json:"date"`
	Notes    string    `json:"notes"`
}

func (es Expenses) FindByID(id string) *Expense {
	for _, e := range es {
		if e.ID == id {
			return &e
		}
	}
	return nil
}

type ExpenseSummary struct {
	Year        int     `json:"year"`
	Month       int     `json:"month"`
	Category    string  `json:"category"`
	TotalAmount float64 `json:"total_amount"`
}
