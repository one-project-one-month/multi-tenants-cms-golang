package authentication

import (
	"context"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type AuthenticationService struct {
	authenticationpb.UnimplementedAuthenticationServiceServer
	databaseCtx context.Context
	store       *repo.Store
	logger      *logrus.Logger
	//redisClient *redis.Client
	//natsConn    *nats.Conn
}

func NewAuthenticationService(
	store *repo.Store,
	logger *logrus.Logger,
	redisClient *redis.Client,
	natsConn *nats.Conn,
) *AuthenticationService {
	return &AuthenticationService{
		databaseCtx: context.Background(),
		store:       store,
		logger:      logger,
		//redisClient: redisClient,
		//natsConn:    natsConn,
	}
}
