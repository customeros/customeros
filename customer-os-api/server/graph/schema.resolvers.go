package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

var contacts []*model.Contact

// CreateContact is the resolver for the createContact field.
func (r *mutationResolver) CreateContact(ctx context.Context, request model.ContactRequest) (*model.Contact, error) {
	var contact = model.Contact{
		ID:          fmt.Sprintf("%d%s", rand.Uint32(), ""),
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		MiddleName:  request.MiddleName,
		Comments:    request.Comments,
		PhoneNumber: request.PhoneNumber,
		Email:       request.Email,
		Address:     request.Address,
	}

	contacts = append(contacts, &contact)

	return &contact, nil
}

// UpdateContact is the resolver for the updateContact field.
func (r *mutationResolver) UpdateContact(ctx context.Context, id string, request model.ContactRequest) (*model.Contact, error) {
	panic(fmt.Errorf("not implemented: UpdateContact - updateContact"))
}

// Contacts is the resolver for the contacts field.
func (r *queryResolver) Contacts(ctx context.Context) ([]*model.Contact, error) {
	return contacts[0:len(contacts)], nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
