package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type BankAccountService interface {
	GetTenantBankAccounts(ctx context.Context) (*neo4jentity.BankAccountEntities, error)
	GetTenantBankAccount(ctx context.Context, id string) (*neo4jentity.BankAccountEntity, error)
}

type bankAccountService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
}

func NewBankAccountService(log logger.Logger, repository *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) BankAccountService {
	return &bankAccountService{
		log:          log,
		repositories: repository,
		grpcClients:  grpcClients,
		services:     services,
	}
}

func (s *bankAccountService) GetTenantBankAccounts(ctx context.Context) (*neo4jentity.BankAccountEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountService.GetTenantBankAccounts")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	dbNodes, err := s.repositories.Neo4jRepositories.BankAccountReadRepository.GetBankAccounts(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantBankAccounts: %s", err.Error())
	}

	tenantBankAccounts := neo4jentity.BankAccountEntities{}
	for _, dbNode := range dbNodes {
		tenantBankAccounts = append(tenantBankAccounts, *neo4jmapper.MapDbNodeToBankAccountEntity(dbNode))
	}

	return &tenantBankAccounts, nil
}

func (s *bankAccountService) GetTenantBankAccount(ctx context.Context, id string) (*neo4jentity.BankAccountEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BankAccountService.GetTenantBankAccount")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("bankAccountId", id))

	dbNode, err := s.repositories.Neo4jRepositories.BankAccountReadRepository.GetBankAccountById(ctx, common.GetTenantFromContext(ctx), id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("GetTenantBankAccount: %s", err.Error())
	}

	return neo4jmapper.MapDbNodeToBankAccountEntity(dbNode), nil
}
