package service

import (
	"context"
	"fmt"
	"github.com/google/go-github/v57/github"
	"github.com/google/uuid"
	"github.com/multi-tenants-cms-golang/file-service/internal/types"
	"github.com/sirupsen/logrus"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileSystemService interface {
	UploadFile(context.Context, []byte, *types.UploadOptions) (*types.FileInfo, error)
	DownloadFile(context.Context, string) ([]byte, error)
	DeleteFile(context.Context, string) error
	ListFiles(context.Context, string) ([]*types.FileInfo, error)
	UpdateFile(context.Context, string, []byte, *types.UploadOptions) (*types.FileInfo, error)
	generateStoragePath(*types.UploadOptions) string
	sanitizeFilename(string) string
	detectContentType(string, []byte) string
	generateFileID(sha string) string
}

type FileSystemServiceImpl struct {
	logger   *logrus.Logger
	ghClient *github.Client
	owner    string
	repo     string
	branch   string
}

func NewFileSystem(
	logger *logrus.Logger,
	ghClient *github.Client,
) *FileSystemServiceImpl {
	return &FileSystemServiceImpl{
		logger:   logger,
		ghClient: ghClient,
		owner:    os.Getenv("GITHUB_OWNER"),
		repo:     os.Getenv("GITHUB_REPO"),
		branch:   getEnvOrDefault("GITHUB_BRANCH", "main"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

var _ FileSystemService = (*FileSystemServiceImpl)(nil)

func (f FileSystemServiceImpl) UploadFile(ctx context.Context, data []byte, options *types.UploadOptions) (*types.FileInfo, error) {
	storagePath := f.generateStoragePath(options)

	if options.ContentType == "" {
		options.ContentType = f.detectContentType(options.Filename, data)
	}

	existing, _, _, err := f.ghClient.Repositories.GetContents(ctx, f.owner, f.repo, storagePath, &github.RepositoryContentGetOptions{
		Ref: f.branch,
	})

	commitMessage := fmt.Sprintf("Upload %s", options.Filename)
	if options.Directory != "" {
		commitMessage = fmt.Sprintf("Upload %s to %s", options.Filename, options.Directory)
	}

	opts := &github.RepositoryContentFileOptions{
		Message: github.String(commitMessage),
		Content: data,
		Branch:  github.String(f.branch),
	}

	var result *github.RepositoryContentResponse

	if err == nil && existing != nil {
		if !options.Overwrite {
			return nil, fmt.Errorf("file already exists: %s", storagePath)
		}
		opts.SHA = existing.SHA
		result, _, err = f.ghClient.Repositories.UpdateFile(ctx, f.owner, f.repo, storagePath, opts)
	} else {
		result, _, err = f.ghClient.Repositories.CreateFile(ctx, f.owner, f.repo, storagePath, opts)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %v", err)
	}

	return &types.FileInfo{
		ID:           f.generateFileID(result.Content.GetSHA()),
		OriginalName: options.Filename,
		StoredName:   filepath.Base(storagePath),
		Path:         result.Content.GetPath(),
		Size:         result.Content.GetSize(),
		ContentType:  options.ContentType,
		URL:          result.Content.GetDownloadURL(),
		SHA:          result.Content.GetSHA(),
		Metadata:     options.Metadata,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (f FileSystemServiceImpl) DownloadFile(ctx context.Context, path string) ([]byte, error) {
	content, _, _, err := f.ghClient.Repositories.GetContents(ctx, f.owner, f.repo, path, &github.RepositoryContentGetOptions{
		Ref: f.branch,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}

	decodedContent, err := content.GetContent()
	if err != nil {
		return nil, fmt.Errorf("failed to decode content: %v", err)
	}

	return []byte(decodedContent), nil
}

func (f FileSystemServiceImpl) DeleteFile(ctx context.Context, path string) error {
	existing, _, _, err := f.ghClient.Repositories.GetContents(ctx, f.owner, f.repo, path, &github.RepositoryContentGetOptions{
		Ref: f.branch,
	})
	if err != nil {
		return fmt.Errorf("failed to get file for deletion: %v", err)
	}

	opts := &github.RepositoryContentFileOptions{
		Message: github.String(fmt.Sprintf("Delete %s", filepath.Base(path))),
		SHA:     existing.SHA,
		Branch:  github.String(f.branch),
	}

	_, _, err = f.ghClient.Repositories.DeleteFile(ctx, f.owner, f.repo, path, opts)
	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

func (f FileSystemServiceImpl) ListFiles(ctx context.Context, directory string) ([]*types.FileInfo, error) {
	_, contents, _, err := f.ghClient.Repositories.GetContents(ctx, f.owner, f.repo, directory, &github.RepositoryContentGetOptions{
		Ref: f.branch,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list directory: %v", err)
	}

	var files []*types.FileInfo
	for _, content := range contents {
		if content.GetType() == "file" {
			files = append(files, &types.FileInfo{
				ID:         f.generateFileID(content.GetSHA()),
				StoredName: content.GetName(),
				Path:       content.GetPath(),
				Size:       content.GetSize(),
				URL:        content.GetDownloadURL(),
				SHA:        content.GetSHA(),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			})
		}
	}

	return files, nil
}

func (f FileSystemServiceImpl) UpdateFile(ctx context.Context, path string, data []byte, options *types.UploadOptions) (*types.FileInfo, error) {
	existing, _, _, err := f.ghClient.Repositories.GetContents(ctx, f.owner, f.repo, path, &github.RepositoryContentGetOptions{
		Ref: f.branch,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get existing file: %v", err)
	}

	if options.ContentType == "" {
		options.ContentType = f.detectContentType(options.Filename, data)
	}

	opts := &github.RepositoryContentFileOptions{
		Message: github.String(fmt.Sprintf("Update %s", filepath.Base(path))),
		Content: data,
		SHA:     existing.SHA,
		Branch:  github.String(f.branch),
	}

	result, _, err := f.ghClient.Repositories.UpdateFile(ctx, f.owner, f.repo, path, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to update file: %v", err)
	}

	return &types.FileInfo{
		ID:           f.generateFileID(result.Content.GetSHA()),
		OriginalName: options.Filename,
		StoredName:   filepath.Base(path),
		Path:         result.Content.GetPath(),
		Size:         result.Content.GetSize(),
		ContentType:  options.ContentType,
		URL:          result.Content.GetDownloadURL(),
		SHA:          result.Content.GetSHA(),
		Metadata:     options.Metadata,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (f FileSystemServiceImpl) generateStoragePath(options *types.UploadOptions) string {
	filename := f.sanitizeFilename(options.Filename)
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	timestamp := time.Now().Format("20060102_150405")
	uniqueID := uuid.New().String()[:8]

	storedName := fmt.Sprintf("%s_%s_%s%s", nameWithoutExt, timestamp, uniqueID, ext)

	if options.Directory != "" {
		return filepath.Join(options.Directory, storedName)
	}

	return storedName
}

func (f FileSystemServiceImpl) sanitizeFilename(filename string) string {
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, ":", "_")
	filename = strings.ReplaceAll(filename, "*", "_")
	filename = strings.ReplaceAll(filename, "?", "_")
	filename = strings.ReplaceAll(filename, "\"", "_")
	filename = strings.ReplaceAll(filename, "<", "_")
	filename = strings.ReplaceAll(filename, ">", "_")
	filename = strings.ReplaceAll(filename, "|", "_")

	return filename
}

func (f FileSystemServiceImpl) detectContentType(filename string, data []byte) string {
	if contentType := mime.TypeByExtension(filepath.Ext(filename)); contentType != "" {
		return contentType
	}
	return http.DetectContentType(data)
}

func (f FileSystemServiceImpl) generateFileID(sha string) string {
	return fmt.Sprintf("gh_%s", sha[:16])
}
