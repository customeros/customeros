package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/model"
)

// Contacts is the resolver for the contacts field.
func (r *contactGroupResolver) Contacts(ctx context.Context, obj *model.ContactGroup, pagination *model.Pagination, where *model.Filter, sort []*model.SortBy) (*model.ContactsPage, error) {
	panic(fmt.Errorf("not implemented: Contacts - contacts"))
}

// ContactGroupCreate is the resolver for the contactGroupCreate field.
func (r *mutationResolver) ContactGroupCreate(ctx context.Context, input model.ContactGroupInput) (*model.ContactGroup, error) {
	panic(fmt.Errorf("not implemented: ContactGroupCreate - contactGroupCreate"))
}

// ContactGroupUpdate is the resolver for the contactGroupUpdate field.
func (r *mutationResolver) ContactGroupUpdate(ctx context.Context, input model.ContactGroupUpdateInput) (*model.ContactGroup, error) {
	panic(fmt.Errorf("not implemented: ContactGroupUpdate - contactGroupUpdate"))
}

// ContactGroupDeleteAndUnlinkAllContacts is the resolver for the contactGroupDeleteAndUnlinkAllContacts field.
func (r *mutationResolver) ContactGroupDeleteAndUnlinkAllContacts(ctx context.Context, id string) (*model.Result, error) {
	panic(fmt.Errorf("not implemented: ContactGroupDeleteAndUnlinkAllContacts - contactGroupDeleteAndUnlinkAllContacts"))
}

// ContactGroupAddContact is the resolver for the contactGroupAddContact field.
func (r *mutationResolver) ContactGroupAddContact(ctx context.Context, contactID string, groupID string) (*model.Result, error) {
	panic(fmt.Errorf("not implemented: ContactGroupAddContact - contactGroupAddContact"))
}

// ContactGroupRemoveContact is the resolver for the contactGroupRemoveContact field.
func (r *mutationResolver) ContactGroupRemoveContact(ctx context.Context, contactID string, groupID string) (*model.Result, error) {
	panic(fmt.Errorf("not implemented: ContactGroupRemoveContact - contactGroupRemoveContact"))
}

// ContactGroup is the resolver for the contactGroup field.
func (r *queryResolver) ContactGroup(ctx context.Context, id string) (*model.ContactGroup, error) {
	panic(fmt.Errorf("not implemented: ContactGroup - contactGroup"))
}

// ContactGroups is the resolver for the contactGroups field.
func (r *queryResolver) ContactGroups(ctx context.Context, pagination *model.Pagination, where *model.Filter, sort []*model.SortBy) (*model.ContactGroupPage, error) {
	panic(fmt.Errorf("not implemented: ContactGroups - contactGroups"))
}

// ContactGroup returns generated.ContactGroupResolver implementation.
func (r *Resolver) ContactGroup() generated.ContactGroupResolver { return &contactGroupResolver{r} }

type contactGroupResolver struct{ *Resolver }
