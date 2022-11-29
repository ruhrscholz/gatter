package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"gatter/internal/environment"
	"log"
	"net/http"
)

type authKey int

const (
	KeyAuthUserId authKey = iota
)

func Auth(env *environment.Env, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionIdCookie, err := r.Cookie("session_id")

		if err == http.ErrNoCookie {
			nctx := context.WithValue(r.Context(), KeyAuthUserId, nil)
			next.ServeHTTP(w, r.WithContext(nctx))
		} else if err != nil {
			log.Printf("Error while reading session_id cookie: %s", err.Error())
			nctx := context.WithValue(r.Context(), KeyAuthUserId, nil)
			next.ServeHTTP(w, r.WithContext(nctx))
		} else {
			stmt := "SELECT user_id FROM sessions WHERE session_id=$1"
			rows := env.Db.QueryRow(stmt, fmt.Sprintf("\\x%s", sessionIdCookie.Value))

			var authUserId int
			if err = rows.Scan(&authUserId); err == sql.ErrNoRows {
				nctx := context.WithValue(r.Context(), KeyAuthUserId, nil)
				next.ServeHTTP(w, r.WithContext(nctx))
			} else if err != nil {
				log.Printf("Error while querying session_id cookie: %s", err.Error())
				nctx := context.WithValue(r.Context(), KeyAuthUserId, nil)
				next.ServeHTTP(w, r.WithContext(nctx))
			} else {
				nctx := context.WithValue(r.Context(), KeyAuthUserId, authUserId)
				next.ServeHTTP(w, r.WithContext(nctx))
			}
		}
	})
}
