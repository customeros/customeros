package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/neo4jutil"
	neo4jrepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	masterplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/master_plan"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type MasterPlanService interface {
	CreateMasterPlan(ctx context.Context, name string) (string, error)
	UpdateMasterPlan(ctx context.Context, id string, name *string, retired *bool) error
	DuplicateMasterPlan(ctx context.Context, sourceMasterPlanId string) (string, error)
	GetMasterPlanById(ctx context.Context, masterPlanId string) (*neo4jentity.MasterPlanEntity, error)
	GetMasterPlans(ctx context.Context, returnRetired *bool) (*neo4jentity.MasterPlanEntities, error)
	CreateMasterPlanMilestone(ctx context.Context, masterPlanId, name string, order, durationHours int64, optional bool, items []string) (string, error)
	UpdateMasterPlanMilestone(ctx context.Context, masterPlanId, masterPlanMilestoneId string, name *string, order, hours *int64, items []string, optional *bool, retired *bool) error
	GetMasterPlanMilestoneById(ctx context.Context, masterPlanMilestoneId string) (*neo4jentity.MasterPlanMilestoneEntity, error)
	GetMasterPlanMilestonesForMasterPlans(ctx context.Context, masterPlanIds []string) (*neo4jentity.MasterPlanMilestoneEntities, error)
	ReorderMasterPlanMilestones(ctx context.Context, masterPlanId string, masterPlanMilestoneIds []string) error
	DuplicateMasterPlanMilestone(ctx context.Context, masterPlanId, sourceMasterPlanMilestoneId string) (string, error)
	CreateDefaultMasterPlan(ctx context.Context) (string, error)
}
type masterPlanService struct {
	log          logger.Logger
	repositories *repository.Repositories
	grpcClients  *grpc_client.Clients
}

func NewMasterPlanService(log logger.Logger, repositories *repository.Repositories, grpcClients *grpc_client.Clients) MasterPlanService {
	return &masterPlanService{
		log:          log,
		repositories: repositories,
		grpcClients:  grpcClients,
	}
}

func (s *masterPlanService) CreateMasterPlan(ctx context.Context, name string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanService.CreateMasterPlan")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("name", name))

	grpcRequest := masterplanpb.CreateMasterPlanGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Name:           name,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: constants.AppSourceCustomerOsApi,
		},
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*masterplanpb.MasterPlanIdGrpcResponse](func() (*masterplanpb.MasterPlanIdGrpcResponse, error) {
		return s.grpcClients.MasterPlanClient.CreateMasterPlan(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, neo4jutil.NodeLabelMasterPlan, span)

	return response.Id, nil
}

func (s *masterPlanService) UpdateMasterPlan(ctx context.Context, masterPlanId string, name *string, retired *bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanService.UpdateMasterPlan")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, masterPlanId)
	span.LogFields(log.Object("name", name), log.Object("retired", retired))

	if name == nil && retired == nil {
		// nothing to update
		return nil
	}

	masterPlanExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), masterPlanId, neo4jutil.NodeLabelMasterPlan)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if !masterPlanExists {
		err = errors.New(fmt.Sprintf("Master plan with id {%s} not found", masterPlanId))
		tracing.TraceErr(span, err)
		return err
	}

	grpcRequest := masterplanpb.UpdateMasterPlanGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		MasterPlanId:   masterPlanId,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		Name:           utils.IfNotNilString(name),
		Retired:        utils.IfNotNilBool(retired),
		AppSource:      constants.AppSourceCustomerOsApi,
	}
	fieldsMask := make([]masterplanpb.MasterPlanFieldMask, 0)
	if name != nil {
		fieldsMask = append(fieldsMask, masterplanpb.MasterPlanFieldMask_MASTER_PLAN_PROPERTY_NAME)
	}
	if retired != nil {
		fieldsMask = append(fieldsMask, masterplanpb.MasterPlanFieldMask_MASTER_PLAN_PROPERTY_RETIRED)
	}
	if len(fieldsMask) == 0 {
		span.LogFields(log.String("result", "No fields to update"))
		return nil
	}
	grpcRequest.FieldsMask = fieldsMask

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*masterplanpb.MasterPlanIdGrpcResponse](func() (*masterplanpb.MasterPlanIdGrpcResponse, error) {
		return s.grpcClients.MasterPlanClient.UpdateMasterPlan(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}
	return nil
}

func (s *masterPlanService) GetMasterPlanById(ctx context.Context, masterPlanId string) (*neo4jentity.MasterPlanEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanService.GetMasterPlanById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, masterPlanId)

	if masterPlanDbNode, err := s.repositories.Neo4jRepositories.MasterPlanReadRepository.GetMasterPlanById(ctx, common.GetContext(ctx).Tenant, masterPlanId); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Master plan with id {%s} not found", masterPlanId))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToMasterPlanEntity(masterPlanDbNode), nil
	}
}

