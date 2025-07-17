package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/multi-tenants-cms-golang/lms-sys/app"
	"github.com/multi-tenants-cms-golang/lms-sys/app/gateway"
	"github.com/multi-tenants-cms-golang/lms-sys/app/rpc"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils/env"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

func initLogger() *logrus.Logger {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.WithError(err).Fatal("Failed to load .env")
	}

	// Initialize logger
	logger := logrus.New()
	logger.Info("Initializing logger with rotation and JSON formatting")

	// Configure Lumberjack for log rotation
	logFile := &lumberjack.Logger{
		Filename:   "logs/lms-system.log",
		MaxSize:    100, // MB
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
		LocalTime:  true,
	}
	logger.Infof("Configured log rotation: file=%s, maxSize=%dMB, maxBackups=%d, maxAge=%ddays, compress=%t",
		logFile.Filename, logFile.MaxSize, logFile.MaxBackups, logFile.MaxAge, logFile.Compress)

	// Set output to both file and stdout
	logger.Info("Setting dual output to stdout and log file")
	logger.SetOutput(io.MultiWriter(os.Stdout, logFile))

	// Configure JSON formatter with enhanced caller info
	logger.Info("Configuring JSON formatter with custom fields")
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999Z07:00",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "caller",
		},
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			file := strings.TrimPrefix(f.File, "/Users/swanhtet/Desktop/multi-tenants-cms-golang/backend/lms-System/")
			funcName := f.Function[strings.LastIndex(f.Function, ".")+1:]
			return file, funcName + ":" + string(rune(f.Line))
		},
		PrettyPrint: true,
	})

	// Enable caller reporting
	logger.Info("Enabling caller reporting in logs")
	logger.SetReportCaller(true)

	// Set log level from environment
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	logger.Infof("Setting log level: %s", logLevel)

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.Warnf("Invalid LOG_LEVEL '%s', defaulting to 'info'", logLevel)
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Verify logger configuration
	logger.WithFields(logrus.Fields{
		"logLevel":    logger.GetLevel(),
		"outputs":     "stdout+file",
		"formatter":   "JSON",
		"callerInfo":  true,
		"compression": logFile.Compress,
	}).Info("Logger initialization complete")

	return logger
}
func main() {
	logger := initLogger()
	defer func() {
		if file, ok := logger.Out.(*lumberjack.Logger); ok {
			err := file.Close()
			if err != nil {
				return
			}
		}
	}()

	logger.Info("Application starting with log rotation enabled")

	grpcServerAddress := env.GetEnv("LMS_GRPC_SERVER_ADDRESS", ":9001")
	grpcGatewayAddress := env.GetEnv("LMS_GRPC_GATEWAY_ADDRESS", ":8086")
	jwtSecret := env.GetEnv("LMS_JWT_SECRET", "default-secret-key-at-least-32-characters-long")
	jwtIssuer := env.GetEnv("LMS_JWT_ISSUER", "lms-system")
	jwtAudience := env.GetEnv("LMS_JWT_AUDIENCE", "lms-client")
	consulAddress := env.GetEnv("LMS_CONSUL_ADDRESS", "localhost:8500")
	//serviceName := env.GetEnv("LMS_SERVICE_NAME", "lms-service")
	serviceID := env.GetEnv("LMS_SERVICE_ID", "lms-srvice-1")

	dbConn := DatabaseConn(logger)
	defer dbConn.Close()

	dbStore := db.NewStore(dbConn)

	consulClient, err := NewConsulClient(consulAddress)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create Consul client")
	}
	defer func(client *api.Client, serviceID string) {
		err := DeregisterService(client, serviceID)
		if err != nil {
			logger.WithError(err).Fatal("Failed to deregister service")
		}
	}(consulClient, serviceID)

	//if err := RegisterService(consulClient, serviceName, serviceID+"-grpc", grpcServerAddress, "grpc"); err != nil {
	//	logger.WithError(err).Fatal("Failed to register gRPC server with Consul")
	//}
	//
	//if err := RegisterService(consulClient, serviceName, serviceID+"-gateway", grpcGatewayAddress, "http"); err != nil {
	//	logger.WithError(err).Fatal("Failed to register gRPC gateway with Consul")
	//}

	grpcServer := rpc.NewServer(
		dbStore,
		logger,
		jwtSecret,
		jwtIssuer,
		jwtAudience,
	)

	grpcGateway := gateway.NewGateway(
		logger,
		grpcServerAddress,
		grpcGatewayAddress,
		jwtSecret,
		jwtIssuer,
		jwtAudience,
	)

	server := app.NewApp(
		grpcServer,
		grpcGateway,
		logger,
	)
	if err := server.Run(); err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}
}

func NewConsulClient(address string) (*api.Client, error) {
	config := api.DefaultConfig()
	config.Address = address
	return api.NewClient(config)
}

func RegisterService(client *api.Client, serviceName, serviceID, address, serviceType string) error {
	registration := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: address,
		Port:    getPortFromAddress(address),
		Tags:    []string{serviceType},
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s/healthz", address),
			Interval: "10s",
			Timeout:  "5s",
		},
	}

	if serviceType == "grpc" {
		registration.Check.GRPC = fmt.Sprintf("%s/grpc.health.v1.Health/Check", address)
		registration.Check.Interval = "15s"
	}

	return client.Agent().ServiceRegister(registration)
}

func DeregisterService(client *api.Client, serviceID string) error {
	return client.Agent().ServiceDeregister(serviceID)
}

func getPortFromAddress(address string) int {
	var port int
	_, err := fmt.Sscanf(address, ":%d", &port)
	if err != nil {

		return 0
	}
	return port
}

func DatabaseConn(logger *logrus.Logger) *pgxpool.Pool {
	user := env.GetEnv("LMS_DB_USER", "lms_user")
	password := env.GetEnv("LMS_DB_PASSWORD", "lms_password")
	dbName := env.GetEnv("LMS_DB_NAME", "lms_db")
	host := env.GetEnv("LMS_DB_HOST", "localhost")
	port := env.GetEnv("LMS_DB_PORT", "5433")
	maxConns := env.GetEnvAsInt("LMS_DB_MAX_CONNS", 10)
	minConns := env.GetEnvAsInt("LMS_DB_MIN_CONNS", 2)
	maxConnLifetime := env.GetEnvAsDuration("LMS_DB_MAX_CONN_LIFETIME", time.Hour)
	maxConnIdleTime := env.GetEnvAsDuration("LMS_DB_MAX_CONN_IDLE_TIME", 30*time.Minute)

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)

	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		logger.WithError(err).Fatal("Failed to parse database configuration")
	}

	config.MaxConns = int32(maxConns)
	config.MinConns = int32(minConns)
	config.MaxConnLifetime = maxConnLifetime
	config.MaxConnIdleTime = maxConnIdleTime

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
