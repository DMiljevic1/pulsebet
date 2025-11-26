package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/DMiljevic1/pulsebet/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectWithRetry(logger *slog.Logger, cfg config.DatabaseConfig, attempts int, delay time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 1; i <= attempts; i++ {
		db, err = connect(cfg)
		if err == nil {
			logger.Info("connected to database", "attempt", i)
			return db, nil
		}

		logger.Error("failed to connect to database",
			"attempt", i,
			"error", err,
		)

		time.Sleep(delay)
	}

	return nil, err
}

func connect(cfg config.DatabaseConfig) (*sql.DB, error) {
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
