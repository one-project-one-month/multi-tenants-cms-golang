package app

import (
	"context"
	gate "github.com/multi-tenants-cms-golang/lms-sys/app/gateway"
	"github.com/multi-tenants-cms-golang/lms-sys/app/rpc"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type App struct {
	grpcServer  *rpc.Server
	grpcGateway *gate.Gateway
	logger      *logrus.Logger
}

func NewApp(
	grpcServer *rpc.Server,
	grpcGateway *gate.Gateway,
	logger *logrus.Logger,
) *App {
	return &App{
		grpcServer:  grpcServer,
		grpcGateway: grpcGateway,
		logger:      logger,
	}
}

func (app *App) Run() error {
	// Create context that cancels on interrupt signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel for errors from goroutines
	errChan := make(chan error, 2)

	// WaitGroup to wait for both servers to shutdown
	var wg sync.WaitGroup
	wg.Add(2)

	// Start gRPC server
	go func() {
		defer wg.Done()
		app.logger.Info("Starting gRPC server...")
		if err := app.grpcServer.Run(); err != nil {
			errChan <- err
			cancel() // Trigger shutdown on error
		}
	}()

	// Brief delay to ensure gRPC server is up before gateway
	time.Sleep(1 * time.Second)

	// Start gRPC gateway
	go func() {
		defer wg.Done()
		app.logger.Info("Starting gRPC gateway...")
		if err := app.grpcGateway.Start(); err != nil {
			errChan <- err
			cancel() // Trigger shutdown on error
		}
	}()

	// Signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		app.logger.Infof("Received signal %v, shutting down...", sig)
		cancel()
	case err := <-errChan:
		app.logger.Errorf("Server error: %v", err)
		cancel()
		return err
	case <-ctx.Done():
		// Normal shutdown
	}

	// Wait for servers to shutdown
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Add shutdown timeout
	select {
	case <-done:
		app.logger.Info("Servers stopped gracefully")
	case <-time.After(10 * time.Second):
		app.logger.Warn("Forced shutdown after timeout")
	}

	return nil
}
