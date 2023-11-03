package helper

import (
	"math"
	"strings"
	"time"
)

func CheckErrMessages(err error, messages ...string) bool {
	for _, message := range messages {
		if strings.Contains(strings.TrimSpace(strings.ToLower(err.Error())), strings.TrimSpace(strings.ToLower(message))) {
			return true
		}
	}
	return false
}

// Implement a backoffDelay function that calculates the delay before the next retry.
func BackoffDelay(attempt int) time.Duration {
	if attempt <= 0 {
		attempt = 1
	}
	// Calculate the delay with a simple exponential backoff formula
	delay := time.Duration(math.Pow(2, float64(attempt))) * time.Millisecond * 50
	// Cap the delay at 5 seconds
	maxDelay := 5 * time.Second
	if delay > maxDelay {
		return maxDelay
	}
	return delay
}
