package main

import (
	"database/sql"
	"encoding/json"
	v1 "gatter/internal/api/v1"
	. "gatter/internal/env"
	"gatter/internal/webfinger"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
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

	db, err := sql.Open("pgx", "postgres://localhost:5432/gatter") // TODO Move to config file
	if err != nil {
		log.Panicf("Could not establish database connection: %s", err.Error())
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver)
	m.Up()

	env := Env{Db: db}

	// Route registration
	mux := http.NewServeMux()

	mux.HandleFunc("/.well-known/webfinger", webfinger.Handle(env))
	mux.Handle("api/v1/", v1.GetRoutes(env))

	err = http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatalf("Could not start http server: %s", err)
	}
}
