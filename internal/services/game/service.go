package game

import "github.com/google/uuid"

type Service interface {
	CreateMatch(home, away string) Match
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (s *service) CreateMatch(home, away string) Match {
	return Match{
		ID:   uuid.NewString(),
		Home: home,
		Away: away,
	}
}
