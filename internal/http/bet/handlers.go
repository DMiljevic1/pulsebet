package bethttp

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/DMiljevic1/pulsebet/internal/services/bet"
)

func RegisterHandlers(mux *http.ServeMux, service bet.Service, logger *slog.Logger) {
	mux.HandleFunc("/bet", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleSaveBet(w, r, service, logger)
			return
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

func handleSaveBet(w http.ResponseWriter, r *http.Request, service bet.Service, logger *slog.Logger) {
	var req SaveBetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode saveBet request", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	b := bet.Bet{
		MatchID:   req.MatchID,
		BettorID:  req.BettorID,
		Selection: req.Selection,
		Stake:     req.Stake,
	}

	created, err := service.SaveBet(b)
	if err != nil {
		logger.Error("Failed to save bet", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := SaveBetResponse{
		ID:        created.ID,
		MatchID:   created.MatchID,
		BettorID:  created.BettorID,
		Selection: created.Selection,
		Stake:     created.Stake,
		Odds:      created.Odds,
		PlacedAt:  created.PlacedAt.Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}
