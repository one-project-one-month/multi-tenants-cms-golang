package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/task"
	"github.com/sirupsen/logrus"
)

type IPFSHandler struct {
	logger *logrus.Logger
	shell  *shell.Shell
	//s      repo.Store
}

func NewIPFSHandler(
	logger *logrus.Logger,
	shell *shell.Shell,
	// store repo.Store,
) *IPFSHandler {
	return &IPFSHandler{
		logger: logger,
		shell:  shell,
		//s:      store,
	}
}
func (h *IPFSHandler) HandleIPFSUpload(ctx context.Context, t *asynq.Task) error {
	var payload *task.TaskFileUpload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		h.logger.Error(err.Error())
		return err
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	h.logger.Info(string(payloadBytes))
	h.logger.Info("uploading to the ipfs ......")
	uploadedFile, err := h.shell.Add(bytes.NewReader(payloadBytes), shell.Pin(true))
	if err != nil {
		h.logger.Error(err.Error())
		return err
	}
	h.logger.Info("successfully uploaded the file")
	h.logger.Info("HASH of the file", uploadedFile)

	return nil
}
