package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/model"
)

// Place is the resolver for the place field.
func (r *locationResolver) Place(ctx context.Context, obj *model.Location) (*model.Place, error) {
	panic(fmt.Errorf("not implemented: Place - place"))
}

// Location returns generated.LocationResolver implementation.
func (r *Resolver) Location() generated.LocationResolver { return &locationResolver{r} }

type locationResolver struct{ *Resolver }
