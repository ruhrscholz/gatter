package v1

import "net/http"

func GetV1Routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("timelines/", getTimelinesRoutes())
	mux.Handle("statuses/", getStatusesRoutes())

	return mux
}
