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
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/common"
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-common/gen/proto/go/api/grpc/v1/contract"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

type ContractService interface {
	Create(ctx context.Context, contract *ContractCreateData) (string, error)
	Update(ctx context.Context, contract *entity.ContractEntity) error
	GetById(ctx context.Context, id string) (*entity.ContractEntity, error)
	GetContractsForOrganizations(ctx context.Context, organizationIds []string) (*entity.ContractEntities, error)
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
	Source            entity.DataSource
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

	contractExists, _ := s.repositories.CommonRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), contract.Id, entity.NodeLabel_Contract)
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
	}
	switch contract.ContractRenewalCycle {
	case entity.ContractRenewalCycleMonthlyRenewal:
		contractUpdateRequest.RenewalCycle = contractpb.RenewalCycle_MONTHLY_RENEWAL
	case entity.ContractRenewalCycleAnnualRenewal:
		contractUpdateRequest.RenewalCycle = contractpb.RenewalCycle_ANNUALLY_RENEWAL
	default:
		contractUpdateRequest.RenewalCycle = contractpb.RenewalCycle_NONE
	}

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
	}

	switch contractDetails.ContractEntity.ContractRenewalCycle {
	case entity.ContractRenewalCycleMonthlyRenewal:
		createContractRequest.RenewalCycle = contractpb.RenewalCycle_MONTHLY_RENEWAL
	case entity.ContractRenewalCycleAnnualRenewal:
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

	response, err := s.grpcClients.ContractClient.CreateContract(ctx, &createContractRequest)

	for i := 1; i <= constants.MaxRetriesCheckDataInNeo4jAfterEventRequest; i++ {
		contractFound, findErr := s.repositories.CommonRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), response.Id, entity.NodeLabel_Contract)
		if contractFound && findErr == nil {
			span.LogFields(log.Bool("contractSavedInGraphDb", true))
			break
		}
		time.Sleep(utils.BackOffIncrementalDelay(i))
	}
	return response.Id, err
}

func (s *contractService) GetById(ctx context.Context, contractId string) (*entity.ContractEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contractId", contractId))

	if contractDbNode, err := s.repositories.ContractRepository.GetById(ctx, common.GetContext(ctx).Tenant, contractId); err != nil {
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

	contracts, err := s.repositories.ContractRepository.GetForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIDs)
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

func (s *contractService) mapDbNodeToContractEntity(dbNode dbtype.Node) *entity.ContractEntity {
	props := utils.GetPropsFromNode(dbNode)
	contractStatus := entity.GetContractStatus(utils.GetStringPropOrEmpty(props, "status"))
	contractRenewalCycle := entity.GetContractRenewalCycle(utils.GetStringPropOrEmpty(props, "renewalCycle"))

	contract := entity.ContractEntity{
		Id:                   utils.GetStringPropOrEmpty(props, "id"),
		Name:                 utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:            utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:            utils.GetTimePropOrEpochStart(props, "updatedAt"),
		ServiceStartedAt:     utils.GetTimePropOrNil(props, "serviceStartedAt"),
		SignedAt:             utils.GetTimePropOrNil(props, "signedAt"),
		EndedAt:              utils.GetTimePropOrNil(props, "endedAt"),
		ContractUrl:          utils.GetStringPropOrEmpty(props, "contractUrl"),
		ContractStatus:       contractStatus,
		ContractRenewalCycle: contractRenewalCycle,
		Source:               entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:            utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &contract
}
