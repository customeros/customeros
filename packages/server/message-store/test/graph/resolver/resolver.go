package resolver

import (
	"context"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/test/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ContactCreate func(ctx context.Context, input model.ContactInput) (*model.Contact, error)
}
