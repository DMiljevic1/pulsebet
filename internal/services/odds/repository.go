package odds

import (
	"context"
	"database/sql"
)

type Repository interface {
	GetByMatchID(ctx context.Context, matchID string) (Odds, error)
	Save(ctx context.Context, o Odds) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetByMatchID(ctx context.Context, matchID string) (Odds, error) {
	const query = `
		SELECT match_id, home_odds, draw_odds, away_odds, updated_at
		FROM odds
		WHERE match_id = $1;
	`

	var o Odds
	err := r.db.QueryRowContext(ctx, query, matchID).Scan(
		&o.MatchID,
		&o.HomeOdds,
		&o.DrawOdds,
		&o.AwayOdds,
		&o.UpdatedAt,
	)
	if err != nil {
		return Odds{}, err
	}
	return o, nil
}

func (r *repository) Save(ctx context.Context, o Odds) error {
	const query = `
		INSERT INTO odds (match_id, home_odds, draw_odds, away_odds)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (match_id) DO UPDATE SET
			home_odds  = EXCLUDED.home_odds,
			draw_odds  = EXCLUDED.draw_odds,
			away_odds  = EXCLUDED.away_odds,
			updated_at = now();
	`

	_, err := r.db.ExecContext(ctx, query,
		o.MatchID,
		o.HomeOdds,
		o.DrawOdds,
		o.AwayOdds,
	)
	return err
}
