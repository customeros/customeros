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

// CustomFields is the resolver for the customFields field.
func (r *fieldSetResolver) CustomFields(ctx context.Context, obj *model.FieldSet) ([]*model.CustomField, error) {
	var customFields []*model.CustomField
	textCustomFieldEntities, err := r.ServiceContainer.CustomFieldService.FindAllForFieldSet(ctx, obj)
	for _, v := range mapper.MapEntitiesToCustomFields(textCustomFieldEntities) {
		customFields = append(customFields, v)
	}
	return customFields, err
}

// Definition is the resolver for the definition field.
func (r *fieldSetResolver) Definition(ctx context.Context, obj *model.FieldSet) (*model.FieldSetDefinition, error) {
	entity, err := r.ServiceContainer.FieldSetDefinitionService.FindLinkedWithFieldSet(ctx, obj.ID)
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to get contact definition for field set %s", obj.ID)
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}
	return mapper.MapEntityToFieldSetDefinition(entity), err
}

// FieldSet returns generated.FieldSetResolver implementation.
func (r *Resolver) FieldSet() generated.FieldSetResolver { return &fieldSetResolver{r} }

type fieldSetResolver struct{ *Resolver }
