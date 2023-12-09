package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
)

type CustomFieldRepository interface {
	MergeCustomFieldToContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, contactId string, entity entity.CustomFieldEntity) (*dbtype.Node, error)
	MergeCustomFieldInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entityType *model.CustomFieldEntityType, entity entity.CustomFieldEntity) (*dbtype.Node, error)
	MergeCustomFieldToFieldSetInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, obj *model.CustomFieldEntityType, fieldSet string, entity entity.CustomFieldEntity) (*dbtype.Node, error)
	LinkWithCustomFieldTemplateForContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, fieldId, contactId, templateId string) error
	LinkWithCustomFieldTemplateInTx(ctx context.Context, tx neo4j.ManagedTransaction, fieldId string, obj *model.CustomFieldEntityType, templateId string) error
	LinkWithCustomFieldTemplateForFieldSetInTx(ctx context.Context, tx neo4j.ManagedTransaction, fieldId, fieldSetId, templateId string) error
	UpdateForContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, entity entity.CustomFieldEntity) (*dbtype.Node, error)
	UpdateForFieldSetInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId, fieldSetId string, entity entity.CustomFieldEntity) (*dbtype.Node, error)
	FindAll(ctx context.Context, session neo4j.SessionWithContext, tenant string, obj *model.CustomFieldEntityType) ([]*neo4j.Record, error)
	FindAllForFieldSet(ctx context.Context, session neo4j.SessionWithContext, tenant, fieldSetId string) ([]*neo4j.Record, error)
	DeleteByNameFromContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId, fieldName string) error
	DeleteByIdFromContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId, fieldId string) error
	DeleteByIdFromFieldSet(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId, fieldSetId, fieldId string) error

	GetCustomFields(ctx context.Context, session neo4j.SessionWithContext, tenant string, obj *model.CustomFieldEntityType) ([]*neo4j.Record, error)
}

type customFieldRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewCustomFieldRepository(driver *neo4j.DriverWithContext, database string) CustomFieldRepository {
	return &customFieldRepository{
		driver:   driver,
		database: database,
	}
}

func (r *customFieldRepository) MergeCustomFieldToContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, entity entity.CustomFieldEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldRepository.MergeCustomFieldToContactInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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

