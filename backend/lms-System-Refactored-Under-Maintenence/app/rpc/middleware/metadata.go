package middleware

import (
	"context"
	"google.golang.org/grpc/metadata"
)

func ExtractMetaData(ctx context.Context, key string) (string, bool) {
	md, condition := metadata.FromIncomingContext(ctx)
	if !condition {
		return "", false
	}
	values := md.Get(key)
	if len(values) != 1 {
		return "", false
	}
	return values[0], true
}

func GetAllMetaData(ctx context.Context, key string) ([]string, bool) {
	md, condition := metadata.FromIncomingContext(ctx)
	if !condition {
		return nil, false
	}
	return md.Get(key), true
}
