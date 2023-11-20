package service

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type OpportunityService interface {
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
		Id:                utils.GetStringPropOrEmpty(props, "id"),
		Name:              utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:         utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:         utils.GetTimePropOrEpochStart(props, "updatedAt"),
		InternalStage:     entity.GetInternalStage(utils.GetStringPropOrEmpty(props, "internalStage")),
		ExternalStage:     utils.GetStringPropOrEmpty(props, "externalStage"),
		InternalType:      entity.GetInternalType(utils.GetStringPropOrEmpty(props, "internalType")),
		ExternalType:      utils.GetStringPropOrEmpty(props, "externalType"),
		Amount:            utils.GetFloatPropOrZero(props, "amount"),
		MaxAmount:         utils.GetFloatPropOrZero(props, "maxAmount"),
		EstimatedClosedAt: utils.GetTimePropOrEpochStart(props, "estimatedClosedAt"),
		NextSteps:         utils.GetStringPropOrEmpty(props, "nextSteps"),
		GeneralNotes:      utils.GetStringPropOrEmpty(props, "generalNotes"),
		Source:            entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth:     entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
		AppSource:         utils.GetStringPropOrEmpty(props, "appSource"),
	}
	return &opportunity
}
