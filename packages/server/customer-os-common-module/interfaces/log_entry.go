package interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
)

type LogEntryService interface {
	Save(ctx context.Context, id *string, logEntryFields data_fields.LogEntryFields) (string, error)
}
