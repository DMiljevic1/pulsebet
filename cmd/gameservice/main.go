package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DMiljevic1/pulsebet/internal/config"
	gamehttp "github.com/DMiljevic1/pulsebet/internal/http/game"
	"github.com/DMiljevic1/pulsebet/internal/httpserver"
	"github.com/DMiljevic1/pulsebet/internal/logging"
	"github.com/DMiljevic1/pulsebet/internal/services/game"
)

func main() {
	// 1) Config
	cfg, err := config.Load("configs/gameservice.yml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2) Logger
	logger := logging.New(cfg.ServiceName)

	// 3) Domain service
	gameService := game.NewService()

	// 4) HTTP router
	mux := http.NewServeMux()

	// health endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	// game HTTP handlers
	gamehttp.RegisterHandlers(mux, gameService)

	// 5) HTTP server
	server := httpserver.New(cfg.HTTPPort, mux)

	if err := server.Start(); err != nil {
		logger.Error("server stopped", "error", err)
	}
}
