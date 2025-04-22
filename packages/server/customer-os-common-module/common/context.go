package common

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

type CustomContext struct {
	AppSource  string
	Tenant     string
	AuthUserId string
	UserId     string
	UserEmail  string
	Roles      []string
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

func WithCustomContextFromGinRequest(c *gin.Context, appSource string) context.Context {
	customContext := &CustomContext{
		AppSource:  appSource,
		AuthUserId: c.GetString("AuthenticatedUserId"),
		Tenant:     c.GetString("TenantName"),
		UserId:     c.GetString("UserId"),
		UserEmail:  strings.ToLower(c.GetString("UserEmail")),
		Roles:      c.GetStringSlice("UserRoles"),
	}
	return WithCustomContext(c.Request.Context(), customContext)
}

func WithCustomContextFromNats(ctx context.Context, msg *nats.Msg) context.Context {
	if msg.Header == nil {
		return ctx
	}

	// Create custom context from message headers
	customContext := &CustomContext{
		Tenant:    msg.Header.Get("X-Tenant"),
		UserId:    msg.Header.Get("X-UserId"),
		AppSource: msg.Header.Get("X-AppSource"),
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

func GetAppSourceFromContext(ctx context.Context) string {
	return GetContext(ctx).AppSource
}

func GetTenantFromContext(ctx context.Context) string {
	return GetContext(ctx).Tenant
}

func GetRolesFromContext(ctx context.Context) []string {
	return GetContext(ctx).Roles
}

func GetAuthUserIdFromContext(ctx context.Context) string {
	return GetContext(ctx).AuthUserId
}

func GetUserIdFromContext(ctx context.Context) string {
	return GetContext(ctx).UserId
}

func GetUserEmailFromContext(ctx context.Context) string {
	return GetContext(ctx).UserEmail
}

func IsImpersonatedUserInContext(ctx context.Context) bool {
	roles := GetRolesFromContext(ctx)
	if len(roles) == 0 {
		return false
	}
	for _, role := range roles {
		if role == "IMPERSONATED" {
			return true
		}
	}
	return false
}

func SetAppSourceInContext(ctx context.Context, appSource string) context.Context {
	customContext := GetContext(ctx)
	customContext.AppSource = appSource
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
		return errors.New("tenant is missing")
	}
	return nil
}
