package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
)

// ContactTypeCreate is the resolver for the contactType_Create field.
func (r *mutationResolver) ContactTypeCreate(ctx context.Context, input model.ContactTypeInput) (*model.ContactType, error) {
	createdContactType, err := r.ServiceContainer.ContactTypeService.Create(ctx, mapper.MapContactTypeInputToEntity(input))
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to create contact type %s", input.Name)
		return nil, err
	}
	return mapper.MapEntityToContactType(createdContactType), nil
}

// ContactTypeUpdate is the resolver for the contactType_Update field.
func (r *mutationResolver) ContactTypeUpdate(ctx context.Context, input model.ContactTypeUpdateInput) (*model.ContactType, error) {
	updatedContactType, err := r.ServiceContainer.ContactTypeService.Update(ctx, mapper.MapContactTypeUpdateInputToEntity(input))
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to update contact type %s", input.ID)
		return nil, err
	}
	return mapper.MapEntityToContactType(updatedContactType), nil
}

// ContactTypeDelete is the resolver for the contactType_Delete field.
func (r *mutationResolver) ContactTypeDelete(ctx context.Context, id string) (*model.Result, error) {
	result, err := r.ServiceContainer.ContactTypeService.Delete(ctx, id)
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to delete contact type %s", id)
		return nil, err
	}
	return &model.Result{
		Result: result,
	}, nil
}

// ContactTypes is the resolver for the contactTypes field.
func (r *queryResolver) ContactTypes(ctx context.Context) ([]*model.ContactType, error) {
	contactTypes, err := r.ServiceContainer.ContactTypeService.GetAll(ctx)
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to fetch contact types")
		return nil, err
	}
	return mapper.MapEntitiesToContactTypes(contactTypes), err
}
