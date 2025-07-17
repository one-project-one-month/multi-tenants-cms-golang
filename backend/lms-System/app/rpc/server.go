package rpc

import (
	"context"
	"fmt"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/types"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/multi-tenants-cms-golang/lms-sys/app/gateway/middleware"
	authSrv "github.com/multi-tenants-cms-golang/lms-sys/app/rpc/authentication"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	authpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/authentication"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
)

type Server struct {
	store       *db.Store
	logger      *logrus.Logger
	jwtSecret   string
	jwtIssuer   string
	jwtAudience string
}

func NewServer(
	db *db.Store,
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
	listener, err := net.Listen("tcp", ":9001")
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
	//grpc.UnaryInterceptor(s.unaryInterceptor()),
	)

	authService := authSrv.NewAuthenticationService(s.store, s.logger, &types.Config{})
	authpb.RegisterAuthenticationServiceServer(grpcServer, authService)
	s.logger.Info("Starting gRPC server on port 9001")
	if err := grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil

}

func (s *Server) unaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if shouldSkipAuth(info.FullMethod) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata not provided")
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token not provided")
		}

		token := strings.TrimPrefix(authHeaders[0], "Bearer ")
		if token == "" {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
		}

		userCtx, err := s.verifyToken(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		if !s.hasRequiredRole(userCtx.Roles, info.FullMethod) {
			return nil, status.Error(codes.PermissionDenied, "insufficient permissions")
		}

		newCtx := context.WithValue(ctx, middleware.UserContextKey, userCtx)
		return handler(newCtx, req)
	}
}

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

	if err := s.validateClaims(claims); err != nil {
		return nil, err
	}

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

func (s *Server) validateClaims(claims jwt.MapClaims) error {
	iss, err := claims.GetIssuer()
	if err != nil || iss != s.jwtIssuer {
		return fmt.Errorf("invalid issuer")
	}

	aud, err := claims.GetAudience()
	if err != nil || !contains(aud, s.jwtAudience) {
		return fmt.Errorf("invalid audience")
	}

	if _, err := claims.GetExpirationTime(); err != nil {
		return fmt.Errorf("invalid expiration")
	}

	return nil
}

func (s *Server) hasRequiredRole(userRoles []middleware.Role, method string) bool {
	requiredRoles := s.getRequiredRoles(method)
	if len(requiredRoles) == 0 {
		return true
	}

	for _, userRole := range userRoles {
		for _, requiredRole := range requiredRoles {
			if userRole == requiredRole {
				return true
			}
		}
	}
	return false
}

func (s *Server) getRequiredRoles(method string) []middleware.Role {
	roleMapping := map[string][]middleware.Role{
		"/lms_user.LMS_USER_SERVICE/AdminMethod":      {middleware.LMSAdmin},
		"/lms_user.LMS_USER_SERVICE/InstructorMethod": {middleware.Instructor, middleware.LMSAdmin},
		"/lms_user.LMS_USER_SERVICE/StudentMethod":    {middleware.Student},
		"/lms_user.LMS_USER_SERVICE/ProtectedMethod":  {middleware.Student, middleware.Instructor, middleware.LMSAdmin},
	}

	if roles, ok := roleMapping[method]; ok {
		return roles
	}
	return nil
}

func shouldSkipAuth(method string) bool {
	skipMethods := []string{
		"/lms_user.LMS_USER_SERVICE/PublicMethod",
		"/grpc.health.v1.Health/Check",
	}

	for _, m := range skipMethods {
		if m == method {
			return true
		}
	}
	return false
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
