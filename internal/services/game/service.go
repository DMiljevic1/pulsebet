package game

import (
	"context"

	"github.com/DMiljevic1/pulsebet/internal/events"
	"github.com/google/uuid"
)

type Service interface {
	CreateMatch(home, away string) (Match, error)
	GetAll() []Match
}

type EventPublisher interface {
	Publish(ctx context.Context, key string, event any) error
}

type service struct {
	repo      Repository
	publisher EventPublisher
}

func NewService(repo Repository, publisher EventPublisher) Service {
	return &service{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *service) CreateMatch(home, away string) (Match, error) {
	match := Match{
		ID:   uuid.NewString(),
		Home: home,
		Away: away,
	}
	err := s.repo.Save(match)
	if err != nil {

		return Match{}, err
	}

	if s.publisher != nil {
		evt := events.MatchCreated{
			ID:   match.ID,
			Home: match.Home,
			Away: match.Away,
		}
		if err := s.publisher.Publish(context.Background(), match.ID, evt); err != nil {
			return Match{}, err
		}
	}

	return match, nil
}

func (s *service) GetAll() []Match {
	return s.repo.GetAll()
}
