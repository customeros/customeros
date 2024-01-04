package service

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	masterplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/master_plan"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

type MasterPlanService interface {
	CreateMasterPlan(ctx context.Context, name string) (string, error)
	GetMasterPlanById(ctx context.Context, masterPlanId string) (*neo4jentity.MasterPlanEntity, error)
	GetMasterPlans(ctx context.Context, returnRetired *bool) (*neo4jentity.MasterPlanEntities, error)

	CreateMasterPlanMilestone(ctx context.Context, masterPlanId, name string, order, durationHours int64, optional bool, items []string) (string, error)
	GetMasterPlanMilestoneById(ctx context.Context, masterPlanMilestoneId string) (*neo4jentity.MasterPlanMilestoneEntity, error)
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
	response, err := s.grpcClients.MasterPlanClient.CreateMasterPlan(ctx, &grpcRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	for i := 1; i <= constants.MaxRetriesCheckDataInNeo4jAfterEventRequest; i++ {
		contractFound, findErr := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), response.Id, neo4jentity.NodeLabel_MasterPlan)
		if contractFound && findErr == nil {
			span.LogFields(log.Bool("response - master plan saved in db", true))
			break
		}
		time.Sleep(utils.BackOffIncrementalDelay(i))
	}

	span.LogFields(log.String("response - created masterPlanId", response.Id))
	return response.Id, nil
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

	masterPlanExists, err := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), masterPlanId, neo4jentity.NodeLabel_MasterPlan)
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
	response, err := s.grpcClients.MasterPlanClient.CreateMasterPlanMilestone(ctx, &grpcRequest)
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Errorf("Error from events processing: %s", err.Error())
		return "", err
	}

	for i := 1; i <= constants.MaxRetriesCheckDataInNeo4jAfterEventRequest; i++ {
		contractFound, findErr := s.repositories.Neo4jRepositories.CommonReadRepository.ExistsById(ctx, common.GetTenantFromContext(ctx), response.Id, neo4jentity.NodeLabel_MasterPlanMilestone)
		if contractFound && findErr == nil {
			span.LogFields(log.Bool("response - master plan milestone saved in db", true))
			break
		}
		time.Sleep(utils.BackOffIncrementalDelay(i))
	}

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
