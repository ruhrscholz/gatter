package wellknown

import (
	"encoding/json"
	"fmt"
	"gatter/internal/environment"
	"gatter/internal/middleware"
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

func SetUpWebfinger(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Header.Get("X-Forwarded-Proto")) != "https" && env.Deployment != environment.Development {
			http.Error(w, "Bad Request: Only https is permitted for this path", http.StatusBadRequest)
			return
		}

		resource := strings.Split(strings.TrimPrefix(r.URL.Query().Get("resource"), "acct:"), "@")

		if len(resource) == 1 && resource[0] != r.Context().Value(middleware.KeyValidUsername).(string) {
			http.NotFound(w, r)
			return
		}

		if len(resource) == 2 && (resource[0] != r.Context().Value(middleware.KeyValidUsername).(string) || resource[1] != r.Context().Value(middleware.KeyDomain).(string)) {
			http.NotFound(w, r)
			return
		}

		if len(resource) > 2 {
			http.Error(w, "Bad Request: Invalid \"acct:\" format", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/jrd+json")

		response := webfingerResponse{
			Subject: fmt.Sprintf("acct:%s@%s", r.Context().Value(middleware.KeyValidUsername), r.Context().Value(middleware.KeyDomain)), // Using the resource[0] directly should be safe since we already checked for existance in the DB,
			Aliases: []string{
				fmt.Sprintf("https://%s/@%s", r.Context().Value(middleware.KeyDomain), r.Context().Value(middleware.KeyValidUsername)),
				fmt.Sprintf("https://%s/users/%s", r.Context().Value(middleware.KeyDomain), r.Context().Value(middleware.KeyValidUsername)), // This forwards to the other link for text/html requests
			},
			Links: []webfingerResponseLink{
				{
					Rel:   "http://webfinger.net/rel/profile-page",
					Type_: "text/html",
					Href:  fmt.Sprintf("https://%s/@%s", r.Context().Value(middleware.KeyDomain), r.Context().Value(middleware.KeyValidUsername)),
				},
				{
					Rel:   "self",
					Type_: "application/activity+json",
					Href:  fmt.Sprintf("https://%s/users/%s", r.Context().Value(middleware.KeyDomain), r.Context().Value(middleware.KeyValidUsername)),
				},
				{
					Rel:      "http://ostatus.org/schema/1.0/subscribe",
					Template: fmt.Sprintf("https://%s/authorize_interaction?uri={uri}", r.Context().Value(middleware.KeyDomain)),
				},
			},
		}

		json.NewEncoder(w).Encode(response)
		return
	}
}
