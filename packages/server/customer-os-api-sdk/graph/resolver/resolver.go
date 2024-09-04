package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Attachment    func(ctx context.Context, id string) (*model.Attachment, error)
	ContactCreate func(ctx context.Context, input model.ContactInput) (string, error)
	SentBy        func(ctx context.Context, obj *model.InteractionEvent) ([]model.InteractionEventParticipant, error)
	SentTo        func(ctx context.Context, obj *model.InteractionEvent) ([]model.InteractionEventParticipant, error)
	RepliesTo     func(ctx context.Context, obj *model.InteractionEvent) (*model.InteractionEvent, error)
	Events        func(ctx context.Context, obj *model.InteractionSession) ([]*model.InteractionEvent, error)
	AttendedBy    func(ctx context.Context, obj *model.InteractionSession) ([]model.InteractionSessionParticipant, error)
}
