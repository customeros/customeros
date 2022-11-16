package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/mapper"
)

// UserCreate is the resolver for the userCreate field.
func (r *mutationResolver) UserCreate(ctx context.Context, input model.UserInput) (*model.User, error) {
	createdTenantEntity, err := r.ServiceContainer.UserService.Create(ctx, mapper.MapUserInputToEntity(input))
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to create user %s %s", input.FirstName, input.LastName)
		return nil, err
	}
	return mapper.MapEntityToUser(createdTenantEntity), nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, paginationFilter *model.PaginationFilter) (*model.UserPage, error) {
	if paginationFilter == nil {
		paginationFilter = &model.PaginationFilter{Page: 0, Limit: 0}
	}
	paginatedResult, err := r.ServiceContainer.UserService.FindAll(ctx, paginationFilter.Page, paginationFilter.Limit)
	return &model.UserPage{
		Content:       mapper.MapEntitiesToUsers(paginatedResult.Rows.(*entity.UserEntities)),
		TotalPages:    paginatedResult.TotalPages,
		TotalElements: paginatedResult.TotalRows,
	}, err
}
