package bet

import (
	"context"
	"database/sql"
	"time"
)

type Repository interface {
	SaveBet(bet Bet) error
	SaveAvailableMatch(ctx context.Context, m Match) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (repo *repository) SaveBet(bet Bet) error {
	bet.PlacedAt = time.Now()
	const query = `
		INSERT INTO bets (id, match_id, bettor_id, selection, stake, odds, placed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
	_, err := repo.db.Exec(query,
		bet.ID,
		bet.MatchID,
		bet.BettorID,
		bet.Selection,
		bet.Stake,
		bet.Odds,
		bet.PlacedAt)
	if err != nil {
		return err
	}
	return nil
}

func (repo *repository) SaveAvailableMatch(ctx context.Context, m Match) error {
	const query = `
		INSERT INTO matches (id, home, away)
		VALUES ($1, $2, $3);
	`
	_, err := repo.db.Exec(query,
		m.ID,
		m.Home,
		m.Away,
	)
	if err != nil {
		return err
	}
	return nil
}
