package main

import (
	"errors"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	loggMiddleware "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/handler"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/repository"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/routes"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/service"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type DISection struct {
	repo    repository.AuthRepository
	srv     service.AuthService
	handler handler.AuthHandle
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

	err := dbConnection.DB.AutoMigrate(&types.CMSWholeSysRole{}, &types.CMSUser{}, &types.CMSCusPurchase{}, &types.UserPageRequest{})
	if err != nil {
		appLogger.WithError(err).Fatal("Failed to migrate database")
		return
	}
	//if err := utils.InitJWTKeysFromVault(); err != nil {
	//	log.Fatalf("Vault key init failed: %v", err)
	//}
	//
	//uid := uuid.New()
	//token, err := utils.GenerateAccessToken(uid, "user@example.com", "admin")
	//if err != nil {
	//	log.Fatalf("token generation failed: %v", err)
	//}

	//fmt.Println("JWT:", token)
	healthChecker := utils.NewHealthChecker(dbConnection.DB, appLogger)

	app := fiber.New(fiber.Config{
		AppName: "CMS Multi-Tenant System",
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

	app.Get("/health", func(c *fiber.Ctx) error {
		health := healthChecker.CheckHealth()

		statusCode := fiber.StatusOK
		if health.Status != "healthy" {
			statusCode = fiber.StatusServiceUnavailable
		}

		return c.Status(statusCode).JSON(health)
	})

	app.Get("/health/database", func(c *fiber.Ctx) error {
		stats := dbConnection.GetStats()
		return c.JSON(fiber.Map{
			"status": "healthy",
			"stats":  stats,
		})
	})

	di := DependencyInjectionSection(appLogger, dbConnection.DB)
	routes.SetupRoutes(app, di.handler)

	port := utils.GetEnv("PORT", "8080")

	go func() {
		appLogger.WithField("port", port).Info("Server starting")
		if err := app.Listen(":" + port); err != nil {
			appLogger.WithError(err).Fatal("Server failed to start")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		appLogger.WithError(err).Error("Server forced to shutdown")
	}

	if err := dbConnection.Close(); err != nil {
		appLogger.WithError(err).Error("Failed to close database connection")
	}

	appLogger.Info("Server exited")
}

func DependencyInjectionSection(logger *logrus.Logger, db *gorm.DB) *DISection {
	repo := repository.NewRepo(logger, db)
	srv := service.NewService(logger, repo)
	handler := handler.NewHandler(srv)

	return &DISection{
		repo:    repo,
		srv:     srv,
		handler: handler,
	}
}
