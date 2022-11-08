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
func (r *fieldSetResolver) TextCustomFields(ctx context.Context, obj *model.FieldSet) ([]*model.TextCustomField, error) {
	textCustomFieldEntities, err := r.ServiceContainer.TextCustomFieldService.FindAllForFieldSet(ctx, obj)
	return mapper.MapEntitiesToTextCustomFields(textCustomFieldEntities), err
}

// FieldSet returns generated.FieldSetResolver implementation.
func (r *Resolver) FieldSet() generated.FieldSetResolver { return &fieldSetResolver{r} }

type fieldSetResolver struct{ *Resolver }
