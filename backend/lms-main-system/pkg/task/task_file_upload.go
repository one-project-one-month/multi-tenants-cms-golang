package task

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

const (
	TypeFileUpload = "ipfs:upload"
)

const (
	COURSE_CONTENT = "COURSE_CONTENT"
	PROFILE        = "PROFILE"
	ROLE           = "ROLE"
)

type MetaData struct {
	Namespace    string    `json:"namespace"`
	Type         string    `json:"type"`
	ROLE         string    `json:"role"`
	FILENAME     string    `json:"filename"`
	EntityToSave string    `json:"entityToSave"`
	Uploader     uuid.UUID `json:"uploader"`
}
type TaskFileUpload struct {
	MetaData MetaData `json:"metaData"`
	Content  []byte   `json:"content"`
}

func NewTaskFileUpload(metaData MetaData, content []byte) *asynq.Task {
	payload := TaskFileUpload{
		MetaData: metaData,
		Content:  content,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		logrus.Error(err.Error())
	}
	return asynq.NewTask(TypeFileUpload, payloadBytes)
}
