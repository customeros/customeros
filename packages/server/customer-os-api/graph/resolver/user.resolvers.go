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
func (r *queryResolver) Users(ctx context.Context, pagination *model.Pagination, where *model.Filter, sort []*model.SortBy) (*model.UserPage, error) {
	if pagination == nil {
		pagination = &model.Pagination{Page: 0, Limit: 0}
	}
	paginatedResult, err := r.ServiceContainer.UserService.FindAll(ctx, pagination.Page, pagination.Limit, where, sort)
	return &model.UserPage{
		Content:       mapper.MapEntitiesToUsers(paginatedResult.Rows.(*entity.UserEntities)),
		TotalPages:    paginatedResult.TotalPages,
		TotalElements: paginatedResult.TotalRows,
	}, err
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	userEntity, err := r.ServiceContainer.UserService.FindUserById(ctx, id)
	if err != nil || userEntity == nil {
		graphql.AddErrorf(ctx, "User with id %s not found", id)
		return nil, err
	}
	return mapper.MapEntityToUser(userEntity), nil
}

// Conversations is the resolver for the conversations field.
func (r *userResolver) Conversations(ctx context.Context, obj *model.User, pagination *model.Pagination, sort []*model.SortBy) (*model.ConversationPage, error) {
	if pagination == nil {
		pagination = &model.Pagination{Page: 0, Limit: 0}
	}
	paginatedResult, err := r.ServiceContainer.ConversationService.GetConversationsForUser(ctx, obj.ID, pagination.Page, pagination.Limit, sort)
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to get user %s conversations", obj.ID)
		return nil, err
	}
	return &model.ConversationPage{
		Content:       mapper.MapEntitiesToConversations(paginatedResult.Rows.(*entity.ConversationEntities)),
		TotalPages:    paginatedResult.TotalPages,
		TotalElements: paginatedResult.TotalRows,
	}, err
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
