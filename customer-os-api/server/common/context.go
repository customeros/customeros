package common

import (
	"context"
	"net/http"
)

type CustomContext struct {
	Tenant string
}

var customContextKey = "CUSTOM_CONTEXT"

func CreateContext(args *CustomContext, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customContext := &CustomContext{
			Tenant: args.Tenant,
		}
		requestWithCtx := r.WithContext(context.WithValue(r.Context(), customContextKey, customContext))
		next.ServeHTTP(w, requestWithCtx)
	})
}

func GetContext(ctx context.Context) *CustomContext {
	customContext, ok := ctx.Value(customContextKey).(*CustomContext)
	if !ok {
		return nil
	}
	return customContext
}
