package interceptor

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/multi-tenants-cms-golang/lms-sys/app/gateway/middleware"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Enabled      bool
	SkipMethods  map[string]bool
	RoleMappings map[string][]middleware.Role
}

// ServiceAuthConfig holds service-specific auth configurations
type ServiceAuthConfig struct {
	ServiceName string
	Config      AuthConfig
}

// InterceptorManager manages gRPC interceptors
type InterceptorManager struct {
	logger             *logrus.Logger
	jwtSecret          string
	jwtIssuer          string
	jwtAudience        string
	serviceAuthConfigs []ServiceAuthConfig
}

func NewInterceptorManager(
	logger *logrus.Logger,
	jwtSecret string,
	jwtIssuer string,
	jwtAudience string,
	serviceAuthConfigs []ServiceAuthConfig,
) *InterceptorManager {
	return &InterceptorManager{
		logger:             logger,
		jwtSecret:          jwtSecret,
		jwtIssuer:          jwtIssuer,
		jwtAudience:        jwtAudience,
		serviceAuthConfigs: serviceAuthConfigs,
	}
}

// UnaryInterceptor returns a chained unary interceptor
func (im *InterceptorManager) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		methods := info.FullMethod
		serviceName, methodName := parseFullMethod(methods)

		// Get auth config for this service
		authConfig := im.getAuthConfig(serviceName)

		// Skip authentication if disabled for this service/method
		if !authConfig.Enabled || authConfig.SkipMethods[methodName] {
			return handler(ctx, req)
		}

		// Authentication logic
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

		// Verify JWT token
		userCtx, err := im.verifyToken(token)
		if err != nil {
			im.logger.WithError(err).Warn("JWT verification failed")
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		// Check authorization
		requiredRoles := authConfig.RoleMappings[methodName]
		if !im.hasRequiredRole(userCtx.Roles, requiredRoles) {
			im.logger.WithFields(logrus.Fields{
				"user_id": userCtx.UserID,
				"method":  info.FullMethod,
			}).Warn("Permission denied")
			return nil, status.Error(codes.PermissionDenied, "insufficient permissions")
		}

		// Add user context and proceed
		newCtx := context.WithValue(ctx, middleware.UserContextKey, userCtx)
		return handler(newCtx, req)
	}
}

// verifyToken validates JWT token and returns user context
func (im *InterceptorManager) verifyToken(tokenString string) (*middleware.UserContext, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(im.jwtSecret), nil
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

	if err := im.validateClaims(claims); err != nil {
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

// validateClaims validates standard JWT claims
func (im *InterceptorManager) validateClaims(claims jwt.MapClaims) error {
	iss, err := claims.GetIssuer()
	if err != nil || iss != im.jwtIssuer {
		return fmt.Errorf("invalid issuer")
	}

	aud, err := claims.GetAudience()
	if err != nil || !containss(aud, im.jwtAudience) {
		return fmt.Errorf("invalid audience")
	}

	if _, err := claims.GetExpirationTime(); err != nil {
		return fmt.Errorf("invalid expiration")
	}

	return nil
}

// hasRequiredRole checks if user has any of the required roles
func (im *InterceptorManager) hasRequiredRole(userRoles []middleware.Role, requiredRoles []middleware.Role) bool {
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

// getAuthConfig returns auth config for the specified service
func (im *InterceptorManager) getAuthConfig(serviceName string) AuthConfig {
	for _, config := range im.serviceAuthConfigs {
		if config.ServiceName == serviceName {
			return config.Config
		}
	}
	return AuthConfig{Enabled: true} // Default to secure
}

// parseFullMethod splits gRPC full method name into service and method
func parseFullMethod(fullMethod string) (service, method string) {

	parts := strings.Split(fullMethod, "/")
	if len(parts) != 3 {
		return "", ""
	}
	return parts[1], parts[2]
}

func containss(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
