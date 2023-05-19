package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/model"
)

// LocationUpdate is the resolver for the location_Update field.
func (r *mutationResolver) LocationUpdate(ctx context.Context, input model.LocationUpdateInput) (*model.Location, error) {
	panic(fmt.Errorf("not implemented: LocationUpdate - location_Update"))
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *locationResolver) Place(ctx context.Context, obj *model.Location) (*model.Place, error) {
	panic(fmt.Errorf("not implemented: Place - place"))
}
func (r *Resolver) Location() generated.LocationResolver { return &locationResolver{r} }

type locationResolver struct{ *Resolver }
