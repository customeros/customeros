package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/customer-os-api/mapper"

	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

// FieldSets is the resolver for the fieldSets field.
func (r *entityDefinitionResolver) FieldSets(ctx context.Context, obj *model.EntityDefinition) ([]*model.FieldSetDefinition, error) {
	result, err := r.ServiceContainer.FieldSetDefinitionService.FindAll(obj.ID)
	return mapper.MapEntitiesToFieldSetDefinitions(result), err
}

// CustomFields is the resolver for the customFields field.
func (r *entityDefinitionResolver) CustomFields(ctx context.Context, obj *model.EntityDefinition) ([]*model.CustomFieldDefinition, error) {
	panic(fmt.Errorf("not implemented: CustomFields - customFields"))
}

// CustomFields is the resolver for the customFields field.
func (r *fieldSetDefinitionResolver) CustomFields(ctx context.Context, obj *model.FieldSetDefinition) ([]*model.CustomFieldDefinition, error) {
	panic(fmt.Errorf("not implemented: CustomFields - customFields"))
}

// EntityDefinition returns generated.EntityDefinitionResolver implementation.
func (r *Resolver) EntityDefinition() generated.EntityDefinitionResolver {
	return &entityDefinitionResolver{r}
}

// FieldSetDefinition returns generated.FieldSetDefinitionResolver implementation.
func (r *Resolver) FieldSetDefinition() generated.FieldSetDefinitionResolver {
	return &fieldSetDefinitionResolver{r}
}

type entityDefinitionResolver struct{ *Resolver }
type fieldSetDefinitionResolver struct{ *Resolver }
