package app

import (
	gate "github.com/multi-tenants-cms-golang/lms-sys/app/gateway"
	rpc "github.com/multi-tenants-cms-golang/lms-sys/app/rpc"
	"github.com/sirupsen/logrus"
	"sync"
)

type App struct {
	grpcServer  *rpc.Server
	grpcGateway *gate.Gateway
	logger      *logrus.Logger
}

func NewApp(
	grpcServer *rpc.Server,
	grpcGateway *gate.Gateway,
) *App {
	return &App{
		grpcServer:  grpcServer,
		grpcGateway: grpcGateway,
	}
}

func (app *App) Run() error {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := app.grpcServer.Run(); err != nil {
			if app.logger != nil {
				app.logger.Error(err.Error())
			}
		}
	}()

	go func() {
		defer wg.Done()
		if err := app.grpcGateway.Start(); err != nil {
			if app.logger != nil {
				app.logger.Error(err.Error())
			}
		}
	}()

	wg.Wait()
	if app.logger != nil {
		app.logger.Info("server stopped")
	}
	return nil
}