func (r *customFieldRepository) MergeCustomFieldInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, obj *model.CustomFieldEntityType, entity entity.CustomFieldEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldRepository.MergeCustomFieldInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	var rel string
	if obj.EntityType == model.EntityTypeContact {
		rel = "CONTACT_BELONGS_TO_TENANT"
	} else {
		rel = "ORGANIZATION_BELONGS_TO_TENANT"
	}
	queryResult, err := tx.Run(ctx,
		fmt.Sprintf("MATCH (c:%s {id:$Id})-[:%s]->(:Tenant {name:$tenant}) "+
			" MERGE (f:%s:CustomField {name: $name, datatype:$datatype})<-[:HAS_PROPERTY]-(c) "+
			" ON CREATE SET f.%s=$value, f.id=randomUUID(), f.createdAt=$now, f.updatedAt=$now, f.source=$source, f.sourceOfTruth=$sourceOfTruth,  f:%s "+
			" ON MATCH SET f.%s=$value, f.sourceOfTruth=$sourceOfTruth, f.updatedAt=$now "+
			" RETURN f", obj.EntityType, rel, entity.NodeLabel(), entity.PropertyName(), "CustomField_"+tenant, entity.PropertyName()),
		map[string]any{
			"tenant":        tenant,
			"Id":            obj.ID,
			"name":          entity.Name,
			"datatype":      entity.DataType,
			"value":         entity.Value.RealValue(),
			"source":        entity.Source,
			"sourceOfTruth": entity.SourceOfTruth,
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *customFieldRepository) MergeCustomFieldToFieldSetInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, obj *model.CustomFieldEntityType, fieldSetId string, entity entity.CustomFieldEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldRepository.MergeCustomFieldToFieldSetInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	var rel string
	if obj.EntityType == model.EntityTypeContact {
		rel = "CONTACT_BELONGS_TO_TENANT"
	} else {
		rel = "ORGANIZATION_BELONGS_TO_TENANT"
	}

	queryResult, err := tx.Run(ctx,
		fmt.Sprintf(" MATCH (s:FieldSet {id:$fieldSetId})<-[:HAS_COMPLEX_PROPERTY]-(c:%s {id:$Id})-[:%s]->(:Tenant {name:$tenant}) "+
			" MERGE (f:%s:CustomField {name: $name, datatype:$datatype})<-[:HAS_PROPERTY]-(s)"+
			" ON CREATE SET f.%s=$value, f.id=randomUUID(), f.createdAt=$now, f.updatedAt=$now, f.source=$source, f.sourceOfTruth=$sourceOfTruth, f:%s "+
			" ON MATCH SET f.%s=$value, f.sourceOfTruth=$sourceOfTruth, f.updatedAt=$now "+
			" RETURN f", obj.EntityType, rel, entity.NodeLabel(), entity.PropertyName(), "CustomField_"+tenant, entity.PropertyName()),
		map[string]any{
			"tenant":        tenant,
			"Id":            obj.ID,
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldRepository.LinkWithCustomFieldTemplateForContactInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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

func (r *customFieldRepository) LinkWithCustomFieldTemplateInTx(ctx context.Context, tx neo4j.ManagedTransaction, fieldId string, obj *model.CustomFieldEntityType, templateId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldRepository.LinkWithCustomFieldTemplateInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	queryResult, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (f:CustomField {id:$fieldId})<-[:HAS_PROPERTY]-(c:%s {id:$Id}),
				  (c)-[:IS_DEFINED_BY]->(e:EntityTemplate),
				  (e)-[:CONTAINS]->(d:CustomFieldTemplate {id:$customFieldTemplateId})
			MERGE (f)-[r:IS_DEFINED_BY]->(d)
			RETURN r`, obj.EntityType),
		map[string]any{
			"customFieldTemplateId": templateId,
			"fieldId":               fieldId,
			"Id":                    obj.ID,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *customFieldRepository) LinkWithCustomFieldTemplateForFieldSetInTx(ctx context.Context, tx neo4j.ManagedTransaction, fieldId, fieldSetId, templateId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldRepository.LinkWithCustomFieldTemplateForFieldSetInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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

func (r *customFieldRepository) FindAll(ctx context.Context, session neo4j.SessionWithContext, tenant string, obj *model.CustomFieldEntityType) ([]*neo4j.Record, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldRepository.FindAll")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	var rel string
	if obj.EntityType == model.EntityTypeContact {
		rel = "CONTACT_BELONGS_TO_TENANT"
	} else {
		rel = "ORGANIZATION_BELONGS_TO_TENANT"
	}

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(`
				MATCH (c:%s {id:$contactId})-[:%s]->(:Tenant {name:$tenant}),
              		  (c)-[:HAS_PROPERTY]->(f:CustomField) 
				RETURN f ORDER BY f.name`, obj.EntityType, rel),
			map[string]any{
				"Id":     obj.ID,
				"tenant": tenant})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	return records.([]*neo4j.Record), err
}

func (r *customFieldRepository) GetCustomFields(ctx context.Context, session neo4j.SessionWithContext, tenant string, obj *model.CustomFieldEntityType) ([]*neo4j.Record, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldRepository.GetCustomFields")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		var rel string
		if obj.EntityType == model.EntityTypeContact {
			rel = "CONTACT_BELONGS_TO_TENANT"
		} else {
			rel = "ORGANIZATION_BELONGS_TO_TENANT"
		}
		query := fmt.Sprintf(`
				MATCH (n:%s {id:$id})-[:%s]->(:Tenant {name:$tenant}),
              		  (n)-[:HAS_PROPERTY]->(f:CustomField) 
				RETURN f ORDER BY f.name`, obj.EntityType, rel)
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"id":     obj.ID,
				"tenant": tenant})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})

	return records.([]*neo4j.Record), err
}

func (r *customFieldRepository) FindAllForFieldSet(ctx context.Context, session neo4j.SessionWithContext, tenant, fieldSetId string) ([]*neo4j.Record, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldRepository.FindAllForFieldSet")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		queryResult, err := tx.Run(ctx, `
				MATCH (s:FieldSet {id:$fieldSetId})<-[:HAS_COMPLEX_PROPERTY]-(:Contact|Organization)-[:CONTACT_BELONGS_TO_TENANT|ORGANIZATION_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldRepository.DeleteByIdFromContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldRepository.DeleteByNameFromContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldRepository.DeleteByIdFromFieldSet")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldRepository.UpdateForContactInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldRepository.UpdateForFieldSetInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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
