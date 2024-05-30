package enummapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
)

var invoiceStatusByModel = map[model.InvoiceStatus]neo4jenum.InvoiceStatus{
	model.InvoiceStatusDraft:       neo4jenum.InvoiceStatusInitialized,
	model.InvoiceStatusInitialized: neo4jenum.InvoiceStatusInitialized,
	model.InvoiceStatusDue:         neo4jenum.InvoiceStatusDue,
	model.InvoiceStatusPaid:        neo4jenum.InvoiceStatusPaid,
	model.InvoiceStatusVoid:        neo4jenum.InvoiceStatusVoid,
	model.InvoiceStatusScheduled:   neo4jenum.InvoiceStatusScheduled,
	model.InvoiceStatusOverdue:     neo4jenum.InvoiceStatusOverdue,
	model.InvoiceStatusOnHold:      neo4jenum.InvoiceStatusOnHold,
	model.InvoiceStatusEmpty:       neo4jenum.InvoiceStatusEmpty,
}

var invoiceStatusByValue = utils.ReverseMap(invoiceStatusByModel)

func MapInvoiceStatusFromModel(input model.InvoiceStatus) neo4jenum.InvoiceStatus {
	return invoiceStatusByModel[input]
}

func MapInvoiceStatusToModel(input neo4jenum.InvoiceStatus) model.InvoiceStatus {
	return invoiceStatusByValue[input]
}
