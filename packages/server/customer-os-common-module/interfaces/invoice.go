package interfaces

import (
	"context"
	"time"

	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/enum"
	neo4jrepository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	invoicepb "github.com/customeros/customeros/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/invoice"
)

type InvoiceService interface {
	SetContractService(contract ContractService)
	SetServiceLineItemService(sli ServiceLineItemService)
	IsInitialized() bool

	GenerateNewRandomInvoiceNumber() string

	GetById(ctx context.Context, invoiceId string) (*neo4jentity.InvoiceEntity, error)
	GetByIdAcrossAllTenants(ctx context.Context, invoiceId string) (*neo4jentity.InvoiceEntity, string, error)
	GetByNumber(ctx context.Context, number string) (*neo4jentity.InvoiceEntity, error)
	GetInvoiceLinesForInvoices(ctx context.Context, invoiceIds []string) (*neo4jentity.InvoiceLineEntities, error)
	GetInvoicesForContracts(ctx context.Context, contractIds []string) (*neo4jentity.InvoiceEntities, error)
	GetNonDryRunInvoicesForOrganization(ctx context.Context, tenant, organizationId string) (*neo4jentity.InvoiceEntities, error)
	SimulateInvoice(ctx context.Context, invoiceData *SimulateInvoiceRequestData) ([]*SimulateInvoiceResponseData, error)
	NextInvoiceDryRun(ctx context.Context, contractId, appSource string) (string, error)
	PayInvoice(ctx context.Context, invoiceId string) error
	VoidInvoice(ctx context.Context, invoiceId, appSource string) error
	UpdateInvoice(ctx context.Context, invoiceId string, data neo4jrepository.InvoiceUpdateFields) error

	FillCycleInvoice(ctx context.Context, invoiceEntity *neo4jentity.InvoiceEntity, sliEntities neo4jentity.ServiceLineItemEntities) (*neo4jentity.InvoiceEntity, []*invoicepb.InvoiceLine, error)
	// Deprecated: Method should be re-worked. DO NOT ENABLE IN PROD
	FillOffCyclePrepaidInvoice(ctx context.Context, invoiceEntity *neo4jentity.InvoiceEntity, sliEntities neo4jentity.ServiceLineItemEntities) (*neo4jentity.InvoiceEntity, []*invoicepb.InvoiceLine, error)
}

type SimulateInvoiceRequestData struct {
	ContractId   string
	ServiceLines []SimulateInvoiceRequestServiceLineData
}

type SimulateInvoiceRequestServiceLineData struct {
	Key               string
	ServiceLineItemID string
	ParentID          string
	Description       string
	Comments          string
	BillingCycle      neo4jenum.BilledType
	Price             float64
	Quantity          int64
	ServiceStarted    time.Time
	ServiceEnded      *time.Time
	TaxRate           *float64
	Canceled          bool
}

type SimulateInvoiceResponseData struct {
	Invoice *neo4jentity.InvoiceEntity
	Lines   []*neo4jentity.InvoiceLineEntity
}
