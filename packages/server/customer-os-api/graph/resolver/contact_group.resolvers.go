package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
)

// Contacts is the resolver for the contacts field.
func (r *contactGroupResolver) Contacts(ctx context.Context, obj *model.ContactGroup, pagination *model.Pagination, where *model.Filter, sort []*model.SortBy) (*model.ContactsPage, error) {
	if pagination == nil {
		pagination = &model.Pagination{Page: 0, Limit: 0}
	}
	paginatedResult, err := r.ServiceContainer.ContactService.FindAllForContactGroup(ctx, pagination.Page, pagination.Limit, where, sort, obj.ID)
	return &model.ContactsPage{
		Content:       mapper.MapEntitiesToContacts(paginatedResult.Rows.(*entity.ContactEntities)),
		TotalPages:    paginatedResult.TotalPages,
		TotalElements: paginatedResult.TotalRows,
	}, err
}

// ContactGroupCreate is the resolver for the contactGroupCreate field.
func (r *mutationResolver) ContactGroupCreate(ctx context.Context, input model.ContactGroupInput) (*model.ContactGroup, error) {
	contactGroupEntityCreated, err := r.ServiceContainer.ContactGroupService.Create(ctx, &entity.ContactGroupEntity{
		Name: input.Name,
	})
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to create contact group %s", input.Name)
		return nil, err
	}
	return mapper.MapEntityToContactGroup(contactGroupEntityCreated), nil
}

// ContactGroupUpdate is the resolver for the contactGroupUpdate field.
func (r *mutationResolver) ContactGroupUpdate(ctx context.Context, input model.ContactGroupUpdateInput) (*model.ContactGroup, error) {
	updatedContactGroup, err := r.ServiceContainer.ContactGroupService.Update(ctx, &entity.ContactGroupEntity{
		Id:   input.ID,
		Name: input.Name,
	})
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to update contact group %s", input.ID)
		return nil, err
	}
	return mapper.MapEntityToContactGroup(updatedContactGroup), nil
}

// ContactGroupDeleteAndUnlinkAllContacts is the resolver for the contactGroupDeleteAndUnlinkAllContacts field.
func (r *mutationResolver) ContactGroupDeleteAndUnlinkAllContacts(ctx context.Context, id string) (*model.Result, error) {
	result, err := r.ServiceContainer.ContactGroupService.Delete(ctx, id)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not delete contact group %s", id)
		return nil, err
	}
	return &model.Result{
		Result: result,
	}, nil
}

// ContactGroupAddContact is the resolver for the contactGroupAddContact field.
func (r *mutationResolver) ContactGroupAddContact(ctx context.Context, contactID string, groupID string) (*model.Result, error) {
	result, err := r.ServiceContainer.ContactGroupService.AddContactToGroup(ctx, contactID, groupID)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not add contact to group")
		return nil, err
	}
	return &model.Result{
		Result: result,
	}, nil
}

// ContactGroupRemoveContact is the resolver for the contactGroupRemoveContact field.
func (r *mutationResolver) ContactGroupRemoveContact(ctx context.Context, contactID string, groupID string) (*model.Result, error) {
	result, err := r.ServiceContainer.ContactGroupService.RemoveContactFromGroup(ctx, contactID, groupID)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not remove contact from group")
		return nil, err
	}
	return &model.Result{
		Result: result,
	}, nil
}

// ContactGroup is the resolver for the contactGroup field.
func (r *queryResolver) ContactGroup(ctx context.Context, id string) (*model.ContactGroup, error) {
	contactGroupEntity, err := r.ServiceContainer.ContactGroupService.FindContactGroupById(ctx, id)
	return mapper.MapEntityToContactGroup(contactGroupEntity), err
}

// ContactGroups is the resolver for the contactGroups field.
func (r *queryResolver) ContactGroups(ctx context.Context, pagination *model.Pagination, where *model.Filter, sort []*model.SortBy) (*model.ContactGroupPage, error) {
	if pagination == nil {
		pagination = &model.Pagination{Page: 0, Limit: 0}
	}
	paginatedResult, err := r.ServiceContainer.ContactGroupService.FindAll(ctx, pagination.Page, pagination.Limit, where, sort)
	if err != nil {
		graphql.AddErrorf(ctx, "Could not fetch contact groups")
		return nil, err
	}
	return &model.ContactGroupPage{
		Content:       mapper.MapEntitiesToContactGroups(paginatedResult.Rows.(*entity.ContactGroupEntities)),
		TotalPages:    paginatedResult.TotalPages,
		TotalElements: paginatedResult.TotalRows,
	}, err
}

// ContactGroup returns generated.ContactGroupResolver implementation.
func (r *Resolver) ContactGroup() generated.ContactGroupResolver { return &contactGroupResolver{r} }

type contactGroupResolver struct{ *Resolver }
