package command

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoicing_cycle/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type UpdateInvoicingCycleCommand struct {
	eventstore.BaseCommand
	UpdatedAt *time.Time
	Type      model.InvoicingCycleType
}

func NewUpdateInvoicingCycleCommand(invoicingCycleId, tenant, loggedInUserId string, sourceFields commonmodel.Source, updatedAt *time.Time, invoicingCycleType model.InvoicingCycleType) *UpdateInvoicingCycleCommand {
	return &UpdateInvoicingCycleCommand{
		BaseCommand: eventstore.NewBaseCommand(invoicingCycleId, tenant, loggedInUserId).WithAppSource(sourceFields.AppSource),
		UpdatedAt:   updatedAt,
		Type:        invoicingCycleType,
	}
}
