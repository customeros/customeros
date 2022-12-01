package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

// EntityDefinitions is the resolver for the entityDefinitions field.
func (r *queryResolver) EntityDefinitions(ctx context.Context, extends *model.EntityDefinitionExtension) ([]*model.EntityDefinition, error) {
	result, err := r.Services.EntityDefinitionService.FindAll(ctx, utils.StringPtr(extends.String()))
	return mapper.MapEntitiesToEntityDefinitions(result), err
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
