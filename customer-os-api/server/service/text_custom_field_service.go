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
	MergeTextCustomFieldToContact(ctx context.Context, id string, entity *entity.TextCustomFieldEntity) (*entity.TextCustomFieldEntity, error)

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

	queryResult, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`MATCH (:Tenant {name:$tenant})<--(:Contact {id:$id})-->(f:TextCustomField) 
				RETURN f `,
			map[string]interface{}{
				"id":     contact.ID,
				"tenant": common.GetContext(ctx).Tenant})
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

	queryResult, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		txResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})--(:Tenant {name:$tenant})
			MERGE (f:TextCustomField {name: $name})<-[:HAS_PROPERTY]-(c)
            ON CREATE SET f.value=$value, f.group=$group
            ON MATCH SET f.value=$value, f.group=$group
			RETURN f`,
			map[string]interface{}{
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

func txAddTextCustomFieldToContact(ctx context.Context, contactId string, input entity.TextCustomFieldEntity, tx neo4j.Transaction) (interface{}, error) {
	result, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})
			CREATE (f:TextCustomField {
				  group: $group,
				  name: $name,
				  value: $value
			})<-[:HAS_PROPERTY]-(c)
			RETURN f`,
		map[string]interface{}{
			"contactId": contactId,
			"group":     input.Group,
			"name":      input.Name,
			"value":     input.Value,
		})

	record, err := result.Single()
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (s *textCustomPropertyService) mapDbNodeToTextCustomFieldEntity(dbContactGroupNode dbtype.Node) *entity.TextCustomFieldEntity {
	props := utils.GetPropsFromNode(dbContactGroupNode)
	result := entity.TextCustomFieldEntity{
		Name:  props["name"].(string),
		Value: props["value"].(string),
		Group: props["group"].(string),
	}
	return &result
}
