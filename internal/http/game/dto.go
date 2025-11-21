package gamehttp

type CreateMatchRequest struct {
	Home string `json:"home"`
	Away string `json:"away"`
}

type CreateMatchResponse struct {
	ID   string `json:"id"`
	Home string `json:"home"`
	Away string `json:"away"`
}

type GetMatchResponse struct {
	ID   string `json:"id"`
	Home string `json:"home"`
	Away string `json:"away"`
}
