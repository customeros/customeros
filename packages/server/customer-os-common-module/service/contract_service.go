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
	contractpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/contract"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
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

	var beforeUpdateContractEntity *neo4jentity.ContractEntity

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
		contractDbNode, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			return contractId, err
		}
		beforeUpdateContractEntity = neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

		err = s.services.Neo4jRepositories.ContractWriteRepository.UpdateContract(ctx, tenant, contractId, dataFields)
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
	}

	// send events
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

	// post save actions
	if createFlow {
		err = s.postCreateContract(ctx, tenant, contractId, dataFields)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error while post create contract %s: %s", contractId, err.Error())
		}
	} else {
		err = s.postUpdateContract(ctx, tenant, contractId, dataFields, beforeUpdateContractEntity)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error while post create contract %s: %s", contractId, err.Error())
		}
		if dataFields.AppSource == nil || *dataFields.AppSource != constants.AppSourceCustomerOsApi {
			utils.EventCompleted(ctx, tenant, model.CONTRACT.String(), contractId, s.services.GrpcClients, utils.NewEventCompletedDetails().WithUpdate())
		}
	}

	return contractId, nil
}

func (s *contractService) postCreateContract(ctx context.Context, tenant, contractId string, dataFields data_fields.ContractSaveFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.postCreateContract")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	_, _, err := s.updateStatus(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error while updating contract %s status: %s", contractId, err.Error())
	}

	if dataFields.LengthInMonths != nil && *dataFields.LengthInMonths > 0 {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err := utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
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

	return nil
}

func (s *contractService) postUpdateContract(ctx context.Context, tenant string, contractId string, dataFields data_fields.ContractSaveFields, beforeUpdateContractEntity *neo4jentity.ContractEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.postCreateContract")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	_, statusChanged, err := s.updateStatus(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error while updating contract %s status: %s", contractId, err.Error())
	}

	if statusChanged {
		err = s.updateOrganizationRelationship(ctx, tenant, contractId, statusChanged)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error while updating organization relationship for contract %s: %s", contractId, err.Error())
		}
	}
	contractDbNode, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	afterUpdateContractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	if beforeUpdateContractEntity.LengthInMonths > 0 && afterUpdateContractEntity.LengthInMonths == 0 {
		err = s.services.Neo4jRepositories.ContractWriteRepository.SuspendActiveRenewalOpportunity(ctx, tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error while suspending renewal opportunity for contract %s: %s", contractId, err.Error())
		}
		organizationDbNode, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByContractId(ctx, tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error while getting organization for contract %s: %s", contractId, err.Error())
			return nil
		}
		if organizationDbNode == nil {
			s.log.Errorf("Organization not found for contract %s", contractId)
			return nil
		}
		organization := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return s.services.GrpcClients.OrganizationClient.RefreshRenewalSummary(ctx, &organizationpb.RefreshRenewalSummaryGrpcRequest{
				Tenant:         tenant,
				OrganizationId: organization.ID,
				AppSource:      common.GetAppSourceFromContext(ctx),
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("RefreshRenewalSummary failed: %v", err.Error())
		}
		_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
			return s.services.GrpcClients.OrganizationClient.RefreshArr(ctx, &organizationpb.OrganizationIdGrpcRequest{
				Tenant:         tenant,
				OrganizationId: organization.ID,
				AppSource:      common.GetAppSourceFromContext(ctx),
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("RefreshArr failed: %v", err.Error())
		}
	} else {
		if beforeUpdateContractEntity.LengthInMonths == 0 && afterUpdateContractEntity.LengthInMonths > 0 {
			err = s.services.Neo4jRepositories.ContractWriteRepository.ActivateSuspendedRenewalOpportunity(ctx, tenant, contractId)
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error while activating renewal opportunity for contract %s: %s", contractId, err.Error())
			}
		}
		err = s.updateActiveRenewalOpportunityRenewDateAndArr(ctx, tenant, contractId)
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("error while updating renewal opportunity for contract %s: %s", contractId, err.Error())
		}
	}

	if beforeUpdateContractEntity.ContractStatus != afterUpdateContractEntity.ContractStatus {
		s.createActionForStatusChange(ctx, tenant, contractId, string(afterUpdateContractEntity.ContractStatus), afterUpdateContractEntity.Name)
	}

	err = s.updateActiveRenewalOpportunityLikelihood(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error while updating renewal opportunity for contract %s: %s", contractId, err.Error())
	}
	s.updateContractLtv(ctx, tenant, contractId)
	return nil
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

