package rpc

import (
	"context"
	"fmt"
	"net"
	"syscall"

	"github.com/golang-jwt/jwt/v5"
	"github.com/multi-tenants-cms-golang/lms-sys/app/gateway/middleware"
	authSrv "github.com/multi-tenants-cms-golang/lms-sys/app/rpc/authentication"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/types"
	authpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/oklog/run"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	grpcAddr = ":9001"
)

type Server struct {
	store       db.Store
	logger      *logrus.Logger
	jwtSecret   string
	jwtIssuer   string
	jwtAudience string
}

func NewServer(
	db db.Store,
	logger *logrus.Logger,
	jwtSecret string,
	jwtIssuer string,
	jwtAudience string,
) *Server {
	return &Server{
		store:       db,
		logger:      logger,
		jwtSecret:   jwtSecret,
		jwtIssuer:   jwtIssuer,
		jwtAudience: jwtAudience,
	}
}

func (s *Server) Run() error {
	grpcServer := grpc.NewServer()

	authService := authSrv.NewAuthenticationService(s.store, s.logger, &types.Config{})
	authpb.RegisterAuthenticationServiceServer(grpcServer, authService)

	var g run.Group

	// gRPC server
	g.Add(func() error {
		listener, err := net.Listen("tcp", grpcAddr)
		if err != nil {
			return fmt.Errorf("failed to listen: %v", err)
		}
		s.logger.Infof("Starting gRPC server on %s", grpcAddr)
		return grpcServer.Serve(listener)
	}, func(err error) {
		grpcServer.GracefulStop()
	})

	// Signal handler
	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	return g.Run()
}

// JWT verification helper
func (s *Server) verifyToken(tokenString string) (*middleware.UserContext, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	//
	//// Validate standard claims
	//if err := claims.VerifyIssuer(s.jwtIssuer, true); err != nil {
	//	return nil, fmt.Errorf("invalid issuer: %w", err)
	//}
	//
	//if err := claims.VerifyAudience(s.jwtAudience, true); err != nil {
	//	return nil, fmt.Errorf("invalid audience: %w", err)
	//}

	// Extract user info
	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("missing or invalid subject claim")
	}

	rolesClaim, ok := claims["roles"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("missing or invalid roles claim")
	}

	var roles []middleware.Role
	for _, r := range rolesClaim {
		if roleStr, ok := r.(string); ok {
			roles = append(roles, middleware.Role(roleStr))
		}
	}

	if len(roles) == 0 {
		return nil, fmt.Errorf("no valid roles found in token")
	}

	return &middleware.UserContext{
		UserID: userID,
		Roles:  roles,
	}, nil
}
