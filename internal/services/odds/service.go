package odds

import (
	"context"
	"errors"
)

type Service interface {
	GetOdds(ctx context.Context, matchID string) (Odds, error)
	SetOdds(ctx context.Context, odds Odds) (Odds, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetOdds(ctx context.Context, matchID string) (Odds, error) {
	if matchID == "" {
		return Odds{}, errors.New("matchID is required")
	}

	return s.repo.GetByMatchID(ctx, matchID)
}

func (s *service) SetOdds(ctx context.Context, odds Odds) (Odds, error) {
	if odds.MatchID == "" {
		return Odds{}, errors.New("matchID is required")
	}
	if odds.HomeOdds <= 1 || odds.AwayOdds <= 1 || odds.DrawOdds <= 1 {
		return Odds{}, errors.New("all odds must be greater than 1")
	}
	if err := s.repo.Save(ctx, odds); err != nil {
		return Odds{}, err
	}
	return odds, nil
}
