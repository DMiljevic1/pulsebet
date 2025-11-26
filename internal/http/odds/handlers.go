package oddshttp

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/DMiljevic1/pulsebet/internal/services/odds"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service odds.Service
	logger  *slog.Logger
}

func RegisterRoutes(r chi.Router, service odds.Service, logger *slog.Logger) {
	h := &Handler{
		service,
		logger.With("component", "odds_http"),
	}

	r.Route("/odds", func(r chi.Router) {
		r.Get("/{matchID}", h.getOdds)
		r.Post("/", h.setOdds)
	})
}

func (h *Handler) getOdds(w http.ResponseWriter, r *http.Request) {
	matchID := chi.URLParam(r, "matchID")
	if matchID == "" {
		http.Error(w, "matchId is required", http.StatusBadRequest)
		return
	}

	o, err := h.service.GetOdds(r.Context(), matchID)
	if err != nil {
		h.logger.Error("Failed to get Odds", "matchID", matchID, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := OddsResponse{
		MatchID:   matchID,
		AwayOdds:  o.AwayOdds,
		DrawOdds:  o.DrawOdds,
		HomeOdds:  o.HomeOdds,
		UpdatedAt: o.UpdatedAt.Format(time.RFC3339),
	}

	writeJSON(w, response, http.StatusOK)
}

func (h *Handler) setOdds(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req SetOddsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to set Odds", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	o, err := h.service.SetOdds(r.Context(), odds.Odds{
		MatchID:  req.MatchID,
		AwayOdds: req.AwayOdds,
		DrawOdds: req.DrawOdds,
		HomeOdds: req.HomeOdds,
	})

	if err != nil {
		h.logger.Error("Failed to set Odds", "matchID", req.MatchID, "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := OddsResponse{
		MatchID:   req.MatchID,
		HomeOdds:  o.HomeOdds,
		DrawOdds:  o.DrawOdds,
		AwayOdds:  o.AwayOdds,
		UpdatedAt: o.UpdatedAt.Format(time.RFC3339),
	}

	writeJSON(w, resp, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
