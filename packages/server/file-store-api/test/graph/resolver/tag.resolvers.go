package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/test/graph/model"
)

// TagCreate is the resolver for the tag_Create field.
func (r *mutationResolver) TagCreate(ctx context.Context, input model.TagInput) (*model.Tag, error) {
	panic(fmt.Errorf("not implemented: TagCreate - tag_Create"))
}

// TagUpdate is the resolver for the tag_Update field.
func (r *mutationResolver) TagUpdate(ctx context.Context, input model.TagUpdateInput) (*model.Tag, error) {
	panic(fmt.Errorf("not implemented: TagUpdate - tag_Update"))
}

// TagDelete is the resolver for the tag_Delete field.
func (r *mutationResolver) TagDelete(ctx context.Context, id string) (*model.Result, error) {
	panic(fmt.Errorf("not implemented: TagDelete - tag_Delete"))
}

// Tags is the resolver for the tags field.
func (r *queryResolver) Tags(ctx context.Context) ([]*model.Tag, error) {
	panic(fmt.Errorf("not implemented: Tags - tags"))
}
