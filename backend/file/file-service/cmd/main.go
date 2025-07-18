package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/multi-tenants-cms-golang/file-service/cmd/app"
	"github.com/multi-tenants-cms-golang/file-service/cmd/provider"
	"github.com/multi-tenants-cms-golang/file-service/config"
	"github.com/multi-tenants-cms-golang/file-service/internal/handler"
	"github.com/multi-tenants-cms-golang/file-service/internal/service"
	utils "github.com/multi-tenants-cms-golang/file-service/pkg"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
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
	token := utils.GetEnv("TOKEN", "")

	logger := utils.NewLogger(utils.GetDefaultLogConfig())
	ghClient := provider.NewGitHubClient(token)
	cfg := config.NewConfiguration()
	fiberApp := provider.NewFiberApp()
	fileService := service.NewFileSystem(logger, ghClient)
	fileHandler := handler.NewFileSystemHandler(logger, fileService)
	application := app.NewServerState(
		cfg,
		fileHandler,
		fiberApp,
		logger,
	)

	consulCfg := loadConsulConfig()

	consulClient, err := newConsulClient(consulCfg, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create Consul client")
	}

	if err := registerService(consulClient, consulCfg, logger); err != nil {
		logger.WithError(err).Fatal("Failed to register service with Consul")
	}

	startHealthUpdateRoutine(consulClient, consulCfg.ServiceID, logger)

	go application.StartServer()

	osChannel := make(chan os.Signal, 1)
	signal.Notify(osChannel, syscall.SIGINT, syscall.SIGTERM)
	<-osChannel

	logger.Info("Shutting down...")

	if err := deregisterService(consulClient, consulCfg.ServiceID, logger); err != nil {
		logger.WithError(err).Error("Failed to deregister service from Consul")
	}

	application.StopServer()
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
		ID:   config.ServiceID,
		Name: config.Name,
		Tags: allTags,
		Port: config.Port,
		//Address: localIP,
		Checks: api.AgentServiceChecks{
			{
				HTTP:                           config.CheckHTTP,
				Interval:                       "10s",
				Timeout:                        "5s",
				DeregisterCriticalServiceAfter: "30s",
			},
			{
				CheckID:                        config.ServiceID + ":ttl", // Ensure this matches the update routine
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
	port, _ := strconv.Atoi(utils.GetEnv("PORT", "8080"))
	checkTTL, _ := time.ParseDuration(utils.GetEnv("CONSUL_CHECK_TTL", "30s"))

	tagsStr := utils.GetEnv("CONSUL_SERVICE_TAGS", "cms-doc,multi-tenant,api")
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
		ServiceID:  utils.GetEnv("CONSUL_SERVICE_ID", fmt.Sprintf("cms-doc-api-%s", hostname)),
		Name:       utils.GetEnv("CONSUL_SERVICE_NAME", "cms-doc-multi-tenant-api"),
		Tags:       tags,
		Port:       port,
		CheckTTL:   checkTTL,
		CheckHTTP:  fmt.Sprintf("http://%s:%d/health", hostname, port),
	}
}

func startHealthUpdateRoutine(client *api.Client, serviceID string, logger *logrus.Logger) chan struct{} {
	ticker := time.NewTicker(10 * time.Second)
	stop := make(chan struct{})

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := client.Agent().UpdateTTL(serviceID+":ttl", "Service is healthy", api.HealthPassing)
				if err != nil {
					logger.WithError(err).Error("Failed to update service health check")
				}
			case <-stop:
				return
			}
		}
	}()

	return stop
}
