package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils/env"
	"github.com/sirupsen/logrus"
	"time"
)

func DatabaseConn(logger *logrus.Logger) *pgxpool.Pool {
	//user := env.GetEnv("LMS_DB_USER", "")
	//password := env.GetEnv("LMS_DB_PASSWORD", "")
	//dbName := env.GetEnv("LMS_DB_NAME", "")
	//host := env.GetEnv("LMS_DB_HOST", "")
	//port := env.GetEnv("LMS_DB_PORT", "")
	maxConns := env.GetEnvAsInt("LMS_DB_MAX_CONNS", 10)
	minConns := env.GetEnvAsInt("LMS_DB_MIN_CONNS", 2)
	maxConnLifetime := env.GetEnvAsDuration("LMS_DB_MAX_CONN_LIFETIME", time.Hour)
	maxConnIdleTime := env.GetEnvAsDuration("LMS_DB_MAX_CONN_IDLE_TIME", 30*time.Minute)

	/*dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)*/
	dbUrl := "postgresql://neondb_owner:npg_MEB4CYJS7TKh@ep-withered-salad-a27oo7gn-pooler.eu-central-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"
	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		logger.WithError(err).Fatal("Failed to parse database configuration")
	}

	config.MaxConns = int32(maxConns)
	config.MinConns = int32(minConns)
	config.MaxConnLifetime = maxConnLifetime
	config.MaxConnIdleTime = maxConnIdleTime
	config.HealthCheckPeriod = 1 * time.Minute

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheStatement

	config.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		err := conn.Ping(ctx)
		if err != nil {
			logger.WithError(err).Warn("Releasing unhealthy connection")
		}
		return err == nil
	}

	var connPool *pgxpool.Pool
	maxRetries := 5
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		connPool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err == nil {
			break
		}

		logger.WithError(err).Warnf("Failed to connect to database (attempt %d/%d)", i+1, maxRetries)
		if i < maxRetries-1 {
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database after retries")
	}

	if err := connPool.Ping(context.Background()); err != nil {
		logger.WithError(err).Fatal("Failed to ping database")
	}

	logger.Info("Successfully connected to database")
	return connPool

}
