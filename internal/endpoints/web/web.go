package web

import (
	"gatter/internal/environment"
	"net/http"
)

func Handle(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
