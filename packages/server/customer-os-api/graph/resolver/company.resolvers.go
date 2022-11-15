package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

// CompaniesByNameLike is the resolver for the companies_ByNameLike field.
func (r *queryResolver) CompaniesByNameLike(ctx context.Context, paginationFilter *model.PaginationFilter, companyName string) (*model.CompaniesPage, error) {
	panic(fmt.Errorf("not implemented: CompaniesByNameLike - companies_ByNameLike"))
}
