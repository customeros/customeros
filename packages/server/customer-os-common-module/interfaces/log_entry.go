package interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
)

type LogEntryService interface {
	SetOrganizationService(org OrganizationService)
	IsInitialized() bool

	Save(ctx context.Context, id *string, logEntryFields data_fields.LogEntryFields) (string, error)
}
