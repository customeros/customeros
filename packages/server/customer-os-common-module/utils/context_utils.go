package utils

import (
	"context"
	"time"
)

const (
	ShortDuration      = 500 * time.Millisecond
	MediumDuration     = 2 * time.Second
	MediumLongDuration = 10 * time.Second
	LongDuration       = 20 * time.Second
)

func GetContextWithTimeout(ctx context.Context, duration time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, duration)
}

func GetShortLivedContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return GetContextWithTimeout(ctx, ShortDuration)
}

func GetMediumLivedContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return GetContextWithTimeout(ctx, MediumDuration)
}

func GetMediumLongLivedContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return GetContextWithTimeout(ctx, MediumLongDuration)
}

func GetLongLivedContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return GetContextWithTimeout(ctx, LongDuration)
}
