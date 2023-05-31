package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type OrganizationRelationshipService interface {
	GetRelationshipsForOrganizations(ctx context.Context, organizationIDs []string) (*entity.OrganizationRelationshipEntities, error)
	GetRelationshipsWithStagesForOrganizations(ctx context.Context, organizationIDs []string) (*entity.OrganizationRelationshipsWithStages, error)
}

type organizationRelationshipService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewOrganizationRelationshipService(log logger.Logger, repositories *repository.Repositories) OrganizationRelationshipService {
	return &organizationRelationshipService{
		log:          log,
		repositories: repositories,
	}
}

func (s *organizationRelationshipService) GetRelationshipsForOrganizations(ctx context.Context, organizationIDs []string) (*entity.OrganizationRelationshipEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRelationshipService.GetRelationshipsForOrganizations")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, common.GetTenantFromContext(ctx))
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)
	span.LogFields(log.Object("organizationIDs", organizationIDs))

	organizationRelationships, err := s.repositories.OrganizationRelationshipRepository.GetOrganizationRelationshipsForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIDs)
	if err != nil {
		return nil, err
	}
	organizationRelationshipEntities := entity.OrganizationRelationshipEntities{}
	for _, v := range organizationRelationships {
		organizationRelationshipEntity := s.mapDbNodeToOrganizationRelationshipEntity(*v.Node)
		organizationRelationshipEntity.DataloaderKey = v.LinkedNodeId
		organizationRelationshipEntities = append(organizationRelationshipEntities, *organizationRelationshipEntity)
	}
	return &organizationRelationshipEntities, nil
}

func (s *organizationRelationshipService) GetRelationshipsWithStagesForOrganizations(ctx context.Context, organizationIDs []string) (*entity.OrganizationRelationshipsWithStages, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OrganizationRelationshipService.GetRelationshipsWithStagesForOrganizations")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, common.GetTenantFromContext(ctx))
	span.SetTag(tracing.SpanTagComponent, constants.ComponentService)
	span.LogFields(log.Object("organizationIDs", organizationIDs))

	dbResults, err := s.repositories.OrganizationRelationshipRepository.GetOrganizationRelationshipsWithStagesForOrganizations(ctx, common.GetTenantFromContext(ctx), organizationIDs)
	if err != nil {
		return nil, err
	}
	organizationRelationshipsWithStages := entity.OrganizationRelationshipsWithStages{}
	for _, v := range dbResults {
		organizationRelationshipWithStage := entity.OrganizationRelationshipWithStage{
			DataloaderKey: v.LinkedNodeId,
		}
		props := utils.GetPropsFromNode(*v.Pair.First)
		organizationRelationshipWithStage.Relationship = entity.OrganizationRelationshipFromString(utils.GetStringPropOrEmpty(props, "name"))
		if v.Pair.Second != nil {
			organizationRelationshipWithStage.Stage = s.mapDbNodeToOrganizationRelationshipStageEntity(*v.Pair.Second)
		}
		organizationRelationshipsWithStages = append(organizationRelationshipsWithStages, organizationRelationshipWithStage)
	}
	return &organizationRelationshipsWithStages, nil
}

func (s *organizationRelationshipService) mapDbNodeToOrganizationRelationshipEntity(node dbtype.Node) *entity.OrganizationRelationshipEntity {
	props := utils.GetPropsFromNode(node)
	return &entity.OrganizationRelationshipEntity{
		ID:    utils.GetStringPropOrEmpty(props, "id"),
		Name:  utils.GetStringPropOrEmpty(props, "name"),
		Group: utils.GetStringPropOrEmpty(props, "group"),
	}
}

func (s *organizationRelationshipService) mapDbNodeToOrganizationRelationshipStageEntity(node dbtype.Node) *entity.OrganizationRelationshipStageEntity {
	props := utils.GetPropsFromNode(node)
	return &entity.OrganizationRelationshipStageEntity{
		Id:   utils.GetStringPropOrEmpty(props, "id"),
		Name: utils.GetStringPropOrEmpty(props, "name"),
	}
}
