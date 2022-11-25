package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
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
	getDriver() neo4j.Driver
}

type fieldSetService struct {
	repository *repository.RepositoryContainer
}

func NewFieldSetService(repository *repository.RepositoryContainer) FieldSetService {
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

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		result, err := tx.Run(`
				MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              			(c)-[r:HAS_COMPLEX_PROPERTY]->(s:FieldSet) 
				RETURN s, r`,
			map[string]any{
				"contactId": contact.ID,
				"tenant":    common.GetContext(ctx).Tenant})
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		return records, nil
	})
	if err != nil {
		return nil, err
	}

	fieldSetEntities := entity.FieldSetEntities{}

	for _, dbRecord := range queryResult.([]*db.Record) {
		fieldSetEntity := s.mapDbNodeToFieldSetEntity(utils.NodePtr(dbRecord.Values[0].(dbtype.Node)))
		s.addDbRelationshipToEntity(utils.RelationshipPtr(dbRecord.Values[1].(dbtype.Relationship)), fieldSetEntity)
		fieldSetEntities = append(fieldSetEntities, *fieldSetEntity)
	}

	return &fieldSetEntities, nil
}

func (s *fieldSetService) MergeFieldSetToContact(ctx context.Context, contactId string, entity *entity.FieldSetEntity) (*entity.FieldSetEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	var fieldSetDbNode *dbtype.Node
	var fieldSetDbRelationship *neo4j.Relationship

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		var err error
		fieldSetDbNode, fieldSetDbRelationship, err = s.repository.FieldSetRepository.MergeFieldSetToContactInTx(tx, common.GetContext(ctx).Tenant, contactId, *entity)
		if err != nil {
			return nil, err
		}
		var fieldSetId = utils.GetPropsFromNode(*fieldSetDbNode)["id"].(string)
		if entity.DefinitionId != nil {
			err := s.repository.FieldSetRepository.LinkWithFieldSetDefinitionInTx(tx, common.GetContext(ctx).Tenant, fieldSetId, *entity.DefinitionId)
			if err != nil {
				return nil, err
			}
		}
		if entity.CustomFields != nil {
			for _, customField := range *entity.CustomFields {
				dbNode, err := s.repository.CustomFieldRepository.MergeCustomFieldToFieldSetInTx(tx, common.GetContext(ctx).Tenant, contactId, fieldSetId, &customField)
				if err != nil {
					return nil, err
				}
				if customField.DefinitionId != nil {
					var fieldId = utils.GetPropsFromNode(*dbNode)["id"].(string)
					err := s.repository.CustomFieldRepository.LinkWithCustomFieldDefinitionForFieldSetInTx(tx, fieldId, fieldSetId, *customField.DefinitionId)
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

	var fieldSetEntity = s.mapDbNodeToFieldSetEntity(fieldSetDbNode)
	s.addDbRelationshipToEntity(fieldSetDbRelationship, fieldSetEntity)
	return fieldSetEntity, nil
}

func (s *fieldSetService) UpdateFieldSetInContact(ctx context.Context, contactId string, entity *entity.FieldSetEntity) (*entity.FieldSetEntity, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	var fieldSetDbNode *dbtype.Node
	var fieldSetDbRelationship *neo4j.Relationship

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		var err error
		fieldSetDbNode, fieldSetDbRelationship, err = s.repository.FieldSetRepository.UpdateForContactInTx(tx, common.GetContext(ctx).Tenant, contactId, *entity)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	var fieldSetEntity = s.mapDbNodeToFieldSetEntity(fieldSetDbNode)
	s.addDbRelationshipToEntity(fieldSetDbRelationship, fieldSetEntity)
	return fieldSetEntity, nil
}

func (s *fieldSetService) DeleteByIdFromContact(ctx context.Context, contactId string, fieldSetId string) (bool, error) {
	session := utils.NewNeo4jWriteSession(s.getDriver())
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c)-[:HAS_COMPLEX_PROPERTY]->(s:FieldSet {id:$fieldSetId}),
				  (s)-[:HAS_PROPERTY]->(f:CustomField)
            DETACH DELETE f, s`,
			map[string]any{
				"contactId":  contactId,
				"fieldSetId": fieldSetId,
				"tenant":     common.GetContext(ctx).Tenant,
			})

		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
}

func (s *fieldSetService) mapDbNodeToFieldSetEntity(node *dbtype.Node) *entity.FieldSetEntity {
	props := utils.GetPropsFromNode(*node)
	result := entity.FieldSetEntity{
		Id:   utils.StringPtr(utils.GetStringPropOrEmpty(props, "id")),
		Name: utils.GetStringPropOrEmpty(props, "name"),
	}
	return &result
}

func (s *fieldSetService) addDbRelationshipToEntity(relationship *dbtype.Relationship, fieldSetEntity *entity.FieldSetEntity) {
	props := utils.GetPropsFromRelationship(*relationship)
	fieldSetEntity.Added = utils.GetTimePropOrNow(props, "added")
}
