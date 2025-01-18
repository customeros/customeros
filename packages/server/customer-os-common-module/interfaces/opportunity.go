package interfaces

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type OpportunityService interface {
	SetContractService(contract ContractService)
	SetOrganizationService(org OrganizationService)
	IsInitialized() bool

	GetById(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, opportunityId string) (*neo4jentity.OpportunityEntity, error)
	GetOpportunitiesForContracts(ctx context.Context, tenant string, contractIds []string) (*neo4jentity.OpportunityEntities, error)
	GetOpportunitiesForOrganizations(ctx context.Context, tenant string, organizationIds []string) (*neo4jentity.OpportunityEntities, error)
	GetPaginatedOrganizationOpportunities(ctx context.Context, tenant string, page int, limit int) (*utils.Pagination, error)

	Save(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, opportunityId *string, input *data_fields.OpportunityFields) (string, error)
	CreateRenewalOpportunity(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, input *data_fields.OpportunityFields) (string, error)
	CloseWon(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, tenant, opportunityId string) error
	CloseLost(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, tenant, opportunityId string) error
	Archive(ctx context.Context, tenant, opportunityId string) error

	RolloutRenewalOpportunity(ctx context.Context, contractId string) error
}
