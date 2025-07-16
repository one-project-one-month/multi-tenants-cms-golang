package authentication

import (
	"context"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/types"
	authenticationpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/sirupsen/logrus"
)

type AuthenticationService struct {
	authenticationpb.UnimplementedAuthenticationServiceServer
	databaseCtx context.Context
	store       *repo.Store
	logger      *logrus.Logger
	cfg         *types.Config
	//redisClient *redis.Client
	//natsConn    *nats.Conn
}

func NewAuthenticationService(
	store *repo.Store,
	logger *logrus.Logger,
	config *types.Config,
) *AuthenticationService {
	return &AuthenticationService{
		databaseCtx: context.Background(),
		store:       store,
		logger:      logger,
		cfg:         config,
		//redisClient: redisClient,
		//natsConn:    natsConn,
	}
}
