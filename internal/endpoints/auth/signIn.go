package auth

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"gatter/internal/environment"
	"gatter/internal/middleware"
	"gatter/internal/util/password"
	"html/template"
	"log"
	"net/http"
)

func HandleSignIn(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectUri := r.URL.Query().Get("redirect_uri") // TODO Very bad practive and phishing possibility

		if redirectUri == "" {
			redirectUri = "/"
		}

		if r.Context().Value(middleware.KeyAuthUserId) != nil {
			http.Redirect(w, r, redirectUri, http.StatusFound)
		}

		switch r.Method {
		case http.MethodGet:
			template, err := template.New("sign_in.html.tmpl").ParseFiles("./internal/endpoints/auth/template/sign_in.html.tmpl")
			if err != nil {
				log.Printf("Could not parse sign_in.html.tmpl: %s", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Currently nil because we didn't implenment the CSRF token (or anything dynamic) yet
			err = template.Execute(w, nil)
			if err != nil {
				log.Printf("Could not execute sign_in.html.tmpl: %s", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			return
		case http.MethodPost:
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			if r.Form.Get("username") == "" {
				http.Error(w, "Username field may not be empty", http.StatusBadRequest)
				return
			}
			if r.Form.Get("password") == "" {
				http.Error(w, "Password field may not be empty", http.StatusBadRequest)
				return
			}

			// Check username exists
			stmt := "SELECT user_id, password FROM users WHERE username=$1"
			rows := env.Db.QueryRow(stmt, r.Form.Get("username"))

			var userId uint
			var correctPassword string
			if err := rows.Scan(&userId, &correctPassword); err == sql.ErrNoRows {
				http.Error(w, "Wrong username or password", http.StatusUnauthorized)
				return
			} else if err != nil {
				log.Printf("Could not query database for sign in: %s", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				if !password.Validate(r.Form.Get("password"), correctPassword) {
					http.Error(w, "Wrong username or password", http.StatusUnauthorized)
					return
				}

				token := make([]byte, 32)
				rand.Read(token)

				stmt := "INSERT INTO sessions (session_id, user_id) VALUES ($1, $2)"
				_, err := env.Db.Exec(stmt, fmt.Sprintf("\\x%x", token), userId)
				if err != nil {
					log.Printf("Could not insert session token into database: %s", err.Error())
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				http.SetCookie(w, &http.Cookie{
					Name:     "session_id",
					Value:    fmt.Sprintf("%x", token),
					Secure:   (env.Deployment != environment.Development),
					HttpOnly: true,
				})
				http.Redirect(w, r, redirectUri, http.StatusFound)
				return
			}

		default:
			http.NotFound(w, r)
			return
		}
	}
}
