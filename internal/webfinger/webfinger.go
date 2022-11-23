package webfinger

import (
	"encoding/json"
	"fmt"
	"gatter/internal/environment"
	"log"
	"net"
	"net/http"
	"strings"
)

type webfingerResponseLink struct {
	Rel      string `json:"rel"`
	Type_    string `json:"type"`
	Href     string `json:"href"`
	Template string `json:"template"`
}

type webfingerResponse struct {
	Subject string                  `json:"subject"`
	Aliases []string                `json:"aliases"`
	Links   []webfingerResponseLink `json:"link"`
}

func Handle(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Header.Get("X-Forwarded-Proto")) != "https" && env.Deployment != environment.Development {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		host, _, err := net.SplitHostPort(r.Host)
		if err != nil {
			log.Printf("Error while parsing HTTP Host header into hostname and port. Host header was: %s. Error: %s", r.Host, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resource := strings.Split(strings.TrimPrefix(r.URL.Query().Get("resource"), "acct:"), "@")

		if err != nil || len(resource) > 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(resource) > 1 && resource[1] != host {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		stmt := "SELECT EXISTS (SELECT username FROM users WHERE username=$1 AND domain=$2)"
		rows := env.Db.QueryRow(stmt, resource[0], host)

		var exists bool
		if err := rows.Scan(&exists); err != nil {
			log.Printf("Error while executing SQL in webfinger request: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.Header().Set("Content-Type", "application/jrd+json")

			response := webfingerResponse{
				Subject: fmt.Sprintf("acct:%s@%s", resource[0], host), // Using the resource[0] directly should be safe since we already checked for existance in the DB,
				Aliases: []string{
					fmt.Sprintf("https://%s/@%s", host, resource[0]),
					fmt.Sprintf("https://%s/users/%s", host, resource[0]), // This forwards to the other link for text/html requests
				},
				Links: []webfingerResponseLink{
					{
						Rel:   "http://webfinger.net/rel/profile-page",
						Type_: "text/html",
						Href:  fmt.Sprintf("https://%s/@%s", host, resource[0]),
					},
					{
						Rel:   "self",
						Type_: "application/activity+json",
						Href:  fmt.Sprintf("https://%s/users/%s", host, resource[0]),
					},
					{
						Rel:      "http://ostatus.org/schema/1.0/subscribe",
						Template: fmt.Sprintf("https://%s/authorize_interaction?uri={uri}", host),
					},
				},
			}

			json.NewEncoder(w).Encode(response)
			return
		}
	}
}
