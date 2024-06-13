package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type OpportunityService interface {
	Create(ctx context.Context, input model.OpportunityCreateInput) (string, error)
	Update(ctx context.Context, input model.OpportunityUpdateInput) error
	UpdateRenewal(ctx context.Context, opportunityId string, renewalLikelihood neo4jenum.RenewalLikelihood, amount *float64, comments *string, ownerUserId *string, adjustedRate *int64, appSource string) error
	GetById(ctx context.Context, id string) (*neo4jentity.OpportunityEntity, error)
	GetOpportunitiesForContracts(ctx context.Context, contractIds []string) (*neo4jentity.OpportunityEntities, error)
	GetOpportunitiesForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.OpportunityEntities, error)
	UpdateRenewalsForOrganization(ctx context.Context, organizationId string, renewalLikelihood neo4jenum.RenewalLikelihood, renewalAdjustedRate *int64) error
	GetPaginatedOrganizationOpportunities(ctx context.Context, page int, limit int) (*utils.Pagination, error)
	CloseWon(ctx context.Context, opportunityId string) error
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

func (s *opportunityService) Create(ctx context.Context, input model.OpportunityCreateInput) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.Create")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	// check organization exists
	orgFound, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), input.OrganizationID, neo4jutil.NodeLabelOrganization)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error checking organization with id {%s} exists: %s", input.OrganizationID, err.Error())
		return "", err
	}
	if !orgFound {
		err := fmt.Errorf("organization with id {%s} not found", input.OrganizationID)
		s.log.Errorf(err.Error())
		tracing.TraceErr(span, err)
		return "", err
	}

	opportunityCreateRequest := opportunitypb.CreateOpportunityGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		OrganizationId: input.OrganizationID,
		Name:           utils.IfNotNilString(input.Name),
		ExternalType:   utils.IfNotNilString(input.ExternalType),
		ExternalStage:  utils.IfNotNilString(input.ExternalStage),
		GeneralNotes:   utils.IfNotNilString(input.GeneralNotes),
		NextSteps:      utils.IfNotNilString(input.NextSteps),
		SourceFields: &commonpb.SourceFields{
			Source:    string(neo4jentity.DataSourceOpenline),
			AppSource: constants.AppSourceCustomerOsApi,
		},
	}
	if input.EstimatedClosedDate != nil {
		opportunityCreateRequest.EstimatedCloseDate = utils.ConvertTimeToTimestampPtr(input.EstimatedClosedDate)
	}
	if input.InternalType != nil {
		switch *input.InternalType {
		case model.InternalTypeNbo:
			opportunityCreateRequest.InternalType = opportunitypb.OpportunityInternalType_NBO
		case model.InternalTypeUpsell:
			opportunityCreateRequest.InternalType = opportunitypb.OpportunityInternalType_UPSELL
		case model.InternalTypeCrossSell:
			opportunityCreateRequest.InternalType = opportunitypb.OpportunityInternalType_CROSS_SELL
		}
	} else {
		opportunityCreateRequest.InternalType = opportunitypb.OpportunityInternalType_NBO
	}
	opportunityCreateRequest.InternalStage = opportunitypb.OpportunityInternalStage_OPEN

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	opportunityIdGrpcResponse, err := utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
		return s.grpcClients.OpportunityClient.CreateOpportunity(ctx, &opportunityCreateRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, opportunityIdGrpcResponse.Id, neo4jutil.NodeLabelOpportunity, span)

	return opportunityIdGrpcResponse.Id, nil
}

func (s *opportunityService) GetById(ctx context.Context, opportunityId string) (*neo4jentity.OpportunityEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("opportunityId", opportunityId))

	if opportunityDbNode, err := s.repositories.Neo4jRepositories.OpportunityReadRepository.GetOpportunityById(ctx, common.GetContext(ctx).Tenant, opportunityId); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("opportunity with id {%s} not found", opportunityId))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode), nil
	}
}

func (s *opportunityService) GetOpportunitiesForContracts(ctx context.Context, contractIDs []string) (*neo4jentity.OpportunityEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.GetOpportunitiesForContracts")
	defer span.Finish()
	span.LogFields(log.Object("contractIDs", contractIDs))

	opportunities, err := s.repositories.Neo4jRepositories.OpportunityReadRepository.GetForContracts(ctx, common.GetTenantFromContext(ctx), contractIDs)
	if err != nil {
		return nil, err
	}
	opportunityEntities := make(neo4jentity.OpportunityEntities, 0, len(opportunities))
	for _, v := range opportunities {
		opportunityEntity := neo4jmapper.MapDbNodeToOpportunityEntity(v.Node)
		opportunityEntity.DataloaderKey = v.LinkedNodeId
		opportunityEntities = append(opportunityEntities, *opportunityEntity)
	}
	return &opportunityEntities, nil
}

