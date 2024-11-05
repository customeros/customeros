package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jmodel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type ContractService interface {
	GetById(ctx context.Context, id string) (*neo4jentity.ContractEntity, error)
	Save(ctx context.Context, id *string, dataFields data_fields.ContractSaveFields) (string, error)
}

type contractService struct {
	log      logger.Logger
	services *Services
}

func NewContractService(log logger.Logger, services *Services) ContractService {
	return &contractService{
		log:      log,
		services: services,
	}
}

func (s *contractService) GetById(ctx context.Context, contractId string) (*neo4jentity.ContractEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("contractId", contractId))

	if contractDbNode, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, common.GetContext(ctx).Tenant, contractId); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Contract with id {%s} not found", contractId))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToContractEntity(contractDbNode), nil
	}
}

func (s *contractService) Save(ctx context.Context, id *string, dataFields data_fields.ContractSaveFields) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.Save")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "dataFields", dataFields)

	// validate tenant
	err := common.ValidateTenant(ctx)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	tenant := common.GetTenantFromContext(ctx)

	createFlow := false
	contractId := ""

	if id == nil || *id == "" {
		createFlow = true
		span.LogKV("flow", "create")
		contractId, err = s.services.Neo4jRepositories.CommonReadRepository.GenerateId(ctx, tenant, model.NodeLabelContract)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		// set missing fields for create flow
		if dataFields.CreatedAt == nil {
			dataFields.CreatedAt = utils.NowPtr()
		} else {
			dataFields.CreatedAt = utils.TimePtr(utils.NowIfZero(*dataFields.CreatedAt))
		}
		if utils.IfNotNilString(dataFields.AppSource) == "" {
			dataFields.AppSource = utils.StringPtr(common.GetAppSourceFromContext(ctx))
		}
		if utils.IfNotNilString(dataFields.Source) == "" {
			dataFields.Source = utils.StringPtr(neo4jentity.DataSourceOpenline.String())
		}
	} else {
		span.LogKV("flow", "update")
		contractId = *id

		// validate contract exists
		exists, err := s.services.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, tenant, contractId, model.NodeLabelContract)
		if err != nil || !exists {
			err = errors.New("contract not found")
			tracing.TraceErr(span, err)
			return "", err
		}
	}
	tracing.TagEntity(span, contractId)

	if createFlow {
		err := s.services.Neo4jRepositories.ContractWriteRepository.CreateForOrganization(ctx, tenant, contractId, dataFields)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}

		if dataFields.ExternalSystem != nil && dataFields.ExternalSystem.Available() {
			externalSystemData := neo4jmodel.ExternalSystem{
				ExternalSystemId: dataFields.ExternalSystem.ExternalSystemId,
				ExternalUrl:      dataFields.ExternalSystem.ExternalUrl,
				ExternalId:       dataFields.ExternalSystem.ExternalId,
				ExternalIdSecond: dataFields.ExternalSystem.ExternalIdSecond,
				ExternalSource:   dataFields.ExternalSystem.ExternalSource,
				SyncDate:         dataFields.ExternalSystem.SyncDate,
			}
			err = s.services.Neo4jRepositories.ExternalSystemWriteRepository.LinkWithEntity(ctx, tenant, contractId, model.NodeLabelContract, externalSystemData)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error while linking contract %s with external system %s: %s", contractId, dataFields.ExternalSystem.ExternalSystemId, err.Error())
				return "", err
			}
		}
	} else {
		err := s.services.Neo4jRepositories.ContractWriteRepository.UpdateContract(ctx, tenant, contractId, dataFields)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
	}

	// send create events
	if createFlow {
		err = s.services.RabbitMQService.PublishEvent(ctx, contractId, model.CONTRACT, dto.CreateContract{dataFields})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message CreateContract"))
		}
		utils.EventCompleted(ctx, tenant, model.CONTRACT.String(), contractId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithCreate())
	} else {
		err = s.services.RabbitMQService.PublishEvent(ctx, contractId, model.CONTRACT, dto.UpdateContract{dataFields})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message UpdateContractOld"))
		}
		if dataFields.AppSource == nil || *dataFields.AppSource != constants.AppSourceCustomerOsApi {
			utils.EventCompleted(ctx, tenant, model.CONTRACT.String(), contractId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
		}
	}

	if createFlow {
		_, _, err = s.updateStatus(ctx, tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error while updating contract %s status: %s", contractId, err.Error())
			return contractId, err
		}

		if dataFields.LengthInMonths != nil && *dataFields.LengthInMonths > 0 {
			ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
			_, err = utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
				return s.services.GrpcClients.OpportunityClient.CreateRenewalOpportunity(ctx, &opportunitypb.CreateRenewalOpportunityGrpcRequest{
					Tenant:     tenant,
					ContractId: contractId,
					SourceFields: &commonpb.SourceFields{
						Source:    *dataFields.Source,
						AppSource: *dataFields.AppSource,
					},
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("CreateRenewalOpportunity failed: %s", err.Error())
			}
		}
	}

	return contractId, nil
}

