package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/model"
)

// GcliSearch is the resolver for the gcli_Search field.
func (r *queryResolver) GcliSearch(ctx context.Context, keyword string, limit *int) ([]*model.GCliSearchResultItem, error) {
	panic(fmt.Errorf("not implemented: GcliSearch - gcli_Search"))
}
