package types

import "time"

type UploadOptions struct {
	Filename    string            `json:"filename"`
	Directory   string            `json:"directory"`
	ContentType string            `json:"content_type"`
	Metadata    map[string]string `json:"metadata"`
	Overwrite   bool              `json:"overwrite"`
}

type FileInfo struct {
	ID           string            `json:"id"`
	OriginalName string            `json:"original_name"`
	StoredName   string            `json:"stored_name"`
	Path         string            `json:"path"`
	Size         int               `json:"size"`
	ContentType  string            `json:"content_type"`
	URL          string            `json:"url"`
	SHA          string            `json:"sha"`
	Metadata     map[string]string `json:"metadata"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}
