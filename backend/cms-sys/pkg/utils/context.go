package utils

import (
	"context"
	"time"
)

func GetContext() context.Context {
	return context.Background()
}

func GetContextWithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

func GetContextWithCancel() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

func GetContextWithDeadline(deadline time.Time) (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), deadline)
}
