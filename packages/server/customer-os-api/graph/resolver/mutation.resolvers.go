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

// EntityDefinitionCreate is the resolver for the entityDefinitionCreate field.
func (r *mutationResolver) EntityDefinitionCreate(ctx context.Context, input model.EntityDefinitionInput) (*model.EntityDefinition, error) {
	entityDefinitionEntity, err := r.ServiceContainer.EntityDefinitionService.Create(ctx, mapper.MapEntityDefinitionInputToEntity(input))
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to create entity definition: %s", input.Name)
		return nil, err
	}
	return mapper.MapEntityToEntityDefinition(entityDefinitionEntity), nil
}

// ConversationCreate is the resolver for the conversationCreate field.
func (r *mutationResolver) ConversationCreate(ctx context.Context, input model.ConversationInput) (*model.Conversation, error) {
	conversationEntity, err := r.ServiceContainer.ConversationService.CreateNewConversation(ctx, input.UserID, input.ContactID, input.ID)
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to create conversation between user: %s and contact: %s", input.UserID, input.ContactID)
		return nil, err
	}
	return mapper.MapEntityToConversation(conversationEntity), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
