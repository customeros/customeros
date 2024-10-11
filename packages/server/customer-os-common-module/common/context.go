package common

import (
	"context"
	"net/http"
)

type CustomContext struct {
	AppSource                string
	Tenant                   string
	UserId                   string
	UserEmail                string
	IdentityId               string
	Roles                    []string
	GraphqlRootOperationName string
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

func GetContext(ctx context.Context) *CustomContext {
	customContext, ok := ctx.Value(customContextKey).(*CustomContext)
	if !ok {
		return new(CustomContext)
	}
	return customContext
}

func GetAppSourceFromContext(ctx context.Context) string {
	return GetContext(ctx).AppSource
}

func GetTenantFromContext(ctx context.Context) string {
	return GetContext(ctx).Tenant
}

func GetRolesFromContext(ctx context.Context) []string {
	return GetContext(ctx).Roles
}

func GetUserIdFromContext(ctx context.Context) string {
	return GetContext(ctx).UserId
}

func GetUserEmailFromContext(ctx context.Context) string {
	return GetContext(ctx).UserEmail
}

func GetIdentityIdFromContext(ctx context.Context) string {
	return GetContext(ctx).IdentityId
}
