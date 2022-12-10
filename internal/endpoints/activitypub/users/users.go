package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gatter/internal/endpoints/activitypub"
	"gatter/internal/environment"
	"net/http"
	"strings"

	"golang.org/x/exp/slices"
)

var env *environment.Env

func HandleUsers(_env *environment.Env) http.HandlerFunc {
	env = _env
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path, "/")

		// We do not check for username exitance here (and thus can not bail early), to reduce the number of DB lookups
		subPath := strings.TrimPrefix(r.URL.Path, path[0])

		switch subPath {
		case "":
			basePath(w, r, path[0])
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

func basePath(w http.ResponseWriter, r *http.Request, username string) {
	if slices.Contains([]string{"application/json", "application/json", "application/ld+json"}, r.Header.Get("Accept")) {

		w.Header().Set("Content-Type", "application/activity+json; charset=utf-8")

		var displayName string
		stmt := "SELECT accounts.display_name FROM accounts INNER JOIN users ON users.account_id=accounts.account_id WHERE users.user_id=$1"
		rows := env.Db.QueryRow(stmt, username)

		if err := rows.Scan(&displayName); err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}

		response := activitypub.ActorResponse{
			Context: []string{"https://www.w3.org/ns/activitystreams"}, Type: "Person",
			Name:      displayName,
			Id:        fmt.Sprintf("https://%s/users/%s", env.WebDomain, username),
			Inbox:     fmt.Sprintf("https://%s/users/%s/inbox", env.WebDomain, username),
			Outbox:    fmt.Sprintf("https://%s/users/%s/outbox", env.WebDomain, username),
			Following: fmt.Sprintf("https://%s/users/%s/following", env.WebDomain, username),
			Followers: fmt.Sprintf("https://%s/users/%s/followers", env.WebDomain, username),
			Likes:     fmt.Sprintf("https://%s/users/%s/likes", env.WebDomain, username),
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Redirect(w, r, fmt.Sprintf("https://%s/@%s", env.WebDomain, username), http.StatusMovedPermanently)
	}
}
