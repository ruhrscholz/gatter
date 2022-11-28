package htmx

import (
	"gatter/internal/environment"
	"net/http"
)

var env *environment.Env

// This gets registered from the main.go, *NOT* from web.go
func Handle(_env *environment.Env) *http.ServeMux {
	env = _env
	mux := http.NewServeMux()

	return mux
}
