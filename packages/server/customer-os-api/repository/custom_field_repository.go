package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"golang.org/x/net/context"
)

type CustomFieldRepository interface {
	MergeCustomFieldToContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, entity entity.CustomFieldEntity) (*dbtype.Node, error)
	MergeCustomFieldToFieldSetInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, fieldSet string, entity entity.CustomFieldEntity) (*dbtype.Node, error)
	LinkWithCustomFieldTemplateForContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, fieldId, contactId, templateId string) error
	LinkWithCustomFieldTemplateForFieldSetInTx(ctx context.Context, tx neo4j.ManagedTransaction, fieldId, fieldSetId, templateId string) error
	UpdateForContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, entity entity.CustomFieldEntity) (*dbtype.Node, error)
	UpdateForFieldSetInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, fieldSetId string, entity entity.CustomFieldEntity) (*dbtype.Node, error)
	FindAllForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string) ([]*neo4j.Record, error)
	FindAllForFieldSet(ctx context.Context, session neo4j.SessionWithContext, tenant, fieldSetId string) ([]*neo4j.Record, error)
	DeleteByNameFromContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId, fieldName string) error
	DeleteByIdFromContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId, fieldId string) error
	DeleteByIdFromFieldSet(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId, fieldSetId, fieldId string) error
}

type customFieldRepository struct {
	driver *neo4j.DriverWithContext
}

func NewCustomFieldRepository(driver *neo4j.DriverWithContext) CustomFieldRepository {
	return &customFieldRepository{
		driver: driver,
	}
}

