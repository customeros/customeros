package service

import (
	"context"
	"fmt"
	model2 "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"time"

	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	orgplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/org_plan"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type OrganizationPlanService interface {
	CreateOrganizationPlan(ctx context.Context, name, masterPlanId, orgId string) (string, error)
	UpdateOrganizationPlan(ctx context.Context, organizationPlanId, orgId string, name *string, retired *bool, statusDetails *model.OrganizationPlanStatusDetailsInput) error
	DuplicateOrganizationPlan(ctx context.Context, sourceOrganizationPlanId, orgId string) (string, error)
	GetOrganizationPlanById(ctx context.Context, organizationPlanId string) (*neo4jentity.OrganizationPlanEntity, error)
	GetOrganizationPlans(ctx context.Context, returnRetired *bool) (*neo4jentity.OrganizationPlanEntities, error)
	CreateOrganizationPlanMilestone(ctx context.Context, organizationPlanId, orgId, name string, order *int64, dueDate *time.Time, optional, adhoc bool, items []string) (string, error)
	UpdateOrganizationPlanMilestone(ctx context.Context, orgId, organizationPlanId, organizationPlanMilestoneId string, name *string, order *int64, dueDate *time.Time, items []*model.OrganizationPlanMilestoneItemInput, optional, adhoc, retired *bool, statusDetails *model.OrganizationPlanMilestoneStatusDetailsInput) error
	GetOrganizationPlanMilestoneById(ctx context.Context, organizationPlanMilestoneId string) (*neo4jentity.OrganizationPlanMilestoneEntity, error)
	GetOrganizationPlanMilestonesForOrganizationPlans(ctx context.Context, organizationPlanIds []string) (*neo4jentity.OrganizationPlanMilestoneEntities, error)
	ReorderOrganizationPlanMilestones(ctx context.Context, organizationPlanId, orgId string, organizationPlanMilestoneIds []string) error
	DuplicateOrganizationPlanMilestone(ctx context.Context, organizationPlanId, orgId, sourceOrganizationPlanMilestoneId string) (string, error)
	GetOrganizationPlansForOrganization(ctx context.Context, organizationId string) (*neo4jentity.OrganizationPlanEntities, error)
}
type organizationPlanService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewOrganizationPlanService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) OrganizationPlanService {
	return &organizationPlanService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
	}
}

