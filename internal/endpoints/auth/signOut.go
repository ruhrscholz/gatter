package auth

import (
	"fmt"
	"gatter/internal/environment"
	"gatter/internal/middleware"
	"log"
	"net/http"
)

func HandleSignOut(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectUri := r.URL.Query().Get("redirect_uri") // TODO Very bad practive and phishing possibility

		if redirectUri == "" {
			redirectUri = "/"
		}

		if r.Context().Value(middleware.KeyAuthUserId) == nil {
			http.Redirect(w, r, redirectUri, http.StatusFound)
		}

		sessionIdCookie, err := r.Cookie("session_id")
		if err != nil {
			log.Printf("Error while reading session_id cookie: %s", err.Error())
			http.Error(w, "Bad Request", http.StatusBadRequest)
		} else {
			stmt := "DELETE FROM sessions WHERE session_id=$1"
			_, err := env.Db.Exec(stmt, fmt.Sprintf("\\x%x", sessionIdCookie))
			if err != nil {
				log.Printf("Could not remove session from database: %s", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    "",
				Secure:   (env.Deployment != environment.Development),
				HttpOnly: true,
				MaxAge:   -1,
			})
			http.Redirect(w, r, redirectUri, http.StatusFound)
			return
		}
	}

}
