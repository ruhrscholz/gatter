package webfinger

import (
	. "gatter/internal/env"
	"net/http"
)

func Handle(env *Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		domain := r.Host
	}
}