func (s *opportunityService) GetOpportunitiesForOrganizations(ctx context.Context, organizationIds []string) (*neo4jentity.OpportunityEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.GetOpportunitiesForOrganizations")
	defer span.Finish()
	span.LogFields(log.Object("organizationIds", organizationIds))

	opportunities, err := s.repositories.Neo4jRepositories.OpportunityReadRepository.GetForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIds)
	if err != nil {
		return nil, err
	}
	opportunityEntities := make(neo4jentity.OpportunityEntities, 0, len(opportunities))
	for _, v := range opportunities {
		opportunityEntity := neo4jmapper.MapDbNodeToOpportunityEntity(v.Node)
		opportunityEntity.DataloaderKey = v.LinkedNodeId
		opportunityEntities = append(opportunityEntities, *opportunityEntity)
	}
	return &opportunityEntities, nil
}

func (s *opportunityService) Update(ctx context.Context, input model.OpportunityUpdateInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	opportunityExists, _ := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), input.OpportunityID, neo4jutil.NodeLabelOpportunity)
	if !opportunityExists {
		err := fmt.Errorf("(OpportunityService.Update) opportunity with id {%s} not found", input.OpportunityID)
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	fieldsMask := make([]opportunitypb.OpportunityMaskField, 0)
	opportunityUpdateRequest := opportunitypb.UpdateOpportunityGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Id:             input.OpportunityID,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: constants.AppSourceCustomerOsApi,
		},
	}
	if input.Name != nil {
		opportunityUpdateRequest.Name = *input.Name
		fieldsMask = append(fieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_NAME)
	}
	if input.Amount != nil {
		opportunityUpdateRequest.Amount = *input.Amount
		fieldsMask = append(fieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_AMOUNT)
	}
	if input.ExternalType != nil {
		opportunityUpdateRequest.ExternalType = *input.ExternalType
		fieldsMask = append(fieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_EXTERNAL_TYPE)
	}
	if input.ExternalStage != nil {
		opportunityUpdateRequest.ExternalStage = *input.ExternalStage
		fieldsMask = append(fieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_EXTERNAL_STAGE)
	}
	if input.EstimatedClosedDate != nil {
		opportunityUpdateRequest.EstimatedCloseDate = utils.ConvertTimeToTimestampPtr(input.EstimatedClosedDate)
		fieldsMask = append(fieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_ESTIMATED_CLOSE_DATE)
	}

	if len(fieldsMask) == 0 {
		span.LogFields(log.String("result", "no fields to update"))
		return nil
	}
	opportunityUpdateRequest.FieldsMask = fieldsMask

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
		return s.grpcClients.OpportunityClient.UpdateOpportunity(ctx, &opportunityUpdateRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}