func (s *contractService) updateOrganizationRelationship(ctx context.Context, tenant, contractId string, statusChanged bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.UpdateOrganizationRelationship")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("contractId", contractId), log.Bool("statusChanged", statusChanged))

	if !statusChanged {
		return nil
	}

	// get contract
	contractDbNode, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error while getting contract %s: %s", contractId, err.Error())
		return err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	// get organization for contract
	organizationDbNode, err := s.services.Neo4jRepositories.OrganizationReadRepository.GetOrganizationByContractId(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error while getting organization for contract %s: %s", contractId, err.Error())
		return err
	}
	orgEntity := neo4jmapper.MapDbNodeToOrganizationEntity(organizationDbNode)

	// get all contracts for organization
	orgContracts, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractsForOrganizations(ctx, tenant, []string{orgEntity.ID})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error while getting contracts for organization %s: %s", orgEntity.ID, err.Error())
		return err
	}
	orgContractEntities := []neo4jentity.ContractEntity{}
	for _, orgContract := range orgContracts {
		orgContractEntities = append(orgContractEntities, *neo4jmapper.MapDbNodeToContractEntity(orgContract.Node))
	}

	if contractEntity.ContractStatus == neo4jenum.ContractStatusEnded {
		// check no other contract is active
		activeContractFound := false
		for _, orgContract := range orgContractEntities {
			if orgContract.ContractStatus != neo4jenum.ContractStatusDraft && orgContract.ContractStatus != neo4jenum.ContractStatusEnded {
				activeContractFound = true
				break
			}
		}

		if !activeContractFound {
			_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
				return s.services.GrpcClients.OrganizationClient.UpdateOrganization(ctx, &organizationpb.UpdateOrganizationGrpcRequest{
					Tenant:         tenant,
					OrganizationId: orgEntity.ID,
					Relationship:   neo4jenum.FormerCustomer.String(),
					Stage:          neo4jenum.Target.String(),
					FieldsMask: []organizationpb.OrganizationMaskField{
						organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_RELATIONSHIP,
						organizationpb.OrganizationMaskField_ORGANIZATION_PROPERTY_STAGE,
					},
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("UpdateOrganization failed: %s", err.Error())
				return errors.Wrap(err, "UpdateOrganization")
			}
			_, err = utils.CallEventsPlatformGRPCWithRetry[*organizationpb.OrganizationIdGrpcResponse](func() (*organizationpb.OrganizationIdGrpcResponse, error) {
				return s.services.GrpcClients.OrganizationClient.RefreshDerivedData(ctx, &organizationpb.RefreshDerivedDataGrpcRequest{
					Tenant:         tenant,
					OrganizationId: orgEntity.ID,
					AppSource:      common.GetAppSourceFromContext(ctx),
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("RefreshDerivedData failed: %s", err.Error())
				return errors.Wrap(err, "RefreshDerivedData")
			}
		}
	}

	return nil
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

func (s *contractService) updateActiveRenewalOpportunityRenewDateAndArr(ctx context.Context, tenant, contractId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.updateActiveRenewalOpportunityRenewDateAndArr")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("contractId", contractId))

	contract, renewalOpportunity, done := s.assertContractAndRenewalOpportunity(ctx, tenant, contractId)
	if done {
		return nil
	}

	err := s.updateRenewalOpportunityRenewedAt(ctx, tenant, contract, renewalOpportunity)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	err = s.updateRenewalArr(ctx, tenant, contract, renewalOpportunity, span)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	return nil
}

func (s *contractService) assertContractAndRenewalOpportunity(ctx context.Context, tenant, contractId string) (*neo4jentity.ContractEntity, *neo4jentity.OpportunityEntity, bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.assertContractAndRenewalOpportunity")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("contractId", contractId))

	contractDbNode, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error while getting contract %s: %s", contractId, err.Error())
		return nil, nil, true
	}
	contract := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)

	// if contract is not frequency based, return
	if contract.LengthInMonths == 0 {
		return nil, nil, true
	}

	currentRenewalOpportunityDbNode, err := s.services.Neo4jRepositories.OpportunityReadRepository.GetActiveRenewalOpportunityForContract(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error while getting renewal opportunity for contract %s: %s", contractId, err.Error())
		return nil, nil, true
	}

	// if there is no renewal opportunity, create one
	if currentRenewalOpportunityDbNode == nil {
		if !contract.IsEnded() {
			ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
			_, err = utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
				return s.services.GrpcClients.OpportunityClient.CreateRenewalOpportunity(ctx, &opportunitypb.CreateRenewalOpportunityGrpcRequest{
					Tenant:     tenant,
					ContractId: contractId,
					SourceFields: &commonpb.SourceFields{
						AppSource: common.GetAppSourceFromContext(ctx),
					},
				})
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("CreateRenewalOpportunity command failed: %v", err.Error())
				return nil, nil, true
			}
			span.LogFields(log.Bool("renewal opportunity create requested", true))
		}
		return nil, nil, true
	}

	currentRenewalOpportunity := neo4jmapper.MapDbNodeToOpportunityEntity(currentRenewalOpportunityDbNode)

	return contract, currentRenewalOpportunity, false
}

func (s *contractService) updateRenewalOpportunityRenewedAt(ctx context.Context, tenant string, contractEntity *neo4jentity.ContractEntity, renewalOpportunityEntity *neo4jentity.OpportunityEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.updateRenewalOpportunityRenewedAt")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)

	if renewalOpportunityEntity == nil {
		err := fmt.Errorf("renewalOpportunityEntity is nil")
		tracing.TraceErr(span, err)
		return nil
	}

	// IF contract already ended, close the renewal opportunity
	if contractEntity.IsEnded() {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err := utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
			return s.services.GrpcClients.OpportunityClient.CloseLooseOpportunity(ctx, &opportunitypb.CloseLooseOpportunityGrpcRequest{
				Tenant:    tenant,
				Id:        renewalOpportunityEntity.Id,
				AppSource: common.GetAppSourceFromContext(ctx),
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("CloseLooseOpportunity failed: %s", err.Error())
			return errors.Wrap(err, "CloseLooseOpportunity")
		}
		return nil
	}

	// Choose starting date for renewal calculation
	if renewalOpportunityEntity.RenewalDetails.RenewedAt != nil {
		return nil
	}
	startRenewalDateCalculation := contractEntity.ServiceStartedAt
	previousClosedWonRenewalDbNode, err := s.services.Neo4jRepositories.OpportunityReadRepository.GetPreviousClosedWonRenewalOpportunityForContract(ctx, tenant, contractEntity.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil
	}
	if previousClosedWonRenewalDbNode != nil {
		previousRenewalOpportunityEntity := neo4jmapper.MapDbNodeToOpportunityEntity(previousClosedWonRenewalDbNode)
		if previousRenewalOpportunityEntity.RenewalDetails.RenewedAt != nil {
			startRenewalDateCalculation = previousRenewalOpportunityEntity.RenewalDetails.RenewedAt
		}
	}
	span.LogFields(log.Object("startRenewalDateCalculation", startRenewalDateCalculation))

	// Calculate until first future date if auto-renew is enabled or renewal is approved
	calculateUntilFirstFutureDate := contractEntity.AutoRenew
	span.LogFields(log.Bool("calculateUntilFirstFutureDate", calculateUntilFirstFutureDate))

	renewedAt := calculateNextCycleDate(startRenewalDateCalculation, contractEntity.LengthInMonths, calculateUntilFirstFutureDate)
	span.LogFields(log.Object("result.renewedAt", renewedAt))
	if !utils.IsEqualTimePtr(renewedAt, renewalOpportunityEntity.RenewalDetails.RenewedAt) {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err := utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
			return s.services.GrpcClients.OpportunityClient.UpdateRenewalOpportunityNextCycleDate(ctx, &opportunitypb.UpdateRenewalOpportunityNextCycleDateGrpcRequest{
				OpportunityId: renewalOpportunityEntity.Id,
				Tenant:        tenant,
				AppSource:     common.GetAppSourceFromContext(ctx),
				RenewedAt:     utils.ConvertTimeToTimestampPtr(renewedAt),
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("UpdateRenewalOpportunityNextCycleDate failed: %s", err.Error())
			return errors.Wrap(err, "UpdateRenewalOpportunityNextCycleDate")
		}
	}

	return nil
}

func (s *contractService) updateRenewalArr(ctx context.Context, tenant string, contract *neo4jentity.ContractEntity, renewalOpportunity *neo4jentity.OpportunityEntity, span opentracing.Span) error {
	// if contract already ended, return
	if contract.IsEnded() {
		span.LogFields(log.Bool("contract ended", true))
		return nil
	}

	maxArr, err := s.calculateMaxArr(ctx, tenant, contract, renewalOpportunity, span)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error while calculating ARR for contract %s: %s", contract.Id, err.Error())
		return nil
	}
	// adjust with likelihood
	currentArr := calculateCurrentArrByAdjustedRate(maxArr, renewalOpportunity.RenewalDetails.RenewalAdjustedRate)

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
		return s.services.GrpcClients.OpportunityClient.UpdateOpportunity(ctx, &opportunitypb.UpdateOpportunityGrpcRequest{
			Tenant:    tenant,
			Id:        renewalOpportunity.Id,
			Amount:    currentArr,
			MaxAmount: maxArr,
			SourceFields: &commonpb.SourceFields{
				AppSource: common.GetAppSourceFromContext(ctx),
				Source:    neo4jentity.DataSourceOpenline.String(),
			},
			FieldsMask: []opportunitypb.OpportunityMaskField{
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_AMOUNT,
				opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_MAX_AMOUNT,
			},
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("UpdateOpportunity failed: %s", err.Error())
	}

	return nil
}

func (s *contractService) calculateMaxArr(ctx context.Context, tenant string, contract *neo4jentity.ContractEntity, renewalOpportunity *neo4jentity.OpportunityEntity, span opentracing.Span) (float64, error) {
	var arr float64

	// Fetch service line items for the contract from the database
	sliDbNodes, err := s.services.Neo4jRepositories.ServiceLineItemReadRepository.GetServiceLineItemsForContract(ctx, tenant, contract.Id)
	if err != nil {
		tracing.TraceErr(span, err)
		return 0, err
	}
	serviceLineItems := neo4jentity.ServiceLineItemEntities{}
	for _, sliDbNode := range sliDbNodes {
		sli := neo4jmapper.MapDbNodeToServiceLineItemEntity(sliDbNode)
		serviceLineItems = append(serviceLineItems, *sli)
	}

	span.LogFields(log.Int("service line items count", len(serviceLineItems)))
	for _, sli := range serviceLineItems {
		if sli.IsEnded() {
			span.LogFields(log.Bool(fmt.Sprintf("service line item {%s} ended", sli.ID), true))
			continue
		}
		span.LogFields(log.Object(fmt.Sprintf("service line item {%s}:", sli.ID), sli))
		annualPrice := float64(0)
		if sli.Billed == neo4jenum.BilledTypeAnnually {
			annualPrice = float64(sli.Price) * float64(sli.Quantity)
		} else if sli.Billed == neo4jenum.BilledTypeMonthly {
			annualPrice = float64(sli.Price) * float64(sli.Quantity)
			annualPrice *= 12
		} else if sli.Billed == neo4jenum.BilledTypeQuarterly {
			annualPrice = float64(sli.Price) * float64(sli.Quantity)
			annualPrice *= 4
		}
		span.LogFields(log.Float64(fmt.Sprintf("service line item {%s} added ARR value:", sli.ID), annualPrice))
		// Add to total ARR
		arr += annualPrice
	}

	// Adjust with end date
	if contract.EndedAt != nil {
		span.LogFields(log.Bool("ARR prorated with contract end date", true))
		arr = prorateArr(arr, monthsUntilContractEnd(utils.Now(), *contract.EndedAt))
	}

	return utils.RoundHalfUpFloat64(arr, 2), nil
}

func calculateCurrentArrByAdjustedRate(maxAmount float64, rate int64) float64 {
	if rate == 0 {
		return 0
	} else if rate == 100 {
		return maxAmount
	}
	return utils.RoundHalfUpFloat64(maxAmount*float64(rate)/100, 2)
}

type ActionStatusMetadata struct {
	Status       string `json:"status"`
	ContractName string `json:"contract-name"`
	Comment      string `json:"comment"`
}

func (s *contractService) createActionForStatusChange(ctx context.Context, tenant, contractId, status, contractName string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "contractService.createActionForStatusChange")
	defer span.Finish()
	var name string
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("contractId", contractId), log.String("status", status), log.String("contractName", contractName))

	if contractName != "" {
		name = contractName
	} else {
		name = "Unnamed contract"
	}
	actionStatusMetadata := ActionStatusMetadata{
		Status:       status,
		ContractName: name,
		Comment:      name + " is now " + status,
	}
	message := ""

	switch status {
	case string(neo4jenum.ContractStatusLive):
		message = contractName + " is now live"
		actionStatusMetadata.Comment = contractName + " is now live"
	case string(neo4jenum.ContractStatusEnded):
		message = contractName + " has ended"
		actionStatusMetadata.Comment = contractName + " has ended"
	case string(neo4jenum.ContractStatusOutOfContract):
		message = contractName + " is now out of contract"
		actionStatusMetadata.Comment = contractName + " is now out of contract"
	}
	metadata, err := utils.ToJson(actionStatusMetadata)
	_, err = s.services.Neo4jRepositories.ActionWriteRepository.Create(ctx, tenant, contractId, model.CONTRACT, neo4jenum.ActionContractStatusUpdated, message, metadata, utils.Now(), common.GetAppSourceFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Failed creating status update action for contract %s: %s", contractId, err.Error())
	}
}

func (s *contractService) updateActiveRenewalOpportunityLikelihood(ctx context.Context, tenant, contractId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.updateActiveRenewalOpportunityLikelihood")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("contractId", contractId))

	opportunityDbNode, err := s.services.Neo4jRepositories.OpportunityReadRepository.GetActiveRenewalOpportunityForContract(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error while getting renewal opportunity for contract %s: %s", contractId, err.Error())
		return err
	}
	if opportunityDbNode == nil {
		s.log.Infof("No open renewal opportunity found for contract %s", contractId)
		return nil
	}
	contractDbNode, err := s.services.Neo4jRepositories.ContractReadRepository.GetContractById(ctx, tenant, contractId)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error while getting contract %s: %s", contractId, err.Error())
		return err
	}
	contractEntity := neo4jmapper.MapDbNodeToContractEntity(contractDbNode)
	opportunityEntity := neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode)

	var renewalLikelihood neo4jenum.RenewalLikelihood
	renewalAdjustedRate := opportunityEntity.RenewalDetails.RenewalAdjustedRate
	if contractEntity.EndedAt != nil &&
		opportunityEntity.RenewalDetails.RenewalLikelihood != neo4jenum.RenewalLikelihoodZero &&
		opportunityEntity.RenewalDetails.RenewedAt != nil &&
		contractEntity.EndedAt.Before(*opportunityEntity.RenewalDetails.RenewedAt) {
		// check if likelihood should be set to Zero
		renewalLikelihood = neo4jenum.RenewalLikelihoodZero
		renewalAdjustedRate = int64(0)
	} else if opportunityEntity.RenewalDetails.RenewalLikelihood == neo4jenum.RenewalLikelihoodZero &&
		opportunityEntity.RenewalDetails.RenewedAt != nil &&
		(contractEntity.EndedAt == nil || contractEntity.EndedAt.After(*opportunityEntity.RenewalDetails.RenewedAt)) {
		// check if likelihood should be set to Medium
		renewalLikelihood = neo4jenum.RenewalLikelihoodMedium
		renewalAdjustedRate = int64(50)
	}

	if renewalLikelihood != "" {
		ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
			return s.services.GrpcClients.OpportunityClient.UpdateRenewalOpportunity(ctx, &opportunitypb.UpdateRenewalOpportunityGrpcRequest{
				Tenant:              tenant,
				Id:                  opportunityEntity.Id,
				RenewalLikelihood:   renewalLikelihoodForGrpcRequest(renewalLikelihood),
				RenewalAdjustedRate: renewalAdjustedRate,
				SourceFields: &commonpb.SourceFields{
					AppSource: common.GetAppSourceFromContext(ctx),
				},
				FieldsMask: []opportunitypb.OpportunityMaskField{
					opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_RENEWAL_LIKELIHOOD,
					opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_ADJUSTED_RATE,
				},
			})
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("UpdateRenewalOpportunity failed: %s", err.Error())
			return errors.Wrap(err, "UpdateRenewalOpportunity")
		}
	}

	return nil
}

