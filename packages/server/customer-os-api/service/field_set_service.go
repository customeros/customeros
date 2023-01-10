package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type FieldSetService interface {
	FindAllForContact(ctx context.Context, contact *model.Contact) (*entity.FieldSetEntities, error)
	MergeFieldSetToContact(ctx context.Context, contactId string, input *entity.FieldSetEntity) (*entity.FieldSetEntity, error)
	UpdateFieldSetInContact(ctx context.Context, contactId string, input *entity.FieldSetEntity) (*entity.FieldSetEntity, error)
	DeleteByIdFromContact(ctx context.Context, contactId string, fieldSetId string) (bool, error)
}

type fieldSetService struct {
	repository *repository.Repositories
}

func NewFieldSetService(repository *repository.Repositories) FieldSetService {
	return &fieldSetService{
		repository: repository,
	}
}

func (s *fieldSetService) getDriver() neo4j.Driver {
	return *s.repository.Drivers.Neo4jDriver
}

func (s *fieldSetService) FindAllForContact(ctx context.Context, contact *model.Contact) (*entity.FieldSetEntities, error) {
	session := utils.NewNeo4jReadSession(s.getDriver())
	defer session.Close()

	dbRecords, err := s.repository.FieldSetRepository.FindAllForContact(session, common.GetContext(ctx).Tenant, contact.ID)
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
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	var fieldSetDbNode *dbtype.Node

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		var err error
		fieldSetDbNode, err = s.repository.FieldSetRepository.MergeFieldSetToContactInTx(tx, common.GetContext(ctx).Tenant, contactId, *entity)
		if err != nil {
			return nil, err
		}
		var fieldSetId = utils.GetPropsFromNode(*fieldSetDbNode)["id"].(string)
		if entity.TemplateId != nil {
			err := s.repository.FieldSetRepository.LinkWithFieldSetTemplateInTx(tx, common.GetContext(ctx).Tenant, fieldSetId, *entity.TemplateId)
			if err != nil {
				return nil, err
			}
		}
		if entity.CustomFields != nil {
			for _, customField := range *entity.CustomFields {
				dbNode, err := s.repository.CustomFieldRepository.MergeCustomFieldToFieldSetInTx(tx, common.GetContext(ctx).Tenant, contactId, fieldSetId, customField)
				if err != nil {
					return nil, err
				}
				if customField.TemplateId != nil {
					var fieldId = utils.GetPropsFromNode(*dbNode)["id"].(string)
					err := s.repository.CustomFieldRepository.LinkWithCustomFieldTemplateForFieldSetInTx(tx, fieldId, fieldSetId, *customField.TemplateId)
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
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	var fieldSetDbNode *dbtype.Node

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		var err error
		fieldSetDbNode, err = s.repository.FieldSetRepository.UpdateFieldSetForContactInTx(tx, common.GetContext(ctx).Tenant, contactId, *entity)
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
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	err := s.repository.FieldSetRepository.DeleteByIdFromContact(session, common.GetContext(ctx).Tenant, contactId, fieldSetId)
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
		CreatedAt:     utils.GetTimePropOrNow(props, "createdAt"),
		Source:        entity.GetDataSource(utils.GetStringPropOrEmpty(props, "source")),
		SourceOfTruth: entity.GetDataSource(utils.GetStringPropOrEmpty(props, "sourceOfTruth")),
	}
	return &result
}
