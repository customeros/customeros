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

// Contact is the resolver for the contact field.
func (r *conversationResolver) Contact(ctx context.Context, obj *model.Conversation) (*model.Contact, error) {
	contactEntity, err := r.ServiceContainer.ContactService.FindContactById(ctx, obj.ContactID)
	if err != nil || contactEntity == nil {
		graphql.AddErrorf(ctx, "Contact with id %s not found", obj.ContactID)
		return nil, err
	}
	return mapper.MapEntityToContact(contactEntity), nil
}

// User is the resolver for the user field.
func (r *conversationResolver) User(ctx context.Context, obj *model.Conversation) (*model.User, error) {
	userEntity, err := r.ServiceContainer.UserService.FindUserById(ctx, obj.UserID)
	if err != nil || userEntity == nil {
		graphql.AddErrorf(ctx, "User with id %s not found", obj.UserID)
		return nil, err
	}
	return mapper.MapEntityToUser(userEntity), nil
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

// Conversation returns generated.ConversationResolver implementation.
func (r *Resolver) Conversation() generated.ConversationResolver { return &conversationResolver{r} }

type conversationResolver struct{ *Resolver }
