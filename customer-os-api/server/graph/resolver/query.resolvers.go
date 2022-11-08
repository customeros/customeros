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

// TenantUsers is the resolver for the tenantUsers field.
func (r *queryResolver) TenantUsers(ctx context.Context, paginationFilter *model.PaginationFilter) (*model.TenantUsersPage, error) {
	if paginationFilter == nil {
		paginationFilter = &model.PaginationFilter{Page: 0, Limit: 0}
	}
	paginatedResult, err := r.ServiceContainer.TenantUserService.FindAll(ctx, paginationFilter.Page, paginationFilter.Limit)
	return &model.TenantUsersPage{
		Content:       mapper.MapEntitiesToTenantUsers(paginatedResult.Rows.(*entity.TenantUserEntities)),
		TotalPages:    paginatedResult.TotalPages,
		TotalElements: paginatedResult.TotalRows,
	}, err
}

// Contact is the resolver for the contact field.
func (r *queryResolver) Contact(ctx context.Context, id string) (*model.Contact, error) {
	contactEntity, err := r.ServiceContainer.ContactService.FindContactById(ctx, id)
	if err != nil || contactEntity == nil {
		graphql.AddErrorf(ctx, "Contact with id %s not found", id)
		return nil, err
	}
	return mapper.MapEntityToContact(contactEntity), nil
}

// Contacts is the resolver for the contacts field.
func (r *queryResolver) Contacts(ctx context.Context, paginationFilter *model.PaginationFilter) (*model.ContactsPage, error) {
	if paginationFilter == nil {
		paginationFilter = &model.PaginationFilter{Page: 0, Limit: 0}
	}
	paginatedResult, err := r.ServiceContainer.ContactService.FindAll(ctx, paginationFilter.Page, paginationFilter.Limit)
	return &model.ContactsPage{
		Content:       mapper.MapEntitiesToContacts(paginatedResult.Rows.(*entity.ContactEntities)),
		TotalPages:    paginatedResult.TotalPages,
		TotalElements: paginatedResult.TotalRows,
	}, err
}

// ContactByEmail is the resolver for the contactByEmail field.
func (r *queryResolver) ContactByEmail(ctx context.Context, email string) (*model.Contact, error) {
	contactEntity, err := r.ServiceContainer.ContactService.FindContactByEmail(ctx, email)
	if err != nil || contactEntity == nil {
		graphql.AddErrorf(ctx, "Contact with email %s not identified", email)
		return nil, err
	}
	return mapper.MapEntityToContact(contactEntity), nil
}

// ContactByPhone is the resolver for the contactByPhone field.
func (r *queryResolver) ContactByPhone(ctx context.Context, number string) (*model.Contact, error) {
	contactEntity, err := r.ServiceContainer.ContactService.FindContactByPhoneNumber(ctx, number)
	if err != nil || contactEntity == nil {
		graphql.AddErrorf(ctx, "Contact with phone number %s not identified", number)
		return nil, err
	}
	return mapper.MapEntityToContact(contactEntity), nil
}

// ContactGroup is the resolver for the contactGroup field.
func (r *queryResolver) ContactGroup(ctx context.Context, id string) (*model.ContactGroup, error) {
	contactGroupEntity, err := r.ServiceContainer.ContactGroupService.FindContactGroupById(ctx, id)
	return mapper.MapEntityToContactGroup(contactGroupEntity), err
}

// ContactGroups is the resolver for the contactGroups field.
func (r *queryResolver) ContactGroups(ctx context.Context, paginationFilter *model.PaginationFilter) (*model.ContactGroupsPage, error) {
	if paginationFilter == nil {
		paginationFilter = &model.PaginationFilter{Page: 0, Limit: 0}
	}
	paginatedResult, err := r.ServiceContainer.ContactGroupService.FindAll(ctx, paginationFilter.Page, paginationFilter.Limit)
	return &model.ContactGroupsPage{
		Content:       mapper.MapEntitiesToContactGroups(paginatedResult.Rows.(*entity.ContactGroupEntities)),
		TotalPages:    paginatedResult.TotalPages,
		TotalElements: paginatedResult.TotalRows,
	}, err
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
