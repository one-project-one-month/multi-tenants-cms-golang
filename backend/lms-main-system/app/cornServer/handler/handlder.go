package handler

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

type HandlerInterface interface {
	HandleDatabaseBackUp(ctx context.Context, req *asynq.Task) error
}

var _ HandlerInterface = (*CornHandler)(nil)

type CornHandler struct {
	logger *logrus.Logger
}

func NewCornHandler(
	logger *logrus.Logger,
) *CornHandler {
	return &CornHandler{
		logger: logger,
	}
}
