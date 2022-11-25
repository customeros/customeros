package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/message-store/test/graph/model"
)

// EmailMergeToContact is the resolver for the emailMergeToContact field.
func (r *mutationResolver) EmailMergeToContact(ctx context.Context, contactID string, input model.EmailInput) (*model.Email, error) {
	panic(fmt.Errorf("not implemented: EmailMergeToContact - emailMergeToContact"))
}

// EmailUpdateInContact is the resolver for the emailUpdateInContact field.
func (r *mutationResolver) EmailUpdateInContact(ctx context.Context, contactID string, input model.EmailUpdateInput) (*model.Email, error) {
	panic(fmt.Errorf("not implemented: EmailUpdateInContact - emailUpdateInContact"))
}

// EmailRemoveFromContact is the resolver for the emailRemoveFromContact field.
func (r *mutationResolver) EmailRemoveFromContact(ctx context.Context, contactID string, email string) (*model.Result, error) {
	panic(fmt.Errorf("not implemented: EmailRemoveFromContact - emailRemoveFromContact"))
}

// EmailRemoveFromContactByID is the resolver for the emailRemoveFromContactById field.
func (r *mutationResolver) EmailRemoveFromContactByID(ctx context.Context, contactID string, id string) (*model.Result, error) {
	panic(fmt.Errorf("not implemented: EmailRemoveFromContactByID - emailRemoveFromContactById"))
}
