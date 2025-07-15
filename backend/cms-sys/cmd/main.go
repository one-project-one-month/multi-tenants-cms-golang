package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/aws"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	loggMiddleware "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hashicorp/consul/api"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/handler"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/repository"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/routes"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/service"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type dISection struct {
	repo               repository.AuthRepository
	srv                service.AuthService
	handler            handler.AuthHandle
	ownerHandler       handler.OwnerHandle
	pageRequestHandler handler.PageRequestHandle
	pageHandler        handler.PageHandle
	consulClient       *api.Client
}

type consulConfig struct {
	Address    string
	Datacenter string
	Token      string
	Scheme     string
	ServiceID  string
	Name       string
	Tags       []string
	Port       int
	CheckTTL   time.Duration
	CheckHTTP  string
}

func newConsulClient(config consulConfig, logger *logrus.Logger) (*api.Client, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = config.Address
	consulConfig.Datacenter = config.Datacenter
	consulConfig.Token = config.Token
	consulConfig.Scheme = config.Scheme

	client, err := api.NewClient(consulConfig)
	if err != nil {
		logger.WithError(err).Error("Failed to create Consul client")
		return nil, err
	}

	logger.WithField("address", config.Address).Info("Consul client created successfully")
	return client, nil
}

func registerService(client *api.Client, config consulConfig, logger *logrus.Logger) error {
	localIP, err := getLocalIP()
	if err != nil {
		logger.WithError(err).Error("Failed to get local IP address")
		return err
	}

	allTags := config.Tags
	traefikTags := []string{
		"traefik.enable=true",
		"traefik.http.routers.cms-doc-api.rule=Host(`api.localhost`) && PathPrefix(`/api/v1`)",
		"traefik.http.routers.cms-doc-api.service=cms-doc-multi-tenant-api",
		fmt.Sprintf("traefik.http.services.cms-doc-multi-tenant-api.loadbalancer.server.port=%d", config.Port),
	}
	allTags = append(allTags, traefikTags...)

	cmsService := &api.AgentServiceRegistration{
		ID:      config.ServiceID,
		Name:    config.Name,
		Tags:    allTags,
		Port:    config.Port,
		Address: localIP,
		Meta: map[string]string{
			"health-check-path": "/cms-doc/health",
		},
		Checks: api.AgentServiceChecks{
			{
				HTTP:                           fmt.Sprintf("http://%s:%d/cms/health", localIP, config.Port),
				Interval:                       "10s",
				Timeout:                        "5s",
				DeregisterCriticalServiceAfter: "30s",
			},
			{
				CheckID:                        config.ServiceID + ":ttl",
				TTL:                            config.CheckTTL.String(),
				DeregisterCriticalServiceAfter: "1m",
			},
		},
	}

	err = client.Agent().ServiceRegister(cmsService)
	if err != nil {
		logger.WithError(err).Error("Failed to register service with Consul")
		return err
	}

	logger.WithFields(logrus.Fields{
		"service_id":   config.ServiceID,
		"service_name": config.Name,
		"address":      localIP,
		"port":         config.Port,
		"metadata":     "health-check-path=/cms-doc/health",
	}).Info("Service registered with Consul successfully")

	return nil
}

func deregisterService(client *api.Client, serviceID string, logger *logrus.Logger) error {
	err := client.Agent().ServiceDeregister(serviceID)
	if err != nil {
		logger.WithError(err).Error("Failed to deregister service from Consul")
		return err
	}

	logger.WithField("service_id", serviceID).Info("Service deregistered from Consul")
	return nil
}

