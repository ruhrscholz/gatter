package wellknown

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gatter/internal/environment"
	"log"
	"net/http"
	"strings"
)

type webfingerResponseLink struct {
	Rel      string `json:"rel"`
	Type_    string `json:"type,omitempty"`
	Href     string `json:"href,omitempty"`
	Template string `json:"template,omitempty"`
}

type webfingerResponse struct {
	Subject string                  `json:"subject"`
	Aliases []string                `json:"aliases"`
	Links   []webfingerResponseLink `json:"link"`
}

func Webfinger(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Header.Get("X-Forwarded-Proto")) != "https" && env.Deployment != environment.Development {
			http.Error(w, "Bad Request: Only https is permitted for this path", http.StatusBadRequest)
			return
		}

		resource := strings.Split(strings.TrimPrefix(r.URL.Query().Get("resource"), "acct:"), "@")

		// Domain not recognized
		if len(resource) == 2 && (resource[1] != env.LocalDomain && resource[1] != env.WebDomain) {
			http.NotFound(w, r)
			return
		}

		if len(resource) > 2 {
			http.Error(w, "Bad Request: Invalid \"acct:\" format", http.StatusBadRequest)
			return
		}

		// Check username exists
		stmt := "SELECT username FROM users WHERE username=$1"
		rows := env.Db.QueryRow(stmt, resource[0])

		// TODO Exists
		var username string
		if err := rows.Scan(&username); err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		} else if err != nil {
			log.Printf("Could not query database for webfinger: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/jrd+json")

		response := webfingerResponse{
			Subject: fmt.Sprintf("acct:%s@%s", username, env.LocalDomain),
			Aliases: []string{
				fmt.Sprintf("https://%s/@%s", env.WebDomain, username),
				fmt.Sprintf("https://%s/users/%s", env.WebDomain, username),
			},
			Links: []webfingerResponseLink{
				{
					Rel:   "http://webfinger.net/rel/profile-page",
					Type_: "text/html",
					Href:  fmt.Sprintf("https://%s/@%s", env.WebDomain, username),
				},
				{
					Rel:   "self",
					Type_: "application/activity+json",
					Href:  fmt.Sprintf("https://%s/users/%s", env.WebDomain, username),
				},
				{
					Rel:      "http://ostatus.org/schema/1.0/subscribe",
					Template: fmt.Sprintf("https://%s/authorize_interaction?uri={uri}", env.WebDomain),
				},
			},
		}

		json.NewEncoder(w).Encode(response)
	}
}
