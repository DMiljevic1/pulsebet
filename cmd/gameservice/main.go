package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DMiljevic1/pulsebet/internal/config"
	"github.com/DMiljevic1/pulsebet/internal/db"
	gamehttp "github.com/DMiljevic1/pulsebet/internal/http/game"
	"github.com/DMiljevic1/pulsebet/internal/httpserver"
	kafkapkg "github.com/DMiljevic1/pulsebet/internal/kafka"
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

	// 3) Db connection
	dbConn, err := db.Connect(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database: %v", err)
		return
	}
	defer dbConn.Close()

	// 4) Repository
	gameRepo := game.NewRepository(dbConn)

	// 5) Kafka producer
	producer := kafkapkg.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topics.MatchCreated, logger)
	defer producer.Close()

	// 6) Domain service
	gameService := game.NewService(gameRepo, producer)

	// 7) HTTP router
	mux := http.NewServeMux()

	// health endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	// game HTTP handlers
	gamehttp.RegisterHandlers(mux, gameService, logger)

	// 8) HTTP server
	server := httpserver.New(cfg.HTTPPort, mux)

	if err := server.Start(); err != nil {
		logger.Error("server stopped", "error", err)
	}
}
