package main

import (
	"github.com/hibiken/asynq"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/joho/godotenv"
	"github.com/multi-tenants-cms-golang/lms-sys/app/middleware"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/infra/consul"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/infra/ipfs"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/multi-tenants-cms-golang/lms-sys/app"
	"github.com/multi-tenants-cms-golang/lms-sys/app/cornServer"
	"github.com/multi-tenants-cms-golang/lms-sys/app/cornServer/handler"
	"github.com/multi-tenants-cms-golang/lms-sys/app/gateway"
	"github.com/multi-tenants-cms-golang/lms-sys/app/rpc"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/infra/nats"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/infra/postgres"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/infra/redis"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils/env"
)

func main() {
	logger := initLogger()
	defer closeLogger(logger)

	loadEnv(logger)

	dbPool := initDatabase(logger)
	defer dbPool.Close()

	initRedis(logger)
	initNATS(logger)

	serviceID := env.GetEnv("LMS_SERVICE_ID", "lms-srvice-1")
	consulClient := initConsul(logger, serviceID)

	app := initApp(logger, dbPool, consulClient)

	if err := app.Run(); err != nil {
		logger.WithError(err).Fatal("Failed to start server")
	}
}

func loadEnv(logger *logrus.Logger) {
	if err := godotenv.Load(); err != nil {
		logger.WithError(err).Fatal("Failed to load .env file")
	}
	logger.Info("Environment variables loaded")
}

func initLogger() *logrus.Logger {
	logger := logrus.New()

	logFile := &lumberjack.Logger{
		Filename:   "logs/lms-system.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
		LocalTime:  true,
	}
	logger.SetOutput(io.MultiWriter(os.Stdout, logFile))

	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999Z07:00",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "caller",
		},
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			funcName := f.Function[strings.LastIndex(f.Function, ".")+1:]
			return f.File, funcName + ":" + string(rune(f.Line))
		},
		ForceColors: true,
	})
	logger.SetReportCaller(true)
	return logger
}

func closeLogger(logger *logrus.Logger) {
	if file, ok := logger.Out.(*lumberjack.Logger); ok {
		_ = file.Close()
	}
}

func initDatabase(logger *logrus.Logger) *pgxpool.Pool {
	connectionPool := postgres.DatabaseConn(logger)
	logger.Info("Database connection initialized")
	return connectionPool
}

func initRedis(logger *logrus.Logger) {
	redisAddr := env.GetEnv("LMS_REDIS_ADDRESS", "localhost:6379")
	redisPass := env.GetEnv("LMS_REDIS_PASSWORD", "")
	dbNum := env.GetEnvAsInt("LMS_DB_REDIS_ADDRESS", 0)

	redis.InitRedis(redisAddr, redisPass, dbNum)
	logger.Info("Redis initialized")
}

