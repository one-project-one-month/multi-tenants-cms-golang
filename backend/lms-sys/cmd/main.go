package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/multi-tenants-cms-golang/lms-sys/app"
	"github.com/multi-tenants-cms-golang/lms-sys/app/api"
	"github.com/multi-tenants-cms-golang/lms-sys/app/rpc"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"github.com/multi-tenants-cms-golang/lms-sys/pkg/utils"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	grpcServerAddress := utils.GetEnv("GRPC_SERVER_ADDRESS", ":9001")
	gprcGatewayAddress := utils.GetEnv("GRPC_GATEWAY_ADDRESS", ":8082")
	dbConn := DatabaseConn(logger)
	dbStore := db.NewStore(dbConn)
	grpcServer := api.NewServer(dbStore, logger)
	grpcGateway := gateway.NewGateway(logger, grpcServerAddress, gprcGatewayAddress)
	server := app.NewApp(grpcServer, grpcGateway)
	err := server.Run()
	if err != nil {
		logger.WithError(err).Fatal("failed to start server")
		return
	}
}
func DatabaseConn(logger *logrus.Logger) *pgxpool.Pool {
	user := utils.GetEnv("DB_USER", "postgres")
	password := utils.GetEnv("DB_PASSWORD", "postgres")
	dbName := utils.GetEnv("DB_NAME", "postgres")
	host := utils.GetEnv("DB_HOST", "localhost")
	port := utils.GetEnv("DB_PORT", "5432")

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbName,
	)

	connPool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		logger.Fatal(err)
	}
	return connPool
}
