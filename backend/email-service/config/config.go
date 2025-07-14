package config

import (
	"os"
	"strconv"

	"time"

	"github.com/pkg/errors"
)

type Config struct {
	Server   ServerConfig
	NATS     NATSConfig
	SMTP     SMTPConfig
	Security SecurityConfig
}

type ServerConfig struct {
	Port            string        `env:"SERVER_PORT" env-default:"8080"`
	ShutdownTimeout time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT" env-default:"30s"`
}

type NATSConfig struct {
	URL          string `env:"NATS_URL" env-default:"nats://localhost:4222"`
	StreamName   string `env:"STREAM_NAME" env-default:"email-stream"`
	Subject      string `env:"SUBJECT_TO" env-default:"email.send"`
	ConsumerName string `env:"CONSUMER_NAME" env-default:"email-consumer"`
}

type SMTPConfig struct {
	Host        string        `env:"SMTP_HOST" env-default:"smtp.gmail.com"`
	Port        int           `env:"SMTP_PORT" env-default:"587"`
	User        string        `env:"SMTP_USER"`
	Password    string        `env:"SMTP_PASSWORD"`
	FromAddr    string        `env:"FROM_ADDR"`
	FromName    string        `env:"FROM_NAME"`
	TemplateDir string        `env:"TEMPLATE_DIR" env-default:"./templates"`
	Timeout     time.Duration `env:"SMTP_TIMEOUT" env-default:"30s"`
	MaxRetries  int           `env:"SMTP_MAX_RETRIES" env-default:"3"`
	RetryDelay  time.Duration `env:"SMTP_RETRY_DELAY" env-default:"5s"`
	PoolSize    int           `env:"SMTP_POOL_SIZE" env-default:"5"`
}

type SecurityConfig struct {
	TLSEnabled bool   `env:"SMTP_TLS_ENABLED" env-default:"true"`
	Insecure   bool   `env:"SMTP_INSECURE" env-default:"false"`
	ServerName string `env:"SMTP_SERVER_NAME"`
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8080"),
			ShutdownTimeout: parseDuration(getEnv("SERVER_SHUTDOWN_TIMEOUT", "30s")),
		},
		NATS: NATSConfig{
			URL:          getEnv("NATS_URL", "nats://localhost:4222"),
			StreamName:   getEnv("STREAM_NAME", "email-stream"),
			Subject:      getEnv("SUBJECT_TO", "email.send"),
			ConsumerName: getEnv("CONSUMER_NAME", "email-consumer"),
		},
		SMTP: SMTPConfig{
			Host:        getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:        parseInt(getEnv("SMTP_PORT", "587")),
			User:        getEnv("SMTP_USER", ""),
			Password:    getEnv("SMTP_PASSWORD", ""),
			FromAddr:    getEnv("FROM_ADDR", ""),
			FromName:    getEnv("FROM_NAME", ""),
			TemplateDir: getEnv("TEMPLATE_DIR", "./templates"),
			Timeout:     parseDuration(getEnv("SMTP_TIMEOUT", "30s")),
			MaxRetries:  parseInt(getEnv("SMTP_MAX_RETRIES", "3")),
			RetryDelay:  parseDuration(getEnv("SMTP_RETRY_DELAY", "5s")),
			PoolSize:    parseInt(getEnv("SMTP_POOL_SIZE", "5")),
		},
		Security: SecurityConfig{
			TLSEnabled: parseBool(getEnv("SMTP_TLS_ENABLED", "true")),
			Insecure:   parseBool(getEnv("SMTP_INSECURE", "false")),
			ServerName: getEnv("SMTP_SERVER_NAME", ""),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, errors.Wrap(err, "configuration validation failed")
	}

	return cfg, nil
}

func parseBool(env string) bool {
	if v := os.Getenv(env); v != "" {
		return v == "true"
	}
	return false
}

func parseDuration(env string) time.Duration {
	if v := os.Getenv(env); v != "" {
		duration, err := time.ParseDuration(v)
		if err != nil {
			panic(err)
		}
		return duration
	}
	return 0
}

func parseInt(env string) int {

	if v := os.Getenv(env); v != "" {
		number, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		return number
	}
	return 0
}

func (c *Config) validate() error {
	if c.SMTP.User == "" {
		return errors.New("SMTP_USER is required")
	}
	if c.SMTP.Password == "" {
		return errors.New("SMTP_PASSWORD is required")
	}
	if c.SMTP.FromAddr == "" {
		return errors.New("FROM_ADDR is required")
	}
	return nil
}
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
