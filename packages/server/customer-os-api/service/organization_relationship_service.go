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

func (s *organizationRelationshipService) mapDbNodeToOrganizationRelationshipEntity(node dbtype.Node) *entity.OrganizationRelationshipEntity {
	props := utils.GetPropsFromNode(node)
	return &entity.OrganizationRelationshipEntity{
		ID:    utils.GetStringPropOrEmpty(props, "id"),
		Name:  utils.GetStringPropOrEmpty(props, "name"),
		Group: utils.GetStringPropOrEmpty(props, "group"),
	}
}
