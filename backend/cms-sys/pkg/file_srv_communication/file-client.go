package file_srv_communication

import (
	"context"
	"fmt"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

type FileClient struct {
	client  *resty.Client
	baseURL string
	timeout time.Duration
}

// NewFileClient creates a new file service client
func NewFileClient(baseURL string, timeout time.Duration) *FileClient {
	return &FileClient{
		client: resty.New().
			SetBaseURL(baseURL).
			SetTimeout(timeout).
			SetHeader("Accept", "application/json"),
		baseURL: baseURL,
		timeout: timeout,
	}
}

func (fc *FileClient) UploadFile(ctx context.Context, req types.FileUploadRequest) (*types.FileUploadResponse, error) {
	if req.FilePath == "" {
		return nil, fmt.Errorf("file path is required")
	}

	file, err := os.Open(req.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	err = file.Close()
	if err != nil {
		fmt.Printf("failed to close file: %v\n", err)
		return nil, err
	}

	resp, err := fc.client.R().
		SetContext(ctx).
		SetFile("file", req.FilePath).
		SetFormData(map[string]string{
			"directory": req.Directory,
			"user_id":   req.UserID,
			"overwrite": fmt.Sprintf("%t", req.Overwrite),
		}).
		SetResult(&types.FileUploadResponse{}).
		Post("/api/v1/files")

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("upload failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	result, ok := resp.Result().(*types.FileUploadResponse)
	if !ok {
		return nil, fmt.Errorf("failed to parse response")
	}

	return result, nil
}
