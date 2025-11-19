package httpserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

func New(port int, handler http.Handler) *Server {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{server: server}
}

func (s *Server) Start() error {
	log.Printf("Starting server on port %d", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	log.Printf("Stopping server on port %d", s.server.Addr)
	return s.server.Shutdown(ctx)
}
