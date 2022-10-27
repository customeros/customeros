package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

// Contacts is the resolver for the contacts field.
func (r *queryResolver) Contacts(ctx context.Context) ([]*model.Contact, error) {
	//contacts, _ := r.ServiceContainer.ContactService.FindAll()
	//return contactsDtoFromNodes(contacts), nil
	return nil, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
