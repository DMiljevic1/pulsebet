package game

import "github.com/google/uuid"

type Service interface {
	CreateMatch(home, away string) (Match, error)
	GetAll() []Match
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s *service) CreateMatch(home, away string) (Match, error) {
	match := Match{
		ID:   uuid.NewString(),
		Home: home,
		Away: away,
	}
	err := s.repository.Save(match)
	if err != nil {
		return Match{}, err
	}
	return match, nil
}

func (s *service) GetAll() []Match {
	return s.repository.GetAll()
}
