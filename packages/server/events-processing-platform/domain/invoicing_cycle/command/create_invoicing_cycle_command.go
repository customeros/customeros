package command

import (
	commonmodel "github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/common/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/invoicing_cycle/model"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/eventstore"
	"time"
)

type CreateInvoicingCycleTypeCommand struct {
	eventstore.BaseCommand
	SourceFields      commonmodel.Source
	CreatedAt         *time.Time
	InvoicingDateType model.InvoicingCycleType
}

func NewCreateInvoicingCycleTypeCommand(invoicingCycleId, tenant, loggedInUserId string, sourceFields commonmodel.Source, createdAt *time.Time, invoicingDateType model.InvoicingCycleType) *CreateInvoicingCycleTypeCommand {
	return &CreateInvoicingCycleTypeCommand{
		BaseCommand:       eventstore.NewBaseCommand(invoicingCycleId, tenant, loggedInUserId).WithAppSource(sourceFields.AppSource),
		SourceFields:      sourceFields,
		CreatedAt:         createdAt,
		InvoicingDateType: invoicingDateType,
	}
}
