package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/generated"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
)

// CreateContact is the resolver for the createContact field.
func (r *mutationResolver) CreateContact(ctx context.Context, request model.ContactRequest) (*model.Contact, error) {
	contactNode := entity.ContactNode{
		FirstName:   *request.FirstName,
		LastName:    *request.LastName,
		Label:       *request.Label,
		ContactType: *request.ContactType,
	}

	contactNodeCreated, _ := r.ServiceContainer.ContactService.Create(contactNode)

	contactCreatedResponse := model.Contact{
		ID:          contactNodeCreated.Id,
		FirstName:   &contactNodeCreated.FirstName,
		LastName:    &contactNodeCreated.LastName,
		Label:       &contactNodeCreated.Label,
		ContactType: &contactNodeCreated.ContactType,
	}
	return &contactCreatedResponse, nil
}

// UpdateContact is the resolver for the updateContact field.
func (r *mutationResolver) UpdateContact(ctx context.Context, id string, request model.ContactRequest) (*model.Contact, error) {
	panic(fmt.Errorf("not implemented: UpdateContact - updateContact"))
}

// Contacts is the resolver for the contacts field.
func (r *queryResolver) Contacts(ctx context.Context) ([]*model.Contact, error) {
	contacts, _ := r.ServiceContainer.ContactService.FindAll()
	return contactsDtoFromNodes(contacts), nil
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
func contactsDtoFromNodes(nodes []entity.ContactNode) []*model.Contact {
	contacts := make([]*model.Contact, 0, len(nodes))
	for _, node := range nodes {
		contact := model.Contact{
			ID:          node.Id,
			FirstName:   &node.FirstName,
			LastName:    &node.LastName,
			Label:       &node.Label,
			ContactType: &node.ContactType,
		}
		contacts = append(contacts, &contact)
	}
	return contacts
}

var contacts []*model.Contact
