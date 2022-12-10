package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gatter/internal/environment"
	"log"
	"os"
)

type configOptions struct {
	Language    string                     `json:"language"`
	Database    string                     `json:"database"`
	Deployment  environment.DeploymentType `json:"deployment"`
	LocalDomain string                     `json:"local_domain"`
	WebDomain   string                     `json:"web_domain"`
}

func ReadConfig(path string) (*environment.Env, error) {
	// Config file parsing
	configFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var env environment.Env

	var config configOptions
	var db *sql.DB
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

		if config.Language != "" {
			env.Language = config.Language
			log.Printf("Language Code: %s", config.Language)
		} else {
			env.Language = "en"
			log.Printf("Warning: Language code unset, defaulting to en")
		}

		if config.LocalDomain != "" {
			env.LocalDomain = config.LocalDomain
			log.Printf("Local Domain: %s", config.LocalDomain)
		} else {
			log.Fatal()
			return nil, fmt.Errorf("Local domain not specified")
		}

		if config.WebDomain != "" {
			env.WebDomain = config.WebDomain
			log.Printf("Web Domain: %s", config.WebDomain)
		} else {
			log.Print("Warning: Web domain not specified, defaulting to local domain")
			env.WebDomain = env.LocalDomain
		}

		if config.Database != "" {
			db, err = sql.Open("pgx", config.Database)
			if err != nil {
				return nil, err
			} else {
				err = db.Ping()
				if err != nil {
					return nil, err
				}
				env.Db = db
			}
		} else {
			return nil, fmt.Errorf("Database connection not specified")
		}
	}
	return &env, nil
}
