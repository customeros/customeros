package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/message-store/test/graph/model"
)

// OrganizationTypeCreate is the resolver for the organizationType_Create field.
func (r *mutationResolver) OrganizationTypeCreate(ctx context.Context, input model.OrganizationTypeInput) (*model.OrganizationType, error) {
	panic(fmt.Errorf("not implemented: OrganizationTypeCreate - organizationType_Create"))
}

// OrganizationTypeUpdate is the resolver for the organizationType_Update field.
func (r *mutationResolver) OrganizationTypeUpdate(ctx context.Context, input model.OrganizationTypeUpdateInput) (*model.OrganizationType, error) {
	panic(fmt.Errorf("not implemented: OrganizationTypeUpdate - organizationType_Update"))
}

// OrganizationTypeDelete is the resolver for the organizationType_Delete field.
func (r *mutationResolver) OrganizationTypeDelete(ctx context.Context, id string) (*model.Result, error) {
	panic(fmt.Errorf("not implemented: OrganizationTypeDelete - organizationType_Delete"))
}

// OrganizationTypes is the resolver for the organizationTypes field.
func (r *queryResolver) OrganizationTypes(ctx context.Context) ([]*model.OrganizationType, error) {
	panic(fmt.Errorf("not implemented: OrganizationTypes - organizationTypes"))
}
