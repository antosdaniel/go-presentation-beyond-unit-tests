package main

// This is simple implementation of third party Bank API that we want to connect with for our expenses tracker app.

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/get-transactions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(transactions))
	})

	log.Println("bank API running...")
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatalf("serving HTTP failed: %v", err)
	}
}

const transactions = `[
  {
    "id": "1a630ebc-f713-4f17-bcba-955292a2490d",
    "amount": 500.00,
    "category": "",
    "created_at": "2020-01-01T00:00:00Z"
  },
  {
    "id": "a9b4f2a9-b76f-40c7-acd8-0553a138819b",
    "amount": 100.00,
    "category": "food",
    "created_at": "2020-01-01T06:00:00Z"
  },
  {
    "id": "c1b8e7d2-8f7a-4d3e-aa91-1a2b3c4d5e6f",
    "amount": 45.50,
    "category": "transportation",
    "created_at": "2020-01-02T08:30:15Z"
  },
  {
    "id": "d3e4f5a6-b7c8-4d5e-8f9a-0b1c2a3d4e5f",
    "amount": 75.25,
    "category": "shopping",
    "created_at": "2020-01-03T14:20:30Z"
  },
  {
    "id": "f6a7b8c9-d0e1-4f2a-9b3c-5e6f7a8b9c0d",
    "amount": 30.00,
    "category": "entertainment",
    "created_at": "2020-01-04T19:45:10Z"
  }
]`
