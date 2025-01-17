package cosapi_interfaces

import (
	"context"
	"time"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graphql/model"
)

type ContractService interface {
	Create(ctx context.Context, contractDetails *ContractCreateData) (string, error)
	Update(ctx context.Context, input model.ContractUpdateInput) error
	SoftDeleteContract(ctx context.Context, contractId string) (bool, error)
	GetById(ctx context.Context, id string) (*entity.ContractEntity, error)
	GetContractsForOrganizations(ctx context.Context, organizationIds []string) (*entity.ContractEntities, error)
	GetContractsForInvoices(ctx context.Context, invoiceIds []string) (*entity.ContractEntities, error)
	GetContractByServiceLineItem(ctx context.Context, serviceLineItemId string) (*entity.ContractEntity, error)
	ContractsExistForTenant(ctx context.Context) (bool, error)
	CountContracts(ctx context.Context, tenant string) (int64, error)
	RenewContract(ctx context.Context, contractId string, renewalDate *time.Time) error
	GetPaginatedContracts(ctx context.Context, page int, limit int) (*utils.Pagination, error)
}

type ContractCreateData struct {
	Input             model.ContractInput
	ExternalReference *entity.ExternalSystemEntity
	Source            entity.DataSource
	AppSource         string
}
