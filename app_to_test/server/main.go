package main

import (
	"log"
	"net/http"

	"github.com/antosdaniel/go-presentation-beyond-unit-tests/app_to_test/server/api"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	setup, err := api.NewSetup()
	defer setup.DB.Close()
	if err != nil {
		log.Fatalf("setup failed: %v", err)
	}

	log.Println("running...")
	err = http.ListenAndServe(setup.APIAddress, setup.APIMux)
	if err != nil {
		log.Fatalf("serving HTTP failed: %v", err)
	}
}
