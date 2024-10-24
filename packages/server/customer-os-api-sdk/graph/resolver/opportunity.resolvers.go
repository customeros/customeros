package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api-sdk/graph/model"
)

// OpportunitySave is the resolver for the opportunity_save field.
func (r *mutationResolver) OpportunitySave(ctx context.Context, input model.OpportunitySaveInput) (*model.Opportunity, error) {
	panic(fmt.Errorf("not implemented: OpportunitySave - opportunity_save"))
}

// OpportunityArchive is the resolver for the opportunity_Archive field.
func (r *mutationResolver) OpportunityArchive(ctx context.Context, id string) (*model.ActionResponse, error) {
	panic(fmt.Errorf("not implemented: OpportunityArchive - opportunity_Archive"))
}

// OpportunityRenewalUpdate is the resolver for the opportunityRenewalUpdate field.
func (r *mutationResolver) OpportunityRenewalUpdate(ctx context.Context, input model.OpportunityRenewalUpdateInput, ownerUserID *string) (*model.Opportunity, error) {
	panic(fmt.Errorf("not implemented: OpportunityRenewalUpdate - opportunityRenewalUpdate"))
}

// OpportunityRenewalUpdateAllForOrganization is the resolver for the opportunityRenewal_UpdateAllForOrganization field.
func (r *mutationResolver) OpportunityRenewalUpdateAllForOrganization(ctx context.Context, input model.OpportunityRenewalUpdateAllForOrganizationInput) (*model.Organization, error) {
	panic(fmt.Errorf("not implemented: OpportunityRenewalUpdateAllForOrganization - opportunityRenewal_UpdateAllForOrganization"))
}

// OpportunityCreate is the resolver for the opportunity_Create field.
func (r *mutationResolver) OpportunityCreate(ctx context.Context, input model.OpportunityCreateInput) (*model.Opportunity, error) {
	panic(fmt.Errorf("not implemented: OpportunityCreate - opportunity_Create"))
}

// OpportunityUpdate is the resolver for the opportunityUpdate field.
func (r *mutationResolver) OpportunityUpdate(ctx context.Context, input model.OpportunityUpdateInput) (*model.Opportunity, error) {
	panic(fmt.Errorf("not implemented: OpportunityUpdate - opportunityUpdate"))
}

// OpportunitySetOwner is the resolver for the opportunity_SetOwner field.
func (r *mutationResolver) OpportunitySetOwner(ctx context.Context, opportunityID string, userID string) (*model.ActionResponse, error) {
	panic(fmt.Errorf("not implemented: OpportunitySetOwner - opportunity_SetOwner"))
}

// OpportunityRemoveOwner is the resolver for the opportunity_RemoveOwner field.
func (r *mutationResolver) OpportunityRemoveOwner(ctx context.Context, opportunityID string) (*model.ActionResponse, error) {
	panic(fmt.Errorf("not implemented: OpportunityRemoveOwner - opportunity_RemoveOwner"))
}

// Organization is the resolver for the organization field.
func (r *opportunityResolver) Organization(ctx context.Context, obj *model.Opportunity) (*model.Organization, error) {
	panic(fmt.Errorf("not implemented: Organization - organization"))
}

// CreatedBy is the resolver for the createdBy field.
func (r *opportunityResolver) CreatedBy(ctx context.Context, obj *model.Opportunity) (*model.User, error) {
	panic(fmt.Errorf("not implemented: CreatedBy - createdBy"))
}

// Owner is the resolver for the owner field.
func (r *opportunityResolver) Owner(ctx context.Context, obj *model.Opportunity) (*model.User, error) {
	panic(fmt.Errorf("not implemented: Owner - owner"))
}

// ExternalLinks is the resolver for the externalLinks field.
func (r *opportunityResolver) ExternalLinks(ctx context.Context, obj *model.Opportunity) ([]*model.ExternalSystem, error) {
	panic(fmt.Errorf("not implemented: ExternalLinks - externalLinks"))
}

// Opportunity is the resolver for the opportunity field.
func (r *queryResolver) Opportunity(ctx context.Context, id string) (*model.Opportunity, error) {
	panic(fmt.Errorf("not implemented: Opportunity - opportunity"))
}

// OpportunitiesLinkedToOrganizations is the resolver for the opportunities_LinkedToOrganizations field.
func (r *queryResolver) OpportunitiesLinkedToOrganizations(ctx context.Context, pagination *model.Pagination) (*model.OpportunityPage, error) {
	panic(fmt.Errorf("not implemented: OpportunitiesLinkedToOrganizations - opportunities_LinkedToOrganizations"))
}

// Opportunity returns generated.OpportunityResolver implementation.
func (r *Resolver) Opportunity() generated.OpportunityResolver { return &opportunityResolver{r} }

type opportunityResolver struct{ *Resolver }
