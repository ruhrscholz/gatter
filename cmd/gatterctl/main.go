package main

// TODO Consider spf13 CLI parser

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"gatter/internal/util/config"
	"gatter/internal/util/password"
	"log"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	env, err := config.ReadConfig("./configs/config.json")
	if err != nil {
		log.Fatalf("Error while parsing config file: %s", err.Error())
		return
	}

	flag.Parse()

	log.Println(flag.Arg(0))
	var command string = flag.Arg(0)

	switch strings.ToLower(command) {
	case "adduser":
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		fmt.Printf("Password: ")
		passwordRaw, _ := reader.ReadString('\n')
		passwordRaw = strings.TrimSpace(passwordRaw)
		password, err := password.GenerateFromPlaintext("argon2", passwordRaw)

		if err != nil {
			log.Panic(err.Error())
		}

		ctx := context.Background()
		tx, err := env.Db.BeginTx(ctx, nil)
		if err != nil {
			log.Panicf("Could not begin transaction: %s", err.Error())
			return
		}

		statement := "INSERT INTO accounts (display_name, fed_username, fed_domain, id) VALUES ($1, $2, $3, $4)"
		_, err = tx.ExecContext(ctx, statement,
			username,
			username,
			env.WebDomain,
			fmt.Sprintf("https://%s/@%s", env.WebDomain, username))

		if err != nil {
			tx.Rollback()
			log.Panicf("Could not insert values into accounts table: %s", err.Error())
			return
		}

		query := tx.QueryRowContext(ctx, "SELECT account_id FROM accounts WHERE fed_username = $1", username)

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

		statement = "INSERT INTO users(username, password, account_id) VALUES ($1, $2, $3)"
		_, err = tx.ExecContext(ctx, statement, username, password, account_id)

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
