package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/file-store-api/test/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Attachment       func(ctx context.Context, id string) (*model.Attachment, error)
	AttachmentCreate func(ctx context.Context, input model.AttachmentInput) (*model.Attachment, error)
}
