package service

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/customer-os-api/common"
	"github.com/openline-ai/openline-customer-os/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type TextCustomFieldService interface {
	FindAllForContact(ctx context.Context, obj *model.Contact) (*entity.TextCustomFieldEntities, error)
	FindAllForFieldsSet(ctx context.Context, obj *model.FieldsSet) (*entity.TextCustomFieldEntities, error)

	MergeTextCustomFieldToContact(ctx context.Context, contactId string, entity *entity.TextCustomFieldEntity) (*entity.TextCustomFieldEntity, error)
	MergeTextCustomFieldToFieldsSet(ctx context.Context, contactId string, fieldsSetId string, entity *entity.TextCustomFieldEntity) (*entity.TextCustomFieldEntity, error)

	UpdateTextCustomFieldInContact(ctx context.Context, contactId string, entity *entity.TextCustomFieldEntity) (*entity.TextCustomFieldEntity, error)
	UpdateTextCustomFieldInFieldsSet(ctx context.Context, contactId string, fieldsSetId string, entity *entity.TextCustomFieldEntity) (*entity.TextCustomFieldEntity, error)

	Delete(ctx context.Context, contactId string, fieldName string) (bool, error)
	DeleteById(ctx context.Context, contactId string, fieldId string) (bool, error)
	DeleteByIdFromFieldsSet(ctx context.Context, contactId string, fieldsSetId string, fieldId string) (bool, error)

	mapDbNodeToTextCustomFieldEntity(dbContactGroupNode dbtype.Node) *entity.TextCustomFieldEntity
}

type textCustomPropertyService struct {
	driver *neo4j.Driver
}

func NewTextCustomFieldService(driver *neo4j.Driver) TextCustomFieldService {
	return &textCustomPropertyService{
		driver: driver,
	}
}

func (s *textCustomPropertyService) FindAllForContact(ctx context.Context, contact *model.Contact) (*entity.TextCustomFieldEntities, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
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

func (s *textCustomPropertyService) FindAllForFieldsSet(ctx context.Context, fieldsSet *model.FieldsSet) (*entity.TextCustomFieldEntities, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		result, err := tx.Run(`
				MATCH (s:FieldsSet {id:$fieldsSetId})<-[:HAS_COMPLEX_PROPERTY]-(:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              		  (s)-[:HAS_TEXT_PROPERTY]->(f:TextCustomField) 
				RETURN f`,
			map[string]any{
				"fieldsSetId": fieldsSet.ID,
				"tenant":      common.GetContext(ctx).Tenant})
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
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		txResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			MERGE (f:TextCustomField {name: $name})<-[:HAS_TEXT_PROPERTY]-(c)
            ON CREATE SET f.value=$value, f.group=$group, f.id=randomUUID()
            ON MATCH SET f.value=$value, f.group=$group
			RETURN f`,
			map[string]any{
				"tenant":    common.GetContext(ctx).Tenant,
				"contactId": contactId,
				"name":      entity.Name,
				"value":     entity.Value,
				"group":     entity.Group,
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

func (s *textCustomPropertyService) MergeTextCustomFieldToFieldsSet(ctx context.Context, contactId string, fieldsSetId string, entity *entity.TextCustomFieldEntity) (*entity.TextCustomFieldEntity, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		txResult, err := tx.Run(`
			MATCH (s:FieldsSet {id:$fieldsSetId})<-[:HAS_COMPLEX_PROPERTY]-(c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			MERGE (f:TextCustomField {name: $name})<-[:HAS_TEXT_PROPERTY]-(s)
            ON CREATE SET f.value=$value, f.group=$group, f.id=randomUUID()
            ON MATCH SET f.value=$value, f.group=$group
			RETURN f`,
			map[string]any{
				"tenant":      common.GetContext(ctx).Tenant,
				"contactId":   contactId,
				"fieldsSetId": fieldsSetId,
				"name":        entity.Name,
				"value":       entity.Value,
				"group":       entity.Group,
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
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		txResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
			  (c)-[:HAS_TEXT_PROPERTY]->(f:TextCustomField {id:$fieldId})
			SET	f.name=$name,
				f.value=$value,
				f.group=$group
			RETURN f`,
			map[string]any{
				"tenant":    common.GetContext(ctx).Tenant,
				"contactId": contactId,
				"fieldId":   entity.Id,
				"name":      entity.Name,
				"value":     entity.Value,
				"group":     entity.Group,
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

func (s *textCustomPropertyService) UpdateTextCustomFieldInFieldsSet(ctx context.Context, contactId string, fieldsSetId string, entity *entity.TextCustomFieldEntity) (*entity.TextCustomFieldEntity, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		txResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              (c)-[:HAS_COMPLEX_PROPERTY]->(s:FieldsSet {id:$fieldsSetId}),
			  (s)-[:HAS_TEXT_PROPERTY]->(f:TextCustomField {id:$fieldId})
			SET	f.name=$name,
				f.value=$value,
				f.group=$group
			RETURN f`,
			map[string]any{
				"tenant":      common.GetContext(ctx).Tenant,
				"contactId":   contactId,
				"fieldsSetId": fieldsSetId,
				"fieldId":     entity.Id,
				"name":        entity.Name,
				"value":       entity.Value,
				"group":       entity.Group,
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

func (s *textCustomPropertyService) Delete(ctx context.Context, contactId string, fieldName string) (bool, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$id})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c:Contact {id:$id})-[:HAS_TEXT_PROPERTY]->(f:TextCustomField {name:$name})
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

func (s *textCustomPropertyService) DeleteById(ctx context.Context, contactId string, fieldId string) (bool, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
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

func (s *textCustomPropertyService) DeleteByIdFromFieldsSet(ctx context.Context, contactId string, fieldsSetId string, fieldId string) (bool, error) {
	session := (*s.driver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c)-[:HAS_COMPLEX_PROPERTY]->(s:FieldsSet {id:$fieldsSetId}),
                  (s)-[:HAS_TEXT_PROPERTY]->(f:TextCustomField {id:$fieldId})
            DETACH DELETE f`,
			map[string]any{
				"contactId":   contactId,
				"fieldsSetId": fieldsSetId,
				"fieldId":     fieldId,
				"tenant":      common.GetContext(ctx).Tenant,
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
				group: $group,
				name: $name,
				value: $value
			})<-[:HAS_TEXT_PROPERTY]-(c)
			RETURN f`,
		map[string]any{
			"contactId": contactId,
			"group":     input.Group,
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
		Group: utils.GetStringPropOrEmpty(props, "group"),
	}
	return &result
}
