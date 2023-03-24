package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/test/graph/model"
)

// Tags is the resolver for the tags field.
func (r *issueResolver) Tags(ctx context.Context, obj *model.Issue) ([]*model.Tag, error) {
	panic(fmt.Errorf("not implemented: Tags - tags"))
}

// Tags is the resolver for the tags field.
func (r *ticketResolver) Tags(ctx context.Context, obj *model.Ticket) ([]*model.Tag, error) {
	panic(fmt.Errorf("not implemented: Tags - tags"))
}

// Notes is the resolver for the notes field.
func (r *ticketResolver) Notes(ctx context.Context, obj *model.Ticket) ([]*model.Note, error) {
	panic(fmt.Errorf("not implemented: Notes - notes"))
}

// Issue returns generated.IssueResolver implementation.
func (r *Resolver) Issue() generated.IssueResolver { return &issueResolver{r} }

// Ticket returns generated.TicketResolver implementation.
func (r *Resolver) Ticket() generated.TicketResolver { return &ticketResolver{r} }

type issueResolver struct{ *Resolver }
type ticketResolver struct{ *Resolver }
