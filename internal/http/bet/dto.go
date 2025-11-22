package bethttp

type SaveBetRequest struct {
	MatchID   string  `json:"match_id"`
	BettorID  string  `json:"bettor_id"`
	Selection string  `json:"selection"`
	Stake     float64 `json:"stake"`
}

type SaveBetResponse struct {
	ID        string  `json:"id"`
	MatchID   string  `json:"match_id"`
	BettorID  string  `json:"bettor_id"`
	Selection string  `json:"selection"`
	Stake     float64 `json:"stake"`
	Odds      float64 `json:"odds"`
	PlacedAt  string  `json:"placed_at"`
}
