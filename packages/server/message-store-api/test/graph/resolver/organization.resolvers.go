package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/model"
)

// OrganizationCreate is the resolver for the organization_Create field.
func (r *mutationResolver) OrganizationCreate(ctx context.Context, input model.OrganizationInput) (*model.Organization, error) {
	panic(fmt.Errorf("not implemented: OrganizationCreate - organization_Create"))
}

// OrganizationUpdate is the resolver for the organization_Update field.
func (r *mutationResolver) OrganizationUpdate(ctx context.Context, id string, input model.OrganizationInput) (*model.Organization, error) {
	panic(fmt.Errorf("not implemented: OrganizationUpdate - organization_Update"))
}

// OrganizationDelete is the resolver for the organization_Delete field.
func (r *mutationResolver) OrganizationDelete(ctx context.Context, id string) (*model.Result, error) {
	panic(fmt.Errorf("not implemented: OrganizationDelete - organization_Delete"))
}

// OrganizationType is the resolver for the organizationType field.
func (r *organizationResolver) OrganizationType(ctx context.Context, obj *model.Organization) (*model.OrganizationType, error) {
	panic(fmt.Errorf("not implemented: OrganizationType - organizationType"))
}

// Addresses is the resolver for the addresses field.
func (r *organizationResolver) Addresses(ctx context.Context, obj *model.Organization) ([]*model.Place, error) {
	panic(fmt.Errorf("not implemented: Addresses - addresses"))
}

// ContactRoles is the resolver for the contactRoles field.
func (r *organizationResolver) ContactRoles(ctx context.Context, obj *model.Organization) ([]*model.ContactRole, error) {
	panic(fmt.Errorf("not implemented: ContactRoles - contactRoles"))
}

// Organizations is the resolver for the organizations field.
func (r *queryResolver) Organizations(ctx context.Context, pagination *model.Pagination, where *model.Filter, sort []*model.SortBy) (*model.OrganizationPage, error) {
	panic(fmt.Errorf("not implemented: Organizations - organizations"))
}

// Organization is the resolver for the organization field.
func (r *queryResolver) Organization(ctx context.Context, id string) (*model.Organization, error) {
	panic(fmt.Errorf("not implemented: Organization - organization"))
}

// Organization returns generated.OrganizationResolver implementation.
func (r *Resolver) Organization() generated.OrganizationResolver { return &organizationResolver{r} }

type organizationResolver struct{ *Resolver }
