package app

import (
	"context"
	"github.com/multi-tenants-cms-golang/lms-sys/app/cornServer"
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
	cronServer  *cornServer.Server
}

func NewApp(
	grpcServer *rpc.Server,
	grpcGateway *gate.Gateway,
	backupServer *cornServer.Server,
	logger *logrus.Logger,
) *App {
	return &App{
		grpcServer:  grpcServer,
		grpcGateway: grpcGateway,
		cronServer:  backupServer,
		logger:      logger,
	}
}

func (app *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 3)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		app.logger.Info("Starting cron  cron server...")
		if err := app.cronServer.RunServer(); err != nil {
			errChan <- err
			cancel()
		}
	}()
	go func() {
		defer wg.Done()
		app.logger.Info("Starting gRPC server...")
		if err := app.grpcServer.Run(); err != nil {
			errChan <- err
			cancel()
		}
	}()

	time.Sleep(1 * time.Second)

	go func() {
		defer wg.Done()
		app.logger.Info("Starting gRPC gateway...")
		if err := app.grpcGateway.Start(); err != nil {
			errChan <- err
			cancel()
		}
	}()

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
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		app.logger.Info("Servers stopped gracefully")
	case <-time.After(10 * time.Second):
		app.logger.Warn("Forced shutdown after timeout")
	}

	return nil
}
