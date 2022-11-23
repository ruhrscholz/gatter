package s2s

import (
	"fmt"
	"gatter/internal/environment"
	"gatter/internal/middleware"
	"net/http"
	"strings"
)

var env *environment.Env

func SetUp(_env *environment.Env) http.HandlerFunc {
	env = _env

	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path, "/")

		username := r.Context().Value(middleware.KeyValidUsername).(string)

		if path[0] != username {
			http.NotFound(w, r)
			return
		}

		subPath := strings.TrimPrefix(r.URL.Path, username)

		switch subPath {
		case "/":
			basePath(w, r)
			return
		default:
			http.NotFound(w, r)
			return
		}
	}
}

func basePath(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") == "application/json" || r.Header.Get("Accept") == "application/activity+json" {
		// TODO
	} else {
		http.Redirect(w, r, fmt.Sprintf("https://%s/@%s", r.Context().Value(middleware.KeyDomain), r.Context().Value(middleware.KeyValidUsername)), http.StatusMovedPermanently)
	}
}
