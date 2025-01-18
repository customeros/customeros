package interfaces

import (
	"context"
	"time"

	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type ServiceLineItemService interface {
	SetContractService(contract ContractService)
	IsInitialized() bool

	GetById(ctx context.Context, id string) (*neo4jentity.ServiceLineItemEntity, error)
	GetServiceLineItemsByParentId(ctx context.Context, sliParentId string) (*neo4jentity.ServiceLineItemEntities, error)
	GetServiceLineItemsForContract(ctx context.Context, contractId string) (*neo4jentity.ServiceLineItemEntities, error)
	GetServiceLineItemsForContracts(ctx context.Context, contractIds []string) (*neo4jentity.ServiceLineItemEntities, error)
	GetServiceLineItemsForInvoiceLines(ctx context.Context, invoiceLineIds []string) (*neo4jentity.ServiceLineItemEntities, error)
	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, id *string, dataFields data_fields.SLIFields) (string, error)
	Pause(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, serviceLineItemId string) error
	Resume(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, serviceLineItemId string) error
	Delete(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, serviceLineItemId string) error
	Close(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, serviceLineItemId string, endedAt time.Time) error
}
