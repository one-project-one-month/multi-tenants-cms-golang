package middleware

import (
	"context"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strings"
)

type GatewayAuthMiddleware struct {
	authMiddleware *AuthMiddleware
	publicPaths    map[string]bool
}

func NewGatewayAuthMiddleware(jwtSecret string) *GatewayAuthMiddleware {
	publicPaths := map[string]bool{
		"/api/v1/auth/login":           true,
		"/api/v1/auth/register":        true,
		"/api/v1/auth/refresh-token":   true,
		"/api/v1/auth/forgot-password": true,
		"/api/v1/auth/reset-password":  true,
		"/api/v1/health":               true,
		// Add other public REST endpoints here
	}

	return &GatewayAuthMiddleware{
		authMiddleware: NewAuthMiddleware(jwtSecret),
		publicPaths:    publicPaths,
	}
}

// HTTPMiddleware for REST API gateway
func (g *GatewayAuthMiddleware) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the path is public
		if g.publicPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		// Extract JWT token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := g.authMiddleware.validateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_claims", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (g *GatewayAuthMiddleware) MetadataAnnotator(ctx context.Context, req *http.Request) metadata.MD {
	md := metadata.New(nil)

	if auth := req.Header.Get("Authorization"); auth != "" {
		md.Set("authorization", auth)
	}

	if userAgent := req.Header.Get("User-Agent"); userAgent != "" {
		md.Set("user-agent", userAgent)
	}

	if clientIP := req.Header.Get("X-Forwarded-For"); clientIP != "" {
		md.Set("x-forwarded-for", clientIP)
	} else if clientIP := req.Header.Get("X-Real-IP"); clientIP != "" {
		md.Set("x-real-ip", clientIP)
	}

	return md
}
