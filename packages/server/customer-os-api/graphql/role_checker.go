package graphql

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

func GetRoleChecker() func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []model.Role) (res interface{}, err error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []model.Role) (res interface{}, err error) {
		// Start a span for role checking that will be a child of the current root
		spans, _ := telemetry.StartSpan(ctx, "RoleChecker")
		defer spans.Finish()
		currentRoles := common.GetRolesFromContext(ctx)

		spans.LogKV("requiredRoles", roles)
		spans.LogKV("currentRoles", currentRoles)

		// Check if the current role is in the list of allowed roles
		for _, allowedRole := range roles {
			for _, currentRole := range currentRoles {
				if currentRole == allowedRole.String() {
					// If the role is in the list of allowed roles, call the next resolver
					spans.LogKV("result", "Access granted")
					// Pass the original context to next resolver to avoid nesting spans
					return next(ctx)
				}
			}
		}
		spans.LogKV("result", "Access denied")
		// If the role is not in the list of allowed roles, return an error
		return nil, coserrors.ErrAccessDenied
	}
}
