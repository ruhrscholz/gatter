package api

import (
	v1 "gatter/internal/api/v1"
	"net/http"
)

func GetApiRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("v1/", v1.GetV1Routes())

	return mux
}
