package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/mapper"
)

// CreateContact is the resolver for the createContact field.
func (r *mutationResolver) CreateContact(ctx context.Context, input model.ContactInput) (*model.Contact, error) {
	contactNodeCreated, err := r.ServiceContainer.ContactService.Create(mapper.MapContactInputToEntity(input))
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to create contact %s %s", input.FirstName, input.LastName)
		return nil, nil
	}

	return mapper.MapEntityToContact(contactNodeCreated), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
