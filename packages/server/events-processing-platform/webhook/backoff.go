package webhook

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

type callable func(*http.Request) (*http.Response, error)

// RetryWithBackoff retries the given operation with exponential backoff
func RetryWithBackoff(operation callable, request *http.Request, maxRetries int, baseDelay time.Duration) callable {
	return func(*http.Request) (*http.Response, error) {
		var lastError error

		for i := 0; i <= maxRetries; i++ {
			val, err := operation(request)
			if err == nil {
				return val, nil
			}
			secRetry := math.Pow(2, float64(i))
			fmt.Printf("Retrying operation in %f seconds\n", secRetry)
			delay := time.Duration(secRetry) * baseDelay
			time.Sleep(delay)
			lastError = err
		}

		return nil, lastError
	}
}
