package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/message-store/test/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/test/graph/model"
)

// Organization is the resolver for the organization field.
func (r *contactRoleResolver) Organization(ctx context.Context, obj *model.ContactRole) (*model.Organization, error) {
	panic(fmt.Errorf("not implemented: Organization - organization"))
}

// Contact is the resolver for the contact field.
func (r *contactRoleResolver) Contact(ctx context.Context, obj *model.ContactRole) (*model.Contact, error) {
	panic(fmt.Errorf("not implemented: Contact - contact"))
}

// ContactRoleDelete is the resolver for the contactRole_Delete field.
func (r *mutationResolver) ContactRoleDelete(ctx context.Context, contactID string, roleID string) (*model.Result, error) {
	panic(fmt.Errorf("not implemented: ContactRoleDelete - contactRole_Delete"))
}

// ContactRoleCreate is the resolver for the contactRole_Create field.
func (r *mutationResolver) ContactRoleCreate(ctx context.Context, contactID string, input model.ContactRoleInput) (*model.ContactRole, error) {
	panic(fmt.Errorf("not implemented: ContactRoleCreate - contactRole_Create"))
}

// ContactRoleUpdate is the resolver for the contactRole_Update field.
func (r *mutationResolver) ContactRoleUpdate(ctx context.Context, contactID string, roleID string, input model.ContactRoleInput) (*model.ContactRole, error) {
	panic(fmt.Errorf("not implemented: ContactRoleUpdate - contactRole_Update"))
}

// ContactRole returns generated.ContactRoleResolver implementation.
func (r *Resolver) ContactRole() generated.ContactRoleResolver { return &contactRoleResolver{r} }

type contactRoleResolver struct{ *Resolver }
