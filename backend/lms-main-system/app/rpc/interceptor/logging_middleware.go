package interceptor

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

func MetadataLoggerInterceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logIncomingMetadata(ctx, logger, info.FullMethod)

		return handler(ctx, req)
	}
}

func logIncomingMetadata(ctx context.Context, logger *logrus.Logger, method string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.WithFields(logrus.Fields{
			"method": method,
		}).Debug("No metadata found in context")
		return
	}

	fields := logrus.Fields{
		"method": method,
	}

	for key, values := range md {
		if isSensitiveHeader(key) {
			fields[key] = "[REDACTED]"
			continue
		}

		if len(values) > 0 {
			fields[key] = strings.Join(values, ", ")
		}
	}

	logger.WithFields(fields).Info("Incoming request metadata")
}

func isSensitiveHeader(key string) bool {
	key = strings.ToLower(key)
	return strings.Contains(key, "auth") ||
		strings.Contains(key, "token") ||
		strings.Contains(key, "password") ||
		strings.Contains(key, "secret") ||
		strings.Contains(key, "cookie")
}
