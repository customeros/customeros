package interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
)

type ContractService interface {
	SetOpportunityService(opportunity OpportunityService)
	SetOrganizationService(org OrganizationService)
	IsInitialized() bool

	GetById(ctx context.Context, contactId string) (*entity.ContractEntity, error)
	Save(ctx context.Context, contactId *string, dataFields data_fields.ContractSaveFields) (string, error)
	SoftDelete(ctx context.Context, contractId string) error
	RefreshContractStatus(ctx context.Context, contractId string) error
	RecalculateContractLtv(ctx context.Context, contractId string) error
	UpdateActiveRenewalOpportunityArr(ctx context.Context, contractId string) error
	UpdateActiveRenewalOpportunityRenewDateAndArr(ctx context.Context, tenant, contractId string) error
	UpdateActiveRenewalOpportunityLikelihood(ctx context.Context, tenant, contractId string) error
}
