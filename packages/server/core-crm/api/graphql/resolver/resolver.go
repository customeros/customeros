// api/graphql/resolver/resolver.go
package resolver

import (
	service "github.com/customeros/customeros/packages/server/customer-os-common-module/services"

	"github.com/customeros/customeros/packages/server/core-crm/api/graphql/generated"
)

// Resolver is the base resolver structure
type Resolver struct {
	services *service.CommonServices
}

func NewResolver(services *service.CommonServices) *Resolver {
	return &Resolver{
		services: services,
	}
}

// Query returns the QueryResolver implementation
func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

type mutationResolver struct {
	*Resolver
}

type (
	queryResolver struct{ *Resolver }
)
