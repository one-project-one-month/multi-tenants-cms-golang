package utils

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type HealthChecker struct {
	db     *gorm.DB
	logger *logrus.Logger
}

type HealthStatus struct {
	Status      string                   `json:"status"`
	Timestamp   time.Time                `json:"timestamp"`
	Version     string                   `json:"version"`
	Environment string                   `json:"environment"`
	Services    map[string]ServiceHealth `json:"services"`
}

type ServiceHealth struct {
	Status       string                 `json:"status"`
	Message      string                 `json:"message,omitempty"`
	ResponseTime time.Duration          `json:"response_time"`
	Details      map[string]interface{} `json:"details,omitempty"`
}

func NewHealthChecker(db *gorm.DB, logger *logrus.Logger) *HealthChecker {
	return &HealthChecker{
		db:     db,
		logger: logger,
	}
}

func (h *HealthChecker) CheckHealth() HealthStatus {
	status := HealthStatus{
		Status:      "healthy",
		Timestamp:   time.Now().UTC(),
		Version:     "1.0.0",
		Environment: GetEnv("APP_ENV", "development"),
		Services:    make(map[string]ServiceHealth),
	}

	dbHealth := h.checkDatabase()
	status.Services["database"] = dbHealth

	if dbHealth.Status != "healthy" {
		status.Status = "unhealthy"
	}

	return status
}

func (h *HealthChecker) checkDatabase() ServiceHealth {
	start := time.Now()

	ctx, cancel := GetContextWithTimeout(5 * time.Second)
	defer cancel()

	var result int
	err := h.db.WithContext(ctx).Raw("SELECT 1").Scan(&result).Error

	responseTime := time.Since(start)

	if err != nil {
		h.logger.WithError(err).Error("Database health check failed")
		return ServiceHealth{
			Status:       "unhealthy",
			Message:      err.Error(),
			ResponseTime: responseTime,
		}
	}

	sqlDB, err := h.db.DB()
	if err != nil {
		return ServiceHealth{
			Status:       "unhealthy",
			Message:      "Failed to get underlying database connection",
			ResponseTime: responseTime,
		}
	}

	stats := sqlDB.Stats()
	details := map[string]interface{}{
		"open_connections": stats.OpenConnections,
		"in_use":           stats.InUse,
		"idle":             stats.Idle,
		"max_open":         stats.MaxOpenConnections,
	}

	return ServiceHealth{
		Status:       "healthy",
		ResponseTime: responseTime,
		Details:      details,
	}
}
