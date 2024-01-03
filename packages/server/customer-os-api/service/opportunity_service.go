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
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type OpportunityService interface {
	Update(ctx context.Context, opportunity *entity.OpportunityEntity) error
	UpdateRenewal(ctx context.Context, opportunityId string, renewalLikelihood entity.OpportunityRenewalLikelihood, amount *float64, comments *string, ownerUserId *string, appSource string) error
	GetById(ctx context.Context, id string) (*entity.OpportunityEntity, error)
	GetOpportunitiesForContracts(ctx context.Context, contractIds []string) (*entity.OpportunityEntities, error)
}
type opportunityService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
	services     *Services
}

func NewOpportunityService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients, services *Services) OpportunityService {
	return &opportunityService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
		services:     services,
	}
}

func (s *opportunityService) GetById(ctx context.Context, opportunityId string) (*entity.OpportunityEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("opportunityId", opportunityId))

	if opportunityDbNode, err := s.repositories.OpportunityRepository.GetById(ctx, common.GetContext(ctx).Tenant, opportunityId); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("opportunity with id {%s} not found", opportunityId))
		return nil, wrappedErr
	} else {
		return s.mapDbNodeToOpportunityEntity(*opportunityDbNode), nil
	}
}

func (s *opportunityService) GetOpportunitiesForContracts(ctx context.Context, contractIDs []string) (*entity.OpportunityEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.GetOpportunitiesForContracts")
	defer span.Finish()
	span.LogFields(log.Object("contractIDs", contractIDs))

	opportunities, err := s.repositories.OpportunityRepository.GetForContracts(ctx, common.GetTenantFromContext(ctx), contractIDs)
	if err != nil {
		return nil, err
	}
	opportunityEntities := make(entity.OpportunityEntities, 0, len(opportunities))
	for _, v := range opportunities {
		opportunityEntity := s.mapDbNodeToOpportunityEntity(*v.Node)
		opportunityEntity.DataloaderKey = v.LinkedNodeId
		opportunityEntities = append(opportunityEntities, *opportunityEntity)
	}
	return &opportunityEntities, nil
}

