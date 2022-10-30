package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/mapper"
)

// CreateContact is the resolver for the createContact field.
func (r *mutationResolver) CreateContact(ctx context.Context, input model.ContactInput) (*model.Contact, error) {
	contactNodeCreated, err := r.ServiceContainer.ContactService.Create(ctx, mapper.MapContactInputToEntity(input))
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to create contact %s %s", input.FirstName, input.LastName)
		return nil, err
	}

	return mapper.MapEntityToContact(contactNodeCreated), nil
}

// AddContactToGroup is the resolver for the addContactToGroup field.
func (r *mutationResolver) AddContactToGroup(ctx context.Context, contactID string, groupID string) (*model.BooleanResult, error) {
	result, err := r.ServiceContainer.ContactWithContactGroupRelationshipService.AddContactToGroup(ctx, contactID, groupID)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not add contact to group")
		return nil, err
	}
	return &model.BooleanResult{
		Result: result,
	}, nil
}

// RemoveContactFromGroup is the resolver for the removeContactFromGroup field.
func (r *mutationResolver) RemoveContactFromGroup(ctx context.Context, contactID string, groupID string) (*model.BooleanResult, error) {
	result, err := r.ServiceContainer.ContactWithContactGroupRelationshipService.RemoveContactFromGroup(ctx, contactID, groupID)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not remove contact from group")
		return nil, err
	}
	return &model.BooleanResult{
		Result: result,
	}, nil
}

// MergeTextCustomFieldToContact is the resolver for the mergeTextCustomFieldToContact field.
func (r *mutationResolver) MergeTextCustomFieldToContact(ctx context.Context, contactID string, input model.TextCustomFieldInput) (*model.TextCustomField, error) {
	result, err := r.ServiceContainer.TextCustomFieldService.MergeTextCustomFieldToContact(ctx, contactID, mapper.MapTextCustomFieldInputToEntity(input))
	if err != nil {
		graphql.AddErrorf(ctx, "Could not add custom field %s to contact %s", input.Name, contactID)
		return nil, err
	}
	return mapper.MapEntityToTextCustomField(result), nil
}

// RemoveTextCustomFieldFromContact is the resolver for the removeTextCustomFieldFromContact field.
func (r *mutationResolver) RemoveTextCustomFieldFromContact(ctx context.Context, contactID string, fieldName string) (*model.BooleanResult, error) {
	result, err := r.ServiceContainer.TextCustomFieldService.Delete(ctx, contactID, fieldName)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not remove property %s from contact %s", fieldName, contactID)
		return nil, err
	}
	return &model.BooleanResult{
		Result: result,
	}, nil
}

// CreateContactGroup is the resolver for the createContactGroup field.
func (r *mutationResolver) CreateContactGroup(ctx context.Context, input model.ContactGroupInput) (*model.ContactGroup, error) {
	contactGroupEntityCreated, err := r.ServiceContainer.ContactGroupService.Create(ctx, &entity.ContactGroupEntity{
		Name: input.Name,
	})
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to create contact group %s", input.Name)
		return nil, err
	}

	return mapper.MapEntityToContactGroup(contactGroupEntityCreated), nil
}

// DeleteContactGroupAndUnlinkAllContacts is the resolver for the deleteContactGroupAndUnlinkAllContacts field.
func (r *mutationResolver) DeleteContactGroupAndUnlinkAllContacts(ctx context.Context, id string) (*model.BooleanResult, error) {
	result, err := r.ServiceContainer.ContactGroupService.Delete(ctx, id)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not delete contact group %s", id)
		return nil, err
	}
	return &model.BooleanResult{
		Result: result,
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
