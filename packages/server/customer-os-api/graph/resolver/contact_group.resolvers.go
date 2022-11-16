package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/mapper"

	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

// ContactGroup is the resolver for the contactGroup field.
func (r *queryResolver) ContactGroup(ctx context.Context, id string) (*model.ContactGroup, error) {
	contactGroupEntity, err := r.ServiceContainer.ContactGroupService.FindContactGroupById(ctx, id)
	return mapper.MapEntityToContactGroup(contactGroupEntity), err
}

// ContactGroups is the resolver for the contactGroups field.
func (r *queryResolver) ContactGroups(ctx context.Context, paginationFilter *model.PaginationFilter) (*model.ContactGroupPage, error) {
	if paginationFilter == nil {
		paginationFilter = &model.PaginationFilter{Page: 0, Limit: 0}
	}
	paginatedResult, err := r.ServiceContainer.ContactGroupService.FindAll(ctx, paginationFilter.Page, paginationFilter.Limit)
	return &model.ContactGroupPage{
		Content:       mapper.MapEntitiesToContactGroups(paginatedResult.Rows.(*entity.ContactGroupEntities)),
		TotalPages:    paginatedResult.TotalPages,
		TotalElements: paginatedResult.TotalRows,
	}, err
}