func (s *contractService) updateContractLtv(ctx context.Context, tenant, contractId string) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.updateContractLtv")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.SetTag(tracing.SpanTagEntityId, contractId)

	// request contract LTV refresh
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := utils.CallEventsPlatformGRPCWithRetry[*contractpb.ContractIdGrpcResponse](func() (*contractpb.ContractIdGrpcResponse, error) {
		return s.services.GrpcClients.ContractClient.RefreshContractLtv(ctx, &contractpb.RefreshContractLtvGrpcRequest{
			Tenant:    tenant,
			Id:        contractId,
			AppSource: common.GetAppSourceFromContext(ctx),
		})
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("RefreshContractLtv failed: %s", err.Error())
	}
}

func renewalLikelihoodForGrpcRequest(renewalLikelihood neo4jenum.RenewalLikelihood) opportunitypb.RenewalLikelihood {
	switch renewalLikelihood {
	case neo4jenum.RenewalLikelihoodHigh:
		return opportunitypb.RenewalLikelihood_HIGH_RENEWAL
	case neo4jenum.RenewalLikelihoodMedium:
		return opportunitypb.RenewalLikelihood_MEDIUM_RENEWAL
	case neo4jenum.RenewalLikelihoodLow:
		return opportunitypb.RenewalLikelihood_LOW_RENEWAL
	case neo4jenum.RenewalLikelihoodZero:
		return opportunitypb.RenewalLikelihood_ZERO_RENEWAL
	default:
		return opportunitypb.RenewalLikelihood_HIGH_RENEWAL
	}
}

