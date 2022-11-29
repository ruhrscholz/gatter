package client

import (
	"gatter/internal/environment"
	"net/http"
)

func HandleApps(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}
