package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/mapper"
)

// TextCustomFields is the resolver for the textCustomFields field.
func (r *fieldsSetResolver) TextCustomFields(ctx context.Context, obj *model.FieldsSet) ([]*model.TextCustomField, error) {
	textCustomFieldEntities, err := r.ServiceContainer.TextCustomFieldService.FindAllForFieldsSet(ctx, obj)
	return mapper.MapEntitiesToTextCustomFields(textCustomFieldEntities), err
}

// FieldsSet returns generated.FieldsSetResolver implementation.
func (r *Resolver) FieldsSet() generated.FieldsSetResolver { return &fieldsSetResolver{r} }

type fieldsSetResolver struct{ *Resolver }
