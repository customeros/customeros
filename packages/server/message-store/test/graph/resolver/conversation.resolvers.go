package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/message-store/test/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/test/graph/model"
)

// Contact is the resolver for the contact field.
func (r *conversationResolver) Contact(ctx context.Context, obj *model.Conversation) (*model.Contact, error) {
	panic(fmt.Errorf("not implemented: Contact - contact"))
}

// User is the resolver for the user field.
func (r *conversationResolver) User(ctx context.Context, obj *model.Conversation) (*model.User, error) {
	panic(fmt.Errorf("not implemented: User - user"))
}

// ConversationCreate is the resolver for the conversationCreate field.
func (r *mutationResolver) ConversationCreate(ctx context.Context, input model.ConversationInput) (*model.Conversation, error) {
	if r.Resolver.ConversationCreate != nil {
		return r.Resolver.ConversationCreate(ctx, input)
	}
	panic(fmt.Errorf("not implemented: ConversationCreate - conversationCreate"))
}

// ConversationAddMessage is the resolver for the conversationAddMessage field.
func (r *mutationResolver) ConversationAddMessage(ctx context.Context, input model.MessageInput) (*model.Message, error) {
	panic(fmt.Errorf("not implemented: ConversationAddMessage - conversationAddMessage"))
}

// Conversation returns generated.ConversationResolver implementation.
func (r *Resolver) Conversation() generated.ConversationResolver { return &conversationResolver{r} }

type conversationResolver struct{ *Resolver }
