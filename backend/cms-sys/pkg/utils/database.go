package utils

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	RetryAttempts   int
	RetryDelay      time.Duration
	LogLevel        logger.LogLevel
}

type DatabaseConnection struct {
	DB     *gorm.DB
	Config DatabaseConfig
	Logger *logrus.Logger
}

func NewDatabaseConnection(config DatabaseConfig, log *logrus.Logger) *DatabaseConnection {
	return &DatabaseConnection{
		Config: config,
		Logger: log,
	}
}

func (dc *DatabaseConnection) Connect() error {
	var db *gorm.DB
	var err error

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dc.Config.Host,
		dc.Config.Port,
		dc.Config.User,
		dc.Config.Password,
		dc.Config.DBName,
		dc.Config.SSLMode,
	)

	gormConfig := &gorm.Config{
		Logger: logger.New(
			dc.Logger,
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  dc.Config.LogLevel,
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	for attempt := 1; attempt <= dc.Config.RetryAttempts; attempt++ {
		dc.Logger.WithFields(logrus.Fields{
			"attempt":      attempt,
			"max_attempts": dc.Config.RetryAttempts,
			"host":         dc.Config.Host,
			"port":         dc.Config.Port,
			"database":     dc.Config.DBName,
		}).Info("Attempting to connect to database")

		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err == nil {
			dc.Logger.Info("Successfully connected to database")
			break
		}

		dc.Logger.WithFields(logrus.Fields{
			"attempt":     attempt,
			"error":       err.Error(),
			"retry_delay": dc.Config.RetryDelay,
		}).Error("Failed to connect to database")

		if attempt < dc.Config.RetryAttempts {
			dc.Logger.WithField("delay", dc.Config.RetryDelay).Info("Retrying database connection")
			time.Sleep(dc.Config.RetryDelay)
		}
	}

	if err != nil {
		dc.Logger.WithError(err).Fatal("Failed to connect to database after all retry attempts")
		return fmt.Errorf("failed to connect to database after %d attempts: %w", dc.Config.RetryAttempts, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		dc.Logger.WithError(err).Error("Failed to get underlying sql.DB")
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(dc.Config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dc.Config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(dc.Config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(dc.Config.ConnMaxIdleTime)

	dc.DB = db

	if err := dc.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	dc.Logger.WithFields(logrus.Fields{
		"max_open_conns":     dc.Config.MaxOpenConns,
		"max_idle_conns":     dc.Config.MaxIdleConns,
		"conn_max_lifetime":  dc.Config.ConnMaxLifetime,
		"conn_max_idle_time": dc.Config.ConnMaxIdleTime,
	}).Info("Database connection pool configured successfully")

	return nil
}

func (dc *DatabaseConnection) Ping() error {
	sqlDB, err := dc.DB.DB()
	if err != nil {
		return err
	}

	ctx, cancel := GetContextWithTimeout(5 * time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}

func (dc *DatabaseConnection) Close() error {
	if dc.DB != nil {
		sqlDB, err := dc.DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func (dc *DatabaseConnection) GetStats() map[string]interface{} {
	if dc.DB == nil {
		return nil
	}

	sqlDB, err := dc.DB.DB()
	if err != nil {
		dc.Logger.WithError(err).Error("Failed to get database stats")
		return nil
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration,
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}

func GetDefaultDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:            "localhost",
		Port:            5432,
		User:            "postgres",
		Password:        "password",
		DBName:          "cms_db",
		SSLMode:         "disable",
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 2 * time.Minute,
		RetryAttempts:   5,
		RetryDelay:      2 * time.Second,
		LogLevel:        logger.Info,
	}
}
