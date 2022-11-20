package v1

import (
	. "gatter/internal/env"
	"net/http"
)

var env *Env

func GetRoutes(_env *Env) *http.ServeMux {
	env = _env
	mux := http.NewServeMux()

	mux.Handle("timelines/", getTimelinesRoutes())
	mux.Handle("statuses/", getStatusesRoutes())

	return mux
}