func calculateNextCycleDate(from *time.Time, lengthInMonths int64, calculateUntilFirstFutureDate bool) *time.Time {
	if from == nil || lengthInMonths <= 0 {
		return nil
	}

	renewalCycleNext := *from
	for {
		renewalCycleNext = renewalCycleNext.AddDate(0, int(lengthInMonths), 0)
		// Break the loop either when the next cycle date is in the future
		// or if we are not calculating until the first future date.
		if renewalCycleNext.After(utils.Now()) || !calculateUntilFirstFutureDate {
			break
		}
	}
	return &renewalCycleNext
}

func prorateArr(arr float64, monthsRemaining int) float64 {
	if monthsRemaining > 12 {
		return arr
	}
	monthlyRate := arr / 12
	return utils.RoundHalfUpFloat64(monthlyRate*float64(monthsRemaining), 2)
}

func monthsUntilContractEnd(start, end time.Time) int {
	yearDiff := end.Year() - start.Year()
	monthDiff := int(end.Month()) - int(start.Month())

	// Total difference in months
	totalMonths := yearDiff*12 + monthDiff

	// If the end day is before the start day in the month, subtract a month
	if end.Day() < start.Day() {
		totalMonths--
	}

	if totalMonths < 0 {
		totalMonths = 0
	}

	return totalMonths
}
