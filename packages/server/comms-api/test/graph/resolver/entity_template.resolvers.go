package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/model"
)

// FieldSets is the resolver for the fieldSets field.
func (r *entityTemplateResolver) FieldSets(ctx context.Context, obj *model.EntityTemplate) ([]*model.FieldSetTemplate, error) {
	panic(fmt.Errorf("not implemented: FieldSets - fieldSets"))
}

// CustomFields is the resolver for the customFields field.
func (r *entityTemplateResolver) CustomFields(ctx context.Context, obj *model.EntityTemplate) ([]*model.CustomFieldTemplate, error) {
	panic(fmt.Errorf("not implemented: CustomFields - customFields"))
}

// CustomFields is the resolver for the customFields field.
func (r *fieldSetTemplateResolver) CustomFields(ctx context.Context, obj *model.FieldSetTemplate) ([]*model.CustomFieldTemplate, error) {
	panic(fmt.Errorf("not implemented: CustomFields - customFields"))
}

// EntityTemplateCreate is the resolver for the entityTemplateCreate field.
func (r *mutationResolver) EntityTemplateCreate(ctx context.Context, input model.EntityTemplateInput) (*model.EntityTemplate, error) {
	panic(fmt.Errorf("not implemented: EntityTemplateCreate - entityTemplateCreate"))
}

// EntityTemplate returns generated.EntityTemplateResolver implementation.
func (r *Resolver) EntityTemplate() generated.EntityTemplateResolver {
	return &entityTemplateResolver{r}
}

// FieldSetTemplate returns generated.FieldSetTemplateResolver implementation.
func (r *Resolver) FieldSetTemplate() generated.FieldSetTemplateResolver {
	return &fieldSetTemplateResolver{r}
}

type entityTemplateResolver struct{ *Resolver }
type fieldSetTemplateResolver struct{ *Resolver }
