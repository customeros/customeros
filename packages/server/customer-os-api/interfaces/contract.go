package cosapi_interfaces

import (
	"context"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	neo4j_entity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-api/graphql/model"
)

type ContractService interface {
	Create(ctx context.Context, contractDetails *ContractCreateData) (string, error)
	Update(ctx context.Context, input model.ContractUpdateInput) error
	SoftDeleteContract(ctx context.Context, contractId string) (bool, error)
	GetById(ctx context.Context, id string) (*neo4j_entity.ContractEntity, error)
	GetContractsForOrganizations(ctx context.Context, organizationIds []string) (*neo4j_entity.ContractEntities, error)
	GetContractsForInvoices(ctx context.Context, invoiceIds []string) (*neo4j_entity.ContractEntities, error)
	GetContractByServiceLineItem(ctx context.Context, serviceLineItemId string) (*neo4j_entity.ContractEntity, error)
	ContractsExistForTenant(ctx context.Context) (bool, error)
	CountContracts(ctx context.Context, tenant string) (int64, error)
	RenewContract(ctx context.Context, contractId string, renewalDate *time.Time) error
	GetPaginatedContracts(ctx context.Context, page int, limit int) (*utils.Pagination, error)
}

type ContractCreateData struct {
	Input             model.ContractInput
	ExternalReference *neo4j_entity.ExternalSystemEntity
	Source            neo4j_entity.DataSource
	AppSource         string
}
