package health

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	server *http.Server
	logger *zap.Logger
}

func NewServer(port string, logger *zap.Logger) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/ready", readyHandler)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		server: server,
		logger: logger,
	}
}

func (s *Server) Start() error {
	s.logger.Info("starting health server", zap.String("addr", s.server.Addr))
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("stopping health server")
	return s.server.Shutdown(ctx)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
}
