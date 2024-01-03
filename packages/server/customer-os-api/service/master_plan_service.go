package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	commonpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/common"
	masterplanpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/master_plan"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"time"
)

type MasterPlanService interface {
	Create(ctx context.Context, name string) (string, error)
	GetById(ctx context.Context, masterPlanId string) (*neo4jentity.MasterPlanEntity, error)
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

func (s *masterPlanService) Create(ctx context.Context, name string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MasterPlanService.Create")
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
			span.LogFields(log.Bool("output - master plan saved in db", true))
			break
		}
		time.Sleep(utils.BackOffIncrementalDelay(i))
	}

	span.LogFields(log.String("output - createdMasterPlanId", response.Id))
	return response.Id, nil
}

func (s *masterPlanService) GetById(ctx context.Context, masterPlanId string) (*neo4jentity.MasterPlanEntity, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContractService.GetById")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.SetTag(tracing.SpanTagEntityId, masterPlanId)

	if masterPlanDbNode, err := s.repositories.Neo4jRepositories.MasterPlanReadRepository.GetMasterPlanById(ctx, common.GetContext(ctx).Tenant, masterPlanId); err != nil {
		tracing.TraceErr(span, err)
		wrappedErr := errors.Wrap(err, fmt.Sprintf("Master plan with id {%s} not found", masterPlanId))
		return nil, wrappedErr
	} else {
		return s.mapDbNodeToMasterPlanEntity(masterPlanDbNode), nil
	}
}

func (s *masterPlanService) mapDbNodeToMasterPlanEntity(dbNode *dbtype.Node) *neo4jentity.MasterPlanEntity {
	if dbNode == nil {
		return nil
	}
	props := utils.GetPropsFromNode(*dbNode)
	masterPlan := neo4jentity.MasterPlanEntity{
		Id:            utils.GetStringPropOrEmpty(props, "id"),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:     utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &masterPlan
}
