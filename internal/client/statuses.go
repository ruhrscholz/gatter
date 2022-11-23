package client

import "net/http"

func getStatusesRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	return mux
}
