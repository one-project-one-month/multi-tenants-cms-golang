package app

//import (
//	"github.com/multi-tenants-cms-golang/lms-sys/app/api"
//	rpcGate "github.com/multi-tenants-cms-golang/lms-sys/app/rpc"
//	"github.com/sirupsen/logrus"
//	"sync"
//)
//
//type App struct {
//	grpcServer  *api.Server
//	grpcGateway *rpcGate.Gateway
//	logger      *logrus.Logger
//}
//
//func NewApp(
//	grpcServer *api.Server,
//	grpcGateway *rpcGate.Gateway,
//) *App {
//	return &App{
//		grpcServer:  grpcServer,
//		grpcGateway: grpcGateway,
//	}
//}
//
//func (app *App) Run() error {
//	var wg sync.WaitGroup
//	wg.Add(2)
//
//	go func() {
//		defer wg.Done()
//		if err := app.grpcServer.Run(); err != nil {
//			if app.logger != nil {
//				app.logger.Error(err.Error())
//			}
//		}
//	}()
//
//	go func() {
//		defer wg.Done()
//		if err := app.grpcGateway.Start(); err != nil {
//			if app.logger != nil {
//				app.logger.Error(err.Error())
//			}
//		}
//	}()
//
//	wg.Wait()
//	if app.logger != nil {
//		app.logger.Info("server stopped")
//	}
//	return nil
//}
