package middleware

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Logger        *logrus.Logger
	TokenVerifier TokenVerifier
}

type TokenVerifier interface {
	VerifyToken(token string) (*UserContext, error)
}

func (c *AuthConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Normalize the path for consistent matching
		path := strings.ToLower(strings.TrimRight(r.URL.Path, "/"))

		if shouldSkipAuth(path) {
			next.ServeHTTP(w, r)
			return
		}

		token, err := extractBearerToken(r)
		if err != nil {
			c.Logger.WithError(err).WithField("path", r.URL.Path).Warn("Failed to extract token")
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		userCtx, err := c.TokenVerifier.VerifyToken(token)
		if err != nil {
			c.Logger.WithError(err).WithField("path", r.URL.Path).Warn("Token verification failed")
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, userCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

var (
	// Exact paths that skip auth
	exactSkipPaths = map[string]bool{
		"/healthz": true,
	}

	// Path prefixes that skip auth
	prefixSkipPaths = []string{
		"/swagger",
		"/grpc.reflection.v1alpha.serverreflection",
		"/lms.authentication.authenticationservice",
	}

	// Regex patterns for auth-skippable paths
	skipPatterns = []*regexp.Regexp{
		regexp.MustCompile(`^/lms/v1/[^/]+/register/?$`),
		regexp.MustCompile(`^/lms/v1/[^/]+/login/?$`),
		regexp.MustCompile(`^/lms/v1/[^/]+/verify-email/?`),
		regexp.MustCompile(`^/lms/v1/[^/]+/resend/?`),
		regexp.MustCompile(`^/lms/v1/auth/[^/]+/?$`),
	}
)

// shouldSkipAuth checks if a path should skip authentication
func shouldSkipAuth(path string) bool {
	normalizedPath := strings.ToLower(strings.TrimRight(path, "/"))
	if exactSkipPaths[normalizedPath] {
		return true
	}

	// Check prefix matches
	for _, prefix := range prefixSkipPaths {
		if strings.HasPrefix(normalizedPath, prefix) {
			return true
		}
	}

	// Check regex patterns
	for _, pattern := range skipPatterns {
		if pattern.MatchString(normalizedPath) {
			return true
		}
	}

	return false
}

func extractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	return parts[1], nil
}
