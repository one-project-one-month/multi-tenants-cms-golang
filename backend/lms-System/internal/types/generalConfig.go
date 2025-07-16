package types

import "time"

type Config struct {
	BaseURL           string
	TokenLength       int
	TokenTTL          time.Duration
	MaxRetries        int
	RetryDelay        time.Duration
	RateLimitAttempts int
	RateLimitWindow   time.Duration
}
