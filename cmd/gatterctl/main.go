package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"gatter/internal/util/password"
	"log"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	flag.Parse()

	log.Println(flag.Arg(0))
	var command string = flag.Arg(0)

	switch strings.ToLower(command) {
	case "adduser":
		db, err := sql.Open("pgx", "postgres://localhost:5432/gatter") // TODO Move to config file
		if err != nil {
			log.Panicf("Error while connecting to database: %s", err)
		}
		defer db.Close()

		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		fmt.Printf("Domain: ")
		domain, _ := reader.ReadString('\n')
		domain = strings.TrimSpace(domain)

		fmt.Printf("Password: ")
		passwordRaw, _ := reader.ReadString('\n')
		passwordRaw = strings.TrimSpace(passwordRaw)
		password, err := password.GenerateFromPlaintext("argon2", passwordRaw)

		if err != nil {
			log.Panic(err.Error())
		}

		ctx := context.Background()
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			log.Panicf("Could not begin transaction: %s", err.Error())
			return
		}

		statement := "INSERT INTO accounts (display_name, fed_username, fed_domain, id) VALUES ($1, $2, $3, $4)"
		_, err = tx.ExecContext(ctx, statement,
			username,
			username,
			domain,
			fmt.Sprintf("https://%s/@%s", domain, username))

		if err != nil {
			tx.Rollback()
			log.Panicf("Could not insert values into accounts table: %s", err.Error())
			return
		}

		query := tx.QueryRowContext(ctx, "SELECT account_id FROM accounts WHERE fed_username = $1 AND fed_domain = $2", username, domain)

		var account_id int
		switch err := query.Scan(&account_id); err {
		case sql.ErrNoRows:
			tx.Rollback()
			log.Panicf("Could not find newly created account")
			return
		case nil:

		default:
			tx.Rollback()
			log.Panicf("Could not query for newly created account: %s", err.Error())
			return
		}

		statement = "INSERT INTO users(username, password, domain, account_id) VALUES ($1, $2, $3, $4)"
		_, err = tx.ExecContext(ctx, statement, username, password, domain, account_id)

		if err != nil {
			tx.Rollback()
			log.Panicf("Could not insert values into users table: %s", err.Error())
			return
		}

		err = tx.Commit()
		if err != nil {
			log.Panicf("Could not execute transaction: %s", err.Error())
			return
		}
	default:
		log.Println("You have to specify a valid command")
		return
	}
}