func (s *masterPlanService) CreateMasterPlanMilestone(ctx context.Context, masterPlanId, name string, order, durationHours int64, optional bool, items []string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanService.CreateMasterPlanMilestone")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("masterPlanId", masterPlanId), log.String("name", name), log.Int64("order", order), log.Int64("durationHours", durationHours), log.Bool("optional", optional), log.Object("items", items))

	masterPlanExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), masterPlanId, neo4jutil.NodeLabelMasterPlan)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	if !masterPlanExists {
		err = errors.New(fmt.Sprintf("Master plan with id {%s} not found", masterPlanId))
		tracing.TraceErr(span, err)
		return "", err
	}

	grpcRequest := masterplanpb.CreateMasterPlanMilestoneGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		MasterPlanId:   masterPlanId,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		Name:           name,
		Order:          order,
		DurationHours:  durationHours,
		Optional:       optional,
		Items:          items,
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: constants.AppSourceCustomerOsApi,
		},
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*masterplanpb.MasterPlanMilestoneIdGrpcResponse](func() (*masterplanpb.MasterPlanMilestoneIdGrpcResponse, error) {
		return s.grpcClients.MasterPlanClient.CreateMasterPlanMilestone(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, neo4jutil.NodeLabelMasterPlanMilestone, span)

	span.LogFields(log.String("response - created masterPlanMilestoneId", response.Id))
	return response.Id, nil
}

func (s *masterPlanService) GetMasterPlanMilestoneById(ctx context.Context, masterPlanMilestoneId string) (*neo4jentity.MasterPlanMilestoneEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanService.GetMasterPlanMilestoneById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, masterPlanMilestoneId)

	if masterPlanMilestoneDbNode, err := s.repositories.Neo4jRepositories.MasterPlanReadRepository.GetMasterPlanMilestoneById(ctx, common.GetContext(ctx).Tenant, masterPlanMilestoneId); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Master plan milestone with id {%s} not found", masterPlanMilestoneId))
		return nil, wrappedErr
	} else {
		return neo4jmapper.MapDbNodeToMasterPlanMilestoneEntity(masterPlanMilestoneDbNode), nil
	}
}

func (s *masterPlanService) GetMasterPlans(ctx context.Context, returnRetired *bool) (*neo4jentity.MasterPlanEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanService.GetMasterPlans")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("returnRetired", returnRetired))

	masterPlanDbNodes, err := s.repositories.Neo4jRepositories.MasterPlanReadRepository.GetMasterPlansOrderByCreatedAt(ctx, common.GetTenantFromContext(ctx), returnRetired)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	masterPlanEntities := make(neo4jentity.MasterPlanEntities, 0, len(masterPlanDbNodes))
	for _, v := range masterPlanDbNodes {
		masterPlanEntities = append(masterPlanEntities, *neo4jmapper.MapDbNodeToMasterPlanEntity(v))
	}
	return &masterPlanEntities, nil
}

func (s *masterPlanService) GetMasterPlanMilestonesForMasterPlans(ctx context.Context, masterPlanIds []string) (*neo4jentity.MasterPlanMilestoneEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanService.GetMasterPlanMilestonesForMasterPlans")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("masterPlanIds", masterPlanIds))

	masterPlanMilestoneDbNodes, err := s.repositories.Neo4jRepositories.MasterPlanReadRepository.GetMasterPlanMilestonesForMasterPlans(ctx, common.GetTenantFromContext(ctx), masterPlanIds)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	masterPlanMilestoneEntities := make(neo4jentity.MasterPlanMilestoneEntities, 0, len(masterPlanMilestoneDbNodes))
	for _, v := range masterPlanMilestoneDbNodes {
		masterPlanMilestoneEntity := neo4jmapper.MapDbNodeToMasterPlanMilestoneEntity(v.Node)
		masterPlanMilestoneEntity.DataloaderKey = v.LinkedNodeId
		masterPlanMilestoneEntities = append(masterPlanMilestoneEntities, *masterPlanMilestoneEntity)
	}
	return &masterPlanMilestoneEntities, nil
}

