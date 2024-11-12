package utils

import (
	"github.com/cenkalti/backoff/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func CallEventsPlatformGRPCWithRetry[T any](operation func() (T, error)) (T, error) {
	operationWithData := func() (T, error) {
		result, opErr := operation()
		if opErr != nil {
			grpcError, ok := status.FromError(opErr)
			if ok && (grpcError.Code() == codes.Unavailable || grpcError.Code() == codes.DeadlineExceeded) {
				return result, opErr
			}
			return result, backoff.Permanent(opErr)
		}
		return result, nil
	}

	response, err := backoff.RetryWithData(operationWithData, BackOffForInvoking())
	return response, err
}

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

func BackOffForInvoking() backoff.BackOff {
	return BackOffConfig(50*time.Millisecond, 2, 2*time.Second, 7*time.Second, 10)
}
