package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
)

// EntityDefinitionCreate is the resolver for the entityDefinitionCreate field.
func (r *mutationResolver) EntityDefinitionCreate(ctx context.Context, input model.EntityDefinitionInput) (*model.EntityDefinition, error) {
	entityDefinitionEntity, err := r.Services.EntityDefinitionService.Create(ctx, mapper.MapEntityDefinitionInputToEntity(input))
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to create entity definition: %s", input.Name)
		return nil, err
	}
	return mapper.MapEntityToEntityDefinition(entityDefinitionEntity), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
