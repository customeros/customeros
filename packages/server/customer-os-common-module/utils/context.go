package utils

import (
	"context"

	"github.com/nats-io/nats.go"
)

type CustomContext struct {
	Tenant    string
	UserId    string
	UserEmail string
	RequestID string
}

var customContextKey = "CUSTOM_CONTEXT"

func WithCustomContext(ctx context.Context, customContext *CustomContext) context.Context {
	return context.WithValue(ctx, customContextKey, customContext)
}

func WithCustomContextFromNats(ctx context.Context, msg *nats.Msg) context.Context {
	if msg.Header == nil {
		return ctx
	}

	// Create custom context from message headers
	customContext := &CustomContext{
		Tenant: msg.Header.Get("X-Tenant"),
		UserId: msg.Header.Get("X-UserId"),
	}
	return WithCustomContext(ctx, customContext)
}
