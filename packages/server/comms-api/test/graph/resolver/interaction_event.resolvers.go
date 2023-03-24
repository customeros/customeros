package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/model"
)

// InteractionSession is the resolver for the interactionSession field.
func (r *interactionEventResolver) InteractionSession(ctx context.Context, obj *model.InteractionEvent) (*model.InteractionSession, error) {
	panic(fmt.Errorf("not implemented: InteractionSession - interactionSession"))
}

// SentBy is the resolver for the sentBy field.
func (r *interactionEventResolver) SentBy(ctx context.Context, obj *model.InteractionEvent) ([]model.InteractionEventParticipant, error) {
	if r.Resolver.SentBy != nil {
		return r.Resolver.SentBy(ctx, obj)
	}
	panic(fmt.Errorf("not implemented: SentBy - sentBy"))
}

// SentTo is the resolver for the sentTo field.
func (r *interactionEventResolver) SentTo(ctx context.Context, obj *model.InteractionEvent) ([]model.InteractionEventParticipant, error) {
	if r.Resolver.SentTo != nil {
		return r.Resolver.SentTo(ctx, obj)
	}
	panic(fmt.Errorf("not implemented: SentTo - sentTo"))
}

// RepliesTo is the resolver for the repliesTo field.
func (r *interactionEventResolver) RepliesTo(ctx context.Context, obj *model.InteractionEvent) (*model.InteractionEvent, error) {
	if r.Resolver.RepliesTo != nil {
		return r.Resolver.RepliesTo(ctx, obj)
	}
	panic(fmt.Errorf("not implemented: RepliesTo - repliesTo"))
}

// Events is the resolver for the events field.
func (r *interactionSessionResolver) Events(ctx context.Context, obj *model.InteractionSession) ([]*model.InteractionEvent, error) {
	if r.Resolver.Events != nil {
		return r.Resolver.Events(ctx, obj)
	}
	panic(fmt.Errorf("not implemented: Events - events"))
}

// AttendedBy is the resolver for the attendedBy field.
func (r *interactionSessionResolver) AttendedBy(ctx context.Context, obj *model.InteractionSession) ([]model.InteractionSessionParticipant, error) {
	if r.Resolver.AttendedBy != nil {
		return r.Resolver.AttendedBy(ctx, obj)
	}
	panic(fmt.Errorf("not implemented: AttendedBy - attendedBy"))
}

// InteractionSessionCreate is the resolver for the interactionSession_Create field.
func (r *mutationResolver) InteractionSessionCreate(ctx context.Context, session model.InteractionSessionInput) (*model.InteractionSession, error) {
	if r.Resolver.InteractionSessionCreate != nil {
		return r.Resolver.InteractionSessionCreate(ctx, session)
	}
	panic(fmt.Errorf("not implemented: InteractionSessionCreate - interactionSession_Create"))
}

// InteractionEventCreate is the resolver for the interactionEvent_Create field.
func (r *mutationResolver) InteractionEventCreate(ctx context.Context, event model.InteractionEventInput) (*model.InteractionEvent, error) {
	if r.Resolver.InteractionEventCreate != nil {
		return r.Resolver.InteractionEventCreate(ctx, event)
	}
	panic(fmt.Errorf("not implemented: InteractionEventCreate - interactionEvent_Create"))
}

// InteractionSession is the resolver for the interactionSession field.
func (r *queryResolver) InteractionSession(ctx context.Context, id string) (*model.InteractionSession, error) {
	panic(fmt.Errorf("not implemented: InteractionSession - interactionSession"))
}

// InteractionSessionBySessionIdentifier is the resolver for the interactionSession_BySessionIdentifier field.
func (r *queryResolver) InteractionSessionBySessionIdentifier(ctx context.Context, sessionIdentifier string) (*model.InteractionSession, error) {
	if r.Resolver.InteractionSessionBySessionIdentifier != nil {
		return r.Resolver.InteractionSessionBySessionIdentifier(ctx, sessionIdentifier)
	}
	panic(fmt.Errorf("not implemented: InteractionSessionBySessionIdentifier - interactionSession_BySessionIdentifier"))
}

// InteractionEvent is the resolver for the interactionEvent field.
func (r *queryResolver) InteractionEvent(ctx context.Context, id string) (*model.InteractionEvent, error) {
	panic(fmt.Errorf("not implemented: InteractionEvent - interactionEvent"))
}

// InteractionEventByEventIdentifier is the resolver for the interactionEvent_ByEventIdentifier field.
func (r *queryResolver) InteractionEventByEventIdentifier(ctx context.Context, eventIdentifier string) (*model.InteractionEvent, error) {
	panic(fmt.Errorf("not implemented: InteractionEventByEventIdentifier - interactionEvent_ByEventIdentifier"))
}

// InteractionEvent returns generated.InteractionEventResolver implementation.
func (r *Resolver) InteractionEvent() generated.InteractionEventResolver {
	return &interactionEventResolver{r}
}

// InteractionSession returns generated.InteractionSessionResolver implementation.
func (r *Resolver) InteractionSession() generated.InteractionSessionResolver {
	return &interactionSessionResolver{r}
}

type interactionEventResolver struct{ *Resolver }
type interactionSessionResolver struct{ *Resolver }
