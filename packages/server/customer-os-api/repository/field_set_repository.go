package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
)

type FieldSetRepository interface {
	LinkWithFieldSetTemplateInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, fieldSetId, templateId string, entityType model.EntityType) error
	MergeFieldSetToContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, error)
	MergeFieldSetInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, obj *model.CustomFieldEntityType, entity entity.FieldSetEntity) (*dbtype.Node, error)

	UpdateFieldSetForContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, error)
	DeleteByIdFromContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId, fieldSetId string) error
	FindAll(ctx context.Context, session neo4j.SessionWithContext, tenant string, obj *model.CustomFieldEntityType) ([]*neo4j.Record, error)
}

type fieldSetRepository struct {
	driver *neo4j.DriverWithContext
}

func NewFieldSetRepository(driver *neo4j.DriverWithContext) FieldSetRepository {
	return &fieldSetRepository{
		driver: driver,
	}
}

func (r *fieldSetRepository) LinkWithFieldSetTemplateInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, fieldSetId, templateId string, entityType model.EntityType) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FieldSetRepository.LinkWithFieldSetTemplateInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	var rel string
	if entityType == model.EntityTypeContact {
		rel = "CONTACT_BELONGS_TO_TENANT"
	} else {
		rel = "ORGANIZATION_BELONGS_TO_TENANT"
	}
	txResult, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (f:FieldSet {id:$fieldSetId})<-[:HAS_COMPLEX_PROPERTY]-(c:%s)-[:%s]->(:Tenant {name:$tenant}),
					(c)-[:IS_DEFINED_BY]->(e:EntityTemplate)-[:CONTAINS]->(d:FieldSetTemplate {id:$templateId})
			MERGE (f)-[r:IS_DEFINED_BY]->(d)
			RETURN r`, entityType, rel),
		map[string]any{
			"templateId": templateId,
			"fieldSetId": fieldSetId,
			"tenant":     tenant,
		})
	if err != nil {
		return err
	}
	_, err = txResult.Single(ctx)
	return err
}

func (r *fieldSetRepository) MergeFieldSetToContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FieldSetRepository.MergeFieldSetToContactInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}) " +
		" MERGE (f:FieldSet {name: $name})<-[r:HAS_COMPLEX_PROPERTY]-(c) " +
		" ON CREATE SET f.id=randomUUID(), " +
		"				f.createdAt=$now, " +
		"				f.createdAt=$now, " +
		"				f.source=$source, " +
		"				f.sourceOfTruth=$sourceOfTruth, " +
		"				f:%s " +
		" RETURN f"
	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "FieldSet_"+tenant),
		map[string]interface{}{
			"tenant":        tenant,
			"contactId":     contactId,
			"name":          entity.Name,
			"source":        entity.Source,
			"sourceOfTruth": entity.SourceOfTruth,
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *fieldSetRepository) MergeFieldSetInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, obj *model.CustomFieldEntityType, entity entity.FieldSetEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FieldSetRepository.MergeFieldSetInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	var rel string
	if obj.EntityType == model.EntityTypeContact {
		rel = "CONTACT_BELONGS_TO_TENANT"
	} else {
		rel = "ORGANIZATION_BELONGS_TO_TENANT"
	}

	query := "MATCH (c:%s {id:$Id})-[:%s]->(:Tenant {name:$tenant}) " +
		" MERGE (f:FieldSet {name: $name})<-[r:HAS_COMPLEX_PROPERTY]-(c) " +
		" ON CREATE SET f.id=randomUUID(), " +
		"				f.createdAt=$now, " +
		"				f.updatedAt=datetime(), " +
		"				f.source=$source, " +
		"				f.sourceOfTruth=$sourceOfTruth, " +
		"				f:%s " +
		" RETURN f"
	queryResult, err := tx.Run(ctx, fmt.Sprintf(query, obj.EntityType, rel, "FieldSet_"+tenant),
		map[string]interface{}{
			"tenant":        tenant,
			"Id":            obj.ID,
			"name":          entity.Name,
			"source":        entity.Source,
			"sourceOfTruth": entity.SourceOfTruth,
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *fieldSetRepository) UpdateFieldSetForContactInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, contactId string, entity entity.FieldSetEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FieldSetRepository.UpdateFieldSetForContactInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	queryResult, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
					(c)-[r:HAS_COMPLEX_PROPERTY]->(s:FieldSet {id:$fieldSetId})
            SET s.name=$name, s.sourceOfTruth=$sourceOfTruth, s.updatedAt=datetime()
			RETURN s`,
		map[string]interface{}{
			"tenant":        tenant,
			"contactId":     contactId,
			"fieldSetId":    entity.Id,
			"name":          entity.Name,
			"sourceOfTruth": entity.SourceOfTruth,
			"now":           utils.Now(),
		})
	return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
}

func (r *fieldSetRepository) DeleteByIdFromContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId, fieldSetId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FieldSetRepository.DeleteByIdFromContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, `
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
                  (c)-[:HAS_COMPLEX_PROPERTY]->(s:FieldSet {id:$fieldSetId}),
				  (s)-[:HAS_PROPERTY]->(f:CustomField)
            DETACH DELETE f, s`,
			map[string]any{
				"contactId":  contactId,
				"fieldSetId": fieldSetId,
				"tenant":     tenant,
			})
		return nil, err
	})
	return err
}

func (r *fieldSetRepository) FindAll(ctx context.Context, session neo4j.SessionWithContext, tenant string, obj *model.CustomFieldEntityType) ([]*neo4j.Record, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FieldSetRepository.FindAll")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		var rel string
		if obj.EntityType == model.EntityTypeContact {
			rel = "CONTACT_BELONGS_TO_TENANT"
		} else {
			rel = "ORGANIZATION_BELONGS_TO_TENANT"
		}

		queryResult, err := tx.Run(ctx, fmt.Sprintf(`
				MATCH (c:%s {id:$Id})-[:%s]->(:Tenant {name:$tenant}),
              			(c)-[:HAS_COMPLEX_PROPERTY]->(s:FieldSet) 
				RETURN s ORDER BY s.name`, obj.EntityType, rel),
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
