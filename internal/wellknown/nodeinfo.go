package wellknown

import (
	"gatter/internal/environment"
	"net/http"
)

func SetUpNodeinfo(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
