package bet

import "time"

type Bet struct {
	ID        string
	MatchID   string
	BettorID  string
	Selection string
	Stake     float64
	Odds      float64
	PlacedAt  time.Time
}

type Match struct {
	ID   string
	Home string
	Away string
}
