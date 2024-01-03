package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
)

type FieldSetService interface {
	FindAll(ctx context.Context, obj *model.CustomFieldEntityType) (*entity.FieldSetEntities, error)
	MergeFieldSetToContact(ctx context.Context, contactId string, input *entity.FieldSetEntity) (*entity.FieldSetEntity, error)
	UpdateFieldSetInContact(ctx context.Context, contactId string, input *entity.FieldSetEntity) (*entity.FieldSetEntity, error)
	DeleteByIdFromContact(ctx context.Context, contactId string, fieldSetId string) (bool, error)
}

type fieldSetService struct {
	log        logger.Logger
	repository *repository.Repositories
}

func NewFieldSetService(log logger.Logger, repository *repository.Repositories) FieldSetService {
	return &fieldSetService{
		log:        log,
		repository: repository,
	}
}

func (s *fieldSetService) getDriver() neo4j.DriverWithContext {
	return *s.repository.Drivers.Neo4jDriver
}

func (s *fieldSetService) FindAll(ctx context.Context, obj *model.CustomFieldEntityType) (*entity.FieldSetEntities, error) {
	session := utils.NewNeo4jReadSession(ctx, s.getDriver())
	defer session.Close(ctx)

	dbRecords, err := s.repository.FieldSetRepository.FindAll(ctx, session, common.GetContext(ctx).Tenant, obj)
	if err != nil {
		return nil, err
	}

	fieldSetEntities := entity.FieldSetEntities{}

	for _, dbRecord := range dbRecords {
		fieldSetEntity := s.mapDbNodeToFieldSetEntity(dbRecord.Values[0].(dbtype.Node))
		fieldSetEntities = append(fieldSetEntities, *fieldSetEntity)
	}

	return &fieldSetEntities, nil
}

func (s *fieldSetService) MergeFieldSetToContact(ctx context.Context, contactId string, entity *entity.FieldSetEntity) (*entity.FieldSetEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)

	var fieldSetDbNode *dbtype.Node
	entityType := &model.CustomFieldEntityType{
		ID:         contactId,
		EntityType: model.EntityTypeOrganization,
	}
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		var err error
		fieldSetDbNode, err = s.repository.FieldSetRepository.MergeFieldSetToContactInTx(ctx, tx, common.GetContext(ctx).Tenant, contactId, *entity)
		if err != nil {
			return nil, err
		}
		var fieldSetId = utils.GetPropsFromNode(*fieldSetDbNode)["id"].(string)
		if entity.TemplateId != nil {
			err := s.repository.FieldSetRepository.LinkWithFieldSetTemplateInTx(ctx, tx, common.GetContext(ctx).Tenant, fieldSetId, *entity.TemplateId, entityType.EntityType)
			if err != nil {
				return nil, err
			}
		}
		if entity.CustomFields != nil {
			for _, customField := range *entity.CustomFields {
				dbNode, err := s.repository.CustomFieldRepository.MergeCustomFieldToFieldSetInTx(ctx, tx, common.GetContext(ctx).Tenant, entityType, fieldSetId, customField)
				if err != nil {
					return nil, err
				}
				if customField.TemplateId != nil {
					var fieldId = utils.GetPropsFromNode(*dbNode)["id"].(string)
					err := s.repository.CustomFieldRepository.LinkWithCustomFieldTemplateForFieldSetInTx(ctx, tx, fieldId, fieldSetId, *customField.TemplateId)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	var fieldSetEntity = s.mapDbNodeToFieldSetEntity(*fieldSetDbNode)
	return fieldSetEntity, nil
}

func (s *fieldSetService) UpdateFieldSetInContact(ctx context.Context, contactId string, entity *entity.FieldSetEntity) (*entity.FieldSetEntity, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)

	var fieldSetDbNode *dbtype.Node

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		var err error
		fieldSetDbNode, err = s.repository.FieldSetRepository.UpdateFieldSetForContactInTx(ctx, tx, common.GetContext(ctx).Tenant, contactId, *entity)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	var fieldSetEntity = s.mapDbNodeToFieldSetEntity(*fieldSetDbNode)
	return fieldSetEntity, nil
}

func (s *fieldSetService) DeleteByIdFromContact(ctx context.Context, contactId string, fieldSetId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(ctx, s.getDriver())
	defer session.Close(ctx)

	err := s.repository.FieldSetRepository.DeleteByIdFromContact(ctx, session, common.GetContext(ctx).Tenant, contactId, fieldSetId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *fieldSetService) mapDbNodeToFieldSetEntity(node dbtype.Node) *entity.FieldSetEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.FieldSetEntity{
		Id:            utils.StringPtr(utils.GetStringPropOrEmpty(props, "id")),
		Name:          utils.GetStringPropOrEmpty(props, "name"),
		CreatedAt:     utils.GetTimePropOrEpochStart(props, "createdAt"),
		UpdatedAt:     utils.GetTimePropOrEpochStart(props, "updatedAt"),
		Source:        neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: neo4jentity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &result
}
