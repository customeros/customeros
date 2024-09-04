package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"strings"
	"time"
)

type CommonWriteRepository interface {
	LinkEntityWithEntity(ctx context.Context, tenant, entityId string, entityType model.EntityType, relationship model.EntityRelation, relationshipProperties *map[string]interface{}, withEntityId string, withEntityType model.EntityType) error
	LinkEntityWithEntityInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, entityId string, entityType model.EntityType, relationship model.EntityRelation, relationshipProperties *map[string]interface{}, withEntityId string, withEntityType model.EntityType) error
	UnlinkEntityWithEntity(ctx context.Context, tenant, entityId string, entityType model.EntityType, relationship model.EntityRelation, withEntityId string, withEntityType model.EntityType) error
	UnlinkEntityWithEntityInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, entityId string, entityType model.EntityType, relationship model.EntityRelation, withEntityId string, withEntityType model.EntityType) error
	UpdateTimeProperty(ctx context.Context, tenant, nodeLabel, entityId, property string, value *time.Time) error
}

type commonWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewCommonWriteRepository(driver *neo4j.DriverWithContext, database string) CommonWriteRepository {
	return &commonWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *commonWriteRepository) LinkEntityWithEntity(ctx context.Context, tenant, entityId string, entityType model.EntityType, relationship model.EntityRelation, relationshipProperties *map[string]interface{}, withEntityId string, withEntityType model.EntityType) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonWriteRepository.LinkEntityWithEntity")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		err := r.LinkEntityWithEntityInTx(ctx, tx, tenant, entityId, entityType, relationship, relationshipProperties, withEntityId, withEntityType)
		return nil, err
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (r *commonWriteRepository) LinkEntityWithEntityInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, entityId string, entityType model.EntityType, relationship model.EntityRelation, relationshipProperties *map[string]interface{}, withEntityId string, withEntityType model.EntityType) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonWriteRepository.LinkEntityWithEntityInTx")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	params := map[string]any{
		"tenant":       tenant,
		"entityId":     entityId,
		"withEntityId": withEntityId,
	}

	cypher := fmt.Sprintf(`MATCH (parent:%s_%s {id:$entityId}) `, entityType.Neo4jLabel(), tenant)
	cypher += fmt.Sprintf(`MATCH (child:%s_%s {id:$withEntityId}) `, withEntityType.Neo4jLabel(), tenant)
	cypher += fmt.Sprintf(`MERGE (parent)-[rel:%s]->(child)`, relationship.String())

	// If there are relationship properties, add a SET clause to the Cypher query
	if relationshipProperties != nil && len(*relationshipProperties) > 0 {
		cypher += " SET "
		props := []string{}
		for k, v := range *relationshipProperties {
			props = append(props, fmt.Sprintf("rel.%s = $rel_%s", k, k))
			params["rel_"+k] = v
		}
		cypher += strings.Join(props, ", ")
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := tx.Run(ctx, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (r *commonWriteRepository) UnlinkEntityWithEntity(ctx context.Context, tenant, entityId string, entityType model.EntityType, relationship model.EntityRelation, withEntityId string, withEntityType model.EntityType) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonWriteRepository.UnlinkEntityWithEntity")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		err := r.UnlinkEntityWithEntityInTx(ctx, tx, tenant, entityId, entityType, relationship, withEntityId, withEntityType)
		return nil, err
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (r *commonWriteRepository) UnlinkEntityWithEntityInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, entityId string, entityType model.EntityType, relationship model.EntityRelation, withEntityId string, withEntityType model.EntityType) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonWriteRepository.UnlinkEntityWithEntityInTx")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := fmt.Sprintf(`MATCH (parent:%s_%s {id:$entityId})-[r:%s]-(child:%s_%s {id:$withEntityId}) `, entityType.Neo4jLabel(), tenant, relationship.String(), withEntityType.Neo4jLabel(), tenant)
	cypher += fmt.Sprintf(`DELETE r`)

	params := map[string]any{
		"tenant":       tenant,
		"entityId":     entityId,
		"withEntityId": withEntityId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := tx.Run(ctx, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (r *commonWriteRepository) UpdateTimeProperty(ctx context.Context, tenant, nodeLabel, entityId, property string, value *time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWriteRepository.UpdateTimeProperty")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, entityId)
	span.LogFields(log.String("property", string(property)), log.String("nodeLabel", nodeLabel), log.Object("value", value))

	cypher := fmt.Sprintf(`MATCH (n:%s:%s_%s {id: $entityId}) SET n.%s = $value`, nodeLabel, nodeLabel, tenant, property)
	params := map[string]any{
		"entityId": entityId,
		"value":    utils.TimePtrAsAny(value),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