func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func loadConsulConfig() consulConfig {
	port, _ := strconv.Atoi(utils.GetEnv("PORT", "8081"))
	checkTTL, _ := time.ParseDuration(utils.GetEnv("CONSUL_CHECK_TTL", "30s"))

	tagsStr := utils.GetEnv("CONSUL_SERVICE_TAGS", "cms-doc,multi-tenant,api")
	var tags []string
	if tagsStr != "" {
		for _, tag := range strings.Split(tagsStr, ",") {
			tags = append(tags, strings.TrimSpace(tag))
		}
	}

	localIP, err := getLocalIP()
	if err != nil {
		localIP, _ = os.Hostname()
	}

	return consulConfig{
		Address:    utils.GetEnv("CONSUL_ADDRESS", "consul:8500"),
		Datacenter: utils.GetEnv("CONSUL_DATACENTER", "dc1"),
		Token:      utils.GetEnv("CONSUL_TOKEN", ""),
		Scheme:     utils.GetEnv("CONSUL_SCHEME", "http"),
		ServiceID:  utils.GetEnv("CONSUL_SERVICE_ID", fmt.Sprintf("cms-doc-api-%s", localIP)),
		Name:       utils.GetEnv("CONSUL_SERVICE_NAME", "cms-doc-service"),
		Tags:       tags,
		Port:       port,
		CheckTTL:   checkTTL,
		CheckHTTP:  fmt.Sprintf("http://%s:%d/health", localIP, port),
	}
}
func startHealthUpdateRoutine(client *api.Client, serviceID string, logger *logrus.Logger) {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			err := client.Agent().UpdateTTL(serviceID+":ttl", "Service is healthy", api.HealthPassing)
			if err != nil {
				logger.WithError(err).Error("Failed to update service health check")
			}
		}
	}()

}

