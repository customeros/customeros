package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type CustomFieldService interface {
	FindAllForContact(ctx context.Context, obj *model.Contact) (*entity.CustomFieldEntities, error)
	FindAllForFieldSet(ctx context.Context, obj *model.FieldSet) (*entity.CustomFieldEntities, error)

	MergeCustomFieldToContact(ctx context.Context, contactId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error)
	MergeCustomFieldToFieldSet(ctx context.Context, contactId string, fieldSetId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error)

	UpdateCustomFieldForContact(ctx context.Context, contactId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error)
	UpdateCustomFieldForFieldSet(ctx context.Context, contactId string, fieldSetId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error)

	DeleteByNameFromContact(ctx context.Context, contactId, fieldName string) (bool, error)
	DeleteByIdFromContact(ctx context.Context, contactId, fieldId string) (bool, error)
	DeleteByIdFromFieldSet(ctx context.Context, contactId, fieldSetId, fieldId string) (bool, error)

	mapDbNodeToCustomFieldEntity(node dbtype.Node) *entity.CustomFieldEntity
	getDriver() neo4j.Driver
}

type customFieldService struct {
	repository *repository.RepositoryContainer
}

func NewCustomFieldService(repository *repository.RepositoryContainer) CustomFieldService {
	return &customFieldService{
		repository: repository,
	}
}

func (s *customFieldService) getDriver() neo4j.Driver {
	return *s.repository.Drivers.Neo4jDriver
}

func (s *customFieldService) FindAllForContact(ctx context.Context, contact *model.Contact) (*entity.CustomFieldEntities, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	dbRecords, err := s.repository.CustomFieldRepository.FindAllForContact(session, common.GetContext(ctx).Tenant, contact.ID)
	if err != nil {
		return nil, err
	}

	customFieldEntities := entity.CustomFieldEntities{}

	for _, dbRecord := range dbRecords {
		customFieldEntity := s.mapDbNodeToCustomFieldEntity(dbRecord.Values[0].(dbtype.Node))
		customFieldEntities = append(customFieldEntities, *customFieldEntity)
	}

	return &customFieldEntities, nil
}

func (s *customFieldService) FindAllForFieldSet(ctx context.Context, fieldSet *model.FieldSet) (*entity.CustomFieldEntities, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	dbRecords, err := s.repository.CustomFieldRepository.FindAllForContact(session, common.GetContext(ctx).Tenant, fieldSet.ID)
	if err != nil {
		return nil, err
	}

	customFieldEntities := entity.CustomFieldEntities{}

	for _, dbRecord := range dbRecords {
		customFieldEntity := s.mapDbNodeToCustomFieldEntity(dbRecord.Values[0].(dbtype.Node))
		customFieldEntities = append(customFieldEntities, *customFieldEntity)
	}

	return &customFieldEntities, nil
}

func (s *customFieldService) MergeCustomFieldToContact(ctx context.Context, contactId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	customFieldNode, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		customFieldDbNode, err := s.repository.CustomFieldRepository.MergeCustomFieldToContactInTx(tx, common.GetContext(ctx).Tenant, contactId, entity)
		if err != nil {
			return nil, err
		}
		if entity.DefinitionId != nil {
			var fieldId = utils.GetPropsFromNode(customFieldDbNode)["id"].(string)
			if err = s.repository.CustomFieldRepository.LinkWithCustomFieldDefinitionForContactInTx(tx, fieldId, contactId, *entity.DefinitionId); err != nil {
				return nil, err
			}
		}
		return customFieldDbNode, err
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToCustomFieldEntity(customFieldNode.(dbtype.Node)), nil
}

func (s *customFieldService) MergeCustomFieldToFieldSet(ctx context.Context, contactId string, fieldSetId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	customFieldNode, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		customFieldNode, err := s.repository.CustomFieldRepository.MergeCustomFieldToFieldSetInTx(tx, common.GetContext(ctx).Tenant, contactId, fieldSetId, entity)
		if err != nil {
			return nil, err
		}
		if entity.DefinitionId != nil {
			var fieldId = utils.GetPropsFromNode(customFieldNode)["id"].(string)
			if err = s.repository.CustomFieldRepository.LinkWithCustomFieldDefinitionForFieldSetInTx(tx, fieldId, fieldSetId, *entity.DefinitionId); err != nil {
				return nil, err
			}
		}
		return customFieldNode, err
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToCustomFieldEntity(customFieldNode.(dbtype.Node)), nil
}

func (s *customFieldService) UpdateCustomFieldForContact(ctx context.Context, contactId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	customFieldDbNode, err := s.repository.CustomFieldRepository.UpdateForContact(session, common.GetContext(ctx).Tenant, contactId, entity)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToCustomFieldEntity(customFieldDbNode), nil
}

func (s *customFieldService) UpdateCustomFieldForFieldSet(ctx context.Context, contactId string, fieldSetId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	customFieldDbNode, err := s.repository.CustomFieldRepository.UpdateForFieldSet(session, common.GetContext(ctx).Tenant, contactId, fieldSetId, entity)
	if err != nil {
		return nil, err
	}
	return s.mapDbNodeToCustomFieldEntity(customFieldDbNode), nil
}

func (s *customFieldService) DeleteByNameFromContact(ctx context.Context, contactId, fieldName string) (bool, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	err := s.repository.CustomFieldRepository.DeleteByNameFromContact(session, common.GetContext(ctx).Tenant, contactId, fieldName)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *customFieldService) DeleteByIdFromContact(ctx context.Context, contactId, fieldId string) (bool, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	err := s.repository.CustomFieldRepository.DeleteByIdFromContact(session, common.GetContext(ctx).Tenant, contactId, fieldId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *customFieldService) DeleteByIdFromFieldSet(ctx context.Context, contactId, fieldSetId, fieldId string) (bool, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	err := s.repository.CustomFieldRepository.DeleteByIdFromFieldSet(session, common.GetContext(ctx).Tenant, contactId, fieldSetId, fieldId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *customFieldService) mapDbNodeToCustomFieldEntity(node dbtype.Node) *entity.CustomFieldEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.CustomFieldEntity{
		Id:       utils.GetStringPropOrEmpty(props, "id"),
		Name:     utils.GetStringPropOrEmpty(props, "name"),
		DataType: utils.GetStringPropOrEmpty(props, "datatype"),
		Value: model.AnyTypeValue{
			Str:   utils.GetStringPropOrNil(props, entity.CustomFieldTextProperty.String()),
			Time:  utils.GetTimePropOrNil(props, entity.CustomFieldTimeProperty.String()),
			Int:   utils.GetIntPropOrNil(props, entity.CustomFieldIntProperty.String()),
			Float: utils.GetFloatPropOrNil(props, entity.CustomFieldFloatProperty.String()),
			Bool:  utils.GetBoolPropOrNil(props, entity.CustomFieldBoolProperty.String()),
		},
	}
	return &result
}
