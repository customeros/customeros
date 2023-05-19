package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/model"
)

// EntityTemplates is the resolver for the entityTemplates field.
func (r *queryResolver) EntityTemplates(ctx context.Context, extends *model.EntityTemplateExtension) ([]*model.EntityTemplate, error) {
	panic(fmt.Errorf("not implemented: EntityTemplates - entityTemplates"))
}

// DashboardViewContacts is the resolver for the dashboardView_Contacts field.
func (r *queryResolver) DashboardViewContacts(ctx context.Context, pagination model.Pagination, where *model.Filter) (*model.ContactsPage, error) {
	panic(fmt.Errorf("not implemented: DashboardViewContacts - dashboardView_Contacts"))
}

// DashboardViewOrganizations is the resolver for the dashboardView_Organizations field.
func (r *queryResolver) DashboardViewOrganizations(ctx context.Context, pagination model.Pagination, where *model.Filter) (*model.OrganizationPage, error) {
	panic(fmt.Errorf("not implemented: DashboardViewOrganizations - dashboardView_Organizations"))
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) DashboardView(ctx context.Context, pagination model.Pagination, searchTerm *string) (*model.DashboardViewItemPage, error) {
	panic(fmt.Errorf("not implemented: DashboardView - dashboardView"))
}
