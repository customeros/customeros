package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/mapper"
)

// CreateContact is the resolver for the createContact field.
func (r *mutationResolver) CreateContact(ctx context.Context, input model.ContactInput) (*model.Contact, error) {
	contactNodeCreated, err := r.ServiceContainer.ContactService.Create(mapper.MapContactInputToEntity(input))
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to create contact %s %s", input.FirstName, input.LastName)
		return nil, err
	}

	return mapper.MapEntityToContact(contactNodeCreated), nil
}

// CreateContactGroup is the resolver for the createContactGroup field.
func (r *mutationResolver) CreateContactGroup(ctx context.Context, name string) (*model.ContactGroup, error) {
	contactGroupNodeCreated, err := r.ServiceContainer.ContactGroupService.Create(&entity.ContactGroupNode{
		Name: name,
	})
	if err != nil {
		graphql.AddErrorf(ctx, "Failed to create contact group %s", name)
		return nil, err
	}

	return mapper.MapEntityToContactGroup(contactGroupNodeCreated), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
