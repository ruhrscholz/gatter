package main

import (
	"database/sql"
	"flag"
	"log"
	"strings"

	_ "github.com/jackc/pgx/v5"
)

func main() {
	var command string = flag.Arg(0)

	switch strings.ToLower(command) {
	case "adduser":
		db, err := sql.Open("pgx", "postgres://localhost:5432/gatter") // TODO Move to config file

		if err != nil {
			log.Panicf("Error while connecting to database: %s", err)

		}

		statement := "INSERT INTO accounts(display_name, fed_username, fed_domain, uri, url) VALUES (?, ?, ?, ?, ?)"
		_, err = db.Exec(statement)

		if err != nil {
			log.Panicf("Could not insert values: %s", err)
		}
	}
}
