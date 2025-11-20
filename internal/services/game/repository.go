package game

import (
	"database/sql"
)

type Repository interface {
	Save(match Match) error
	GetAll() []Match
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Save(match Match) error {
	_, err := r.db.Exec(
		`INSERT INTO matches (id, home, away) VALUES ($1, $2, $3)`,
		match.ID,
		match.Home,
		match.Away,
	)
	return err
}

func (r *repository) GetAll() []Match {
	rows, err := r.db.Query(`SELECT id, home, away FROM matches`)
	if err != nil {
		return []Match{}
	}
	defer rows.Close()

	var result []Match
	for rows.Next() {
		var m Match
		if err := rows.Scan(&m.ID, &m.Home, &m.Away); err != nil {
			continue
		}
		result = append(result, m)
	}

	return result
}
