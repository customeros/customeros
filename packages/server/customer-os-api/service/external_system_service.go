package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ExternalSystemService interface {
	GetExternalSystemsForEntities(ctx context.Context, ids []string, entityType commonModel.EntityType) (*neo4jentity.ExternalSystemEntities, error)
	GetAllExternalSystemInstances(ctx context.Context) (*neo4jentity.ExternalSystemEntities, error)
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

func (s *externalSystemService) GetExternalSystemsForEntities(ctx context.Context, ids []string, entityType commonModel.EntityType) (*neo4jentity.ExternalSystemEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemService.GetExternalSystemsForEntities")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	span.LogFields(log.Object("ids", ids))

	dbRecords, err := s.repositories.ExternalSystemRepository.GetFor(ctx, common.GetTenantFromContext(ctx), ids, entityType.Neo4jLabel())
	if err != nil {
		return nil, err
	}
	externalSystemEntities := make(neo4jentity.ExternalSystemEntities, 0, len(dbRecords))
	for _, v := range dbRecords {
		externalSystemEntity := neo4jmapper.MapDbNodeToExternalSystem(v.Node)
		s.addDbRelationshipToExternalSystemEntity(*v.Relationship, externalSystemEntity)
		externalSystemEntity.DataloaderKey = v.LinkedNodeId
		externalSystemEntities = append(externalSystemEntities, *externalSystemEntity)
	}
	return &externalSystemEntities, nil
}

func (s *externalSystemService) GetAllExternalSystemInstances(ctx context.Context) (*neo4jentity.ExternalSystemEntities, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemService.GetAllExternalSystemInstances")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	dbRecords, err := s.repositories.Neo4jRepositories.ExternalSystemReadRepository.GetAllForTenant(ctx, common.GetTenantFromContext(ctx))
	if err != nil {
		tracing.TraceErr(span, err)
		s.log.Error(ctx, "Error getting all external system instances", err)
		return nil, err
	}
	var externalSystemEntities neo4jentity.ExternalSystemEntities
	for _, v := range dbRecords {
		externalSystemEntity := neo4jmapper.MapDbNodeToExternalSystem(v)
		externalSystemEntities = append(externalSystemEntities, *externalSystemEntity)
	}
	return &externalSystemEntities, nil

}

func (s *externalSystemService) addDbRelationshipToExternalSystemEntity(relationship dbtype.Relationship, entity *neo4jentity.ExternalSystemEntity) {
	props := utils.GetPropsFromRelationship(relationship)
	entity.Relationship.SyncDate = utils.GetTimePropOrNil(props, "syncDate")
	entity.Relationship.ExternalId = utils.GetStringPropOrEmpty(props, "externalId")
	entity.Relationship.ExternalUrl = utils.GetStringPropOrNil(props, "externalUrl")
	entity.Relationship.ExternalSource = utils.GetStringPropOrNil(props, "externalSource")
	entity.Relationship.Primary = utils.GetBoolPropOrFalse(props, "primary")
}
