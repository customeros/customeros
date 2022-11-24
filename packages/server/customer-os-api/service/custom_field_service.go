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

	MergeAndUpdateCustomFieldsForContact(ctx context.Context, contactId string, customFields *entity.CustomFieldEntities, fieldSets *entity.FieldSetEntities) error

	MergeCustomFieldToContact(ctx context.Context, contactId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error)
	MergeCustomFieldToFieldSet(ctx context.Context, contactId string, fieldSetId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error)

	UpdateCustomFieldForContact(ctx context.Context, contactId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error)
	UpdateCustomFieldForFieldSet(ctx context.Context, contactId string, fieldSetId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error)

	DeleteByNameFromContact(ctx context.Context, contactId, fieldName string) (bool, error)
	DeleteByIdFromContact(ctx context.Context, contactId, fieldId string) (bool, error)
	DeleteByIdFromFieldSet(ctx context.Context, contactId, fieldSetId, fieldId string) (bool, error)

	mapDbNodeToCustomFieldEntity(node dbtype.Node) *entity.CustomFieldEntity
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

func (s *customFieldService) MergeAndUpdateCustomFieldsForContact(ctx context.Context, contactId string, customFields *entity.CustomFieldEntities, fieldSets *entity.FieldSetEntities) error {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		tenant := common.GetContext(ctx).Tenant
		if customFields != nil {
			for _, customField := range *customFields {
				if customField.Id == nil {
					dbNode, err := s.repository.CustomFieldRepository.MergeCustomFieldToContactInTx(tx, tenant, contactId, &customField)
					if err != nil {
						return nil, err
					}
					if customField.DefinitionId != nil {
						var fieldId = utils.GetPropsFromNode(*dbNode)["id"].(string)
						err := s.repository.CustomFieldRepository.LinkWithCustomFieldDefinitionForContactInTx(tx, fieldId, contactId, *customField.DefinitionId)
						if err != nil {
							return nil, err
						}
					}
				} else {
					_, err := s.repository.CustomFieldRepository.UpdateForContactInTx(tx, tenant, contactId, &customField)
					if err != nil {
						return nil, err
					}
				}
			}
		}
		if fieldSets != nil {
			for _, fieldSet := range *fieldSets {
				var fieldSetId string
				if fieldSet.Id == nil {
					setDbNode, _, err := s.repository.FieldSetRepository.MergeFieldSetToContactInTx(tx, tenant, contactId, fieldSet)
					if err != nil {
						return nil, err
					}
					fieldSetId = utils.GetPropsFromNode(*setDbNode)["id"].(string)
					if fieldSet.DefinitionId != nil {
						err := s.repository.FieldSetRepository.LinkWithFieldSetDefinitionInTx(tx, tenant, fieldSetId, *fieldSet.DefinitionId)
						if err != nil {
							return nil, err
						}
					}
				} else {
					fieldSetDbNode, _, err := s.repository.FieldSetRepository.UpdateForContactInTx(tx, common.GetContext(ctx).Tenant, contactId, fieldSet)
					if err != nil {
						return nil, err
					}
					fieldSetId = utils.GetPropsFromNode(*fieldSetDbNode)["id"].(string)
				}
				if fieldSet.CustomFields != nil {
					for _, customField := range *fieldSet.CustomFields {
						if customField.Id == nil {
							fieldDbNode, err := s.repository.CustomFieldRepository.MergeCustomFieldToFieldSetInTx(tx, tenant, contactId, fieldSetId, &customField)
							if err != nil {
								return nil, err
							}
							if customField.DefinitionId != nil {
								var fieldId = utils.GetPropsFromNode(*fieldDbNode)["id"].(string)
								err := s.repository.CustomFieldRepository.LinkWithCustomFieldDefinitionForFieldSetInTx(tx, fieldId, fieldSetId, *customField.DefinitionId)
								if err != nil {
									return nil, err
								}
							}
						} else {
							_, err := s.repository.CustomFieldRepository.UpdateForFieldSetInTx(tx, tenant, contactId, fieldSetId, &customField)
							if err != nil {
								return nil, err
							}
						}
					}
				}
			}
		}
		return nil, nil
	})

	return err
}

