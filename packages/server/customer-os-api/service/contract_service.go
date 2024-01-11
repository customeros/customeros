package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type ContractService interface {
	Create(ctx context.Context, contract *ContractCreateData) (string, error)
	Update(ctx context.Context, contract *entity.ContractEntity) error
	GetById(ctx context.Context, id string) (*entity.ContractEntity, error)
	GetContractsForOrganizations(ctx context.Context, organizationIds []string) (*entity.ContractEntities, error)
	ContractsExistForTenant(ctx context.Context) (bool, error)
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

type ContractCreateData struct {
	ContractEntity    *entity.ContractEntity
	OrganizationId    string
	ExternalReference *entity.ExternalSystemEntity
	Source            neo4jentity.DataSource
	AppSource         string
}

func (s *contractService) Create(ctx context.Context, contractDetails *ContractCreateData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contractDetails", contractDetails))

	if contractDetails.ContractEntity == nil {
		err := fmt.Errorf("contract entity is nil")
		tracing.TraceErr(span, err)
		return "", err
	}

	contractId, err := s.createContractWithEvents(ctx, contractDetails)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	span.LogFields(log.String("output - createdContractId", contractId))
	return contractId, nil
}

func (s *contractService) Update(ctx context.Context, contract *entity.ContractEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("contract", contract))

	if contract == nil {
		err := fmt.Errorf("(ContractService.Update) contract entity is nil")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	} else if contract.Id == "" {
		err := fmt.Errorf("(ContractService.Update) contract id is missing")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	contractExists, _ := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), contract.Id, neo4jentity.NodeLabelContract)
	if !contractExists {
		err := fmt.Errorf("(ContractService.Update) contract with id {%s} not found", contract.Id)
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	contractUpdateRequest := contractpb.UpdateContractGrpcRequest{
		Tenant:           common.GetTenantFromContext(ctx),
		Id:               contract.Id,
		LoggedInUserId:   common.GetUserIdFromContext(ctx),
		Name:             contract.Name,
		ContractUrl:      contract.ContractUrl,
		SignedAt:         utils.ConvertTimeToTimestampPtr(contract.SignedAt),
		ServiceStartedAt: utils.ConvertTimeToTimestampPtr(contract.ServiceStartedAt),
		EndedAt:          utils.ConvertTimeToTimestampPtr(contract.EndedAt),
		SourceFields: &commonpb.SourceFields{
			Source:    string(contract.Source),
			AppSource: utils.StringFirstNonEmpty(contract.AppSource, constants.AppSourceCustomerOsApi),
		},
		RenewalPeriods: contract.RenewalPeriods,
	}
	switch contract.RenewalCycle {
	case entity.RenewalCycleMonthlyRenewal:
		contractUpdateRequest.RenewalCycle = contractpb.RenewalCycle_MONTHLY_RENEWAL
	case entity.RenewalCycleQuarterlyRenewal:
		contractUpdateRequest.RenewalCycle = contractpb.RenewalCycle_QUARTERLY_RENEWAL
	case entity.RenewalCycleAnnualRenewal:
		contractUpdateRequest.RenewalCycle = contractpb.RenewalCycle_ANNUALLY_RENEWAL
	default:
		contractUpdateRequest.RenewalCycle = contractpb.RenewalCycle_NONE
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := s.grpcClients.ContractClient.UpdateContract(ctx, &contractUpdateRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}

func (s *contractService) createContractWithEvents(ctx context.Context, contractDetails *ContractCreateData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.createContractWithEvents")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	createContractRequest := contractpb.CreateContractGrpcRequest{
		Tenant:           common.GetTenantFromContext(ctx),
		OrganizationId:   contractDetails.OrganizationId,
		Name:             contractDetails.ContractEntity.Name,
		ContractUrl:      contractDetails.ContractEntity.ContractUrl,
		SignedAt:         utils.ConvertTimeToTimestampPtr(contractDetails.ContractEntity.SignedAt),
		ServiceStartedAt: utils.ConvertTimeToTimestampPtr(contractDetails.ContractEntity.ServiceStartedAt),
		LoggedInUserId:   common.GetUserIdFromContext(ctx),
		SourceFields: &commonpb.SourceFields{
			Source:    string(contractDetails.Source),
			AppSource: utils.StringFirstNonEmpty(contractDetails.AppSource, constants.AppSourceCustomerOsApi),
		},
		RenewalPeriods: contractDetails.ContractEntity.RenewalPeriods,
	}

	switch contractDetails.ContractEntity.RenewalCycle {
	case entity.RenewalCycleMonthlyRenewal:
		createContractRequest.RenewalCycle = contractpb.RenewalCycle_MONTHLY_RENEWAL
	case entity.RenewalCycleQuarterlyRenewal:
		createContractRequest.RenewalCycle = contractpb.RenewalCycle_QUARTERLY_RENEWAL
	case entity.RenewalCycleAnnualRenewal:
		createContractRequest.RenewalCycle = contractpb.RenewalCycle_ANNUALLY_RENEWAL
	default:
		createContractRequest.RenewalCycle = contractpb.RenewalCycle_NONE
	}

	if contractDetails.ExternalReference != nil && contractDetails.ExternalReference.ExternalSystemId != "" {
		createContractRequest.ExternalSystemFields = &commonpb.ExternalSystemFields{
			ExternalSystemId: string(contractDetails.ExternalReference.ExternalSystemId),
			ExternalId:       contractDetails.ExternalReference.Relationship.ExternalId,
			ExternalUrl:      utils.IfNotNilString(contractDetails.ExternalReference.Relationship.ExternalUrl),
			ExternalSource:   utils.IfNotNilString(contractDetails.ExternalReference.Relationship.ExternalSource),
		}
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := s.grpcClients.ContractClient.CreateContract(ctx, &createContractRequest)

	WaitForObjectCreationAndLogSpan(ctx, s.repositories, response.Id, neo4jentity.NodeLabelContact, span)
	return response.Id, err
}

func (s *contractService) GetById(ctx context.Context, contractId string) (*entity.ContractEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contractId", contractId))

	if contractDbNode, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, common.GetContext(ctx).Tenant, contractId); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Contract with id {%s} not found", contractId))
		return nil, wrappedErr
	} else {
		return s.mapDbNodeToContractEntity(*contractDbNode), nil
	}
}

