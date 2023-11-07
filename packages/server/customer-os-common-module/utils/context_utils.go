package utils

import (
	"context"
	"time"
)

const (
	ShortDuration  = 500 * time.Millisecond
	MediumDuration = 2 * time.Second
	LongDuration   = 20 * time.Second
)

func getContextWithTimeout(ctx context.Context, duration time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, duration)
}

func GetShortLivedContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return getContextWithTimeout(ctx, ShortDuration)
}

func GetMediumLivedContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return getContextWithTimeout(ctx, MediumDuration)
}

func GetLongLivedContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return getContextWithTimeout(ctx, LongDuration)
}
