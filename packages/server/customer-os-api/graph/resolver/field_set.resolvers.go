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
	customFieldEntities, err := r.ServiceContainer.CustomFieldService.FindAllForFieldSet(ctx, obj)
	for _, v := range mapper.MapEntitiesToCustomFields(customFieldEntities) {
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

// FieldSetMergeToContact is the resolver for the fieldSetMergeToContact field.
func (r *mutationResolver) FieldSetMergeToContact(ctx context.Context, contactID string, input model.FieldSetInput) (*model.FieldSet, error) {
	result, err := r.ServiceContainer.FieldSetService.MergeFieldSetToContact(ctx, contactID, mapper.MapFieldSetInputToEntity(&input), input.DefinitionID)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not merge fields set <%s> to contact %s", input.Name, contactID)
		return nil, err
	}
	return mapper.MapEntityToFieldSet(result), nil
}

// FieldSetUpdateInContact is the resolver for the fieldSetUpdateInContact field.
func (r *mutationResolver) FieldSetUpdateInContact(ctx context.Context, contactID string, input model.FieldSetUpdateInput) (*model.FieldSet, error) {
	result, err := r.ServiceContainer.FieldSetService.UpdateFieldSetInContact(ctx, contactID, mapper.MapFieldSetUpdateInputToEntity(&input))
	if err != nil {
		graphql.AddErrorf(ctx, "Could not update fields set %s in contact %s", input.ID, contactID)
		return nil, err
	}
	return mapper.MapEntityToFieldSet(result), nil
}

// FieldSetDeleteFromContact is the resolver for the fieldSetDeleteFromContact field.
func (r *mutationResolver) FieldSetDeleteFromContact(ctx context.Context, contactID string, id string) (*model.Result, error) {
	result, err := r.ServiceContainer.FieldSetService.DeleteByIdFromContact(ctx, contactID, id)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not remove fields set %s from contact %s", id, contactID)
		return nil, err
	}
	return &model.Result{
		Result: result,
	}, nil
}

// FieldSet returns generated.FieldSetResolver implementation.
func (r *Resolver) FieldSet() generated.FieldSetResolver { return &fieldSetResolver{r} }

type fieldSetResolver struct{ *Resolver }