func initNATS(logger *logrus.Logger) {
	natsURL := "nats://localhost:4222"
	logger.WithField("url", natsURL).Info("Attempting to connect to NATS")

	if err := nats.InitNATS(natsURL); err != nil {
		logger.WithError(err).WithField("url", natsURL).Fatal("Failed to connect to NATS")
	}
	logger.Info("NATS initialized successfully")
}
func initConsul(logger *logrus.Logger, serviceID string) *api.Client {
	consulAddr := env.GetEnv("LMS_CONSUL_ADDRESS", "localhost:8500")
	client, err := consul.NewConsulClient(consulAddr)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create Consul client")
	}
	logger.Info("Consul client initialized")

	serviceName := env.GetEnv("LMS_SERVICE_NAME", "lms-service")
	serviceTags := env.GetEnv("LMS_SERVICE_TAGS", "lms,learning,api,v1,environment:development")
	serviceTagsSlice := strings.Split(serviceTags, ",")
	serviceAddress := env.GetEnv("LMS_SERVICE_ADDRESS", "localhost:8086")
	serviceType := "http"

	err = consul.RegisterService(client, serviceName, serviceID, serviceAddress, serviceType, serviceTagsSlice)
	if err != nil {
		logger.WithError(err).Fatal("Failed to register service with Consul")
	}
	logger.Infof("Service %s registered with Consul", serviceName)

	go func() {
		<-time.After(time.Second * 1)
		if err := consul.DeregisterService(client, serviceID); err != nil {
			logger.WithError(err).Error("Failed to deregister service")
		} else {
			logger.Infof("Service %s deregistered from Consul", serviceName)
		}
	}()

	return client
}
func initApp(logger *logrus.Logger, dbPool *pgxpool.Pool, consulClient *api.Client) *app.App {
	grpcAddress := env.GetEnv("LMS_GRPC_SERVER_ADDRESS", ":9001")
	jwtSecret := env.GetEnv("LMS_JWT_SECRET", "default-secret-key-at-least-32-characters-long")
	jwtIssuer := env.GetEnv("LMS_JWT_ISSUER", "lms-system")
	jwtAudience := env.GetEnv("LMS_JWT_AUDIENCE", "lms-client")

	dbUrl := "postgresql://neondb_owner:npg_MEB4CYJS7TKh@ep-withered-salad-a27oo7gn-pooler.eu-central-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"

	backupHandler := handler.NewCornHandler(logger)
	ipfsShell := shell.NewShell("localhost:5001")
	ipfsHanlder := handler.NewIPFSHandler(logger, ipfsShell)
	cronServer := cornServer.NewServer(
		logger,
		env.GetEnv("LMS_REDIS_ADDRESS", ""),
		backupHandler,
		ipfsHanlder,
		dbUrl,
	)
	ipfsCfg := ipfs.Config{
		Host:    env.GetEnv("IPFS_HOST", "localhost"),
		Port:    env.GetEnv("IPFS_PORT", "5001"),
		Timeout: 10 * time.Minute,
	}
	ipfsClient := ipfs.NewClient(&ipfsCfg, logger)

	redisOpt := asynq.RedisClientOpt{
		Addr:     env.GetEnv("REDIS_ADDR", "localhost:6379"),
		Password: env.GetEnv("REDIS_PASSWORD", ""),
		DB:       env.GetEnvAsInt("REDIS_DB", 0),
	}
	asynqClient := asynq.NewClient(redisOpt)
	defer func(asynqClient *asynq.Client) {
		err := asynqClient.Close()
		if err != nil {
			logger.WithError(err).Error("Failed to close asynq client connection")
		}
	}(asynqClient)
	//taskCreator := &task.TaskCreator{}

	//fileService := files.NewFileService(ipfsClient, asynqClient, logger, taskCreator)
	//startFileUploadWorker(logger, ipfsClient, fileService)
	store := db.NewStore(logger, dbPool)
	mfaConfig := utils.DefaultMFAConfig()
	mfaManger := utils.NewMFAManager(mfaConfig)
	authMiddlware := middleware.NewAuthMiddleware(jwtSecret)
	grpcSrv := rpc.NewServer(
		store,
		logger,
		jwtSecret,
		jwtIssuer,
		jwtAudience,
		grpcAddress,
		mfaManger,
		authMiddlware,
		ipfsClient,
		asynqClient,
	)

	grpcGateway := gateway.NewGateway(logger, grpcAddress, ":8086", "localhost:6379")

	return app.NewApp(grpcSrv, grpcGateway, cronServer, logger)
}

//func startFileUploadWorker(logger *logrus.Logger, ipfsClient *ipfs.Client, fileService *files.FileService) {
//	redisOpt := asynq.RedisClientOpt{
//		Addr:     env.GetEnv("REDIS_ADDR", "localhost:6379"),
//		Password: env.GetEnv("REDIS_PASSWORD", ""),
//		DB:       env.GetEnvAsInt("REDIS_DB", 0),
//	}
//
//	worker := asynq.NewServer(
//		redisOpt,
//		asynq.Config{
//			Concurrency: 10,
//			Queues: map[string]int{
//				"file-upload": 5,
//				"default":     3,
//			},
//		},
//	)
//
//	//processor := task.NewFileUploadProcessor(fileService, ipfsClient, logger)
//	mux := asynq.NewServeMux()
//	mux.Handle(task.TypeFileUpload, processor)
//
//	go func() {
//		if err := worker.Run(mux); err != nil {
//			logger.WithError(err).Fatal("failed to start file upload worker")
//		}
//	}()
//}
