package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper"
	enummapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/mapper/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jenum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/enum"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	opportunitypb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/opportunity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type OpportunityService interface {
	Create(ctx context.Context, input model.OpportunityCreateInput) (string, error)
	Update(ctx context.Context, input model.OpportunityUpdateInput) error
	UpdateRenewal(ctx context.Context, opportunityId string, renewalLikelihood neo4jenum.RenewalLikelihood, amount *float64, comments *string, ownerUserId *string, adjustedRate *int64, appSource string) error
	UpdateRenewalsForOrganization(ctx context.Context, organizationId string, renewalLikelihood neo4jenum.RenewalLikelihood, renewalAdjustedRate *int64) error
	CloseWon(ctx context.Context, opportunityId string) error
	CloseLost(ctx context.Context, opportunityId string) error
	ReplaceOwner(ctx context.Context, opportunityId, userId string) error
	RemoveOwner(ctx context.Context, opportunityId string) error
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
	orgFound, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), input.OrganizationID, model2.NodeLabelOrganization)
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
		MaxAmount:      utils.IfNotNilFloat64(input.MaxAmount),
		LikelihoodRate: utils.IfNotNilInt64(input.LikelihoodRate),
		SourceFields: &commonpb.SourceFields{
			Source:    string(neo4jentity.DataSourceOpenline),
			AppSource: constants.AppSourceCustomerOsApi,
		},
	}
	if input.Currency != nil {
		opportunityCreateRequest.Currency = enummapper.MapCurrencyFromModel(*input.Currency).String()
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

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, opportunityIdGrpcResponse.Id, model2.NodeLabelOpportunity, span)

	return opportunityIdGrpcResponse.Id, nil
}

