package command

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/service_line_item/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

// UpdateServiceLineItemCommand contains the data needed to update a service line item.
type UpdateServiceLineItemCommand struct {
	eventstore.BaseCommand
	DataFields model.ServiceLineItemDataFields
	Source     commonmodel.Source
	UpdatedAt  *time.Time
}

// NewUpdateServiceLineItemCommand creates a new UpdateServiceLineItemCommand.
func NewUpdateServiceLineItemCommand(serviceLineItemId, tenant, loggedInUserId string, dataFields model.ServiceLineItemDataFields, source commonmodel.Source, updatedAt *time.Time) *UpdateServiceLineItemCommand {
	return &UpdateServiceLineItemCommand{
		BaseCommand: eventstore.NewBaseCommand(serviceLineItemId, tenant, loggedInUserId),
		DataFields:  dataFields,
		Source:      source,
		UpdatedAt:   updatedAt,
	}
}
