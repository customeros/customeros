package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type CustomFieldRepository interface {
	MergeCustomFieldToContactInTx(tx neo4j.Transaction, tenant, contactId string, entity *entity.CustomFieldEntity) (*dbtype.Node, error)
	MergeCustomFieldToFieldSetInTx(tx neo4j.Transaction, tenant, contactId, fieldSet string, entity *entity.CustomFieldEntity) (*dbtype.Node, error)
	LinkWithCustomFieldDefinitionForContactInTx(tx neo4j.Transaction, fieldId, contactId, definitionId string) error
	LinkWithCustomFieldDefinitionForFieldSetInTx(tx neo4j.Transaction, fieldId, fieldSetId, definitionId string) error
	UpdateForContactInTx(tx neo4j.Transaction, tenant, contactId string, entity *entity.CustomFieldEntity) (*dbtype.Node, error)
	UpdateForFieldSetInTx(tx neo4j.Transaction, tenant, contactId, fieldSetId string, entity *entity.CustomFieldEntity) (*dbtype.Node, error)
	FindAllForContact(session neo4j.Session, tenant, contactId string) ([]*neo4j.Record, error)
	FindAllForFieldSet(session neo4j.Session, tenant, fieldSetId string) ([]*neo4j.Record, error)
	DeleteByNameFromContact(session neo4j.Session, tenant, contactId, fieldName string) error
	DeleteByIdFromContact(session neo4j.Session, tenant, contactId, fieldId string) error
	DeleteByIdFromFieldSet(session neo4j.Session, tenant, contactId, fieldSetId, fieldId string) error
}

type customFieldRepository struct {
	driver *neo4j.Driver
}

func NewCustomFieldRepository(driver *neo4j.Driver) CustomFieldRepository {
	return &customFieldRepository{
		driver: driver,
	}
}

