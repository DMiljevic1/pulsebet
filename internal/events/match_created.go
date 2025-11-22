package events

type MatchCreated struct {
	ID   string `json:"id"`
	Home string `json:"home"`
	Away string `json:"away"`
}
