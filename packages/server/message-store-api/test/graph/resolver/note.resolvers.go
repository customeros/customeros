package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/model"
)

// NoteMergeToContact is the resolver for the note_MergeToContact field.
func (r *mutationResolver) NoteMergeToContact(ctx context.Context, contactID string, input model.NoteInput) (*model.Note, error) {
	panic(fmt.Errorf("not implemented: NoteMergeToContact - note_MergeToContact"))
}

// NoteUpdateInContact is the resolver for the note_UpdateInContact field.
func (r *mutationResolver) NoteUpdateInContact(ctx context.Context, contactID string, input model.NoteUpdateInput) (*model.Note, error) {
	panic(fmt.Errorf("not implemented: NoteUpdateInContact - note_UpdateInContact"))
}

// NoteDeleteFromContact is the resolver for the note_DeleteFromContact field.
func (r *mutationResolver) NoteDeleteFromContact(ctx context.Context, contactID string, noteID string) (*model.Result, error) {
	panic(fmt.Errorf("not implemented: NoteDeleteFromContact - note_DeleteFromContact"))
}

// CreatedBy is the resolver for the createdBy field.
func (r *noteResolver) CreatedBy(ctx context.Context, obj *model.Note) (*model.User, error) {
	panic(fmt.Errorf("not implemented: CreatedBy - createdBy"))
}

// Note returns generated.NoteResolver implementation.
func (r *Resolver) Note() generated.NoteResolver { return &noteResolver{r} }

type noteResolver struct{ *Resolver }
