package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
)

type CustomFieldTemplateRepository interface {
	Merge(ctx context.Context, tenant string, inputEntity entity.CustomFieldTemplateEntity) (*dbtype.Node, error)
	createCustomFieldTemplateForEntityInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, entityTemplateId string, entity *entity.CustomFieldTemplateEntity) error
	createCustomFieldTemplateForFieldSetInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, fieldSetTemplateId string, entity *entity.CustomFieldTemplateEntity) error
	FindAllByEntityTemplateId(ctx context.Context, entityTemplateId string) (any, error)
	FindAllByEntityFieldSetTemplateId(ctx context.Context, fieldSetTemplateId string) (any, error)
	FindByCustomFieldId(ctx context.Context, fieldSetId string) (any, error)
}

type customFieldTemplateRepository struct {
	driver *neo4j.DriverWithContext
}

func NewCustomFieldTemplateRepository(driver *neo4j.DriverWithContext) CustomFieldTemplateRepository {
	return &customFieldTemplateRepository{
		driver: driver,
	}
}

func (r *customFieldTemplateRepository) Merge(ctx context.Context, tenant string, inputEntity entity.CustomFieldTemplateEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldTemplateRepository.Create")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant}) 
		 MERGE (t)<-[:CUSTOM_FIELD_TEMPLATE_BELONGS_TO_TENANT]-(cft:CustomFieldTemplate {name:$name}) 
		 ON CREATE SET 
		  	cft.id=randomUUID(),
		  	cft.createdAt=$now,
		  	cft.updatedAt=$now,
		  	cft.order=$order,
			cft.mandatory=$mandatory,
			cft.type=$type,
			cft.length=$length,
			cft.min=$min,
			cft.max=$max,
		  	cft:CustomFieldTemplate_%s
		 RETURN cft`, tenant)

	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant":    tenant,
				"name":      inputEntity.Name,
				"order":     inputEntity.Order,
				"mandatory": inputEntity.Mandatory,
				"type":      inputEntity.Type,
				"length":    inputEntity.Length,
				"min":       inputEntity.Min,
				"max":       inputEntity.Max,
				"now":       utils.Now(),
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *customFieldTemplateRepository) createCustomFieldTemplateForEntityInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, entityTemplateId string, entity *entity.CustomFieldTemplateEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldTemplateRepository.createCustomFieldTemplateForEntityInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (e:EntityTemplate {id:$entityTemplateId}) " +
		" MERGE (e)-[:CONTAINS]->(f:CustomFieldTemplate {id:randomUUID(), name:$name}) " +
		" ON CREATE SET f:%s, " +
		"				f.createdAt=$now, " +
		"				f.updated=$now, " +
		"  				f.order=$order, " +
		"				f.mandatory=$mandatory, " +
		"				f.type=$type, " +
		"				f.length=$length, " +
		"  				f.min=$min, " +
		"				f.max=$max"

	_, err := tx.Run(ctx, fmt.Sprintf(query, "CustomFieldTemplate_"+tenant),
		map[string]any{
			"entityTemplateId": entityTemplateId,
			"name":             entity.Name,
			"order":            entity.Order,
			"mandatory":        entity.Mandatory,
			"type":             entity.Type,
			"length":           entity.Length,
			"min":              entity.Min,
			"max":              entity.Max,
			"now":              utils.Now(),
		})

	return err
}

func (r *customFieldTemplateRepository) createCustomFieldTemplateForFieldSetInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, fieldSetTemplateId string, entity *entity.CustomFieldTemplateEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldTemplateRepository.createCustomFieldTemplateForFieldSetInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := "MATCH (d:FieldSetTemplate {id:$fieldSetTemplateId}) " +
		" MERGE (d)-[:CONTAINS]->(f:CustomFieldTemplate {id:randomUUID(), name:$name}) " +
		" ON CREATE SET f:%s, " +
		"				f.createdAt=$now, " +
		"				f.updatedAt=$now, " +
		"				f.order=$order, " +
		"				f.mandatory=$mandatory, " +
		"				f.type=$type, " +
		"				f.length=$length, " +
		"				f.min=$min, " +
		"				f.max=$max"
	_, err := tx.Run(ctx, fmt.Sprintf(query, "CustomFieldTemplate_"+tenant),
		map[string]any{
			"fieldSetTemplateId": fieldSetTemplateId,
			"name":               entity.Name,
			"order":              entity.Order,
			"mandatory":          entity.Mandatory,
			"type":               entity.Type,
			"length":             entity.Length,
			"min":                entity.Min,
			"max":                entity.Max,
			"now":                utils.Now(),
		})

	return err
}

func (r *customFieldTemplateRepository) FindAllByEntityTemplateId(ctx context.Context, entityTemplateId string) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldTemplateRepository.FindAllByEntityTemplateId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	return session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, `
				MATCH (:EntityTemplate {id:$entityTemplateId})-[:CONTAINS]->(f:CustomFieldTemplate) RETURN f ORDER BY f.order`,
			map[string]any{
				"entityTemplateId": entityTemplateId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
}

func (r *customFieldTemplateRepository) FindAllByEntityFieldSetTemplateId(ctx context.Context, fieldSetTemplateId string) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldTemplateRepository.FindAllByEntityFieldSetTemplateId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	return session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, `
				MATCH (:FieldSetTemplate {id:$fieldSetTemplateId})-[:CONTAINS]->(f:CustomFieldTemplate) RETURN f ORDER BY f.order`,
			map[string]any{
				"fieldSetTemplateId": fieldSetTemplateId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
}

func (r *customFieldTemplateRepository) FindByCustomFieldId(ctx context.Context, customFieldId string) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CustomFieldTemplateRepository.FindByCustomFieldId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	return session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		queryResult, err := tx.Run(ctx, `
				MATCH (:CustomField {id:$customFieldId})-[:IS_DEFINED_BY]->(d:CustomFieldTemplate)
					RETURN d`,
			map[string]any{
				"customFieldId": customFieldId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
}
