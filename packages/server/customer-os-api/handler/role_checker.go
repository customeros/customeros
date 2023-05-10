package handler

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
)

func GetRoleChecker() func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []model.Role) (res interface{}, err error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []model.Role) (res interface{}, err error) {
		currentRole := common.GetRoleFromContext(ctx)
		// Check if the current role is in the list of allowed roles
		for _, allowedRole := range roles {
			if currentRole == allowedRole {
				// If the role is in the list of allowed roles, call the next resolver
				return next(ctx)
			}
		}
		// If the role is not in the list of allowed roles, return an error
		return nil, fmt.Errorf("Access denied")
	}
}
