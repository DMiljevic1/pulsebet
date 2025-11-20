package gamehttp

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/DMiljevic1/pulsebet/internal/services/game"
)

func RegisterHandlers(mux *http.ServeMux, service game.Service, logger *slog.Logger) {
	mux.HandleFunc("/matches", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		handleCreateMatch(w, req, service, logger)
	})
}

func handleCreateMatch(w http.ResponseWriter, r *http.Request, service game.Service, logger *slog.Logger) {
	defer r.Body.Close()
	var req CreateMatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	match, err := service.CreateMatch(req.Home, req.Away)
	if err != nil {
		logger.Error("Failed to create a match", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	response := CreateMatchResponse{
		ID:   match.ID,
		Home: match.Home,
		Away: match.Away,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
