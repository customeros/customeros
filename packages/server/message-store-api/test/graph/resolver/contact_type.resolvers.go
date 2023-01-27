package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/model"
)

// ContactTypeCreate is the resolver for the contactType_Create field.
func (r *mutationResolver) ContactTypeCreate(ctx context.Context, input model.ContactTypeInput) (*model.ContactType, error) {
	panic(fmt.Errorf("not implemented: ContactTypeCreate - contactType_Create"))
}

// ContactTypeUpdate is the resolver for the contactType_Update field.
func (r *mutationResolver) ContactTypeUpdate(ctx context.Context, input model.ContactTypeUpdateInput) (*model.ContactType, error) {
	panic(fmt.Errorf("not implemented: ContactTypeUpdate - contactType_Update"))
}

// ContactTypeDelete is the resolver for the contactType_Delete field.
func (r *mutationResolver) ContactTypeDelete(ctx context.Context, id string) (*model.Result, error) {
	panic(fmt.Errorf("not implemented: ContactTypeDelete - contactType_Delete"))
}

// ContactTypes is the resolver for the contactTypes field.
func (r *queryResolver) ContactTypes(ctx context.Context) ([]*model.ContactType, error) {
	panic(fmt.Errorf("not implemented: ContactTypes - contactTypes"))
}
