package bank

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// This pulls transactions from the bank API and changes them into expenses.

type Transaction struct {
	ID        string    `json:"id"`
	Amount    float64   `json:"amount"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

type Transactions []Transaction

type API struct {
	apiAddress string
}

func New(apiAddress string) *API {
	return &API{apiAddress: apiAddress}
}

func (b *API) GetTransactions() (Transactions, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(fmt.Sprintf("%s/get-transactions", b.apiAddress))
	if err != nil {
		return nil, fmt.Errorf("could not retrieve bank transactions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	var transactions []Transaction
	err = json.NewDecoder(resp.Body).Decode(&transactions)
	if err != nil {
		return nil, fmt.Errorf("could not decode bank transactions: %w", err)
	}

	return transactions, nil
}
