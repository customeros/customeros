package utils

import (
	"context"
)

type CustomContext struct {
	Tenant string
	Source string
	RunId  string
}

var customContextKey = "CUSTOM_CONTEXT"

func WithCustomContext(ctx context.Context, customContext *CustomContext) context.Context {
	return context.WithValue(ctx, customContextKey, customContext)
}

func GetContext(ctx context.Context) *CustomContext {
	customContext, ok := ctx.Value(customContextKey).(*CustomContext)
	if !ok {
		return new(CustomContext)
	}
	return customContext
}

func GetTenantFromContext(ctx context.Context) string {
	return GetContext(ctx).Tenant
}

func GetSourceFromContext(ctx context.Context) string {
	return GetContext(ctx).Source
}

func GetRunIdFromContext(ctx context.Context) string {
	return GetContext(ctx).RunId
}
