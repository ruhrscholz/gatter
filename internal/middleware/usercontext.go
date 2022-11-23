package middleware

import (
	"context"
	"database/sql"
	"gatter/internal/environment"
	"log"
	"net"
	"net/http"
)

type contextKey int

const (
	KeyValidUsername contextKey = iota
	KeyDomain
)

func UserContext(env *environment.Env, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			log.Printf("Error while parsing HTTP Host header into hostname and port. Host header was: %s. Error: %s", r.Host, err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		stmt := "SELECT username FROM users WHERE domain=$1"
		rows := env.Db.QueryRow(stmt, host)

		var validUsername string
		if err := rows.Scan(&validUsername); err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		} else if err != nil {
			log.Printf("Error while executing SQL in UserContext middleware: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		nctx := context.WithValue(r.Context(), KeyValidUsername, validUsername)
		nctx = context.WithValue(nctx, KeyDomain, host)

		next.ServeHTTP(w, r.WithContext(nctx))
	})
}
