package main

import (
	"database/sql"
	"encoding/json"
	"gatter/internal/endpoints/activitypub/users"
	"gatter/internal/endpoints/client"
	"gatter/internal/endpoints/web"
	"gatter/internal/endpoints/web/htmx"
	wellknown2 "gatter/internal/endpoints/wellknown"
	"gatter/internal/environment"
	"gatter/internal/middleware"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type configOptions struct {
	Language     string                     `json:"language"`
	Database     string                     `json:"database"`
	Deployment   environment.DeploymentType `json:"deployment"`
	SessionToken string                     `json:"sessionToken"`
}

func main() {
	// Config file parsing
	configFile, err := os.ReadFile("./configs/config.json")
	if err != nil {
		log.Fatalf("Could not open config file: %s", err.Error())
	}

	var env environment.Env

	var config configOptions
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Could not parse config file: %s", err.Error())
	} else {
		log.Printf("Successfully read config file")

		switch config.Deployment {
		case environment.Development:
			env.Deployment = environment.Development
			log.Println("Deployment Type: Development")
		case environment.Production:
			env.Deployment = environment.Production
			log.Println("Deployment Type: Production")
		}

		env.Language = config.Language
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
	if err != nil {
		log.Printf("Error while connecting to database: %s", err.Error())
		return
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres",
		driver)
	if err != nil {
		log.Printf("Error while preparing database migration: %s", err.Error())
		return
	}
	_ = m.Up()

	env.Db = db

	mux := http.NewServeMux()

	// .well-known
	mux.Handle("/.well-known/webfinger", middleware.UserContext(&env, wellknown2.Webfinger(&env)))
	mux.Handle("/.well-known/nodeinfo", wellknown2.Nodeinfo(&env))

	// ActivityPub
	mux.Handle("/users/", middleware.UserContext(&env, http.StripPrefix("/users/", users.HandleUsers(&env))))

	// Client
	mux.Handle("/api/v1/accounts/", middleware.UserContext(&env, http.StripPrefix("/api/v1/accounts", client.HandleAccounts(&env))))
	mux.Handle("/api/v1/statuses/", middleware.UserContext(&env, http.StripPrefix("/api/v1/statuses", client.HandleStatuses(&env))))
	mux.Handle("/api/v1/timelines/", middleware.UserContext(&env, http.StripPrefix("/api/v1/timelines", client.HandleTimelines(&env))))

	// Web (and HTMX)
	mux.Handle("/", middleware.UserContext(&env, web.Handle(&env)))
	mux.Handle("/htmx/", middleware.UserContext(&env, http.StripPrefix("/htmx", htmx.Handle(&env))))

	err = http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatalf("Could not start http server: %s", err)
	}
}
