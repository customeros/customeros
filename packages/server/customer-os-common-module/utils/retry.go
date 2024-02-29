package utils

import (
	"github.com/cenkalti/backoff/v4"
	"time"
)

func BackOffConfig(initialInterval time.Duration, multiplier float64, maxInterval time.Duration, maxElapsedTime time.Duration, maxRetries uint64) backoff.BackOff {
	// Customize the backoff configuration
	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.InitialInterval = initialInterval
	backoffConfig.Multiplier = multiplier         // Exponential factor
	backoffConfig.MaxInterval = maxInterval       // Ensure individual retry does not exceed this
	backoffConfig.MaxElapsedTime = maxElapsedTime // Total time for all retries

	// Apply a maximum number of retries by wrapping the backoff with a MaxRetries policy
	backoffWithMaxRetries := backoff.WithMaxRetries(backoffConfig, maxRetries)
	return backoffWithMaxRetries
}

func BackOffForInvokingEventsPlatformGrpcClient() backoff.BackOff {
	return BackOffConfig(50*time.Millisecond, 2, 2*time.Second, 7*time.Second, 10)
}
