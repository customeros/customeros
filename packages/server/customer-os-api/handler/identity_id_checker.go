package handler

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
)

func GetIdentityIdChecker() func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		identityId := common.GetIdentityIdFromContext(ctx)

		if identityId == "" {
			return nil, fmt.Errorf("IdentityId is required")
		}
		return next(ctx)
	}
}
