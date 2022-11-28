package client

import (
	"gatter/internal/environment"
	"net/http"
)

func HandleAccounts(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}
