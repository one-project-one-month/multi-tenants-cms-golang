package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type JWTClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	TenantID string `json:"tenant_id"`
	jwt.RegisteredClaims
}

type AuthMiddleware struct {
	jwtSecret   []byte
	publicPaths map[string]bool
}

func NewAuthMiddleware(jwtSecret string) *AuthMiddleware {
	publicPaths := map[string]bool{
		"/authentication.AuthenticationService/": true,
		"/auth/":                                 true,
	}

	return &AuthMiddleware{
		jwtSecret:   []byte(jwtSecret),
		publicPaths: publicPaths,
	}
}

func (a *AuthMiddleware) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if a.publicPaths[info.FullMethod] {
			return handler(ctx, req)
		}

		newCtx, err := a.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		return handler(newCtx, req)
	}
}

func (a *AuthMiddleware) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if a.publicPaths[info.FullMethod] {
			return handler(srv, stream)
		}

		newCtx, err := a.authenticate(stream.Context())
		if err != nil {
			return err
		}

		wrapped := &wrappedStream{
			ServerStream: stream,
			ctx:          newCtx,
		}

		return handler(srv, wrapped)
	}
}

func (a *AuthMiddleware) authenticate(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization header")
	}

	authHeader := authHeaders[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := a.validateToken(tokenString)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	newCtx := context.WithValue(ctx, "user_claims", claims)
	return newCtx, nil
}

func (a *AuthMiddleware) validateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return a.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenMalformed
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, jwt.ErrTokenMalformed
	}

	// Check if token is expired
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, jwt.ErrTokenExpired
	}

	return claims, nil
}

type wrappedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedStream) Context() context.Context {
	return w.ctx
}

func GetUserClaims(ctx context.Context) (*JWTClaims, error) {
	claims, ok := ctx.Value("user_claims").(*JWTClaims)
	if !ok {
		return nil, status.Error(codes.Internal, "failed to get user claims from context")
	}
	return claims, nil
}
