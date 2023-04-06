package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/generated"
)

// PhoneNumberUpsertInEventStore is the resolver for the phoneNumberUpsertInEventStore field.
func (r *mutationResolver) PhoneNumberUpsertInEventStore(ctx context.Context, size int) (int, error) {
	panic(fmt.Errorf("not implemented: PhoneNumberUpsertInEventStore - phoneNumberUpsertInEventStore"))
}

// ContactUpsertInEventStore is the resolver for the contactUpsertInEventStore field.
func (r *mutationResolver) ContactUpsertInEventStore(ctx context.Context, size int) (int, error) {
	panic(fmt.Errorf("not implemented: ContactUpsertInEventStore - contactUpsertInEventStore"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
