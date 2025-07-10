package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server ServerConfig
	NATS   NATSConfig
	SMTP   SMTPConfig
}

type ServerConfig struct {
	Port            string
	ShutdownTimeout time.Duration
}

type NATSConfig struct {
	URL          string
	StreamName   string
	Subject      string
	ConsumerName string
}

type SMTPConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	FromAddr    string
	TemplateDir string
}

func Load() (*Config, error) {
	smtpPortStr := os.Getenv("SMTP_PORT")
	if smtpPortStr == "" {
		smtpPortStr = "587"
	}

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_PORT: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:            getEnvOrDefault("SERVER_PORT", "8080"),
			ShutdownTimeout: 30 * time.Second,
		},
		NATS: NATSConfig{
			URL:          getEnvOrDefault("NATS_URL", "nats://localhost:4222"),
			StreamName:   getEnvOrDefault("STREAM_NAME", "email-stream"),
			Subject:      getEnvOrDefault("SUBJECT_TO", "email.send"),
			ConsumerName: getEnvOrDefault("CONSUMER_NAME", "email-consumer"),
		},
		SMTP: SMTPConfig{
			Host:        getEnvOrDefault("SMTP_HOST", "localhost"),
			Port:        smtpPort,
			User:        os.Getenv("SMTP_USER"),
			Password:    os.Getenv("SMTP_PASSWORD"),
			FromAddr:    os.Getenv("FROM_ADDR"),
			TemplateDir: getEnvOrDefault("TEMPLATE_DIR", "./templates"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("convert validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.SMTP.User == "" {
		return fmt.Errorf("SMTP_USER is required")
	}
	if c.SMTP.Password == "" {
		return fmt.Errorf("SMTP_PASSWORD is required")
	}
	if c.SMTP.FromAddr == "" {
		return fmt.Errorf("FROM_ADDR is required")
	}
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
