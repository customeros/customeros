package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/model"
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
