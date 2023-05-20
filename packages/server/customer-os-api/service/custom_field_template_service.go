package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type CustomFieldTemplateService interface {
	FindAllForEntityTemplate(ctx context.Context, entityTemplateId string) (*entity.CustomFieldTemplateEntities, error)
	FindAllForFieldSetTemplate(ctx context.Context, fieldSetTemplateId string) (*entity.CustomFieldTemplateEntities, error)
	FindLinkedWithCustomField(ctx context.Context, customFieldId string) (*entity.CustomFieldTemplateEntity, error)
}

type customFieldTemplateService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewCustomFieldTemplateService(log logger.Logger, repositories *repository.Repositories) CustomFieldTemplateService {
	return &customFieldTemplateService{
		log:          log,
		repositories: repositories,
	}
}

func (s *customFieldTemplateService) FindAllForEntityTemplate(ctx context.Context, entityTemplateId string) (*entity.CustomFieldTemplateEntities, error) {
	all, err := s.repositories.CustomFieldTemplateRepository.FindAllByEntityTemplateId(ctx, entityTemplateId)
	if err != nil {
		return nil, err
	}
	customFieldTemplateEntities := entity.CustomFieldTemplateEntities{}
	for _, dbRecord := range all.([]*db.Record) {
		customFieldTemplateEntities = append(customFieldTemplateEntities, *s.mapDbNodeToCustomFieldTemplate(dbRecord.Values[0].(dbtype.Node)))
	}
	return &customFieldTemplateEntities, nil
}

func (s *customFieldTemplateService) FindAllForFieldSetTemplate(ctx context.Context, fieldSetTemplateId string) (*entity.CustomFieldTemplateEntities, error) {
	all, err := s.repositories.CustomFieldTemplateRepository.FindAllByEntityFieldSetTemplateId(ctx, fieldSetTemplateId)
	if err != nil {
		return nil, err
	}
	customFieldTemplateEntities := entity.CustomFieldTemplateEntities{}
	for _, dbRecord := range all.([]*db.Record) {
		customFieldTemplateEntities = append(customFieldTemplateEntities, *s.mapDbNodeToCustomFieldTemplate(dbRecord.Values[0].(dbtype.Node)))
	}
	return &customFieldTemplateEntities, nil
}

func (s *customFieldTemplateService) FindLinkedWithCustomField(ctx context.Context, customFieldId string) (*entity.CustomFieldTemplateEntity, error) {
	queryResult, err := s.repositories.CustomFieldTemplateRepository.FindByCustomFieldId(ctx, customFieldId)
	if err != nil {
		return nil, err
	}
	if len(queryResult.([]*db.Record)) == 0 {
		return nil, nil
	}
	return s.mapDbNodeToCustomFieldTemplate((queryResult.([]*db.Record))[0].Values[0].(dbtype.Node)), nil
}

func (s *customFieldTemplateService) mapDbNodeToCustomFieldTemplate(dbNode dbtype.Node) *entity.CustomFieldTemplateEntity {
	props := utils.GetPropsFromNode(dbNode)
	customFieldTemplate := entity.CustomFieldTemplateEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Name:      utils.GetStringPropOrEmpty(props, "name"),
		Order:     utils.GetIntPropOrMinusOne(props, "order"),
		Mandatory: utils.GetBoolPropOrFalse(props, "mandatory"),
		Type:      utils.GetStringPropOrEmpty(props, "type"),
		Length:    utils.GetIntPropOrNil(props, "length"),
		Min:       utils.GetIntPropOrNil(props, "min"),
		Max:       utils.GetIntPropOrNil(props, "max"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
	}
	return &customFieldTemplate
}
