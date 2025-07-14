package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/multi-tenants-cms-golang/file-service/config"
	"github.com/multi-tenants-cms-golang/file-service/internal/handler"
	"github.com/multi-tenants-cms-golang/file-service/internal/routes"
	"github.com/sirupsen/logrus"
)

type ServerState struct {
	cfg       *config.Configuration
	fsHandler handler.FileSystemHandler
	app       *fiber.App
	logger    *logrus.Logger
}

func NewServerState(
	cfg *config.Configuration,

	fsHandler handler.FileSystemHandler,
	app *fiber.App,
	logger *logrus.Logger,
) *ServerState {
	return &ServerState{
		cfg,
		fsHandler,
		app,
		logger,
	}
}

func (state *ServerState) StartServer() {
	state.app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusOK)
	})
	state.logger.Info("Setting up the route")
	routes.FileRoute(state.app, state.fsHandler)
	state.logger.Info("Starting server on http://localhost:" + state.cfg.Port)
	if err := state.app.Listen(":" + state.cfg.Port); err != nil {
		state.logger.Fatal(err.Error())
	}
}

func (state *ServerState) StopServer() {
	state.logger.Info("Shutting down the server")
	err := state.app.Shutdown()
	if err != nil {
		state.logger.Fatal(err.Error())
		return
	}
}
