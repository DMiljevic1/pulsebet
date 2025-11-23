package bet

import (
	"context"
	"errors"

	"github.com/DMiljevic1/pulsebet/internal/events"
	"github.com/google/uuid"
)

type Service interface {
	SaveBet(bet Bet) (Bet, error)
	HandleMatchCreated(ctx context.Context, key string, evt events.MatchCreated) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) SaveBet(bet Bet) (Bet, error) {
	if bet.MatchID == "" {
		return Bet{}, errors.New("match_id id is required")
	}
	if bet.BettorID == "" {
		return Bet{}, errors.New("bettor_id is required")
	}
	if bet.Selection == "" {
		return Bet{}, errors.New("selection is required")
	}
	if bet.Stake <= 0 {
		return Bet{}, errors.New("stake must be greater than zero")
	}

	if bet.ID == "" {
		bet.ID = uuid.NewString()
	}
	if bet.Odds <= 0 {
		bet.Odds = 1
	}

	err := s.repo.SaveBet(bet)
	if err != nil {
		return Bet{}, err
	}
	return bet, nil
}

func (s *service) HandleMatchCreated(ctx context.Context, key string, evt events.MatchCreated) error {
	m := Match{
		ID:   evt.ID,
		Home: evt.Home,
		Away: evt.Away,
	}

	if err := s.repo.SaveAvailableMatch(ctx, m); err != nil {
		return err
	}
	return nil
}