func (s *opportunityService) Update(ctx context.Context, input model.OpportunityUpdateInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.Update")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "input", input)

	tenant := common.GetTenantFromContext(ctx)

	opportunity, err := s.services.CommonServices.OpportunityService.GetById(ctx, tenant, input.OpportunityID)
	if err != nil {
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
	if input.MaxAmount != nil {
		opportunityUpdateRequest.MaxAmount = *input.MaxAmount
		fieldsMask = append(fieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_MAX_AMOUNT)
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
	if input.NextSteps != nil {
		opportunityUpdateRequest.NextSteps = *input.NextSteps
		fieldsMask = append(fieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_NEXT_STEPS)
	}
	if input.LikelihoodRate != nil {
		opportunityUpdateRequest.LikelihoodRate = *input.LikelihoodRate
		fieldsMask = append(fieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_LIKELIHOOD_RATE)
	}
	if input.Currency != nil {
		opportunityUpdateRequest.Currency = enummapper.MapCurrencyFromModel(*input.Currency).String()
		fieldsMask = append(fieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_CURRENCY)
	}
	if input.InternalStage != nil && opportunity.InternalStage != mapper.MapInternalStageFromModel(*input.InternalStage) {
		switch *input.InternalStage {
		case model.InternalStageOpen:
			opportunityUpdateRequest.InternalStage = opportunitypb.OpportunityInternalStage_OPEN
		case model.InternalStageClosedWon,
			model.InternalStageClosedLost:
			err := fmt.Errorf("final internal stage should be set with dedicated APIs")
			s.log.Error(err.Error())
			tracing.TraceErr(span, err)
			return err
		}
		fieldsMask = append(fieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_INTERNAL_STAGE)
	}
	// Changing external stage should set internal stage back to OPEN
	if input.ExternalStage != nil && *input.ExternalStage != "" && opportunity.ExternalStage != opportunity.ExternalStage && opportunity.InternalStage != neo4jenum.OpportunityInternalStageOpen {
		opportunityUpdateRequest.InternalStage = opportunitypb.OpportunityInternalStage_OPEN
		fieldsMask = append(fieldsMask, opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_INTERNAL_STAGE)
	}

	if len(fieldsMask) == 0 {
		span.LogFields(log.String("result", "no fields to update"))
		return nil
	}
	opportunityUpdateRequest.FieldsMask = fieldsMask

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
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

	opportunityExists, _ := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), opportunityId, model2.NodeLabelOpportunity)
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

	opportunityDbNodes, err := s.repositories.Neo4jRepositories.OpportunityReadRepository.GetActiveRenewalOpportunitiesForOrganization(ctx, common.GetTenantFromContext(ctx), organizationId, true)
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

	tenant := common.GetTenantFromContext(ctx)

	opportunity, err := s.services.CommonServices.OpportunityService.GetById(ctx, tenant, opportunityId)
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

func (s *opportunityService) CloseLost(ctx context.Context, opportunityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.CloseLost")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	tenant := common.GetTenantFromContext(ctx)

	opportunity, err := s.services.CommonServices.OpportunityService.GetById(ctx, tenant, opportunityId)
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

	// check opportunity is not already closed lost
	if opportunity.InternalStage == neo4jenum.OpportunityInternalStageClosedLost {
		err = fmt.Errorf("opportunity already closed lost")
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	closeLostRequest := opportunitypb.CloseLooseOpportunityGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Id:             opportunity.Id,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		AppSource:      constants.AppSourceCustomerOsApi,
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
		return s.grpcClients.OpportunityClient.CloseLooseOpportunity(ctx, &closeLostRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error from events processing: %s", err.Error())
		return err
	}

	return nil
}

func (s *opportunityService) ReplaceOwner(ctx context.Context, opportunityId, userId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.ReplaceOwner")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)
	span.LogFields(log.String("userId", userId))

	opportunityExists, _ := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), opportunityId, model2.NodeLabelOpportunity)
	if !opportunityExists {
		err := fmt.Errorf("(OpportunityService.ReplaceOwner) opportunity with id {%s} not found", opportunityId)
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	userExists, _ := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), userId, model2.NodeLabelUser)
	if !userExists {
		err := fmt.Errorf("(OpportunityService.ReplaceOwner) user with id {%s} not found", userId)
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	updateOpportunityGrpcRequest := opportunitypb.UpdateOpportunityGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Id:             opportunityId,
		OwnerUserId:    userId,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		FieldsMask:     []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_OWNER_USER_ID},
		SourceFields: &commonpb.SourceFields{
			Source:    string(neo4jentity.DataSourceOpenline),
			AppSource: constants.AppSourceCustomerOsApi,
		},
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
		return s.grpcClients.OpportunityClient.UpdateOpportunity(ctx, &updateOpportunityGrpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error from events processing: %s", err.Error())
		return err
	}

	return nil
}

func (s *opportunityService) RemoveOwner(ctx context.Context, opportunityId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityService.RemoveOwner")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, opportunityId)

	opportunityExists, _ := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), opportunityId, model2.NodeLabelOpportunity)
	if !opportunityExists {
		err := fmt.Errorf("(OpportunityService.ReplaceOwner) opportunity with id {%s} not found", opportunityId)
		s.log.Error(err.Error())
		tracing.TraceErr(span, err)
		return err
	}

	updateOpportunityGrpcRequest := opportunitypb.UpdateOpportunityGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Id:             opportunityId,
		OwnerUserId:    "",
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		FieldsMask:     []opportunitypb.OpportunityMaskField{opportunitypb.OpportunityMaskField_OPPORTUNITY_PROPERTY_OWNER_USER_ID},
		SourceFields: &commonpb.SourceFields{
			Source:    string(neo4jentity.DataSourceOpenline),
			AppSource: constants.AppSourceCustomerOsApi,
		},
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err := utils.CallEventsPlatformGRPCWithRetry[*opportunitypb.OpportunityIdGrpcResponse](func() (*opportunitypb.OpportunityIdGrpcResponse, error) {
		return s.grpcClients.OpportunityClient.UpdateOpportunity(ctx, &updateOpportunityGrpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("error from events processing: %s", err.Error())
		return err
	}

	return nil
}