func (s *customFieldService) FindAllForContact(ctx context.Context, contact *model.Contact) (*entity.CustomFieldEntities, error) {
	session := utils.NewNeo4jReadSession(s.getDriver())
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
	session := utils.NewNeo4jReadSession(s.getDriver())
	defer session.Close()

	dbRecords, err := s.repository.CustomFieldRepository.FindAllForFieldSet(session, common.GetContext(ctx).Tenant, fieldSet.ID)
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
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	customFieldNode, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		customFieldDbNode, err := s.repository.CustomFieldRepository.MergeCustomFieldToContactInTx(tx, common.GetContext(ctx).Tenant, contactId, entity)
		if err != nil {
			return nil, err
		}
		if entity.DefinitionId != nil {
			var fieldId = utils.GetPropsFromNode(*customFieldDbNode)["id"].(string)
			if err = s.repository.CustomFieldRepository.LinkWithCustomFieldDefinitionForContactInTx(tx, fieldId, contactId, *entity.DefinitionId); err != nil {
				return nil, err
			}
		}
		return customFieldDbNode, err
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodePtrToCustomFieldEntity(customFieldNode.(*dbtype.Node)), nil
}

func (s *customFieldService) MergeCustomFieldToFieldSet(ctx context.Context, contactId string, fieldSetId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	customFieldNode, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		customFieldNode, err := s.repository.CustomFieldRepository.MergeCustomFieldToFieldSetInTx(tx, common.GetContext(ctx).Tenant, contactId, fieldSetId, entity)
		if err != nil {
			return nil, err
		}
		if entity.DefinitionId != nil {
			var fieldId = utils.GetPropsFromNode(*customFieldNode)["id"].(string)
			if err = s.repository.CustomFieldRepository.LinkWithCustomFieldDefinitionForFieldSetInTx(tx, fieldId, fieldSetId, *entity.DefinitionId); err != nil {
				return nil, err
			}
		}
		return customFieldNode, err
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodePtrToCustomFieldEntity(customFieldNode.(*dbtype.Node)), nil
}

func (s *customFieldService) UpdateCustomFieldForContact(ctx context.Context, contactId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	customFieldDbNode, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return s.repository.CustomFieldRepository.UpdateForContactInTx(tx, common.GetContext(ctx).Tenant, contactId, entity)
	})

	if err != nil {
		return nil, err
	}
	return s.mapDbNodePtrToCustomFieldEntity(customFieldDbNode.(*dbtype.Node)), nil
}

func (s *customFieldService) UpdateCustomFieldForFieldSet(ctx context.Context, contactId string, fieldSetId string, entity *entity.CustomFieldEntity) (*entity.CustomFieldEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	customFieldDbNode, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return s.repository.CustomFieldRepository.UpdateForFieldSetInTx(tx, common.GetContext(ctx).Tenant, contactId, fieldSetId, entity)
	})
	if err != nil {
		return nil, err
	}
	return s.mapDbNodePtrToCustomFieldEntity(customFieldDbNode.(*dbtype.Node)), nil
}

func (s *customFieldService) DeleteByNameFromContact(ctx context.Context, contactId, fieldName string) (bool, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()
	err := s.repository.CustomFieldRepository.DeleteByNameFromContact(session, common.GetContext(ctx).Tenant, contactId, fieldName)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *customFieldService) DeleteByIdFromContact(ctx context.Context, contactId, fieldId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()
	err := s.repository.CustomFieldRepository.DeleteByIdFromContact(session, common.GetContext(ctx).Tenant, contactId, fieldId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *customFieldService) DeleteByIdFromFieldSet(ctx context.Context, contactId, fieldSetId, fieldId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()
	err := s.repository.CustomFieldRepository.DeleteByIdFromFieldSet(session, common.GetContext(ctx).Tenant, contactId, fieldSetId, fieldId)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *customFieldService) mapDbNodeToCustomFieldEntity(node dbtype.Node) *entity.CustomFieldEntity {
	return s.mapDbNodePtrToCustomFieldEntity(&node)
}

func (s *customFieldService) mapDbNodePtrToCustomFieldEntity(node *dbtype.Node) *entity.CustomFieldEntity {
	props := utils.GetPropsFromNode(*node)
	result := entity.CustomFieldEntity{
		Id:       utils.StringPtr(utils.GetStringPropOrEmpty(props, "id")),
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
