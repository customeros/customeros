package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
)

type EntityTemplateRepository interface {
	Create(ctx context.Context, tenant string, entity *entity.EntityTemplateEntity) (any, error)
	FindAllByTenant(ctx context.Context, session neo4j.SessionWithContext, tenant string) ([]*db.Record, error)
	FindAllByTenantAndExtends(ctx context.Context, session neo4j.SessionWithContext, tenant, extends string) ([]*db.Record, error)
	FindById(ctx context.Context, tenant string, obj *model.CustomFieldEntityType) (any, error)
}

type entityTemplateRepository struct {
	driver       *neo4j.DriverWithContext
	repositories *Repositories
}

func NewEntityTemplateRepository(driver *neo4j.DriverWithContext, repositories *Repositories) EntityTemplateRepository {
	return &entityTemplateRepository{
		driver:       driver,
		repositories: repositories,
	}
}

func (r *entityTemplateRepository) Create(ctx context.Context, tenant string, entity *entity.EntityTemplateEntity) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EntityTemplateRepository.Create")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	queryResult, err := session.ExecuteWrite(ctx, r.createFullEntityTemplateInTxWork(ctx, tenant, entity))
	if err != nil {
		return nil, err
	}
	return queryResult, nil
}

func (r *entityTemplateRepository) FindAllByTenant(ctx context.Context, session neo4j.SessionWithContext, tenant string) ([]*db.Record, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EntityTemplateRepository.FindAllByTenant")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	if result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
				MATCH (:Tenant {name:$tenant})<-[r:ENTITY_TEMPLATE_BELONGS_TO_TENANT]-(e:EntityTemplate) RETURN e`,
			map[string]any{
				"tenant": tenant,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	}); err != nil {
		return nil, err
	} else {
		return result.([]*db.Record), nil
	}
}

func (r *entityTemplateRepository) FindAllByTenantAndExtends(ctx context.Context, session neo4j.SessionWithContext, tenant, extends string) ([]*db.Record, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EntityTemplateRepository.FindAllByTenantAndExtends")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	if result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `
				MATCH (:Tenant {name:$tenant})<-[r:ENTITY_TEMPLATE_BELONGS_TO_TENANT]-(e:EntityTemplate) 
					WHERE e.extends=$extends
				RETURN e`,
			map[string]any{
				"tenant":  tenant,
				"extends": extends,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	}); err != nil {
		return nil, err
	} else {
		return result.([]*db.Record), nil
	}
}

func (r *entityTemplateRepository) FindById(ctx context.Context, tenant string, obj *model.CustomFieldEntityType) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EntityTemplateRepository.FindById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	var rel string
	if obj.EntityType == model.EntityTypeContact {
		rel = "CONTACT_BELONGS_TO_TENANT"
	} else {
		rel = "ORGANIZATION_BELONGS_TO_TENANT"
	}

	return session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(`
				MATCH (c:%s {id:$Id})-[:%s]->(t:Tenant {name:$tenant})<-[r:ENTITY_TEMPLATE_BELONGS_TO_TENANT]->(e:EntityTemplate),
					(c)-[:IS_DEFINED_BY]->(e)
					RETURN e`, obj.EntityType, rel),
			map[string]any{
				"tenant": tenant,
				"Id":     obj.ID,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
}

func (r *entityTemplateRepository) createFullEntityTemplateInTxWork(ctx context.Context, tenant string, entity *entity.EntityTemplateEntity) func(tx neo4j.ManagedTransaction) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "EntityTemplateRepository.createFullEntityTemplateInTxWork")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	return func(tx neo4j.ManagedTransaction) (any, error) {
		query := "MATCH (t:Tenant {name:$tenant}) " +
			" MERGE (t)<-[r:ENTITY_TEMPLATE_BELONGS_TO_TENANT]-(e:EntityTemplate {id: randomUUID()}) " +
			" ON CREATE SET e:%s, " +
			" 				e.extends=$extends, " +
			"				e.createdAt=datetime({timezone: 'UTC'}), " +
			"				e.name=$name, " +
			"				e.version=$version" +
			" RETURN e"
		txResult, err := tx.Run(ctx, fmt.Sprintf(query, "EntityTemplate_"+tenant),
			map[string]any{
				"tenant":  tenant,
				"name":    entity.Name,
				"version": entity.Version,
				"extends": entity.Extends,
			})
		if err != nil {
			return nil, err
		}
		records, err := txResult.Collect(ctx)
		if err != nil {
			return nil, err
		}
		entityTemplateId := utils.GetPropsFromNode(records[0].Values[0].(dbtype.Node))["id"].(string)
		for _, v := range entity.FieldSets {
			err := r.repositories.FieldSetTemplateRepository.createFieldSetTemplateInTx(ctx, tx, tenant, entityTemplateId, v)
			if err != nil {
				return nil, err
			}
		}
		for _, v := range entity.CustomFields {
			err := r.repositories.CustomFieldTemplateRepository.createCustomFieldTemplateForEntityInTx(ctx, tx, tenant, entityTemplateId, v)
			if err != nil {
				return nil, err
			}
		}
		return records, nil
	}
}
