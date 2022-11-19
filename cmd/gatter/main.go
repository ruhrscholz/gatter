package main

import (
	"encoding/json"
	"gatter/internal/api"
	"gatter/internal/webfinger"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type configOptions struct {
	SingleUser   bool   `json:"singleUser"`
	Language     string `json:"language"`
	SessionToken string `json:"sessionToken"`
}

func main() {
	// Config file parsing
	configFile, err := os.ReadFile("./configs/config.json")
	if err != nil {
		log.Fatalf("Could not open config file: %s", err)
	}

	var config configOptions
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Could not parse config file: %s", err)
	} else {
		log.Printf("Successfully read config file")
		log.Printf("Single user mode: %t", config.SingleUser)
		log.Printf("Language Code: %s", config.Language)
	}

	if config.SessionToken == "changeme" || config.SessionToken == "" {
		log.Panic("Please change the default session token before attempting to start the server")
	}

	m, err := migrate.New(
		"file://db/migrations",
		"postgres://localhost:5432/gatter") // TODO Move to config file
	m.Steps(2)

	// Route registration
	mux := http.NewServeMux()

	mux.HandleFunc("/.well-known/webfinger", webfinger.Handle)
	mux.Handle("/api", api.GetApiRoutes())

	err = http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatalf("Could not start http server: %s", err)
	}
}
