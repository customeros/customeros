package handler

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
)

func GetTenantChecker() func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		tenant := common.GetTenantFromContext(ctx)

		if tenant == "" {
			return nil, fmt.Errorf("Tenant is required")
		}
		return next(ctx)
	}
}
