package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ContactCreate                         func(ctx context.Context, input model.ContactInput) (*model.Contact, error)
	InteractionEventCreate                func(ctx context.Context, event model.InteractionEventInput) (*model.InteractionEvent, error)
	InteractionSessionBySessionIdentifier func(ctx context.Context, sessionIdentifier string) (*model.InteractionSession, error)
	InteractionSessionCreate              func(ctx context.Context, session model.InteractionSessionInput) (*model.InteractionSession, error)
	InteractionSessionResolver            func(ctx context.Context, obj *model.InteractionEvent) (*model.InteractionSession, error)
	AnalysisCreate                        func(ctx context.Context, analysis model.AnalysisInput) (*model.Analysis, error)
	SentBy                                func(ctx context.Context, obj *model.InteractionEvent) ([]model.InteractionEventParticipant, error)
	SentTo                                func(ctx context.Context, obj *model.InteractionEvent) ([]model.InteractionEventParticipant, error)
	RepliesTo                             func(ctx context.Context, obj *model.InteractionEvent) (*model.InteractionEvent, error)
	Events                                func(ctx context.Context, obj *model.InteractionSession) ([]*model.InteractionEvent, error)
	AttendedBy                            func(ctx context.Context, obj *model.InteractionSession) ([]model.InteractionSessionParticipant, error)
}
