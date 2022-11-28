package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"gatter/internal/environment"
	"log"
	"net"
	"net/http"
)

type contextKey int

const (
	KeyValidUsername contextKey = iota
	KeyDomain
	KeyValidAccount
)

func UserContext(env *environment.Env, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			errText := fmt.Sprintf("Error while parsing HTTP Host header into hostname and port. Host header was: %s. Error: %s", r.Host, err.Error())
			log.Panic(errText)
			if env.Deployment == environment.Development {
				http.Error(w, errText, http.StatusInternalServerError)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		stmt := "SELECT username, account_id FROM users WHERE domain=$1"
		rows := env.Db.QueryRow(stmt, host)

		var validUsername string
		var validAccount int64
		if err := rows.Scan(&validUsername, &validAccount); err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		} else if err != nil {
			errText := fmt.Sprintf("Error while executing SQL in UserContext middleware: %s", err.Error())
			log.Panic(errText)
			if env.Deployment == environment.Development {
				http.Error(w, errText, http.StatusInternalServerError)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		nctx := context.WithValue(r.Context(), KeyValidUsername, validUsername)
		nctx = context.WithValue(nctx, KeyDomain, host)
		nctx = context.WithValue(nctx, KeyValidAccount, validAccount)

		next.ServeHTTP(w, r.WithContext(nctx))
	})
}