func main() {
	appLogger := utils.NewLogger(utils.LogConfig{
		Level:      utils.GetEnv("LOG_LEVEL", "info"),
		FilePath:   utils.GetEnv("LOG_FILE_PATH", "logs/app.log"),
		MaxSize:    utils.GetEnvAsInt("LOG_MAX_SIZE", 100),
		MaxBackups: utils.GetEnvAsInt("LOG_MAX_BACKUPS", 5),
		MaxAge:     utils.GetEnvAsInt("LOG_MAX_AGE", 30),
		Compress:   utils.GetEnvAsBool("LOG_COMPRESS", true),
		Console:    utils.GetEnvAsBool("LOG_CONSOLE", true),
	})

	appLogger.Info("Starting CMS Multi-Tenant System")

	consulConfig := loadConsulConfig()
	consulEnabled := utils.GetEnvAsBool("CONSUL_ENABLED", false)

	var consulClient *api.Client
	if consulEnabled {
		var err error
		consulClient, err = newConsulClient(consulConfig, appLogger)
		if err != nil {
			appLogger.WithError(err).Fatal("Failed to create Consul client")
		}
	}

	dbConfig := utils.DatabaseConfig{
		Host:            utils.GetEnv("DB_HOST", "localhost"),
		Port:            utils.GetEnvAsInt("DB_PORT", 5432),
		User:            utils.GetEnv("DB_USER", "postgres"),
		Password:        utils.GetEnv("DB_PASSWORD", "Swanhtet12@"),
		DBName:          utils.GetEnv("DB_NAME", "cms_db"),
		SSLMode:         utils.GetEnv("DB_SSL_MODE", "disable"),
		MaxOpenConns:    utils.GetEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    utils.GetEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		ConnMaxLifetime: utils.GetEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		ConnMaxIdleTime: utils.GetEnvAsDuration("DB_CONN_MAX_IDLE_TIME", 2*time.Minute),
		RetryAttempts:   utils.GetEnvAsInt("DB_RETRY_ATTEMPTS", 5),
		RetryDelay:      utils.GetEnvAsDuration("DB_RETRY_DELAY", 2*time.Second),
		LogLevel:        logger.Info,
	}

	dbConnection := utils.NewDatabaseConnection(dbConfig, appLogger)
	if err := dbConnection.Connect(); err != nil {
		appLogger.WithError(err).Fatal("Failed to initialize database connection")
	}

	err := dbConnection.DB.AutoMigrate(&types.CMSWholeSysRole{}, &types.CMSUser{}, &types.MFAToken{}, &types.CMSCusPurchase{}, &types.UserPageRequest{}, &types.Page{}, &types.PageRequest{})
	if err != nil {
		appLogger.WithError(err).Fatal("Failed to migrate database")
		return
	}

	if err := utils.InitRedis(); err != nil {
		appLogger.Fatal("Failed to initialize Redis:", err)
	}
	defer func() {
		err := utils.CloseRedis()
		if err != nil {
			appLogger.WithError(err).Fatal("Failed to close Redis connection")
		}
	}()

	healthChecker := utils.NewHealthChecker(dbConnection.DB, appLogger)

	if err := utils.InitNats(); err != nil {
		log.Fatalf("Failed to initialize NATS: %v", err)
	}
	//defer func() {
	//	err := utils.CloseNats()
	//	if err != nil {
	//		appLogger.WithError(err).Fatal("Failed to close NATS")
	//	}
	//}()
	app := fiber.New(fiber.Config{
		AppName: "CMS Multi-Tenant System ",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			appLogger.WithFields(logrus.Fields{
				"method": c.Method(),
				"path":   c.Path(),
				"ip":     c.IP(),
				"error":  err.Error(),
				"code":   code,
			}).Error("Request error")

			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		},
	})

	app.Use(loggMiddleware.New(loggMiddleware.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		Output: appLogger.Writer(),
	}))

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "CMS Multi-Tenant API",
			"status":  "success",
			"version": "1.0.0",
		})
	})

	cmsGroup := app.Group("/cms-doc")

	cmsGroup.Get("/health", func(c *fiber.Ctx) error {
		health := healthChecker.CheckHealth()

		statusCode := fiber.StatusOK
		if health.Status != "healthy" {
			statusCode = fiber.StatusServiceUnavailable
		}

		return c.Status(statusCode).JSON(health)
	})

	cmsGroup.Get("/health/database", func(c *fiber.Ctx) error {
		stats := dbConnection.GetStats()
		return c.JSON(fiber.Map{
			"status": "healthy",
			"stats":  stats,
		})
	})

	cmsGroup.Get("/health/consul", func(c *fiber.Ctx) error {
		if !consulEnabled || consulClient == nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "disabled",
				"message": "Consul is not enabled",
			})
		}

		_, err := consulClient.Status().Leader()
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy",
				"error":  err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"status": "healthy",
			"consul": "connected",
		})
	})

	di := dependencyInjectionSection(appLogger, dbConnection.DB, consulClient, utils.GetRedisClient())
	routes.SetupRoutes(app, di.handler)
	routes.SetupOwnerRoutes(app, di.ownerHandler)
	routes.SetupPageRequestRoutes(app, di.pageRequestHandler)
	routes.SetupPageRoutes(app, di.pageHandler)

	port := utils.GetEnv("PORT", "8080")

	if consulEnabled && consulClient != nil {
		if err := registerService(consulClient, consulConfig, appLogger); err != nil {
			appLogger.WithError(err).Error("Failed to register service with Consul")
		} else {
			startHealthUpdateRoutine(consulClient, consulConfig.ServiceID, appLogger)
		}
	}

	go func() {
		appLogger.WithField("port", port).Info("Server starting")
		if err := app.Listen("0.0.0.0:" + port); err != nil {
			appLogger.WithError(err).Fatal("Server failed to start")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")
	if consulEnabled && consulClient != nil {
		if err := deregisterService(consulClient, consulConfig.ServiceID, appLogger); err != nil {
			appLogger.WithError(err).Error("Failed to deregister service from Consul")
		}
	}

	if err := app.Shutdown(); err != nil {
		appLogger.WithError(err).Error("Server forced to shutdown")
	}

	if err := dbConnection.Close(); err != nil {
		appLogger.WithError(err).Error("Failed to close database connection")
	}

	appLogger.Info("Server exited")
}

func dependencyInjectionSection(
	logger *logrus.Logger,
	db *gorm.DB,
	consulClient *api.Client,
	redisClient *redis.Client,
) *dISection {
	repo := repository.NewRepo(logger, db)
	if err := repo.CreateDefaultRoles(); err != nil {
		logger.Fatalf("Failed to create default roles: %v", err)
	}
	srv := service.NewService(logger, repo, redisClient)
	authHandler := handler.NewHandler(srv)
	bucketName := utils.GetEnv("BUCKET_NAME", "")
	s3, err := aws.NewS3Service(bucketName, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create S3 service")
	}
	ownerRepo := repository.NewOwnerRepository(logger, db)
	ownerService := service.NewOwnerService(logger, ownerRepo, repo)
	ownerHandler := handler.NewOwnerHandler(ownerService)

	pageRequestRepo := repository.NewPageRequestRepository(logger, db)
	pageRequestSrv := service.NewPageRequestService(logger, pageRequestRepo)
	pageRequestHandler := handler.NewPageRequestHandler(pageRequestSrv, s3)

	// Page
	pageRepo := repository.NewPageRepository(logger, db)
	pageService := service.NewPageService(logger, pageRepo, pageRequestRepo, ownerRepo)
	pageHandler := handler.NewPageHandler(pageService)

	return &dISection{

		repo:               repo,
		srv:                srv,
		handler:            authHandler,
		ownerHandler:       ownerHandler,
		pageRequestHandler: pageRequestHandler,
		pageHandler:        pageHandler,
		consulClient:       consulClient,
	}
}
