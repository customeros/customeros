// api/graphql/resolver/resolver.go
package resolver

import (
	"github.com/customeros/customeros/packages/server/leads/api/graphql/generated"
	"github.com/customeros/customeros/packages/server/leads/services"
)

// Resolver is the base resolver structure
type Resolver struct {
	services *services.Services
}

func NewResolver(services *services.Services) *Resolver {
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
