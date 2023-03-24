package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/generated"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/test/graph/model"
)

// Tags is the resolver for the tags field.
func (r *issueResolver) Tags(ctx context.Context, obj *model.Issue) ([]*model.Tag, error) {
	panic(fmt.Errorf("not implemented: Tags - tags"))
}

// Issue returns generated.IssueResolver implementation.
func (r *Resolver) Issue() generated.IssueResolver { return &issueResolver{r} }

type issueResolver struct{ *Resolver }
