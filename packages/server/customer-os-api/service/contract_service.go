package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ContractService interface {
	Create(ctx context.Context, contract *ContractCreateData) (*entity.ContractEntity, error)
}
type contractService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
}

func NewContractService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) ContractService {
	return &contractService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
	}
}

func (s *contractService) getNeo4jDriver() neo4j.DriverWithContext {
	return *s.repositories.Drivers.Neo4jDriver
}

type ContractCreateData struct {
	ContractEntity    *entity.ContractEntity
	Organization      string
	ExternalReference *entity.ExternalSystemEntity
	Source            entity.DataSource
	AppSource         string
}

func (s *contractService) Create(ctx context.Context, contract *ContractCreateData) (*entity.ContractEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "NoteService.CreateNoteForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organization", contract.Organization))

	dbNodePtr, err := s.repositories.ContractRepository.CreateContract(ctx, common.GetContext(ctx).Tenant, contract.Organization, *contract.ContractEntity)
	if err != nil {
		return nil, err
	}
	// set contract creator
	if len(common.GetUserIdFromContext(ctx)) > 0 {
		props := utils.GetPropsFromNode(*dbNodePtr)
		contractId := utils.GetStringPropOrEmpty(props, "id")
		err = s.repositories.ContractRepository.SetContractCreator(ctx, common.GetTenantFromContext(ctx), common.GetUserIdFromContext(ctx), contractId)
		if err != nil {
			s.log.Error("Failed to set contract creator", err)
			return nil, err
		}
	}
	return s.mapDbNodeToContractEntity(*dbNodePtr), nil
}

func (s *contractService) mapDbNodeToContractEntity(node dbtype.Node) *entity.ContractEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.ContractEntity{
		ID:        utils.GetStringPropOrEmpty(props, "id"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:    entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		AppSource: utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &result
}
