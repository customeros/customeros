package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/model"
)

// ContactUpsertInEventStore is the resolver for the contactUpsertInEventStore field.
func (r *mutationResolver) ContactUpsertInEventStore(ctx context.Context, size int) (int, error) {
	panic(fmt.Errorf("not implemented: ContactUpsertInEventStore - contactUpsertInEventStore"))
}

// ContactPhoneNumberRelationUpsertInEventStore is the resolver for the contactPhoneNumberRelationUpsertInEventStore field.
func (r *mutationResolver) ContactPhoneNumberRelationUpsertInEventStore(ctx context.Context, size int) (int, error) {
	panic(fmt.Errorf("not implemented: ContactPhoneNumberRelationUpsertInEventStore - contactPhoneNumberRelationUpsertInEventStore"))
}

// UpsertInEventStore is the resolver for the UpsertInEventStore field.
func (r *mutationResolver) UpsertInEventStore(ctx context.Context, size int) (*model.UpsertToEventStoreResult, error) {
	panic(fmt.Errorf("not implemented: UpsertInEventStore - UpsertInEventStore"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
