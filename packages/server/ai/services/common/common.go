package common

import "time"

const (
	MAX_RETRIES = 3
)

func ShouldRetry(attempt int) bool {
	if attempt >= MAX_RETRIES {
		return false
	}
	return true
}

func HandleRetry(attempt int, err error) {
	backoff := time.Duration(attempt) * time.Second
	time.Sleep(backoff)
}
