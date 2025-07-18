package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/klauspost/compress/gzip"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/task"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"time"
)

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

func (handler *CornHandler) HandleDatabaseBackUp(ctx context.Context, req *asynq.Task) error {
	var payload task.DatabaseBackup
	if err := json.Unmarshal(req.Payload(), &payload); err != nil {
		handler.logger.Errorf("Failed to unmarshal payload: %v", err)
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	backupFileName := fmt.Sprintf("backup_%s.sql", time.Now().Format("2006-01-02"))
	if payload.BackUpPath != "" {
		if err := os.MkdirAll(payload.BackUpPath, 0755); err != nil {
			handler.logger.Errorf("Failed to create backup directory: %v", err)
			return fmt.Errorf("failed to create backup directory: %v", err)
		}
		backupFileName = filepath.Join(payload.BackUpPath, backupFileName)
	}

	// Create backup file
	backupFile, err := os.Create(backupFileName)
	if err != nil {
		handler.logger.Errorf("Failed to create backup file: %v", err)
		return fmt.Errorf("failed to create backup file: %v", err)
	}
	defer func(backupFile *os.File) {
		err := backupFile.Close()
		if err != nil {
			logrus.Errorf("Failed to close backup file: %v", err)
		}
	}(backupFile)

	handler.logger.Infof("Starting PostgreSQL backup to %s", backupFileName)

	// Connect to database
	conn, err := pgx.Connect(ctx, payload.DbUrl)
	if err != nil {
		handler.logger.Errorf("Failed to connect to database: %v", err)
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			handler.logger.Errorf("Failed to close connection: %v", err)
		}
	}(conn, ctx)

	rows, err := conn.Query(ctx,
		"SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'")
	if err != nil {
		handler.logger.Errorf("Failed to get tables: %v", err)
		return fmt.Errorf("failed to get tables: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			handler.logger.Errorf("Failed to scan table name: %v", err)
			continue
		}
		tables = append(tables, table)
	}

	for _, table := range tables {

		var schema string
		err := conn.QueryRow(ctx,
			"SELECT pg_get_tabledef($1)", table).Scan(&schema)
		if err != nil {
			handler.logger.Errorf("Failed to get schema for table %s: %v", table, err)
			continue
		}
		if _, err := backupFile.WriteString(schema + ";\n\n"); err != nil {
			handler.logger.Errorf("Failed to write schema for table %s: %v", table, err)
			continue
		}

		if _, err := backupFile.WriteString(
			fmt.Sprintf("COPY %s FROM stdin;\n", table)); err != nil {
			handler.logger.Errorf("Failed to write COPY header for table %s: %v", table, err)
			continue
		}

		rows, err := conn.Query(ctx, fmt.Sprintf("SELECT * FROM %s", table))
		if err != nil {
			handler.logger.Errorf("Failed to query data from table %s: %v", table, err)
			continue
		}
		for rows.Next() {
			vals, err := rows.Values()
			if err != nil {
				handler.logger.Errorf("Failed to get values from table %s: %v", table, err)
				continue
			}
			line := ""
			for i, val := range vals {
				if i > 0 {
					line += "\t"
				}
				line += fmt.Sprintf("%v", val)
			}
			if _, err := backupFile.WriteString(line + "\n"); err != nil {
				handler.logger.Errorf("Failed to write data for table %s: %v", table, err)
				continue
			}
		}
		if _, err := backupFile.WriteString("\\.\n\n"); err != nil {
			handler.logger.Errorf("Failed to write COPY footer for table %s: %v", table, err)
			continue
		}
	}

	handler.logger.Infof("PostgreSQL backup completed successfully: %s", backupFileName)
	return nil
}

func compressBackupFile(filename string) error {
	original, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func(original *os.File) {
		err := original.Close()
		if err != nil {
			logrus.Errorf("Failed to close gzip writer: %v", err)
		}
	}(original)

	compressed, err := os.Create(filename + ".gz")
	if err != nil {
		return err
	}
	defer func(compressed *os.File) {
		err := compressed.Close()
		if err != nil {
			logrus.Errorf("Failed to close gzip writer: %v", err)
		}
	}(compressed)

	gz := gzip.NewWriter(compressed)
	defer func(gz *gzip.Writer) {
		err := gz.Close()
		if err != nil {
			logrus.Errorf("Failed to close gzip writer: %v", err)
		}
	}(gz)

	if _, err := io.Copy(gz, original); err != nil {
		return err
	}

	err = os.Remove(filename)
	if err != nil {
		logrus.Errorf("Failed to remove backups file: %v", err)
		return err
	}
	return nil
}
