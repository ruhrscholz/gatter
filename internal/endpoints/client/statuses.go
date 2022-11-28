package client

import (
	"gatter/internal/environment"
	"net/http"
)

func HandleStatuses(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}
