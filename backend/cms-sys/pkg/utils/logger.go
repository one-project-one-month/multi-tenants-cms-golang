package utils

import (
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogConfig struct {
	Level      string
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
	Console    bool
}

func NewLogger(config LogConfig) *logrus.Logger {
	logger := logrus.New()

	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "caller",
		},
	})

	if config.FilePath != "" {
		if err := os.MkdirAll(filepath.Dir(config.FilePath), 0755); err != nil {
			logger.WithError(err).Error("Failed to create log directory")
		}

		lumberjackLogger := &lumberjack.Logger{
			Filename:   config.FilePath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}

		if config.Console {
			multiWriter := io.MultiWriter(os.Stdout, lumberjackLogger)
			logger.SetOutput(multiWriter)
		} else {
			logger.SetOutput(lumberjackLogger)
		}
	} else {
		logger.SetOutput(os.Stdout)
	}

	logger.SetReportCaller(true)

	return logger
}

func GetDefaultLogConfig() LogConfig {
	return LogConfig{
		Level:      "info",
		FilePath:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
		Console:    true,
	}
}
