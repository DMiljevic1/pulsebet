package gamehttp

import (
	"encoding/json"
	"net/http"

	"github.com/DMiljevic1/pulsebet/internal/services/game"
)

func RegisterHandlers(mux *http.ServeMux, service game.Service) {
	mux.HandleFunc("/matches", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}

	})
}

func handleCreateMatch(w http.ResponseWriter, r *http.Request, service game.Service) {
	defer r.Body.Close()
	var req CreateMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	match := service.CreateMatch(req.Home, req.Away)

	response := CreateMatchResponse{
		ID:   match.ID,
		Home: match.Home,
		Away: match.Away,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
