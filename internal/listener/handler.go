package listener

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func newEventHandler(eventsCh chan<- string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
			return
		}

		var payload EventPayload
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&payload); err != nil {
			http.Error(w, fmt.Sprintf("Invalid Json: %v", err), http.StatusBadRequest)
			return
		}

		if payload.ID == "" {
			http.Error(w, "required 'id' field", http.StatusBadRequest)
			return
		}

		eventsCh <- payload.ID

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
			"id":     payload.ID,
		})
	}
}
