package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/DMiljevic1/pulsebet/internal/config"
	"github.com/DMiljevic1/pulsebet/internal/httpserver"
	"github.com/DMiljevic1/pulsebet/internal/logging"
)

func main() {
	cfg, err := config.Load("configs/gameservice.yml")
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	logger := logging.New(cfg.ServiceName)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	server := httpserver.New(cfg.HTTPPort, mux)

	if err := server.Start(); err != nil {
		logger.Error("Failed to start server: ", err)
	}
}