func (s *organizationPlanService) CreateOrganizationPlan(ctx context.Context, name, masterPlanId, orgId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanService.CreateOrganizationPlan")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("name", name))

	grpcRequest := orgplanpb.CreateOrganizationPlanGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Name:           name,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: constants.AppSourceCustomerOsApi,
		},
		MasterPlanId: masterPlanId,
		OrgId:        orgId,
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*orgplanpb.OrganizationPlanIdGrpcResponse](func() (*orgplanpb.OrganizationPlanIdGrpcResponse, error) {
		return s.grpcClients.OrganizationPlanClient.CreateOrganizationPlan(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, model2.NodeLabelOrganizationPlan, span)

	return response.Id, nil
}

func (s *organizationPlanService) UpdateOrganizationPlan(ctx context.Context, organizationPlanId, orgId string, name *string, retired *bool, statusDetails *model.OrganizationPlanStatusDetailsInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanService.UpdateOrganizationPlan")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, organizationPlanId)
	span.LogFields(log.Object("name", name), log.Object("retired", retired))

	if name == nil && retired == nil && statusDetails == nil {
		// nothing to update
		return nil
	}

	organizationPlanExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), organizationPlanId, model2.NodeLabelOrganizationPlan)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if !organizationPlanExists {
		err = errors.New(fmt.Sprintf("Organization plan with id {%s} not found", organizationPlanId))
		tracing.TraceErr(span, err)
		return err
	}

	grpcRequest := orgplanpb.UpdateOrganizationPlanGrpcRequest{
		Tenant:             common.GetTenantFromContext(ctx),
		OrganizationPlanId: organizationPlanId,
		LoggedInUserId:     common.GetUserIdFromContext(ctx),
		Name:               utils.IfNotNilString(name),
		Retired:            utils.IfNotNilBool(retired),
		AppSource:          constants.AppSourceCustomerOsApi,
		OrgId:              orgId,
	}
	fieldsMask := make([]orgplanpb.OrganizationPlanFieldMask, 0)
	if name != nil {
		fieldsMask = append(fieldsMask, orgplanpb.OrganizationPlanFieldMask_ORGANIZATION_PLAN_PROPERTY_NAME)
	}
	if retired != nil {
		fieldsMask = append(fieldsMask, orgplanpb.OrganizationPlanFieldMask_ORGANIZATION_PLAN_PROPERTY_RETIRED)
	}
	if statusDetails != nil {
		fieldsMask = append(fieldsMask, orgplanpb.OrganizationPlanFieldMask_ORGANIZATION_PLAN_PROPERTY_STATUS_DETAILS)
		grpcRequest.StatusDetails = statusDetailsInputToProtobuf(statusDetails)
	}
	if len(fieldsMask) == 0 {
		span.LogFields(log.String("result", "No fields to update"))
		return nil
	}
	grpcRequest.FieldsMask = fieldsMask

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*orgplanpb.OrganizationPlanIdGrpcResponse](func() (*orgplanpb.OrganizationPlanIdGrpcResponse, error) {
		return s.grpcClients.OrganizationPlanClient.UpdateOrganizationPlan(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}
	return nil
}

func (s *organizationPlanService) GetOrganizationPlanById(ctx context.Context, organizationPlanId string) (*neo4jentity.OrganizationPlanEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanService.GetOrganizationPlanById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, organizationPlanId)

	if organizationPlanDbNode, err := s.repositories.Neo4jRepositories.OrganizationPlanReadRepository.GetOrganizationPlanById(ctx, common.GetContext(ctx).Tenant, organizationPlanId); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Organization plan with id {%s} not found", organizationPlanId))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToOrganizationPlanEntity(organizationPlanDbNode), nil
	}
}

func (s *organizationPlanService) CreateOrganizationPlanMilestone(ctx context.Context, organizationPlanId, orgId, name string, order *int64, dueDate *time.Time, optional, adhoc bool, items []string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanService.CreateOrganizationPlanMilestone")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationId", orgId), log.String("organizationPlanId", organizationPlanId), log.String("name", name), log.Int64("order", *order), log.Bool("optional", optional), log.Object("items", items))

	organizationPlanExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), organizationPlanId, model2.NodeLabelOrganizationPlan)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	if !organizationPlanExists {
		err = errors.New(fmt.Sprintf("Organization plan with id {%s} not found", organizationPlanId))
		tracing.TraceErr(span, err)
		return "", err
	}

	var due time.Time
	var dueDatePtr *time.Time
	if dueDate != nil {
		due = dueDate.UTC()
		dueDatePtr = &due
	} else {
		dueDatePtr = nil
	}

	grpcRequest := orgplanpb.CreateOrganizationPlanMilestoneGrpcRequest{
		Tenant:             common.GetTenantFromContext(ctx),
		OrganizationPlanId: organizationPlanId,
		LoggedInUserId:     common.GetUserIdFromContext(ctx),
		Name:               name,
		Order:              *order,
		DueDate:            utils.ConvertTimeToTimestampPtr(dueDatePtr),
		Optional:           optional,
		Items:              items,
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: constants.AppSourceCustomerOsApi,
		},
		OrgId: orgId,
		Adhoc: adhoc,
	}

	tracing.LogObjectAsJson(span, "CreateOrganizationPlanMilestoneGrpcRequest", &grpcRequest)

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*orgplanpb.OrganizationPlanMilestoneIdGrpcResponse](func() (*orgplanpb.OrganizationPlanMilestoneIdGrpcResponse, error) {
		return s.grpcClients.OrganizationPlanClient.CreateOrganizationPlanMilestone(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, model2.NodeLabelOrganizationPlanMilestone, span)

	span.LogFields(log.String("response - created organizationPlanMilestoneId", response.Id))
	return response.Id, nil
}

func (s *organizationPlanService) GetOrganizationPlanMilestoneById(ctx context.Context, organizationPlanMilestoneId string) (*neo4jentity.OrganizationPlanMilestoneEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanService.GetOrganizationPlanMilestoneById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, organizationPlanMilestoneId)

	if organizationPlanMilestoneDbNode, err := s.repositories.Neo4jRepositories.OrganizationPlanReadRepository.GetOrganizationPlanMilestoneById(ctx, common.GetContext(ctx).Tenant, organizationPlanMilestoneId); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Organization plan milestone with id {%s} not found", organizationPlanMilestoneId))
		return nil, wrappedErr
	} else {
		opm := neo4jmapper.MapDbNodeToOrganizationPlanMilestoneEntity(organizationPlanMilestoneDbNode)
		return opm, nil
	}
}