func (s *masterPlanService) UpdateMasterPlanMilestone(ctx context.Context, masterPlanId, masterPlanMilestoneId string, name *string, order, durationHours *int64, items []string, optional *bool, retired *bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanService.UpdateMasterPlanMilestone")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, masterPlanMilestoneId)
	span.LogFields(log.Object("name", name), log.Object("order", order), log.Object("durationHours", durationHours), log.Object("items", items), log.Object("optional", optional), log.Object("retired", retired))

	if name == nil && retired == nil && order == nil && durationHours == nil && optional == nil && items == nil {
		// nothing to update
		return nil
	}

	masterPlanMilestoneExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), masterPlanMilestoneId, neo4jutil.NodeLabelMasterPlanMilestone)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if !masterPlanMilestoneExists {
		err = errors.New(fmt.Sprintf("Master plan milestone with id {%s} not found", masterPlanId))
		tracing.TraceErr(span, err)
		return err
	}

	grpcRequest := masterplanpb.UpdateMasterPlanMilestoneGrpcRequest{
		Tenant:                common.GetTenantFromContext(ctx),
		MasterPlanId:          masterPlanId,
		MasterPlanMilestoneId: masterPlanMilestoneId,
		LoggedInUserId:        common.GetUserIdFromContext(ctx),
		AppSource:             constants.AppSourceCustomerOsApi,
	}
	fieldsMask := make([]masterplanpb.MasterPlanMilestoneFieldMask, 0)
	if name != nil {
		grpcRequest.Name = utils.IfNotNilString(name)
		fieldsMask = append(fieldsMask, masterplanpb.MasterPlanMilestoneFieldMask_MASTER_PLAN_MILESTONE_PROPERTY_NAME)
	}
	if retired != nil {
		grpcRequest.Retired = utils.IfNotNilBool(retired)
		fieldsMask = append(fieldsMask, masterplanpb.MasterPlanMilestoneFieldMask_MASTER_PLAN_MILESTONE_PROPERTY_RETIRED)
	}
	if order != nil {
		grpcRequest.Order = utils.IfNotNilInt64(order)
		fieldsMask = append(fieldsMask, masterplanpb.MasterPlanMilestoneFieldMask_MASTER_PLAN_MILESTONE_PROPERTY_ORDER)
	}
	if durationHours != nil {
		grpcRequest.DurationHours = utils.IfNotNilInt64(durationHours)
		fieldsMask = append(fieldsMask, masterplanpb.MasterPlanMilestoneFieldMask_MASTER_PLAN_MILESTONE_PROPERTY_DURATION_HOURS)
	}
	if optional != nil {
		grpcRequest.Optional = utils.IfNotNilBool(optional)
		fieldsMask = append(fieldsMask, masterplanpb.MasterPlanMilestoneFieldMask_MASTER_PLAN_MILESTONE_PROPERTY_OPTIONAL)
	}
	if items != nil {
		grpcRequest.Items = items
		fieldsMask = append(fieldsMask, masterplanpb.MasterPlanMilestoneFieldMask_MASTER_PLAN_MILESTONE_PROPERTY_ITEMS)
	}
	if len(fieldsMask) == 0 {
		span.LogFields(log.String("result", "No fields to update"))
		return nil
	}
	grpcRequest.FieldsMask = fieldsMask

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*masterplanpb.MasterPlanMilestoneIdGrpcResponse](func() (*masterplanpb.MasterPlanMilestoneIdGrpcResponse, error) {
		return s.grpcClients.MasterPlanClient.UpdateMasterPlanMilestone(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}
	return nil
}

func (s *masterPlanService) ReorderMasterPlanMilestones(ctx context.Context, masterPlanId string, masterPlanMilestoneIds []string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanService.ReorderMasterPlanMilestones")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("masterPlanId", masterPlanId), log.Object("masterPlanMilestoneIds", masterPlanMilestoneIds))

	masterPlanMilestoneExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), masterPlanId, neo4jutil.NodeLabelMasterPlan)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	if !masterPlanMilestoneExists {
		err = errors.New(fmt.Sprintf("Master plan with id {%s} not found", masterPlanId))
		tracing.TraceErr(span, err)
		return err
	}

	grpcRequest := masterplanpb.ReorderMasterPlanMilestonesGrpcRequest{
		Tenant:                 common.GetTenantFromContext(ctx),
		MasterPlanId:           masterPlanId,
		LoggedInUserId:         common.GetUserIdFromContext(ctx),
		AppSource:              constants.AppSourceCustomerOsApi,
		MasterPlanMilestoneIds: masterPlanMilestoneIds,
	}
	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	_, err = utils.CallEventsPlatformGRPCWithRetry[*masterplanpb.MasterPlanIdGrpcResponse](func() (*masterplanpb.MasterPlanIdGrpcResponse, error) {
		return s.grpcClients.MasterPlanClient.ReorderMasterPlanMilestones(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return err
	}
	return nil
}

func (s *masterPlanService) DuplicateMasterPlanMilestone(ctx context.Context, masterPlanId, sourceMasterPlanMilestoneId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanService.DuplicateMasterPlanMilestone")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("masterPlanId", masterPlanId), log.String("sourceMasterPlanMilestoneId", sourceMasterPlanMilestoneId))

	masterPlanMilestoneDbNode, err := s.repositories.Neo4jRepositories.MasterPlanReadRepository.GetMasterPlanMilestoneByPlanAndId(ctx, common.GetContext(ctx).Tenant, masterPlanId, sourceMasterPlanMilestoneId)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	if masterPlanMilestoneDbNode == nil {
		err = errors.New(fmt.Sprintf("Master plan milestone with id {%s} not found", sourceMasterPlanMilestoneId))
		tracing.TraceErr(span, err)
		return "", err
	}
	souceMasterPlanMilestoneEntity := neo4jmapper.MapDbNodeToMasterPlanMilestoneEntity(masterPlanMilestoneDbNode)
	maxOrder, err := s.repositories.Neo4jRepositories.MasterPlanReadRepository.GetMaxOrderForMasterPlanMilestones(ctx, common.GetContext(ctx).Tenant, masterPlanId)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	grpcRequest := masterplanpb.CreateMasterPlanMilestoneGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		MasterPlanId:   masterPlanId,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		Name:           souceMasterPlanMilestoneEntity.Name,
		Order:          maxOrder + 1,
		DurationHours:  souceMasterPlanMilestoneEntity.DurationHours,
		Optional:       souceMasterPlanMilestoneEntity.Optional,
		Items:          souceMasterPlanMilestoneEntity.Items,
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: constants.AppSourceCustomerOsApi,
		},
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*masterplanpb.MasterPlanMilestoneIdGrpcResponse](func() (*masterplanpb.MasterPlanMilestoneIdGrpcResponse, error) {
		return s.grpcClients.MasterPlanClient.CreateMasterPlanMilestone(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, neo4jutil.NodeLabelMasterPlanMilestone, span)
	return response.Id, nil
}

func (s *masterPlanService) DuplicateMasterPlan(ctx context.Context, sourceMasterPlanId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanService.DuplicateMasterPlan")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.String("sourceMasterPlanId", sourceMasterPlanId))

	sourceMasterPlanEntity, err := s.GetMasterPlanById(ctx, sourceMasterPlanId)
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}
	if sourceMasterPlanEntity == nil || sourceMasterPlanEntity.Id == "" {
		err = errors.New(fmt.Sprintf("Master plan with id {%s} not found", sourceMasterPlanId))
		tracing.TraceErr(span, err)
		return "", err
	}
	masterPlanMilestoneEntities, err := s.GetMasterPlanMilestonesForMasterPlans(ctx, []string{sourceMasterPlanId})
	if err != nil {
		tracing.TraceErr(span, err)
		return "", err
	}

	grpcRequest := masterplanpb.CreateMasterPlanGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		Name:           sourceMasterPlanEntity.Name,
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: constants.AppSourceCustomerOsApi,
		},
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*masterplanpb.MasterPlanIdGrpcResponse](func() (*masterplanpb.MasterPlanIdGrpcResponse, error) {
		return s.grpcClients.MasterPlanClient.CreateMasterPlan(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, neo4jutil.NodeLabelMasterPlan, span)

	for _, masterPlanMilestoneEntity := range *masterPlanMilestoneEntities {
		if !masterPlanMilestoneEntity.Retired {
			grpcRequestCreateMilestone := masterplanpb.CreateMasterPlanMilestoneGrpcRequest{
				Tenant:         common.GetTenantFromContext(ctx),
				MasterPlanId:   response.Id,
				LoggedInUserId: common.GetUserIdFromContext(ctx),
				Name:           masterPlanMilestoneEntity.Name,
				Order:          masterPlanMilestoneEntity.Order,
				DurationHours:  masterPlanMilestoneEntity.DurationHours,
				Optional:       masterPlanMilestoneEntity.Optional,
				Items:          masterPlanMilestoneEntity.Items,
				SourceFields: &commonpb.SourceFields{
					Source:    neo4jentity.DataSourceOpenline.String(),
					AppSource: constants.AppSourceCustomerOsApi,
				},
			}
			_, err = utils.CallEventsPlatformGRPCWithRetry[*masterplanpb.MasterPlanMilestoneIdGrpcResponse](func() (*masterplanpb.MasterPlanMilestoneIdGrpcResponse, error) {
				return s.grpcClients.MasterPlanClient.CreateMasterPlanMilestone(ctx, &grpcRequestCreateMilestone)
			})
			if err != nil {
				tracing.TraceErr(span, err)
				s.log.Errorf("Error from events processing: %s", err.Error())
			}
		}
	}

	return response.Id, nil
}

func (s *masterPlanService) CreateDefaultMasterPlan(ctx context.Context) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanService.CreateDefaultMasterPlan")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	grpcRequest := masterplanpb.CreateMasterPlanGrpcRequest{
		Tenant:         common.GetTenantFromContext(ctx),
		LoggedInUserId: common.GetUserIdFromContext(ctx),
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: constants.AppSourceCustomerOsApi,
		},
		Name: "Default Master Plan",
	}

	ctx = tracing.InjectSpanContextIntoGrpcMetadata(ctx, span)
	response, err := utils.CallEventsPlatformGRPCWithRetry[*masterplanpb.MasterPlanIdGrpcResponse](func() (*masterplanpb.MasterPlanIdGrpcResponse, error) {
		return s.grpcClients.MasterPlanClient.CreateMasterPlan(ctx, &grpcRequest)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	neo4jrepository.WaitForNodeCreatedInNeo4j(ctx, s.repositories.Neo4jRepositories, response.Id, neo4jutil.NodeLabelMasterPlan, span)
	mid := response.Id

	milestones := []map[string]any{
		{"text": "Kickoff activities", "order": 0, "durationHours": 24, "items": []string{
			"Identify success criteria",
			"Validate onboarding plan",
			"Get list of admins",
			"Get list of users",
			"Validate training plan",
		}},
		{"text": "Add users to organization", "order": 1, "durationHours": 72, "items": []string{}},
		{"text": "Admin training session", "order": 2, "durationHours": 120, "items": []string{}},
		{"text": "User training session", "order": 3, "durationHours": 168, "items": []string{}},
		{"text": "Office hours", "order": 4, "durationHours": 192, "items": []string{}},
	}

	for _, milestone := range milestones {
		text := milestone["text"].(string)
		order := milestone["order"].(int64)
		durationHours := milestone["durationHours"].(int64)
		items := milestone["items"].([]string)
		grpcRequestCreateMilestone := newDefaultMasterPlanMilestone(mid, text, common.GetUserIdFromContext(ctx), common.GetTenantFromContext(ctx), order, items, durationHours)
		_, err = utils.CallEventsPlatformGRPCWithRetry[*masterplanpb.MasterPlanMilestoneIdGrpcResponse](func() (*masterplanpb.MasterPlanMilestoneIdGrpcResponse, error) {
			return s.grpcClients.MasterPlanClient.CreateMasterPlanMilestone(ctx, &grpcRequestCreateMilestone)
		})
		if err != nil {
			tracing.TraceErr(span, err)
			s.log.Errorf("Error from events processing: %s", err.Error())
		}
	}
	return response.Id, nil
}

func newDefaultMasterPlanMilestone(
	mid,
	name,
	loggedInUserId,
	tenant string,
	order int64,
	items []string,
	durationHours int64,
) masterplanpb.CreateMasterPlanMilestoneGrpcRequest {
	return masterplanpb.CreateMasterPlanMilestoneGrpcRequest{
		Name:           name,
		Optional:       false,
		Order:          order,
		DurationHours:  durationHours,
		Items:          items,
		MasterPlanId:   mid,
		LoggedInUserId: loggedInUserId,
		SourceFields: &commonpb.SourceFields{
			Source:    neo4jentity.DataSourceOpenline.String(),
			AppSource: constants.AppSourceCustomerOsApi,
		},
		Tenant: tenant,
	}
}
