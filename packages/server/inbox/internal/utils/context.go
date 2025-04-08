package utils

import (
	"context"
	"net/http"

	"github.com/nats-io/nats.go"

	inbox_errors "github.com/customeros/customeros/packages/server/inbox/errors"
)

type contextKey string

const (
	tenantKey contextKey = "tenant"
	userIDKey contextKey = "user_id"
)

// TenantHeaders is a list of possible header names for tenant identification
var TenantHeaders = []string{
	"X-TENANT",
	"x-tenant",
	"X-Tenant",
	"tenant",
	"Tenant",
	"TENANT",
	"tenantname",
	"TenantName",
	"tenantName",
	"TENANTNAME",
}

var UserIdHeaders = []string{
	"X-USER-ID",
	"X-USERID",
	"X-User-Id",
	"X-UserId",
	"X-User-ID",
	"X-Userid",
	"x-user-id",
	"User-ID",
	"UserId",
	"Userid",
	"USERID",
}

type CustomContext struct {
	Tenant    string
	UserId    string
	UserEmail string
	RequestID string
}

var customContextKey = "CUSTOM_CONTEXT"

func WithContext(customContext *CustomContext, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestWithCtx := r.WithContext(context.WithValue(r.Context(), customContextKey, customContext))
		next.ServeHTTP(w, requestWithCtx)
	})
}

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

func GetUserIdFromContext(ctx context.Context) string {
	return GetContext(ctx).UserId
}

func GetUserEmailFromContext(ctx context.Context) string {
	return GetContext(ctx).UserEmail
}

func SetAppSourceInContext(ctx context.Context, appSource string) context.Context {
	customContext := GetContext(ctx)
	return WithCustomContext(ctx, customContext)
}

func SetUserIdInContext(ctx context.Context, userId string) context.Context {
	customContext := GetContext(ctx)
	customContext.UserId = userId
	return WithCustomContext(ctx, customContext)
}

func SetTenantInContext(ctx context.Context, tenant string) context.Context {
	customContext := GetContext(ctx)
	customContext.Tenant = tenant
	return WithCustomContext(ctx, customContext)
}

func ValidateTenant(ctx context.Context) error {
	if GetTenantFromContext(ctx) == "" {
		return inbox_errors.ErrTenantMissing
	}
	return nil
}

func ValidateUserId(ctx context.Context) error {
	if GetUserIdFromContext(ctx) == "" {
		return inbox_errors.ErrUserIdMissing
	}
	return nil
}

// WithTenantContext creates a new context with the specified tenant
func WithTenantContext(ctx context.Context, tenant string) context.Context {
	return WithCustomContext(ctx, &CustomContext{
		Tenant: tenant,
	})
}

// WithTenant adds tenant to context
func WithTenant(ctx context.Context, tenant string) context.Context {
	return context.WithValue(ctx, tenantKey, tenant)
}

// WithUserId adds user ID to context
func WithUserId(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
