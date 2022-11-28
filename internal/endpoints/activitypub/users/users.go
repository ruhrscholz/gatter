package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gatter/internal/endpoints/activitypub"
	"gatter/internal/environment"
	"gatter/internal/middleware"
	"log"
	"net/http"
	"strings"
)

var env *environment.Env

func HandleUsers(_env *environment.Env) http.HandlerFunc {
	env = _env
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path, "/")

		username := r.Context().Value(middleware.KeyValidUsername).(string)

		if !strings.EqualFold(path[0], username) {
			http.NotFound(w, r)
			return
		}

		subPath := strings.TrimPrefix(r.URL.Path, username)

		switch subPath {
		case "":
			basePath(w, r)
			return
		case "inbox":
			handleInbox(w, r)
			return
		case "outbox":
			handleOutbox(w, r)
			return
		case "followers":
			handleFollowers(w, r)
			return
		case "following":
			handleFollowing(w, r)
			return
		case "likes":
			handleLikes(w, r)
			return
		case "statuses":
			handleStatuses(w, r)
			return
		default:
			http.NotFound(w, r)
			return
		}
	}
}

func basePath(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") == "application/json" ||
		r.Header.Get("Accept") == "application/activity+json" ||
		r.Header.Get("Accept") == "application/ld+json" {

		w.Header().Set("Content-Type", "application/activity+json; charset=utf-8")
		domain := r.Context().Value(middleware.KeyDomain)
		username := r.Context().Value(middleware.KeyValidUsername)
		accountId := r.Context().Value(middleware.KeyValidAccount)

		var displayName string
		stmt := "SELECT display_name FROM accounts WHERE account_id=$1"
		rows := env.Db.QueryRow(stmt, accountId)

		if err := rows.Scan(&displayName); err == sql.ErrNoRows {
			errText := fmt.Sprintf("No rows returned for statement \"%s\", account_id %d even though that row should exist", stmt, accountId)
			log.Panic(errText)
			if env.Deployment == environment.Development {
				http.Error(w, errText, http.StatusInternalServerError)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}

		response := activitypub.ActorResponse{
			Context:   []string{"https://www.w3.org/ns/activitystreams"},
			Type:      "Person",
			Name:      displayName,
			Id:        fmt.Sprintf("https://%s/users/%s", domain, username),
			Inbox:     fmt.Sprintf("https://%s/users/%s/inbox", domain, username),
			Outbox:    fmt.Sprintf("https://%s/users/%s/outbox", domain, username),
			Following: fmt.Sprintf("https://%s/users/%s/following", domain, username),
			Followers: fmt.Sprintf("https://%s/users/%s/followers", domain, username),
			Likes:     fmt.Sprintf("https://%s/users/%s/likes", domain, username),
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Redirect(w, r, fmt.Sprintf("https://%s/@%s", r.Context().Value(middleware.KeyDomain), r.Context().Value(middleware.KeyValidUsername)), http.StatusMovedPermanently)
	}
}