func (r *customFieldRepository) MergeCustomFieldToContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, entity entity.CustomFieldEntity) (*dbtype.Node, error) {
	queryResult, err := tx.Run(ctx,
		fmt.Sprintf("MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) "+
			" MERGE (f:%s:CustomField {name: $name, datatype:$datatype})<-[:HAS_PROPERTY]-(c) "+
			" ON CREATE SET f.%s=$value, f.id=randomUUID(), f.createdAt=$now, f.updatedAt=$now, f.source=$source, f.sourceOfTruth=$sourceOfTruth,  f:%s "+
			" ON MATCH SET f.%s=$value, f.sourceOfTruth=$sourceOfTruth, f.updatedAt=$now "+
			" RETURN f", entity.NodeLabel(), entity.PropertyName(), "CustomField_"+tenant, entity.PropertyName()),
		map[string]any{
			"tenant":        tenant,
			"contactId":     contactId,
			"name":          entity.Name,
			"datatype":      entity.DataType,
			"value":         entity.Value.RealValue(),
			"source":        entity.Source,
			"sourceOfTruth": entity.SourceOfTruth,
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *customFieldRepository) MergeCustomFieldToFieldSetInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, fieldSetId string, entity entity.CustomFieldEntity) (*dbtype.Node, error) {
	queryResult, err := tx.Run(ctx,
		fmt.Sprintf(" MATCH (s:FieldSet {id:$fieldSetId})<-[:HAS_COMPLEX_PROPERTY]-(c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) "+
			" MERGE (f:%s:CustomField {name: $name, datatype:$datatype})<-[:HAS_PROPERTY]-(s)"+
			" ON CREATE SET f.%s=$value, f.id=randomUUID(), f.createdAt=$now, f.updatedAt=$now, f.source=$source, f.sourceOfTruth=$sourceOfTruth, f:%s "+
			" ON MATCH SET f.%s=$value, f.sourceOfTruth=$sourceOfTruth, f.updatedAt=$now "+
			" RETURN f", entity.NodeLabel(), entity.PropertyName(), "CustomField_"+tenant, entity.PropertyName()),
		map[string]any{
			"tenant":        tenant,
			"contactId":     contactId,
			"fieldSetId":    fieldSetId,
			"name":          entity.Name,
			"datatype":      entity.DataType,
			"value":         entity.Value.RealValue(),
			"source":        entity.Source,
			"sourceOfTruth": entity.SourceOfTruth,
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *customFieldRepository) LinkWithCustomFieldTemplateForContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, fieldId, contactId, templateId string) error {
	queryResult, err := tx.Run(ctx, `
			MATCH (f:CustomField {id:$fieldId})<-[:HAS_PROPERTY]-(c:Contact {id:$contactId}),
				  (c)-[:IS_DEFINED_BY]->(e:EntityTemplate),
				  (e)-[:CONTAINS]->(d:CustomFieldTemplate {id:$customFieldTemplateId})
			MERGE (f)-[r:IS_DEFINED_BY]->(d)
			RETURN r`,
		map[string]any{
			"customFieldTemplateId": templateId,
			"fieldId":               fieldId,
			"contactId":             contactId,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *customFieldRepository) LinkWithCustomFieldTemplateForFieldSetInTx(ctx context.Context, tx neo4j.ManagedTransaction, fieldId, fieldSetId, templateId string) error {
	queryResult, err := tx.Run(ctx, `
			MATCH (f:CustomField {id:$fieldId})<-[:HAS_PROPERTY]-(s:FieldSet {id:$fieldSetId}),
				  (s)-[:IS_DEFINED_BY]->(e:FieldSetTemplate),
				  (e)-[:CONTAINS]->(d:CustomFieldTemplate {id:$customFieldTemplateId})
			MERGE (f)-[r:IS_DEFINED_BY]->(d)
			RETURN r`,
		map[string]any{
			"customFieldTemplateId": templateId,
			"fieldId":               fieldId,
			"fieldSetId":            fieldSetId,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *customFieldRepository) FindAllForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string) ([]*neo4j.Record, error) {
	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
				MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              		  (c)-[:HAS_PROPERTY]->(f:CustomField) 
				RETURN f ORDER BY f.name`,
			map[string]any{
				"contactId": contactId,
				"tenant":    tenant})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	return records.([]*neo4j.Record), err
}

func (r *customFieldRepository) FindAllForFieldSet(ctx context.Context, session neo4j.SessionWithContext, tenant, fieldSetId string) ([]*neo4j.Record, error) {
	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
				MATCH (s:FieldSet {id:$fieldSetId})<-[:HAS_COMPLEX_PROPERTY]-(:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
              		  (s)-[:HAS_PROPERTY]->(f:CustomField) 
				RETURN f ORDER BY f.name`,
			map[string]any{
				"fieldSetId": fieldSetId,
				"tenant":     tenant})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	return records.([]*neo4j.Record), err
}

func (r *customFieldRepository) DeleteByIdFromContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId, fieldId string) error {
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
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

func (r *customFieldRepository) DeleteByNameFromContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId, fieldId string) error {
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
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

func (r *customFieldRepository) DeleteByIdFromFieldSet(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId, fieldSetId, fieldId string) error {
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
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

func (r *customFieldRepository) UpdateForContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, entity entity.CustomFieldEntity) (*dbtype.Node, error) {
	queryResult, err := tx.Run(ctx, fmt.Sprintf(
		"MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), "+
			" (c)-[:HAS_PROPERTY]->(f:%s:CustomField {id:$fieldId}) "+
			" SET f.name=$name, "+
			" f.%s=$value, "+
			" f.sourceOfTruth=$sourceOfTruth, "+
			" f.updatedAt=$now "+
			" RETURN f", entity.NodeLabel(), entity.PropertyName()),
		map[string]any{
			"tenant":        tenant,
			"contactId":     contactId,
			"fieldId":       entity.Id,
			"name":          entity.Name,
			"sourceOfTruth": entity.SourceOfTruth,
			"value":         entity.Value.RealValue(),
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *customFieldRepository) UpdateForFieldSetInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, fieldSetId string, entity entity.CustomFieldEntity) (*dbtype.Node, error) {
	queryResult, err := tx.Run(ctx, fmt.Sprintf(
		"MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), "+
			" (c)-[:HAS_COMPLEX_PROPERTY]->(s:FieldSet {id:$fieldSetId}),"+
			" (s)-[:HAS_PROPERTY]->(f:%s:CustomField {id:$fieldId})"+
			" SET f.name=$name, "+
			" f.%s=$value, "+
			" f.sourceOfTruth=$sourceOfTruth, "+
			" f.updatedAt=$now "+
			"RETURN f", entity.NodeLabel(), entity.PropertyName()),
		map[string]any{
			"tenant":        tenant,
			"contactId":     contactId,
			"fieldSetId":    fieldSetId,
			"fieldId":       entity.Id,
			"name":          entity.Name,
			"sourceOfTruth": entity.SourceOfTruth,
			"value":         entity.Value.RealValue(),
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}
