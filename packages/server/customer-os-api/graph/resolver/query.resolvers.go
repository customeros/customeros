package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/mapper"
)

// EntityDefinitions is the resolver for the entityDefinitions field.
func (r *queryResolver) EntityDefinitions(ctx context.Context) ([]*model.EntityDefinition, error) {
	result, err := r.ServiceContainer.EntityDefinitionService.FindAll(ctx)
	return mapper.MapEntitiesToEntityDefinitions(result), err
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
