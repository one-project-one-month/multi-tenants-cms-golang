package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/multi-tenants-cms-golang/email-service/config"
	"github.com/multi-tenants-cms-golang/email-service/internal/email"
	"github.com/multi-tenants-cms-golang/email-service/internal/health"
	natss "github.com/multi-tenants-cms-golang/email-service/internal/nats"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			panic(err.Error())
		}
	}(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load configuration", zap.Error(err))
	}

	// Connect to NATS
	nc, err := nats.Connect(cfg.NATS.URL,
		nats.ReconnectWait(time.Second),
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.Error("NATS disconnected", zap.Error(err))
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Info("NATS reconnected")
		}),
	)
	if err != nil {
		logger.Fatal("failed to connect to NATS", zap.Error(err))
	}
	defer nc.Close()

	emailCh := make(chan natss.EmailRequest, 100)
	emailService, err := email.NewService(
		logger,
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.User,
		cfg.SMTP.Password,
		cfg.SMTP.FromAddr,
		cfg.SMTP.TemplateDir,
	)
	if err != nil {
		logger.Fatal("failed to create email service", zap.Error(err))
	}

	// Create NATS consumer
	natsConsumer, err := natss.NewConsumer(
		nc,
		logger,
		emailCh,
		cfg.NATS.StreamName,
		cfg.NATS.Subject,
		cfg.NATS.ConsumerName,
	)
	if err != nil {
		logger.Fatal("failed to create NATS consumer", zap.Error(err))
	}

	// Create contexts
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// Start worker goroutines
	workerCount := runtime.NumCPU()
	logger.Info("starting email workers", zap.Int("count", workerCount))

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			logger.Info("worker started", zap.Int("id", workerID))

			for {
				select {
				case req, ok := <-emailCh:
					if !ok {
						logger.Info("worker stopping", zap.Int("id", workerID))
						return
					}

					if err := emailService.Send(ctx, req); err != nil {
						logger.Error("failed to send email",
							zap.Int("worker", workerID),
							zap.String("to", req.To),
							zap.String("template", req.Template),
							zap.Error(err))
					}
				case <-ctx.Done():
					logger.Info("worker stopping due to context cancellation", zap.Int("id", workerID))
					return
				}
			}
		}(i)
	}

	// Start NATS consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := natsConsumer.Start(ctx); err != nil {
			logger.Error("NATS consumer failed", zap.Error(err))
		}
	}()

	// Start health server
	healthServer := health.NewServer(cfg.Server.Port, logger)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := healthServer.Start(); !errors.Is(err, http.ErrServerClosed) && err != nil {
			logger.Error("health server failed", zap.Error(err))
		}
	}()

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("shutdown signal received, starting graceful shutdown")

	// Create shutdown context
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	// Stop health server
	if err := healthServer.Stop(shutdownCtx); err != nil {
		logger.Error("failed to stop health server gracefully", zap.Error(err))
	}

	// Cancel main context to stop all goroutines
	cancel()

	close(emailCh)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Info("all goroutines stopped")
	case <-shutdownCtx.Done():
		logger.Warn("shutdown timeout exceeded, forcing exit")
	}

	logger.Info("shutdown complete")
}