func (s *contractService) GetContractsForOrganizations(ctx context.Context, organizationIDs []string) (*entity.ContractEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.GetContractsForOrganizations")
	defer span.Finish()
	span.LogFields(log.Object("organizationIDs", organizationIDs))

	contracts, err := s.repositories.Neo4jRepositories.ContractReadRepository.GetContractsForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIDs)
	if err != nil {
		return nil, err
	}
	contractEntities := make(entity.ContractEntities, 0, len(contracts))
	for _, v := range contracts {
		contractEntity := s.mapDbNodeToContractEntity(*v.Node)
		contractEntity.DataloaderKey = v.LinkedNodeId
		contractEntities = append(contractEntities, *contractEntity)
	}
	return &contractEntities, nil
}

func (s *contractService) ContractsExistForTenant(ctx context.Context) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.ContractsExistForTenant")
	defer span.Finish()

	contractsExistForTenant, err := s.repositories.Neo4jRepositories.ContractReadRepository.TenantsHasAtLeastOneContract(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		return false, err
	}
	return contractsExistForTenant, nil
}

func (s *contractService) mapDbNodeToContractEntity(dbNode dbtype.Node) *entity.ContractEntity {
	props := utils.GetPropsFromNode(dbNode)
	contractStatus := entity.GetContractStatus(utils.GetStringPropOrEmpty(props, "status"))
	contractRenewalCycle := entity.GetRenewalCycle(utils.GetStringPropOrEmpty(props, "renewalCycle"))

	contract := entity.ContractEntity{
		Id:               utils.GetStringPropOrEmpty(props, "id"),
		Name:             utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:        utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:        utils.GetTimePropOrEpochStart(props, "updatedAt"),
		ServiceStartedAt: utils.GetTimePropOrNil(props, "serviceStartedAt"),
		SignedAt:         utils.GetTimePropOrNil(props, "signedAt"),
		EndedAt:          utils.GetTimePropOrNil(props, "endedAt"),
		ContractUrl:      utils.GetStringPropOrEmpty(props, "contractUrl"),
		ContractStatus:   contractStatus,
		RenewalCycle:     contractRenewalCycle,
		RenewalPeriods:   utils.GetInt64PropOrNil(props, "renewalPeriods"),
		Source:           neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:    neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:        utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &contract
}
