package oauth

import (
	"gatter/internal/environment"
	"net/http"
)

func HandleToken(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}