func (s *opportunityService) UpdateRenewal(ctx context.Context, opportunityId string, renewalLikelihood neo4jenum.RenewalLikelihood, amount *float64, comments *string, ownerUserId *string, adjustedRate *int64, appSource string) error {
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

	opportunityExists, _ := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), opportunityId, neo4jutil.NodeLabelOpportunity)
	if !opportunityExists {
		err := fmt.Errorf("(OpportunityService.UpdateRenewal) opportunity with id {%s} not found", opportunityId)
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	fieldsMask := make([]opportunitypb.OpportunityMaskField, 0)
	opportunityRenewalUpdateRequest := opportunitypb.UpdateRenewalOpportunityGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Id:             opportunityId,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		OwnerUserId:    utils.IfNotNilString(ownerUserId),
		SourceFields: &commonpb.SourceFields{
			Source:    string(neo4jentity.DataSourceOpenline),
			AppSource: appSource,
		},
	}
	if amount != nil {
		opportunityRenewalUpdateRequest.Amount = *amount
		fieldsMask = append(fieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_AMOUNT)
	}
	if comments != nil {
		opportunityRenewalUpdateRequest.Comments = *comments
		fieldsMask = append(fieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_COMMENTS)
	}
	if adjustedRate != nil {
		opportunityRenewalUpdateRequest.RenewalAdjustedRate = *adjustedRate
		fieldsMask = append(fieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_ADJUSTED_RATE)
	}

	switch renewalLikelihood {
	case neo4jenum.RenewalLikelihoodHigh:
		opportunityRenewalUpdateRequest.RenewalLikelihood = opportunitypb.RenewalLikelihood_HIGH_RENEWAL
	case neo4jenum.RenewalLikelihoodMedium:
		opportunityRenewalUpdateRequest.RenewalLikelihood = opportunitypb.RenewalLikelihood_MEDIUM_RENEWAL
	case neo4jenum.RenewalLikelihoodLow:
		opportunityRenewalUpdateRequest.RenewalLikelihood = opportunitypb.RenewalLikelihood_LOW_RENEWAL
	case neo4jenum.RenewalLikelihoodZero:
		opportunityRenewalUpdateRequest.RenewalLikelihood = opportunitypb.RenewalLikelihood_ZERO_RENEWAL
	default:
		opportunityRenewalUpdateRequest.RenewalLikelihood = opportunitypb.RenewalLikelihood_ZERO_RENEWAL
	}
	fieldsMask = append(fieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_RENEWAL_LIKELIHOOD)
	opportunityRenewalUpdateRequest.FieldsMask = fieldsMask

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
		return s.grpcClients.OpportunityClient.UpdateRenewalOpportunity(ctx, &opportunityRenewalUpdateRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}

	return nil
}

func (s *opportunityService) UpdateRenewalsForOrganization(ctx context.Context, organizationId string, renewalLikelihood neo4jenum.RenewalLikelihood, renewalAdjustedRate *int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.UpdateRenewalsForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", organizationId), log.String("renewalLikelihood", renewalLikelihood.String()), log.Object("renewalAdjustedRate", renewalAdjustedRate))

	_, err := s.services.OrganizationService.GetById(ctx, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	opportunityDbNodes, err := s.repositories.Neo4jRepositories.OpportunityReadRepository.GetActiveRenewalOpportunitiesForOrganization(ctx, common.GetTenantFromContext(ctx), organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	for _, opportunityDbNode := range opportunityDbNodes {
		opportunity := neo4jmapper.MapDbNodeToOpportunityEntity(opportunityDbNode)
		if err := s.UpdateRenewal(ctx, opportunity.Id, renewalLikelihood, nil, nil, nil, renewalAdjustedRate, constants.AppSourceCustomerOsApi); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (s *opportunityService) GetPaginatedOrganizationOpportunities(ctx context.Context, page int, limit int) (*utils.Pagination, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.GetPaginatedOrganizationOpportunities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Int("page", page), log.Int("limit", limit))

	var paginatedResult = utils.Pagination{
		Limit: limit,
		Page:  page,
	}

	dbNodesWithTotalCount, err := s.repositories.Neo4jRepositories.OpportunityReadRepository.GetPaginatedOpportunitiesLinkedToAnOrganization(ctx, common.GetContext(ctx).Tenant, paginatedResult.GetSkip(), paginatedResult.GetLimit())
	if err != nil {
		return nil, err
	}
	paginatedResult.SetTotalRows(dbNodesWithTotalCount.Count)

	opportunities := neo4jentity.OpportunityEntities{}

	for _, v := range dbNodesWithTotalCount.Nodes {
		opportunities = append(opportunities, *neo4jmapper.MapDbNodeToOpportunityEntity(v))
	}
	paginatedResult.SetRows(&opportunities)
	return &paginatedResult, nil
}

func (s *opportunityService) CloseWon(ctx context.Context, opportunityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.CloseWon")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	opportunity, err := s.GetById(ctx, opportunityId)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if opportunity == nil {
		err = fmt.Errorf("opportunity not found")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	// check opportunity is not already closed won
	if opportunity.InternalStage == neo4jenum.OpportunityInternalStageClosedWon {
		err = fmt.Errorf("opportunity already closed won")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	closeWonRequest := opportunitypb.CloseWinOpportunityGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Id:             opportunity.Id,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		AppSource:      constants.AppSourceCustomerOsApi,
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
		return s.grpcClients.OpportunityClient.CloseWinOpportunity(ctx, &closeWonRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error from events processing: %s", err.Error())
		return err
	}

	return nil
}
