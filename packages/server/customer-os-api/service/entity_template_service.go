package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
)

type EntityTemplateService interface {
	Create(ctx context.Context, entity *neo4jentity.EntityTemplateEntity) (*neo4jentity.EntityTemplateEntity, error)
	FindAll(ctx context.Context, extends string) (*neo4jentity.EntityTemplateEntities, error)
	FindLinked(ctx context.Context, obj *model.CustomFieldEntityType) (*neo4jentity.EntityTemplateEntity, error)
}

type entityTemplateService struct {
	log        logger.Logger
	repository *repository.Repositories
}

func NewEntityTemplateService(log logger.Logger, repository *repository.Repositories) EntityTemplateService {
	return &entityTemplateService{
		log:        log,
		repository: repository,
	}
}

func (s *entityTemplateService) Create(ctx context.Context, entity *neo4jentity.EntityTemplateEntity) (*neo4jentity.EntityTemplateEntity, error) {
	record, err := s.repository.EntityTemplateRepository.Create(ctx, common.GetContext(ctx).Tenant, entity)
	if err != nil {
		return nil, err
	}
	entityTemplate := s.mapDbNodeToEntityTemplate((record.([]*db.Record)[0]).Values[0].(dbtype.Node))
	return entityTemplate, nil
}

func (s *entityTemplateService) FindAll(ctx context.Context, extends *string) (*neo4jentity.EntityTemplateEntities, error) {
	session := utils.NewNeo4jReadSession(ctx, *s.repository.Drivers.Neo4jDriver)
	defer session.Close(ctx)

	var err error
	var entityTemplatesDbRecords []*db.Record
	if extends == nil {
		entityTemplatesDbRecords, err = s.repository.EntityTemplateRepository.FindAllByTenant(ctx, session, common.GetContext(ctx).Tenant)
	} else {
		entityTemplatesDbRecords, err = s.repository.EntityTemplateRepository.FindAllByTenantAndExtends(ctx, session, common.GetContext(ctx).Tenant, *extends)
	}
	if err != nil {
		return nil, err
	}
	entityTemplateEntities := neo4jentity.EntityTemplateEntities{}
	for _, dbRecord := range entityTemplatesDbRecords {
		entityTemplateEntity := neo4jmapper.MapDbNodeToEntityTemplate(utils.NodePtr(dbRecord.Values[0].(dbtype.Node)))
		entityTemplateEntities = append(entityTemplateEntities, *entityTemplateEntity)
	}
	return &entityTemplateEntities, nil
}

func (s *entityTemplateService) FindLinked(ctx context.Context, obj *model.CustomFieldEntityType) (*neo4jentity.EntityTemplateEntity, error) {
	queryResult, err := s.repository.EntityTemplateRepository.FindById(ctx, common.GetContext(ctx).Tenant, obj)
	if err != nil {
		return nil, err
	}
	if len(queryResult.([]*db.Record)) == 0 {
		return nil, nil
	}
	entityTemplateEntity := neo4jmapper.MapDbNodeToEntityTemplate(utils.NodePtr((queryResult.([]*db.Record))[0].Values[0].(dbtype.Node)))
	return entityTemplateEntity, nil
}
