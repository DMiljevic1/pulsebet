package oddshttp

type SetOddsRequest struct {
	MatchID  string  `json:"matchId"`
	HomeOdds float64 `json:"homeOdds"`
	DrawOdds float64 `json:"drawOdds"`
	AwayOdds float64 `json:"awayOdds"`
}

type SetOddsResponse struct {
	MatchID string `json:"matchId"`
	Status  string `json:"status"`
}

type OddsResponse struct {
	MatchID   string  `json:"matchId"`
	HomeOdds  float64 `json:"homeOdds"`
	DrawOdds  float64 `json:"drawOdds"`
	AwayOdds  float64 `json:"awayOdds"`
	UpdatedAt string  `json:"updatedAt"`
}
