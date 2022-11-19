package v1

import (
	"encoding/json"
	"gatter/internal/api/model"
	"net/http"
)

func getTimelinesRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/public", timelinesPublic)
	mux.HandleFunc("/tag", timelinesTag)

	return mux
}

func timelinesPublic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	queryParams := r.URL.Query()
	statuses := []model.Status{}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

func timelinesTag(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}
