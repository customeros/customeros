package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"

	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/mapper"
)

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
		Content:       mapper.MapEntitiesToContacts(paginatedResult.Rows.(*entity.ContactNodes)),
		TotalPages:    paginatedResult.TotalPages,
		TotalElements: paginatedResult.TotalRows,
	}, err
}

// ContactGroups is the resolver for the contactGroups field.
func (r *queryResolver) ContactGroups(ctx context.Context) ([]*model.ContactGroup, error) {
	contactGroupEntities, err := r.ServiceContainer.ContactGroupService.FindAll(ctx)
	return mapper.MapEntitiesToContactGroups(contactGroupEntities), err
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
