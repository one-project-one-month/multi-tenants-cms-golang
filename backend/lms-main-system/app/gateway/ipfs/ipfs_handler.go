package ipfs

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/task"
	"github.com/sirupsen/logrus"
)

type IPFSGateway struct {
	logger      *logrus.Logger
	asynqClient *asynq.Client
	inspector   *asynq.Inspector
}

func NewIPFSGateway(logger *logrus.Logger, redisAddr string) *IPFSGateway {
	return &IPFSGateway{
		logger:      logger,
		asynqClient: asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr}),
	}
}

func (g *IPFSGateway) RegisterHandlers(mux *runtime.ServeMux) {
	mux.HandlePath("POST", "/lms/v1/files/upload", g.handleIPFSUpload)
	mux.HandlePath("GET", "/lms/v1/files/status/{task_id}", g.handleIPFSStatus)
}

func (g *IPFSGateway) handleIPFSUpload(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	r.Body = http.MaxBytesReader(w, r.Body, 50<<20)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		g.logger.WithError(err).Error("Failed to parse multipart form")
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		g.logger.WithError(err).Error("Failed to get file from form")
		http.Error(w, "Invalid file upload", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		g.logger.WithError(err).Error("Failed to read file content")
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	metadata := task.MetaData{
		Namespace:    r.FormValue("namespace"),
		Type:         r.FormValue("type"),
		ROLE:         r.FormValue("role"),
		FILENAME:     header.Filename,
		EntityToSave: r.FormValue("entityToSave"),
	}

	if metadata.Namespace == "" || metadata.Type == "" {
		http.Error(w, "namespace and type are required fields", http.StatusBadRequest)
		return
	}

	metadata.Uploader = uuid.New()

	taskPayload := task.NewTaskFileUpload(metadata, fileContent)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	info, err := g.asynqClient.EnqueueContext(ctx, taskPayload)
	if err != nil {
		g.logger.WithError(err).Error("Failed to enqueue IPFS upload task")
		http.Error(w, "Failed to start upload process", http.StatusInternalServerError)
		return
	}

	g.logger.WithFields(logrus.Fields{
		"task_id":   info.ID,
		"filename":  header.Filename,
		"uploader":  metadata.Uploader.String(),
		"namespace": metadata.Namespace,
	}).Info("IPFS upload task enqueued")

	response := map[string]interface{}{
		"task_id":    info.ID,
		"status":     "enqueued",
		"created_at": info.State,
		"queue":      info.Queue,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}

func (g *IPFSGateway) handleIPFSStatus(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	taskID := pathParams["task_id"]
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	info, err := g.inspector.GetTaskInfo("default", taskID)
	if err != nil {
		if errors.Is(err, asynq.ErrTaskNotFound) {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		g.logger.WithError(err).Error("Failed to get task info")
		http.Error(w, "Failed to get task status", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"task_id":      taskID,
		"status":       info.State.String(),
		"completed_at": info.CompletedAt,
	}

	if info.State == asynq.TaskStateCompleted {
		var result struct {
			CID string `json:"cid"`
		}
		if err := json.Unmarshal(info.Result, &result); err == nil {
			response["cid"] = result.CID
		}
	}

	//// If failed, include error message
	//if info.State == asynq.TaskStateFailed {
	//	response["error"] = info.Error
	//}
	//
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
