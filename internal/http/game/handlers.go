package gamehttp

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/DMiljevic1/pulsebet/internal/services/game"
)

func RegisterHandlers(mux *http.ServeMux, service game.Service, logger *slog.Logger) {
	mux.HandleFunc("/matches", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			handleGetAllMatches(w, req, service)
			return
		case http.MethodPost:
			handleCreateMatch(w, req, service, logger)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
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
		return
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

func handleGetAllMatches(w http.ResponseWriter, r *http.Request, service game.Service) {
	matches := service.GetAll()

	response := make([]GetMatchResponse, 0, len(matches))
	for _, v := range matches {
		response = append(response, GetMatchResponse{
			ID:   v.ID,
			Home: v.Home,
			Away: v.Away,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
