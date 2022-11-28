package wellknown

import (
	"gatter/internal/environment"
	"net/http"
)

type nodeinfoResponse struct {
	// TODO
}

func Nodeinfo(env *environment.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := nodeinfoResponse{}
		_ = response

		//w.Header().Set("Content-Type", "application/jrd+json")

		// TODO Actually implement Nodeinfo
		//
		// See specification: http://nodeinfo.diaspora.software/
		// Example response: https://mastodon.sdf.org/nodeinfo/2.0
		http.NotFound(w, r)
	}
}
