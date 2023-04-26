package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/test/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/test/graph/model"
)

// EntityTemplates is the resolver for the entityTemplates field.
func (r *queryResolver) EntityTemplates(ctx context.Context, extends *model.EntityTemplateExtension) ([]*model.EntityTemplate, error) {
	panic(fmt.Errorf("not implemented: EntityTemplates - entityTemplates"))
}

// DashboardView is the resolver for the dashboardView field.
func (r *queryResolver) DashboardView(ctx context.Context, pagination model.Pagination, searchTerm *string) (*model.DashboardViewItemPage, error) {
	panic(fmt.Errorf("not implemented: DashboardView - dashboardView"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
