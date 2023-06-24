package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type FieldSetTemplateService interface {
	FindAll(ctx context.Context, entityTemplateId string) (*entity.FieldSetTemplateEntities, error)
	FindLinkedWithFieldSet(ctx context.Context, fieldSetId string) (*entity.FieldSetTemplateEntity, error)
}

type fieldSetTemplateService struct {
	log          logger.Logger
	repositories *repository.Repositories
}

func NewFieldSetTemplateService(log logger.Logger, repositories *repository.Repositories) FieldSetTemplateService {
	return &fieldSetTemplateService{
		log:          log,
		repositories: repositories,
	}
}

func (s *fieldSetTemplateService) FindAll(ctx context.Context, entityTemplateId string) (*entity.FieldSetTemplateEntities, error) {
	all, err := s.repositories.FieldSetTemplateRepository.FindAllByEntityTemplateId(ctx, entityTemplateId)
	if err != nil {
		return nil, err
	}
	fieldSetTemplateEntities := entity.FieldSetTemplateEntities{}
	for _, dbRecord := range all.([]*db.Record) {
		fieldSetTemplateEntities = append(fieldSetTemplateEntities, *s.mapDbNodeToFieldSetTemplate(dbRecord.Values[0].(dbtype.Node)))
	}
	return &fieldSetTemplateEntities, nil
}

func (s *fieldSetTemplateService) FindLinkedWithFieldSet(ctx context.Context, fieldSetId string) (*entity.FieldSetTemplateEntity, error) {
	queryResult, err := s.repositories.FieldSetTemplateRepository.FindByFieldSetId(ctx, fieldSetId)
	if err != nil {
		return nil, err
	}
	if len(queryResult.([]*db.Record)) == 0 {
		return nil, nil
	}
	return s.mapDbNodeToFieldSetTemplate((queryResult.([]*db.Record))[0].Values[0].(dbtype.Node)), nil
}

func (s *fieldSetTemplateService) mapDbNodeToFieldSetTemplate(dbNode dbtype.Node) *entity.FieldSetTemplateEntity {
	props := utils.GetPropsFromNode(dbNode)
	fieldSetTemplate := entity.FieldSetTemplateEntity{
		Id:        utils.GetStringPropOrEmpty(props, "id"),
		Name:      utils.GetStringPropOrEmpty(props, "name"),
		Order:     utils.GetIntPropOrMinusOne(props, "order"),
		CreatedAt: utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt: utils.GetTimePropOrEpochStart(props, "updatedAt"),
	}
	return &fieldSetTemplate
}
