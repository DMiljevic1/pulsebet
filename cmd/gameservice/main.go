package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DMiljevic1/pulsebet/internal/config"
	"github.com/DMiljevic1/pulsebet/internal/db"
	gamehttp "github.com/DMiljevic1/pulsebet/internal/http/game"
	"github.com/DMiljevic1/pulsebet/internal/httpserver"
	kafkapkg "github.com/DMiljevic1/pulsebet/internal/kafka"
	"github.com/DMiljevic1/pulsebet/internal/logging"
	"github.com/DMiljevic1/pulsebet/internal/services/game"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	db, err := db.ConnectWithRetry(logger, cfg.Database, 5, 3*time.Second)
	if err != nil {
		logger.Error("Failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	// 4) Repository
	gameRepo := game.NewRepository(db)

	// 5) Kafka producer
	producer := kafkapkg.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topics.MatchCreated, logger)
	fmt.Printf("Kafka topic: %q\n", cfg.Kafka.Topics.MatchCreated)
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

	mux.Handle("/metrics", promhttp.Handler())

	if err := server.Start(); err != nil {
		logger.Error("server stopped", "error", err)
	}
}
