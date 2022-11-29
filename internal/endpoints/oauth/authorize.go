package oauth

import (
	"gatter/internal/environment"
	"net/http"
)

func HandleAuthorize(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
	}
}
