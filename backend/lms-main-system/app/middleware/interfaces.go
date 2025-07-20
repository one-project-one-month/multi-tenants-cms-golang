package middleware

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// AuthInterceptorInterface defines methods for managing authentication and authorization interceptors
type AuthInterceptorInterface interface {
	// UnaryInterceptor returns a gRPC UnaryServerInterceptor for authentication/authorization
	UnaryInterceptor() grpc.UnaryServerInterceptor

	// VerifyToken parses and validates the JWT, returning user context if valid
	VerifyToken(tokenString string) (*UserContext, error)

	// ValidateClaims checks standard JWT claims like issuer, audience, and expiration
	ValidateClaims(claims map[string]interface{}) error

	// HasRequiredRole checks if user roles match required method roles
	HasRequiredRole(userRoles []Role, requiredRoles []Role) bool

	// GetAuthConfig retrieves auth config for a specific gRPC service
	GetAuthConfig(serviceName string) AuthConfig
	//MetadataLoggerInterceptor log every incoming metadata
	MetadataLoggerInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor
}
