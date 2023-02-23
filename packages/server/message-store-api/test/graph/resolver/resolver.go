package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ContactCreate         func(ctx context.Context, input model.ContactInput) (*model.Contact, error)
	ConversationCreate    func(ctx context.Context, input model.ConversationInput) (*model.Conversation, error)
	GetContactByEmail     func(ctx context.Context, email string) (*model.Contact, error)
	GetContactByPhone     func(ctx context.Context, e164 string) (*model.Contact, error)
	GetContactById        func(ctx context.Context, id string) (*model.Contact, error)
	PhoneNumbersByContact func(ctx context.Context, obj *model.Contact) ([]*model.PhoneNumber, error)
	EmailsByContact       func(ctx context.Context, obj *model.Contact) ([]*model.Email, error)
	UserByEmail           func(ctx context.Context, email string) (*model.User, error)
	EmailsByUser          func(ctx context.Context, obj *model.User) ([]*model.Email, error)
}
