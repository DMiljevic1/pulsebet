package db

import (
	"database/sql"
	"fmt"

	"github.com/DMiljevic1/pulsebet/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