func (r *customFieldRepository) MergeCustomFieldToContactInTx(tx neo4j.Transaction, tenant, contactId string, entity *entity.CustomFieldEntity) (*dbtype.Node, error) {
	queryResult, err := tx.Run(
		fmt.Sprintf("MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) "+
			" MERGE (f:%s:CustomField {name: $name, datatype:$datatype})<-[:HAS_PROPERTY]-(c) "+
			" ON CREATE SET f.%s=$value, f.id=randomUUID(), f.source=$source, f:%s "+
			" ON MATCH SET f.%s=$value, f.source=$source "+
			" RETURN f", entity.NodeLabel(), entity.PropertyName(), "CustomField_"+tenant, entity.PropertyName()),
		map[string]any{
			"tenant":    tenant,
			"contactId": contactId,
			"name":      entity.Name,
			"datatype":  entity.DataType,
			"source":    entity.Source,
			"value":     entity.Value.RealValue(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
}

func (r *customFieldRepository) MergeCustomFieldToFieldSetInTx(tx neo4j.Transaction, tenant, contactId, fieldSetId string, entity *entity.CustomFieldEntity) (*dbtype.Node, error) {
	queryResult, err := tx.Run(
		fmt.Sprintf(" MATCH (s:FieldSet {id:$fieldSetId})<-[:HAS_COMPLEX_PROPERTY]-(c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) "+
			" MERGE (f:%s:CustomField {name: $name, datatype:$datatype})<-[:HAS_PROPERTY]-(s)"+
			" ON CREATE SET f.%s=$value, f.id=randomUUID(), f.source=$source, f:%s "+
			" ON MATCH SET f.%s=$value, f.source=$source "+
			" RETURN f", entity.NodeLabel(), entity.PropertyName(), "CustomField_"+tenant, entity.PropertyName()),
		map[string]any{
			"tenant":     tenant,
			"contactId":  contactId,
			"fieldSetId": fieldSetId,
			"name":       entity.Name,
			"datatype":   entity.DataType,
			"source":     entity.Source,
			"value":      entity.Value.RealValue(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
}

func (r *customFieldRepository) LinkWithCustomFieldDefinitionForContactInTx(tx neo4j.Transaction, fieldId, contactId, definitionId string) error {
	queryResult, err := tx.Run(`
			MATCH (f:CustomField {id:$fieldId})<-[:HAS_PROPERTY]-(c:Contact {id:$contactId}),
				  (c)-[:IS_DEFINED_BY]->(e:EntityDefinition),
				  (e)-[:CONTAINS]->(d:CustomFieldDefinition {id:$customFieldDefinitionId})
			MERGE (f)-[r:IS_DEFINED_BY]->(d)
			RETURN r`,
		map[string]any{
			"customFieldDefinitionId": definitionId,
			"fieldId":                 fieldId,
			"contactId":               contactId,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single()
	return err
}

func (r *customFieldRepository) LinkWithCustomFieldDefinitionForFieldSetInTx(tx neo4j.Transaction, fieldId, fieldSetId, definitionId string) error {
	queryResult, err := tx.Run(`
			MATCH (f:CustomField {id:$fieldId})<-[:HAS_PROPERTY]-(s:FieldSet {id:$fieldSetId}),
				  (s)-[:IS_DEFINED_BY]->(e:FieldSetDefinition),
				  (e)-[:CONTAINS]->(d:CustomFieldDefinition {id:$customFieldDefinitionId})
			MERGE (f)-[r:IS_DEFINED_BY]->(d)
			RETURN r`,
		map[string]any{
			"customFieldDefinitionId": definitionId,
			"fieldId":                 fieldId,
			"fieldSetId":              fieldSetId,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single()
	return err
}

func (r *customFieldRepository) FindAllForContact(session neo4j.Session, tenant, contactId string) ([]*neo4j.Record, error) {
	records, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
				MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              		  (c)-[:HAS_PROPERTY]->(f:CustomField) 
				RETURN f ORDER BY f.name`,
			map[string]any{
				"contactId": contactId,
				"tenant":    tenant})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
	return records.([]*neo4j.Record), err
}

func (r *customFieldRepository) FindAllForFieldSet(session neo4j.Session, tenant, fieldSetId string) ([]*neo4j.Record, error) {
	records, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
				MATCH (s:FieldSet {id:$fieldSetId})<-[:HAS_COMPLEX_PROPERTY]-(:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              		  (s)-[:HAS_PROPERTY]->(f:CustomField) 
				RETURN f ORDER BY f.name`,
			map[string]any{
				"fieldSetId": fieldSetId,
				"tenant":     tenant})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect()
	})
	return records.([]*neo4j.Record), err
}

func (r *customFieldRepository) DeleteByIdFromContact(session neo4j.Session, tenant, contactId, fieldId string) error {
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c)-[:HAS_PROPERTY]->(f:CustomField {id:$fieldId})
            DETACH DELETE f`,
			map[string]any{
				"contactId": contactId,
				"fieldId":   fieldId,
				"tenant":    tenant,
			})
		return nil, err
	})
	return err
}

func (r *customFieldRepository) DeleteByNameFromContact(session neo4j.Session, tenant, contactId, fieldId string) error {
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$id})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c)-[:HAS_PROPERTY]->(f:CustomField {name:$name})
            DETACH DELETE f`,
			map[string]any{
				"contactId": contactId,
				"fieldId":   fieldId,
				"tenant":    tenant,
			})
		return nil, err
	})
	return err
}

func (r *customFieldRepository) DeleteByIdFromFieldSet(session neo4j.Session, tenant, contactId, fieldSetId, fieldId string) error {
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c)-[:HAS_COMPLEX_PROPERTY]->(s:FieldSet {id:$fieldSetId}),
                  (s)-[:HAS_PROPERTY]->(f:CustomField {id:$fieldId})
            DETACH DELETE f`,
			map[string]any{
				"contactId":  contactId,
				"fieldSetId": fieldSetId,
				"fieldId":    fieldId,
				"tenant":     tenant,
			})
		return nil, err
	})
	return err
}

func (r *customFieldRepository) UpdateForContactInTx(tx neo4j.Transaction, tenant, contactId string, entity *entity.CustomFieldEntity) (*dbtype.Node, error) {
	queryResult, err := tx.Run(fmt.Sprintf(
		"MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), "+
			" (c)-[:HAS_PROPERTY]->(f:%s:CustomField {id:$fieldId}) "+
			" SET f.name=$name, "+
			" f.%s=$value, "+
			" f.source=$source "+
			" RETURN f", entity.NodeLabel(), entity.PropertyName()),
		map[string]any{
			"tenant":    tenant,
			"contactId": contactId,
			"fieldId":   entity.Id,
			"name":      entity.Name,
			"source":    entity.Source,
			"value":     entity.Value.RealValue(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
}

func (r *customFieldRepository) UpdateForFieldSetInTx(tx neo4j.Transaction, tenant, contactId, fieldSetId string, entity *entity.CustomFieldEntity) (*dbtype.Node, error) {
	queryResult, err := tx.Run(fmt.Sprintf(
		"MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), "+
			" (c)-[:HAS_COMPLEX_PROPERTY]->(s:FieldSet {id:$fieldSetId}),"+
			" (s)-[:HAS_PROPERTY]->(f:%s:CustomField {id:$fieldId})"+
			" SET f.name=$name, "+
			" f.%s=$value, "+
			" f.source=$source "+
			"RETURN f", entity.NodeLabel(), entity.PropertyName()),
		map[string]any{
			"tenant":     tenant,
			"contactId":  contactId,
			"fieldSetId": fieldSetId,
			"fieldId":    entity.Id,
			"name":       entity.Name,
			"source":     entity.Source,
			"value":      entity.Value.RealValue(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
}
