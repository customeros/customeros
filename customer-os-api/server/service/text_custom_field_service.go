package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/repository"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type TextCustomFieldService interface {
	FindAllForContact(ctx context.Context, obj *model.Contact) (*entity.TextCustomFieldEntities, error)
	FindAllForFieldSet(ctx context.Context, obj *model.FieldSet) (*entity.TextCustomFieldEntities, error)

	MergeTextCustomFieldToContact(ctx context.Context, contactId string, entity *entity.TextCustomFieldEntity) (*entity.TextCustomFieldEntity, error)
	MergeTextCustomFieldToFieldSet(ctx context.Context, contactId string, fieldSetId string, entity *entity.TextCustomFieldEntity) (*entity.TextCustomFieldEntity, error)

	UpdateTextCustomFieldInContact(ctx context.Context, contactId string, entity *entity.TextCustomFieldEntity) (*entity.TextCustomFieldEntity, error)
	UpdateTextCustomFieldInFieldSet(ctx context.Context, contactId string, fieldSetId string, entity *entity.TextCustomFieldEntity) (*entity.TextCustomFieldEntity, error)

	DeleteByNameFromContact(ctx context.Context, contactId string, fieldName string) (bool, error)
	DeleteByIdFromContact(ctx context.Context, contactId string, fieldId string) (bool, error)
	DeleteByIdFromFieldSet(ctx context.Context, contactId string, fieldSetId string, fieldId string) (bool, error)

	mapDbNodeToTextCustomFieldEntity(node dbtype.Node) *entity.TextCustomFieldEntity
}

type textCustomPropertyService struct {
	repository *repository.RepositoryContainer
}

func NewTextCustomFieldService(repository *repository.RepositoryContainer) TextCustomFieldService {
	return &textCustomPropertyService{
		repository: repository,
	}
}

func (s *textCustomPropertyService) getDriver() neo4j.Driver {
	return *s.repository.Drivers.Neo4jDriver
}

func (s *textCustomPropertyService) FindAllForContact(ctx context.Context, contact *model.Contact) (*entity.TextCustomFieldEntities, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		result, err := tx.Run(`
				MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              		  (c)-[:HAS_TEXT_PROPERTY]->(f:TextCustomField) 
				RETURN f `,
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

	textCustomFieldEntities := entity.TextCustomFieldEntities{}

	for _, dbRecord := range queryResult.([]*db.Record) {
		textCustomFieldEntity := s.mapDbNodeToTextCustomFieldEntity(dbRecord.Values[0].(dbtype.Node))
		textCustomFieldEntities = append(textCustomFieldEntities, *textCustomFieldEntity)
	}

	return &textCustomFieldEntities, nil
}

func (s *textCustomPropertyService) FindAllForFieldSet(ctx context.Context, fieldSet *model.FieldSet) (*entity.TextCustomFieldEntities, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		result, err := tx.Run(`
				MATCH (s:FieldSet {id:$fieldSetId})<-[:HAS_COMPLEX_PROPERTY]-(:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              		  (s)-[:HAS_TEXT_PROPERTY]->(f:TextCustomField) 
				RETURN f`,
			map[string]any{
				"fieldSetId": fieldSet.ID,
				"tenant":     common.GetContext(ctx).Tenant})
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		return records, nil
	})
	if err != nil {
		return nil, err
	}

	textCustomFieldEntities := entity.TextCustomFieldEntities{}

	for _, dbRecord := range queryResult.([]*db.Record) {
		textCustomFieldEntity := s.mapDbNodeToTextCustomFieldEntity(dbRecord.Values[0].(dbtype.Node))
		textCustomFieldEntities = append(textCustomFieldEntities, *textCustomFieldEntity)
	}

	return &textCustomFieldEntities, nil
}