func (s *opportunityService) mapDbNodeToOpportunityEntity(dbNode dbtype.Node) *entity.OpportunityEntity {
	props := utils.GetPropsFromNode(dbNode)
	opportunity := entity.OpportunityEntity{
		Id:                     utils.GetStringPropOrEmpty(props, "id"),
		Name:                   utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:              utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:              utils.GetTimePropOrEpochStart(props, "updatedAt"),
		InternalStage:          entity.GetInternalStage(utils.GetStringPropOrEmpty(props, "internalStage")),
		ExternalStage:          utils.GetStringPropOrEmpty(props, "externalStage"),
		InternalType:           entity.GetInternalType(utils.GetStringPropOrEmpty(props, "internalType")),
		ExternalType:           utils.GetStringPropOrEmpty(props, "externalType"),
		Amount:                 utils.GetFloatPropOrZero(props, "amount"),
		MaxAmount:              utils.GetFloatPropOrZero(props, "maxAmount"),
		EstimatedClosedAt:      utils.GetTimePropOrNil(props, "estimatedClosedAt"),
		NextSteps:              utils.GetStringPropOrEmpty(props, "nextSteps"),
		GeneralNotes:           utils.GetStringPropOrEmpty(props, "generalNotes"),
		RenewedAt:              utils.GetTimePropOrEpochStart(props, "renewedAt"),
		RenewalLikelihood:      entity.GetOpportunityRenewalLikelihood(utils.GetStringPropOrEmpty(props, "renewalLikelihood")),
		RenewalUpdatedByUserAt: utils.GetTimePropOrEpochStart(props, "renewalUpdatedByUserAt"),
		RenewalUpdatedByUserId: utils.GetStringPropOrEmpty(props, "renewalUpdatedByUserId"),
		Comments:               utils.GetStringPropOrEmpty(props, "comments"),
		Source:                 neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:          neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:              utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &opportunity
}

func (s *opportunityService) Update(ctx context.Context, opportunity *entity.OpportunityEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("opportunity", opportunity))

	if opportunity == nil {
		err := fmt.Errorf("(OpportunityService.Update) opportunity entity is nil")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	} else if opportunity.Id == "" {
		err := fmt.Errorf("(OpportunityService.Update) opportunity id is missing")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	opportunityExists, _ := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), opportunity.Id, entity.NodeLabel_Opportunity)
	if !opportunityExists {
		err := fmt.Errorf("(OpportunityService.Update) opportunity with id {%s} not found", opportunity.Id)
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	opportunityUpdateRequest := opportunitypb.UpdateOpportunityGrpcRequest{
		Tenant:             common.GetTenantFromContext(ctx),
		Id:                 opportunity.Id,
		LoggedInUserId:     common.GetUserIdFromContext(ctx),
		Name:               opportunity.Name,
		Amount:             opportunity.Amount,
		ExternalType:       opportunity.ExternalType,
		ExternalStage:      opportunity.ExternalStage,
		GeneralNotes:       opportunity.GeneralNotes,
		NextSteps:          opportunity.NextSteps,
		EstimatedCloseDate: utils.ConvertTimeToTimestampPtr(opportunity.EstimatedClosedAt),
		SourceFields: &commonpb.SourceFields{
			Source:    string(opportunity.Source),
			AppSource: utils.StringFirstNonEmpty(opportunity.AppSource, constants.AppSourceCustomerOsApi),
		},
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := s.grpcClients.OpportunityClient.UpdateOpportunity(ctx, &opportunityUpdateRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}

func (s *opportunityService) UpdateRenewal(ctx context.Context, opportunityId string, renewalLikelihood entity.OpportunityRenewalLikelihood, amount *float64, comments *string, ownerUserId *string, appSource string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.UpdateRenewal")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("opportunityId", opportunityId), log.Object("renewalLikelihood", renewalLikelihood), log.Object("amount", amount), log.Object("comments", comments), log.String("appSource", appSource))

	if opportunityId == "" {
		err := fmt.Errorf("(OpportunityService.UpdateRenewal) opportunity id is missing")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	opportunityExists, _ := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), opportunityId, entity.NodeLabel_Opportunity)
	if !opportunityExists {
		err := fmt.Errorf("(OpportunityService.UpdateRenewal) opportunity with id {%s} not found", opportunityId)
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	opportunityRenewalUpdateRequest := opportunitypb.UpdateRenewalOpportunityGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Id:             opportunityId,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		Amount:         utils.IfNotNilFloat64(amount),
		Comments:       utils.IfNotNilString(comments),
		OwnerUserId:    utils.IfNotNilString(ownerUserId),
		SourceFields: &commonpb.SourceFields{
			Source:    string(neo4jentity.DataSourceOpenline),
			AppSource: appSource,
		},
	}
	fieldsMask := make([]opportunitypb.OpportunityMaskField, 0)
	if amount != nil {
		fieldsMask = append(fieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_AMOUNT)
	}
	if comments != nil {
		fieldsMask = append(fieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_COMMENTS)
	}
	fieldsMask = append(fieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_RENEWAL_LIKELIHOOD)
	opportunityRenewalUpdateRequest.FieldsMask = fieldsMask

	switch renewalLikelihood {
	case entity.OpportunityRenewalLikelihoodHigh:
		opportunityRenewalUpdateRequest.RenewalLikelihood = opportunitypb.RenewalLikelihood_HIGH_RENEWAL
	case entity.OpportunityRenewalLikelihoodMedium:
		opportunityRenewalUpdateRequest.RenewalLikelihood = opportunitypb.RenewalLikelihood_MEDIUM_RENEWAL
	case entity.OpportunityRenewalLikelihoodLow:
		opportunityRenewalUpdateRequest.RenewalLikelihood = opportunitypb.RenewalLikelihood_LOW_RENEWAL
	case entity.OpportunityRenewalLikelihoodZero:
		opportunityRenewalUpdateRequest.RenewalLikelihood = opportunitypb.RenewalLikelihood_ZERO_RENEWAL
	default:
		opportunityRenewalUpdateRequest.RenewalLikelihood = opportunitypb.RenewalLikelihood_ZERO_RENEWAL
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := s.grpcClients.OpportunityClient.UpdateRenewalOpportunity(ctx, &opportunityRenewalUpdateRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}
