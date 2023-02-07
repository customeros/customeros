package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/model"
)

// SearchBasic is the resolver for the search_Basic field.
func (r *queryResolver) SearchBasic(ctx context.Context, keyword string) ([]*model.SearchBasicResultItem, error) {
	panic(fmt.Errorf("not implemented: SearchBasic - search_Basic"))
}