func (s *textCustomPropertyService) MergeTextCustomFieldToContact(ctx context.Context, contactId string, entity *entity.TextCustomFieldEntity) (*entity.TextCustomFieldEntity, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		txResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			MERGE (f:TextCustomField {name: $name})<-[:HAS_TEXT_PROPERTY]-(c)
            ON CREATE SET f.value=$value, f.id=randomUUID()
            ON MATCH SET f.value=$value
			RETURN f`,
			map[string]any{
				"tenant":    common.GetContext(ctx).Tenant,
				"contactId": contactId,
				"name":      entity.Name,
				"value":     entity.Value,
			})
		record, err := txResult.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToTextCustomFieldEntity(queryResult.(dbtype.Node)), nil
}

func (s *textCustomPropertyService) MergeTextCustomFieldToFieldSet(ctx context.Context, contactId string, fieldSetId string, entity *entity.TextCustomFieldEntity) (*entity.TextCustomFieldEntity, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		txResult, err := tx.Run(`
			MATCH (s:FieldSet {id:$fieldSetId})<-[:HAS_COMPLEX_PROPERTY]-(c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			MERGE (f:TextCustomField {name: $name})<-[:HAS_TEXT_PROPERTY]-(s)
            ON CREATE SET f.value=$value, f.id=randomUUID()
            ON MATCH SET f.value=$value
			RETURN f`,
			map[string]any{
				"tenant":     common.GetContext(ctx).Tenant,
				"contactId":  contactId,
				"fieldSetId": fieldSetId,
				"name":       entity.Name,
				"value":      entity.Value,
			})
		record, err := txResult.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToTextCustomFieldEntity(queryResult.(dbtype.Node)), nil
}

func (s *textCustomPropertyService) UpdateTextCustomFieldInContact(ctx context.Context, contactId string, entity *entity.TextCustomFieldEntity) (*entity.TextCustomFieldEntity, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		txResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
			  (c)-[:HAS_TEXT_PROPERTY]->(f:TextCustomField {id:$fieldId})
			SET	f.name=$name,
				f.value=$value
			RETURN f`,
			map[string]any{
				"tenant":    common.GetContext(ctx).Tenant,
				"contactId": contactId,
				"fieldId":   entity.Id,
				"name":      entity.Name,
				"value":     entity.Value,
			})
		record, err := txResult.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToTextCustomFieldEntity(queryResult.(dbtype.Node)), nil
}

func (s *textCustomPropertyService) UpdateTextCustomFieldInFieldSet(ctx context.Context, contactId string, fieldSetId string, entity *entity.TextCustomFieldEntity) (*entity.TextCustomFieldEntity, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		txResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              (c)-[:HAS_COMPLEX_PROPERTY]->(s:FieldSet {id:$fieldSetId}),
			  (s)-[:HAS_TEXT_PROPERTY]->(f:TextCustomField {id:$fieldId})
			SET	f.name=$name,
				f.value=$value
			RETURN f`,
			map[string]any{
				"tenant":     common.GetContext(ctx).Tenant,
				"contactId":  contactId,
				"fieldSetId": fieldSetId,
				"fieldId":    entity.Id,
				"name":       entity.Name,
				"value":      entity.Value,
			})
		record, err := txResult.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
	if err != nil {
		return nil, err
	}

	return s.mapDbNodeToTextCustomFieldEntity(queryResult.(dbtype.Node)), nil
}

func (s *textCustomPropertyService) DeleteByNameFromContact(ctx context.Context, contactId string, fieldName string) (bool, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$id})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c)-[:HAS_TEXT_PROPERTY]->(f:TextCustomField {name:$name})
            DETACH DELETE f
			`,
			map[string]any{
				"id":     contactId,
				"name":   fieldName,
				"tenant": common.GetContext(ctx).Tenant,
			})

		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
}

func (s *textCustomPropertyService) DeleteByIdFromContact(ctx context.Context, contactId string, fieldId string) (bool, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c)-[:HAS_TEXT_PROPERTY]->(f:TextCustomField {id:$fieldId})
            DETACH DELETE f`,
			map[string]any{
				"contactId": contactId,
				"fieldId":   fieldId,
				"tenant":    common.GetContext(ctx).Tenant,
			})

		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
}

func (s *textCustomPropertyService) DeleteByIdFromFieldSet(ctx context.Context, contactId string, fieldSetId string, fieldId string) (bool, error) {
	session := s.getDriver().NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c)-[:HAS_COMPLEX_PROPERTY]->(s:FieldSet {id:$fieldSetId}),
                  (s)-[:HAS_TEXT_PROPERTY]->(f:TextCustomField {id:$fieldId})
            DETACH DELETE f`,
			map[string]any{
				"contactId":  contactId,
				"fieldSetId": fieldSetId,
				"fieldId":    fieldId,
				"tenant":     common.GetContext(ctx).Tenant,
			})

		return true, err
	})
	if err != nil {
		return false, err
	}

	return queryResult.(bool), nil
}

func addTextCustomFieldToContactInTx(ctx context.Context, contactId string, input entity.TextCustomFieldEntity, tx neo4j.Transaction) error {
	_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})
			CREATE (f:TextCustomField {
				id: randomUUID(),
				name: $name,
				value: $value
			})<-[:HAS_TEXT_PROPERTY]-(c)
			RETURN f`,
		map[string]any{
			"contactId": contactId,
			"name":      input.Name,
			"value":     input.Value,
		})

	return err
}

func (s *textCustomPropertyService) mapDbNodeToTextCustomFieldEntity(node dbtype.Node) *entity.TextCustomFieldEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.TextCustomFieldEntity{
		Id:    utils.GetStringPropOrEmpty(props, "id"),
		Name:  utils.GetStringPropOrEmpty(props, "name"),
		Value: utils.GetStringPropOrEmpty(props, "value"),
	}
	return &result
}
