package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/mapper"
)

// Definition is the resolver for the definition field.
func (r *customFieldResolver) Definition(ctx context.Context, obj *model.CustomField) (*model.CustomFieldDefinition, error) {
	entity, err := r.ServiceContainer.CustomFieldDefinitionService.FindLinkedWithCustomField(ctx, obj.ID)
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to get contact definition for custom field %s", obj.ID)
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}
	return mapper.MapEntityToCustomFieldDefinition(entity), err
}

// CustomField returns generated.CustomFieldResolver implementation.
func (r *Resolver) CustomField() generated.CustomFieldResolver { return &customFieldResolver{r} }

type customFieldResolver struct{ *Resolver }
