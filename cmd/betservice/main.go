package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DMiljevic1/pulsebet/internal/config"
	db2 "github.com/DMiljevic1/pulsebet/internal/db"
	bethttp "github.com/DMiljevic1/pulsebet/internal/http/bet"
	"github.com/DMiljevic1/pulsebet/internal/httpserver"
	"github.com/DMiljevic1/pulsebet/internal/logging"
	"github.com/DMiljevic1/pulsebet/internal/services/bet"
)

func main() {
	cfg, err := config.Load("configs/betservice.yaml")
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.New(cfg.ServiceName)

	db, err := db2.Connect(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		return
	}
	defer db.Close()

	repo := bet.NewRepository(db)

	service := bet.NewService(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	bethttp.RegisterHandlers(mux, service, logger)

	server := httpserver.New(cfg.HTTPPort, mux)

	if err := server.Start(); err != nil {
		logger.Error("Failed to start server", "error", err)
	}
}
