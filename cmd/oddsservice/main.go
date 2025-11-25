package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DMiljevic1/pulsebet/internal/config"
	db2 "github.com/DMiljevic1/pulsebet/internal/db"
	oddshttp "github.com/DMiljevic1/pulsebet/internal/http/odds"
	"github.com/DMiljevic1/pulsebet/internal/httpserver"
	"github.com/DMiljevic1/pulsebet/internal/logging"
	"github.com/DMiljevic1/pulsebet/internal/services/odds"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	cfg, err := config.Load("configs/oddsservice.yml")
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.New(cfg.ServiceName)

	db, err := db2.Connect(cfg.Database)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		return
	}
	defer db.Close()

	repo := odds.NewRepository(db)
	service := odds.NewService(repo)

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	r.Handle("/metrics", promhttp.Handler())

	oddshttp.RegisterRoutes(r, service, logger)

	server := httpserver.New(cfg.HTTPPort, r)

	if err := server.Start(); err != nil {
		logger.Error("server stopped", "error", err)
	}
}
