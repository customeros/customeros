package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ExternalSystemService interface {
	GetExternalSystemsFor(ctx context.Context, ids []string, entityType entity.EntityType) (*entity.ExternalSystemEntities, error)
}

type externalSystemService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewExternalSystemService(log logger.Logger, repositories *repository.Repositories) ExternalSystemService {
	return &externalSystemService{
		log:          log,
		repositories: repositories,
	}
}

func (s *externalSystemService) GetExternalSystemsFor(ctx context.Context, ids []string, entityType entity.EntityType) (*entity.ExternalSystemEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemService.GetExternalSystemsFor")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	dbRecords, err := s.repositories.ExternalSystemRepository.GetFor(ctx, common.GetTenantFromContext(ctx), ids, entityType.Neo4jLabel())
	if err != nil {
		return nil, err
	}
	externalSystemEntities := make(entity.ExternalSystemEntities, 0, len(dbRecords))
	for _, v := range dbRecords {
		externalSystemEntity := s.mapDbNodeToExternalSystemEntity(*v.Node)
		s.addDbRelationshipToExternalSystemEntity(*v.Relationship, externalSystemEntity)
		externalSystemEntity.DataloaderKey = v.LinkedNodeId
		externalSystemEntities = append(externalSystemEntities, *externalSystemEntity)
	}
	return &externalSystemEntities, nil
}

func (s *externalSystemService) mapDbNodeToExternalSystemEntity(dbNode dbtype.Node) *entity.ExternalSystemEntity {
	props := utils.GetPropsFromNode(dbNode)
	return &entity.ExternalSystemEntity{
		ExternalSystemId: entity.ExternalSystemTypeFromString(utils.GetStringPropOrEmpty(props, "id")),
	}
}

func (s *externalSystemService) addDbRelationshipToExternalSystemEntity(relationship dbtype.Relationship, entity *entity.ExternalSystemEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	entity.Relationship.SyncDate = utils.GetTimePropOrNil(props, "syncDate")
	entity.Relationship.ExternalId = utils.GetStringPropOrEmpty(props, "externalId")
	entity.Relationship.ExternalUrl = utils.GetStringPropOrNil(props, "externalUrl")
	entity.Relationship.ExternalSource = utils.GetStringPropOrNil(props, "externalSource")
}
