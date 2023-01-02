package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func TimelinesPublic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		local := false
		limit := 20

		if strings.EqualFold(r.URL.Query().Get("local"), "true") {
			local = true
		}

		if local {
			http.Error(w, "Not yet implemented", http.StatusNotImplemented)
		}

		if r.URL.Query().Get("limit") != "" {
			limitRequest, err := strconv.Atoi(r.URL.Query().Get("limit"))
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			if limitRequest <= 40 {
				limit = limitRequest
			}
		}

		var statuses []Status

		stmt := "SELECT public_id, created_at, content FROM statuses INNER JOIN accounts on statuses.author_id = accounts.account_id ORDER BY created_at DESC LIMIT $1"
		rows, err := env.Db.Query(stmt, limit)

		if err != nil {
			log.Printf("Could not query public timeline: %s", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			var status Status

			if err := rows.Scan(
				&status.Id,
				&status.CreatedAt,
				&status.Content); err != nil {
				log.Printf("Could not scan database statuses into model: %s", err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			statuses = append(statuses, status)
		}

		w.Header().Set("Content-Type", "application/json")
		if len(statuses) == 0 {
			fmt.Fprint(w, "[]")
			return
		}
		json.NewEncoder(w).Encode(statuses)
	}
}
