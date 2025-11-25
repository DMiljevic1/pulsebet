package odds

import "time"

type Odds struct {
	MatchID   string
	HomeOdds  float64
	DrawOdds  float64
	AwayOdds  float64
	UpdatedAt time.Time
}
