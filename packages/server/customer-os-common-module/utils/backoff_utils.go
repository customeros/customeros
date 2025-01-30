package utils

import (
	"math"
	"math/rand"
	"time"
)

type BackoffConfig struct {
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Factor       float64 // The exponential factor (e.g., 2 for doubling)
	Jitter       float64 // Random jitter factor between 0-1 (e.g., 0.1 for 10% jitter)
}

func DefaultBackoffConfig() BackoffConfig {
	return BackoffConfig{
		InitialDelay: 50 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Factor:       2.0,
		Jitter:       0.1,
	}
}

func CalculateExponentialBackoffDelay(attempt int, config BackoffConfig) time.Duration {
	if attempt <= 0 {
		attempt = 1
	}

	// Calculate base delay
	delay := float64(config.InitialDelay) * math.Pow(config.Factor, float64(attempt-1))

	// Add jitter if configured
	if config.Jitter > 0 {
		jitterRange := delay * config.Jitter
		jitter := rand.Float64() * jitterRange
		delay += jitter
	}

	// Convert to time.Duration and cap at max delay
	actualDelay := time.Duration(delay)
	if actualDelay > config.MaxDelay {
		actualDelay = config.MaxDelay
	}

	return actualDelay
}

func BackOffExponentialDelay(attempt int) time.Duration {
	config := DefaultBackoffConfig()
	return CalculateExponentialBackoffDelay(attempt, config)
}

func BackOffIncrementalDelay(attempt int) time.Duration {
	if attempt <= 0 {
		attempt = 1
	}
	// Calculate the delay with a simple exponential backoff formula
	delay := time.Duration(attempt) * time.Millisecond * 50
	// Cap the delay at 2 seconds
	maxDelay := 2 * time.Second
	if delay > maxDelay {
		return maxDelay
	}
	return delay
}
