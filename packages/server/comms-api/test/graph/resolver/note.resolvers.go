package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/model"
)

// NoteCreateForContact is the resolver for the note_CreateForContact field.
func (r *mutationResolver) NoteCreateForContact(ctx context.Context, contactID string, input model.NoteInput) (*model.Note, error) {
	panic(fmt.Errorf("not implemented: NoteCreateForContact - note_CreateForContact"))
}

// NoteCreateForOrganization is the resolver for the note_CreateForOrganization field.
func (r *mutationResolver) NoteCreateForOrganization(ctx context.Context, organizationID string, input model.NoteInput) (*model.Note, error) {
	panic(fmt.Errorf("not implemented: NoteCreateForOrganization - note_CreateForOrganization"))
}

// NoteUpdate is the resolver for the note_Update field.
func (r *mutationResolver) NoteUpdate(ctx context.Context, input model.NoteUpdateInput) (*model.Note, error) {
	panic(fmt.Errorf("not implemented: NoteUpdate - note_Update"))
}

// NoteDelete is the resolver for the note_Delete field.
func (r *mutationResolver) NoteDelete(ctx context.Context, id string) (*model.Result, error) {
	panic(fmt.Errorf("not implemented: NoteDelete - note_Delete"))
}

// CreatedBy is the resolver for the createdBy field.
func (r *noteResolver) CreatedBy(ctx context.Context, obj *model.Note) (*model.User, error) {
	panic(fmt.Errorf("not implemented: CreatedBy - createdBy"))
}

// Noted is the resolver for the noted field.
func (r *noteResolver) Noted(ctx context.Context, obj *model.Note) ([]model.NotedEntity, error) {
	panic(fmt.Errorf("not implemented: Noted - noted"))
}

// Note returns generated.NoteResolver implementation.
func (r *Resolver) Note() generated.NoteResolver { return &noteResolver{r} }

type noteResolver struct{ *Resolver }
