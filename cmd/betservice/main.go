package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DMiljevic1/pulsebet/internal/config"
	"github.com/DMiljevic1/pulsebet/internal/db"
	"github.com/DMiljevic1/pulsebet/internal/events"
	bethttp "github.com/DMiljevic1/pulsebet/internal/http/bet"
	"github.com/DMiljevic1/pulsebet/internal/httpserver"
	"github.com/DMiljevic1/pulsebet/internal/kafka"
	"github.com/DMiljevic1/pulsebet/internal/logging"
	"github.com/DMiljevic1/pulsebet/internal/services/bet"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	segmentio "github.com/segmentio/kafka-go"
)

func main() {
	cfg, err := config.Load("configs/betservice.yml")
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.New(cfg.ServiceName)

	db, err := db.ConnectWithRetry(logger, cfg.Database, 5, 3*time.Second)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		return
	}
	defer db.Close()

	repo := bet.NewRepository(db)

	service := bet.NewService(repo)

	consumer := kafka.NewConsumer[events.MatchCreated](
		cfg.Kafka.Brokers,
		cfg.Kafka.Topics.MatchCreated,
		cfg.Kafka.GroupID,
		4,
		logger,
	)

	defer consumer.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		consumerErr := consumer.Consume(ctx, func(ctx context.Context, key string, evt events.MatchCreated, msg segmentio.Message) error {
			return service.HandleMatchCreated(ctx, key, evt)
		})

		if consumerErr != nil && consumerErr != context.Canceled {
			logger.Error("consumer stopped", "error", consumerErr)
		}
	}()

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	bethttp.RegisterHandlers(mux, service, logger)

	server := httpserver.New(cfg.HTTPPort, mux)

	if err := server.Start(); err != nil {
		logger.Error("Failed to start server", "error", err)
	}
}
