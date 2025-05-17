package api_bank_account

import (
	"fmt"
	cosapi_interfaces "github.com/customeros/customeros/packages/server/customer-os-api/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-api/repository"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	"golang.org/x/net/context"
)

type bankAccountService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewBankAccountService(log logger.Logger, repository *repository.Repositories) cosapi_interfaces.BankAccountService {
	return &bankAccountService{
		log:          log,
		repositories: repository,
	}
}

func (s *bankAccountService) GetTenantBankAccounts(ctx context.Context) (*neo4jentity.BankAccountEntities, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "BankAccountService.GetTenantBankAccounts")
	defer spans.Finish()

	dbNodes, err := s.repositories.Neo4jRepositories.BankAccountReadRepository.GetBankAccounts(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		spans.TraceError(err)
		return nil, fmt.Errorf("GetTenantBankAccounts: %s", err.Error())
	}

	tenantBankAccounts := neo4jentity.BankAccountEntities{}
	for _, dbNode := range dbNodes {
		tenantBankAccounts = append(tenantBankAccounts, *neo4jmapper.MapDbNodeToBankAccountEntity(dbNode))
	}

	return &tenantBankAccounts, nil
}

func (s *bankAccountService) GetTenantBankAccount(ctx context.Context, id string) (*neo4jentity.BankAccountEntity, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "BankAccountService.GetTenantBankAccount")
	defer span.Finish()
	span.LogKV("bankAccountId", id)

	dbNode, err := s.repositories.Neo4jRepositories.BankAccountReadRepository.GetBankAccountById(ctx, common.GetTenantFromContext(ctx), id)
	if err != nil {
		span.TraceError(err)
		return nil, fmt.Errorf("GetTenantBankAccount: %s", err.Error())
	}

	return neo4jmapper.MapDbNodeToBankAccountEntity(dbNode), nil
}