func (s *organizationPlanService) GetOrganizationPlans(ctx context.Context, returnRetired *bool) (*neo4jentity.OrganizationPlanEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanService.GetOrganizationPlans")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("returnRetired", returnRetired))

	organizationPlanDbNodes, err := s.repositories.Neo4jRepositories.OrganizationPlanReadRepository.GetOrganizationPlansOrderByCreatedAt(ctx, common.GetTenantFromContext(ctx), returnRetired)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	organizationPlanEntities := make(neo4jentity.OrganizationPlanEntities, 0, len(organizationPlanDbNodes))
	for _, v := range organizationPlanDbNodes {
		organizationPlanEntities = append(organizationPlanEntities, *neo4jmapper.MapDbNodeToOrganizationPlanEntity(v))
	}
	return &organizationPlanEntities, nil
}

func (s *organizationPlanService) GetOrganizationPlanMilestonesForOrganizationPlans(ctx context.Context, organizationPlanIds []string) (*neo4jentity.OrganizationPlanMilestoneEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanService.GetOrganizationPlanMilestonesForOrganizationPlans")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("organizationPlanIds", organizationPlanIds))

	organizationPlanMilestoneDbNodes, err := s.repositories.Neo4jRepositories.OrganizationPlanReadRepository.GetOrganizationPlanMilestonesForOrganizationPlans(ctx, common.GetTenantFromContext(ctx), organizationPlanIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	organizationPlanMilestoneEntities := make(neo4jentity.OrganizationPlanMilestoneEntities, 0, len(organizationPlanMilestoneDbNodes))
	for _, v := range organizationPlanMilestoneDbNodes {
		organizationPlanMilestoneEntity := neo4jmapper.MapDbNodeToOrganizationPlanMilestoneEntity(v.Node)
		organizationPlanMilestoneEntity.DataloaderKey = v.LinkedNodeId
		organizationPlanMilestoneEntities = append(organizationPlanMilestoneEntities, *organizationPlanMilestoneEntity)
	}
	return &organizationPlanMilestoneEntities, nil
}

func (s *organizationPlanService) UpdateOrganizationPlanMilestone(ctx context.Context, orgId, organizationPlanId, organizationPlanMilestoneId string, name *string, order *int64, dueDate *time.Time, items []*model.OrganizationPlanMilestoneItemInput, optional, adhoc, retired *bool, statusDetails *model.OrganizationPlanMilestoneStatusDetailsInput) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanService.UpdateOrganizationPlanMilestone")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, organizationPlanMilestoneId)
	span.LogFields(log.Object("name", name), log.Object("order", order), log.Object("dueDate", dueDate), log.Object("items", items), log.Object("optional", optional), log.Object("retired", retired))

	if name == nil && retired == nil && order == nil && dueDate == nil && optional == nil && items == nil && adhoc == nil && statusDetails == nil {
		// nothing to update
		return nil
	}

	organizationPlanMilestoneExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), organizationPlanMilestoneId, model2.NodeLabelOrganizationPlanMilestone)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if !organizationPlanMilestoneExists {
		err = errors.New(fmt.Sprintf("Organization plan milestone with id {%s} not found", organizationPlanId))
		tracing.TraceErr(span, err)
		return err
	}

	grpcRequest := orgplanpb.UpdateOrganizationPlanMilestoneGrpcRequest{
		Tenant:                      common.GetTenantFromContext(ctx),
		OrganizationPlanId:          organizationPlanId,
		OrganizationPlanMilestoneId: organizationPlanMilestoneId,
		LoggedInUserId:              common.GetUserIdFromContext(ctx),
		AppSource:                   constants.AppSourceCustomerOsApi,
		OrgId:                       orgId,
		Name:                        utils.IfNotNilString(name),
		Retired:                     utils.IfNotNilBool(retired),
		Order:                       utils.IfNotNilInt64(order),
		Optional:                    utils.IfNotNilBool(optional),
		Adhoc:                       utils.IfNotNilBool(adhoc),
	}
	fieldsMask := make([]orgplanpb.OrganizationPlanMilestoneFieldMask, 0)
	if name != nil {
		fieldsMask = append(fieldsMask, orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_NAME)
	}
	if retired != nil {
		fieldsMask = append(fieldsMask, orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_RETIRED)
	}
	if order != nil {
		fieldsMask = append(fieldsMask, orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_ORDER)
	}
	if dueDate != nil {
		due := utils.IfNotNilTimeWithDefault(dueDate, time.Now().UTC()).UTC()
		grpcRequest.DueDate = utils.ConvertTimeToTimestampPtr(&due)
		fieldsMask = append(fieldsMask, orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_DUE_DATE)
	}
	if optional != nil {
		fieldsMask = append(fieldsMask, orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_OPTIONAL)
	}
	if items != nil {
		grpcRequest.Items = modelItemsToPbItems(items)
		fieldsMask = append(fieldsMask, orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_ITEMS)
	}
	if statusDetails != nil {
		grpcRequest.StatusDetails = milestoneStatusDetailsInputToProtobuf(statusDetails)
		fieldsMask = append(fieldsMask, orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_STATUS_DETAILS)
	}
	if adhoc != nil {
		fieldsMask = append(fieldsMask, orgplanpb.OrganizationPlanMilestoneFieldMask_ORGANIZATION_PLAN_MILESTONE_PROPERTY_ADHOC)
	}
	if len(fieldsMask) == 0 {
		span.LogFields(log.String("result", "No fields to update"))
		return nil
	}
	grpcRequest.FieldsMask = fieldsMask

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*orgplanpb.OrganizationPlanMilestoneIdGrpcResponse](func() (*orgplanpb.OrganizationPlanMilestoneIdGrpcResponse, error) {
		return s.grpcClients.OrganizationPlanClient.UpdateOrganizationPlanMilestone(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}
	return nil
}

func (s *organizationPlanService) ReorderOrganizationPlanMilestones(ctx context.Context, organizationPlanId, orgId string, organizationPlanMilestoneIds []string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanService.ReorderOrganizationPlanMilestones")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationPlanId", organizationPlanId), log.Object("organizationPlanMilestoneIds", organizationPlanMilestoneIds))

	organizationPlanMilestoneExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), organizationPlanId, model2.NodeLabelOrganizationPlan)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if !organizationPlanMilestoneExists {
		err = errors.New(fmt.Sprintf("Organization plan with id {%s} not found", organizationPlanId))
		tracing.TraceErr(span, err)
		return err
	}

	grpcRequest := orgplanpb.ReorderOrganizationPlanMilestonesGrpcRequest{
		Tenant:                       common.GetTenantFromContext(ctx),
		OrganizationPlanId:           organizationPlanId,
		LoggedInUserId:               common.GetUserIdFromContext(ctx),
		AppSource:                    constants.AppSourceCustomerOsApi,
		OrganizationPlanMilestoneIds: organizationPlanMilestoneIds,
		OrgId:                        orgId,
	}
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*orgplanpb.OrganizationPlanIdGrpcResponse](func() (*orgplanpb.OrganizationPlanIdGrpcResponse, error) {
		return s.grpcClients.OrganizationPlanClient.ReorderOrganizationPlanMilestones(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}
	return nil
}

func (s *organizationPlanService) DuplicateOrganizationPlanMilestone(ctx context.Context, organizationPlanId, orgId, sourceOrganizationPlanMilestoneId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanService.DuplicateOrganizationPlanMilestone")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("organizationPlanId", organizationPlanId), log.String("sourceOrganizationPlanMilestoneId", sourceOrganizationPlanMilestoneId))

	organizationPlanMilestoneDbNode, err := s.repositories.Neo4jRepositories.OrganizationPlanReadRepository.GetOrganizationPlanMilestoneByPlanAndId(ctx, common.GetContext(ctx).Tenant, organizationPlanId, sourceOrganizationPlanMilestoneId)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	if organizationPlanMilestoneDbNode == nil {
		err = errors.New(fmt.Sprintf("Organization plan milestone with id {%s} not found", sourceOrganizationPlanMilestoneId))
		tracing.TraceErr(span, err)
		return "", err
	}
	souceOrganizationPlanMilestoneEntity := neo4jmapper.MapDbNodeToOrganizationPlanMilestoneEntity(organizationPlanMilestoneDbNode)
	maxOrder, err := s.repositories.Neo4jRepositories.OrganizationPlanReadRepository.GetMaxOrderForOrganizationPlanMilestones(ctx, common.GetContext(ctx).Tenant, organizationPlanId)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	itemsText := make([]string, 0, len(souceOrganizationPlanMilestoneEntity.Items))
	for _, v := range souceOrganizationPlanMilestoneEntity.Items {
		itemsText = append(itemsText, v.Text)
	}
	grpcRequest := orgplanpb.CreateOrganizationPlanMilestoneGrpcRequest{
		Tenant:             common.GetTenantFromContext(ctx),
		OrganizationPlanId: organizationPlanId,
		LoggedInUserId:     common.GetUserIdFromContext(ctx),
		Name:               souceOrganizationPlanMilestoneEntity.Name,
		Order:              maxOrder + 1,
		DueDate:            utils.ConvertTimeToTimestampPtr(&souceOrganizationPlanMilestoneEntity.DueDate),
		Optional:           souceOrganizationPlanMilestoneEntity.Optional,
		Items:              itemsText,
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: constants.AppSourceCustomerOsApi,
		},
		OrgId: orgId,
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*orgplanpb.OrganizationPlanMilestoneIdGrpcResponse](func() (*orgplanpb.OrganizationPlanMilestoneIdGrpcResponse, error) {
		return s.grpcClients.OrganizationPlanClient.CreateOrganizationPlanMilestone(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, model2.NodeLabelOrganizationPlanMilestone, span)
	return response.Id, nil
}

func (s *organizationPlanService) DuplicateOrganizationPlan(ctx context.Context, sourceOrganizationPlanId, orgId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanService.DuplicateOrganizationPlan")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("sourceOrganizationPlanId", sourceOrganizationPlanId))

	sourceOrganizationPlanEntity, err := s.GetOrganizationPlanById(ctx, sourceOrganizationPlanId)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	if sourceOrganizationPlanEntity == nil {
		err = errors.New(fmt.Sprintf("Organization plan with id {%s} not found", sourceOrganizationPlanId))
		tracing.TraceErr(span, err)
		return "", err
	}
	organizationPlanMilestoneEntities, err := s.GetOrganizationPlanMilestonesForOrganizationPlans(ctx, []string{sourceOrganizationPlanId})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	grpcRequest := orgplanpb.CreateOrganizationPlanGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Name:           sourceOrganizationPlanEntity.Name,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: constants.AppSourceCustomerOsApi,
		},
		OrgId: orgId,
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*orgplanpb.OrganizationPlanIdGrpcResponse](func() (*orgplanpb.OrganizationPlanIdGrpcResponse, error) {
		return s.grpcClients.OrganizationPlanClient.CreateOrganizationPlan(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, model2.NodeLabelOrganizationPlan, span)

	for _, organizationPlanMilestoneEntity := range *organizationPlanMilestoneEntities {
		itemsText := make([]string, 0, len(organizationPlanMilestoneEntity.Items))
		for _, v := range organizationPlanMilestoneEntity.Items {
			itemsText = append(itemsText, v.Text)
		}
		if !organizationPlanMilestoneEntity.Retired {
			grpcRequestCreateMilestone := orgplanpb.CreateOrganizationPlanMilestoneGrpcRequest{
				Tenant:             common.GetTenantFromContext(ctx),
				OrganizationPlanId: response.Id,
				LoggedInUserId:     common.GetUserIdFromContext(ctx),
				Name:               organizationPlanMilestoneEntity.Name,
				Order:              organizationPlanMilestoneEntity.Order,
				DueDate:            utils.ConvertTimeToTimestampPtr(&organizationPlanMilestoneEntity.DueDate),
				Optional:           organizationPlanMilestoneEntity.Optional,
				Items:              itemsText,
				SourceFields: &commonpb.SourceFields{
					Source:    neo4jentity.DataSourceOpenline.String(),
					AppSource: constants.AppSourceCustomerOsApi,
				},
				OrgId: orgId,
			}
			_, err = utils.CallEventsPlatformGRPCWithRetry[*orgplanpb.OrganizationPlanMilestoneIdGrpcResponse](func() (*orgplanpb.OrganizationPlanMilestoneIdGrpcResponse, error) {
				return s.grpcClients.OrganizationPlanClient.CreateOrganizationPlanMilestone(ctx, &grpcRequestCreateMilestone)
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error from events processing: %s", err.Error())
			}
		}
	}

	return response.Id, nil
}

func (s *organizationPlanService) GetOrganizationPlansForOrganization(ctx context.Context, organizationId string) (*neo4jentity.OrganizationPlanEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationPlanService.GetOrganizationPlanForOrganization")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, organizationId)

	organizationPlanDbNodes, err := s.repositories.Neo4jRepositories.OrganizationPlanReadRepository.GetOrganizationPlansForOrganization(ctx, common.GetContext(ctx).Tenant, organizationId)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	organizationPlanEntities := make(neo4jentity.OrganizationPlanEntities, 0, len(organizationPlanDbNodes))
	if len(organizationPlanDbNodes) == 0 {
		span.LogFields(log.String("Warning", fmt.Sprintf("Organization plans for organization with id {%s} not found", organizationId)))
		return &organizationPlanEntities, nil
	}

	for _, v := range organizationPlanDbNodes {
		organizationPlanEntities = append(organizationPlanEntities, *neo4jmapper.MapDbNodeToOrganizationPlanEntity(v))
	}
	return &organizationPlanEntities, nil
}

//////////////////////////////////////////
////// utils just for this service //////
////////////////////////////////////////

func statusDetailsInputToProtobuf(sd *model.OrganizationPlanStatusDetailsInput) *orgplanpb.StatusDetails {
	out := &orgplanpb.StatusDetails{
		Status:    sd.Status.String(),
		UpdatedAt: utils.ConvertTimeToTimestampPtr(&sd.UpdatedAt),
		Comments:  sd.Text,
	}
	return out
}

func milestoneStatusDetailsInputToProtobuf(sd *model.OrganizationPlanMilestoneStatusDetailsInput) *orgplanpb.StatusDetails {
	out := &orgplanpb.StatusDetails{
		Status:    sd.Status.String(),
		UpdatedAt: utils.ConvertTimeToTimestampPtr(&sd.UpdatedAt),
		Comments:  sd.Text,
	}
	return out
}

func modelItemsToPbItems(items []*model.OrganizationPlanMilestoneItemInput) []*orgplanpb.OrganizationPlanMilestoneItem {
	out := make([]*orgplanpb.OrganizationPlanMilestoneItem, 0, len(items))
	for _, v := range items {
		var iUuid string
		if v.UUID == nil {
			iUuid = uuid.New().String()
		} else {
			iUuid = *v.UUID
		}
		out = append(out, &orgplanpb.OrganizationPlanMilestoneItem{
			Status:    v.Status.String(),
			UpdatedAt: utils.ConvertTimeToTimestampPtr(&v.UpdatedAt),
			Text:      v.Text,
			Uuid:      iUuid,
		})
	}
	return out
}
