package auth

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"gatter/internal/environment"
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

		switch r.Method {
		case http.MethodGet:
			template, err := template.New("sign_in.html.tmpl").ParseFiles("./internal/endpoints/auth/template/sign_in.html.tmpl")
			if err != nil {
				errText := fmt.Sprintf("Could not parse sign_in.html.tmpl: %s", err.Error())
				log.Panic(errText)
				if env.Deployment == environment.Development {
					http.Error(w, errText, http.StatusInternalServerError)
				} else {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}

			// Currently nil because we didn't implenment the CSRF token (or anything dynamic) yet
			err = template.Execute(w, nil)
			if err != nil {
				errText := fmt.Sprintf("Could not execute sign_in.html.tmpl: %s", err.Error())
				log.Panic(errText)
				if env.Deployment == environment.Development {
					http.Error(w, errText, http.StatusInternalServerError)
				} else {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}
			return
		case http.MethodPost:
			err := r.ParseForm()
			if err != nil {
				errText := fmt.Sprintf("Could not parse sign_in form data: %s", err.Error())
				log.Print(errText)
				if env.Deployment == environment.Development {
					http.Error(w, errText, http.StatusBadRequest)
					return
				} else {
					http.Error(w, "Bad Request", http.StatusBadRequest)
					return
				}
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
				errText := fmt.Sprintf("Could not query database for sign in: %s", err.Error())
				log.Print(errText)
				if env.Deployment == environment.Development {
					http.Error(w, errText, http.StatusInternalServerError)
					return
				} else {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
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
					errText := fmt.Sprintf("Could not insert session token into database: %s", err.Error())
					if env.Deployment == environment.Development {
						http.Error(w, errText, http.StatusInternalServerError)
					} else {
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					}
				}

				http.SetCookie(w, &http.Cookie{
					Name:     "session_id",
					Value:    fmt.Sprintf("%x", token),
					Secure:   true,
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
