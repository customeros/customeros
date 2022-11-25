package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/message-store/test/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/test/graph/model"
)

// EntityDefinitions is the resolver for the entityDefinitions field.
func (r *queryResolver) EntityDefinitions(ctx context.Context, extends *model.EntityDefinitionExtension) ([]*model.EntityDefinition, error) {
	panic(fmt.Errorf("not implemented: EntityDefinitions - entityDefinitions"))
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