func (s *contractService) updateStatus(ctx context.Context, tenant, contractId string) (string, bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.updateStatus")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	contractDbNode, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error while getting contract %s: %s", contractId, err.Error())
		return "", false, err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	status, err := s.deriveContractStatus(ctx, tenant, *contractEntity)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error while deriving contract %s status: %s", contractId, err.Error())
		return "", false, err
	}
	statusChanged := contractEntity.ContractStatus.String() != status

	if statusChanged {
		err = s.services.Neo4jRepositories.ContractWriteRepository.UpdateStatus(ctx, tenant, contractId, status)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error while updating contract %s status: %s", contractId, err.Error())
			return "", false, err
		}

		// TODO add event for status change
		utils.EventCompleted(ctx, tenant, model.CONTRACT.String(), contractId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())

		err = s.services.RabbitMQService.PublishEvent(ctx, contractId, model.CONTRACT, dto.ChangeStatusForContract{Status: status})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "unable to publish message ChangeStatusForContract"))
		}
	}

	return status, statusChanged, nil
}

func (s *contractService) deriveContractStatus(ctx context.Context, tenant string, contractEntity neo4jentity.ContractEntity) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.deriveContractStatus")
	defer span.Finish()

	now := utils.Now()

	// If endedAt is not nil and is in the past, the contract is considered Ended.
	if contractEntity.IsEnded() {
		span.LogFields(log.String("result.status", neo4jenum.ContractStatusEnded.String()))
		return neo4jenum.ContractStatusEnded.String(), nil
	}

	// check if contract is draft
	if !contractEntity.Approved {
		span.LogFields(log.String("result.status", neo4jenum.ContractStatusDraft.String()))
		return neo4jenum.ContractStatusDraft.String(), nil
	}

	// Check contract is scheduled
	if contractEntity.ServiceStartedAt == nil || contractEntity.ServiceStartedAt.After(now) {
		span.LogFields(log.String("result.status", neo4jenum.ContractStatusScheduled.String()))
		return neo4jenum.ContractStatusScheduled.String(), nil
	}

	// Check if contract is out of contract
	if !contractEntity.AutoRenew {
		// fetch active renewal opportunity for the contract
		opportunityDbNode, err := s.services.Neo4jRepositories.OpportunityReadRepository.GetActiveRenewalOpportunityForContract(ctx, tenant, contractEntity.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			return "", err
		}
		if opportunityDbNode != nil {
			opportunityEntity := neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode)
			if opportunityEntity.RenewalDetails.RenewedAt != nil && opportunityEntity.RenewalDetails.RenewedAt.Before(now) {
				span.LogFields(log.String("result.status", neo4jenum.ContractStatusLive.String()))
				return neo4jenum.ContractStatusOutOfContract.String(), nil
			}
		}
	}

	// Otherwise, the contract is considered Live.
	span.LogFields(log.String("result.status", neo4jenum.ContractStatusLive.String()))
	return neo4jenum.ContractStatusLive.String(), nil
}

func (s *contractService) startOnboardingIfEligible(ctx context.Context, tenant, contractId string, span opentracing.Span) {

	// TODO temporary not eligible for all contracts
	return

	contractDbNode, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}
	if contractDbNode == nil {
		return
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	if contractEntity.IsEligibleToStartOnboarding() {
		organizationDbNode, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByContractId(ctx, tenant, contractEntity.Id)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error while getting organization for contract %s: %s", contractEntity.Id, err.Error())
			return
		}
		if organizationDbNode == nil {
			return
		}
		organization := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return s.services.GrpcClients.OrganizationClient.UpdateOnboardingStatus(ctx, &organizationpb.UpdateOnboardingStatusGrpcRequest{
				Tenant:             tenant,
				OrganizationId:     organization.ID,
				CausedByContractId: contractEntity.Id,
				OnboardingStatus:   organizationpb.OnboardingStatus_ONBOARDING_STATUS_NOT_STARTED,
				AppSource:          constants.AppSourceCustomerOsApi,
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("UpdateOnboardingStatus gRPC request failed: %v", err.Error())
		}
	}
}
