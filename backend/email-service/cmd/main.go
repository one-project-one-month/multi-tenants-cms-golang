package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/multi-tenants-cms-golang/email-service/utils"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/multi-tenants-cms-golang/email-service/config"
	"github.com/multi-tenants-cms-golang/email-service/internal/email"
	"github.com/multi-tenants-cms-golang/email-service/internal/health"
	natss "github.com/multi-tenants-cms-golang/email-service/internal/nats"
)

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

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("failed to load configuration", zap.Error(err))
	}

	nc, err := nats.Connect(cfg.NATS.URL,
		nats.ReconnectWait(time.Second),
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.Error("NATS disconnected", zap.Error(err))
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Info("NATS reconnected")
		}),
	)
	if err != nil {
		logger.Fatal("failed to connect to NATS", zap.Error(err))
	}
	defer nc.Close()

	emailCh := make(chan natss.EmailRequest, 100)
	emailService, err := email.NewService(
		logger,
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.User,
		cfg.SMTP.Password,
		cfg.SMTP.FromAddr,
		cfg.SMTP.TemplateDir,
	)
	if err != nil {
		logger.Fatal("failed to create email service", zap.Error(err))
	}

	natsConsumer, err := natss.NewConsumer(
		nc,
		logger,
		emailCh,
		cfg.NATS.StreamName,
		cfg.NATS.Subject,
		cfg.NATS.ConsumerName,
	)
	if err != nil {
		logger.Fatal("failed to create NATS consumer", zap.Error(err))
	}

	consulCfg := loadConsulConfig()
	consulClient, err := api.NewClient(&api.Config{
		Address:    consulCfg.Address,
		Datacenter: consulCfg.Datacenter,
		Token:      consulCfg.Token,
		Scheme:     consulCfg.Scheme,
	})
	if err != nil {
		logger.Fatal("failed to create Consul client", zap.Error(err))
	}

	if err := registerService(consulClient, consulCfg, logger); err != nil {
		logger.Fatal("failed to register service with Consul", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	workerCount := runtime.NumCPU()
	logger.Info("starting email workers", zap.Int("count", workerCount))

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			logger.Info("worker started", zap.Int("id", workerID))

			for {
				select {
				case req, ok := <-emailCh:
					if !ok {
						logger.Info("worker stopping", zap.Int("id", workerID))
						return
					}

					if err := emailService.Send(ctx, req); err != nil {
						logger.Error("failed to send email",
							zap.Int("worker", workerID),
							zap.String("to", req.To),
							zap.String("template", req.Template),
							zap.Error(err))
					}
				case <-ctx.Done():
					logger.Info("worker stopping due to context cancellation", zap.Int("id", workerID))
					return
				}
			}
		}(i)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := natsConsumer.Start(ctx); err != nil {
			logger.Error("NATS consumer failed", zap.Error(err))
		}
	}()

	healthServer := health.NewServer(cfg.Server.Port, logger)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := healthServer.Start(); !errors.Is(err, http.ErrServerClosed) && err != nil {
			logger.Error("health server failed", zap.Error(err))
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info("shutdown signal received, starting graceful shutdown")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	if err := healthServer.Stop(shutdownCtx); err != nil {
		logger.Error("failed to stop health server gracefully", zap.Error(err))
	}

	cancel()
	close(emailCh)

	defer func() {
		if err := deregisterService(consulClient, consulCfg.ServiceID, logger); err != nil {
			logger.Error("failed to deregister service from Consul", zap.Error(err))
		}
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Info("all goroutines stopped")
	case <-shutdownCtx.Done():
		logger.Warn("shutdown timeout exceeded, forcing exit")
	}

	logger.Info("shutdown complete")
}

func loadConsulConfig() consulConfig {
	port, _ := strconv.Atoi(utils.GetEnv("PORT", "8080"))
	checkTTL, _ := time.ParseDuration(utils.GetEnv("CONSUL_CHECK_TTL", "30s"))

	tagsStr := utils.GetEnv("CONSUL_SERVICE_TAGS", "cms,multi-tenant,api")
	var tags []string
	if tagsStr != "" {
		for _, tag := range strings.Split(tagsStr, ",") {
			tags = append(tags, strings.TrimSpace(tag))
		}
	}

	hostname, _ := os.Hostname()

	return consulConfig{
		Address:    utils.GetEnv("CONSUL_ADDRESS", "consul:8500"),
		Datacenter: utils.GetEnv("CONSUL_DATACENTER", "dc1"),
		Token:      utils.GetEnv("CONSUL_TOKEN", ""),
		Scheme:     utils.GetEnv("CONSUL_SCHEME", "http"),
		ServiceID:  utils.GetEnv("CONSUL_SERVICE_ID", fmt.Sprintf("cms-api-%s", hostname)),
		Name:       utils.GetEnv("CONSUL_SERVICE_NAME", "cms-multi-tenant-api"),
		Tags:       tags,
		Port:       port,
		CheckTTL:   checkTTL,
		CheckHTTP:  fmt.Sprintf("http://%s:%d/health", hostname, port),
	}
}

func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("failed to get local IP: %w", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func registerService(client *api.Client, config consulConfig, logger *zap.Logger) error {
	localIP, err := getLocalIP()
	if err != nil {
		logger.Error("failed to get local IP address", zap.Error(err))
		return err
	}

	service := &api.AgentServiceRegistration{
		ID:      config.ServiceID,
		Name:    config.Name,
		Tags:    config.Tags,
		Port:    config.Port,
		Address: localIP,
		Check: &api.AgentServiceCheck{
			HTTP:                           config.CheckHTTP,
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "30s",
		},
		Meta: map[string]string{
			"version":     "1.0.0",
			"environment": utils.GetEnv("ENV", "development"),
			"region":      utils.GetEnv("REGION", "us-east-1"),
		},
	}

	if err := client.Agent().ServiceRegister(service); err != nil {
		logger.Error("failed to register service with Consul", zap.Error(err))
		return err
	}

	logger.Info("service registered with Consul",
		zap.String("service_id", config.ServiceID),
		zap.String("service_name", config.Name),
		zap.String("address", localIP),
		zap.Int("port", config.Port))

	return nil
}

func deregisterService(client *api.Client, serviceID string, logger *zap.Logger) error {
	if err := client.Agent().ServiceDeregister(serviceID); err != nil {
		logger.Error("failed to deregister service from Consul", zap.Error(err))
		return err
	}

	logger.Info("service deregistered from Consul", zap.String("service_id", serviceID))
	return nil
}
