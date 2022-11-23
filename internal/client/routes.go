package client

import (
	"gatter/internal/environment"
	"net/http"
)

var env *environment.Env

func GetRoutes(_env *environment.Env) *http.ServeMux {
	env = _env
	mux := http.NewServeMux()

	mux.Handle("timelines/", getTimelinesRoutes())
	mux.Handle("statuses/", getStatusesRoutes())

	return mux
}
