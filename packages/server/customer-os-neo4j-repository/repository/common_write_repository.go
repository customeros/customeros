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

type LinkDetails struct {
	FromEntityId   string
	FromEntityType model.EntityType

	Relationship           model.EntityRelation
	RelationshipProperties *map[string]interface{}

	ToEntityId   string
	ToEntityType model.EntityType
}

type CommonWriteRepository interface {
	Link(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, details LinkDetails) error
	Unlink(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, details LinkDetails) error
	Delete(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, id, label string) error
	UpdateTimeProperty(ctx context.Context, tenant, nodeLabel, entityId, property string, value *time.Time) error
	UpdateInt64Property(ctx context.Context, tenant, nodeLabel, entityId, property string, value int64) error
	UpdateBoolProperty(ctx context.Context, tenant, nodeLabel, entityId, property string, value bool) error
	UpdateStringProperty(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, nodeLabel, entityId, property string, value string) error
	IncrementProperty(ctx context.Context, tenant, nodeLabel, entityId, property string) error
	RemoveProperty(ctx context.Context, tenant, nodeLabel, entityId, property string) error
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

func (r *commonWriteRepository) Link(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, details LinkDetails) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonWriteRepository.Link")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	params := map[string]any{
		"tenant":       tenant,
		"entityId":     details.FromEntityId,
		"withEntityId": details.ToEntityId,
	}

	cypher := fmt.Sprintf(`MATCH (parent:%s_%s {id:$entityId}) `, details.FromEntityType.Neo4jLabel(), tenant)
	cypher += fmt.Sprintf(`MATCH (child:%s_%s {id:$withEntityId}) `, details.ToEntityType.Neo4jLabel(), tenant)
	cypher += fmt.Sprintf(`MERGE (parent)-[rel:%s]->(child)`, details.Relationship.String())

	// If there are relationship properties, add a SET clause to the Cypher query
	if details.RelationshipProperties != nil && len(*details.RelationshipProperties) > 0 {
		cypher += " SET "
		props := []string{}
		for k, v := range *details.RelationshipProperties {
			props = append(props, fmt.Sprintf("rel.%s = $rel_%s", k, k))
			params["rel_"+k] = v
		}
		cypher += strings.Join(props, ", ")
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *commonWriteRepository) Unlink(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, details LinkDetails) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonWriteRepository.Unlink")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (parent:%s_%s {id:$entityId})-[r:%s]-(child:%s_%s {id:$withEntityId}) `, details.FromEntityType.Neo4jLabel(), tenant, details.Relationship.String(), details.ToEntityType.Neo4jLabel(), tenant)
	cypher += fmt.Sprintf(`DELETE r`)

	params := map[string]any{
		"tenant":       tenant,
		"entityId":     details.FromEntityId,
		"withEntityId": details.ToEntityId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *commonWriteRepository) Delete(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, id, label string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonWriteRepository.Delete")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})<-[:BELONGS_TO_TENANT]-(n:%s_%s {id:$id}) delete n`, label, tenant)

	params := map[string]any{
		"tenant": tenant,
		"id":     id,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	if tx == nil {
		session := utils.NewNeo4jWriteSession(ctx, *r.driver)
		defer session.Close(ctx)

		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			_, err := tx.Run(ctx, cypher, params)
			if err != nil {
				return nil, err
			}
			return nil, nil
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	} else {
		_, err := (*tx).Run(ctx, cypher, params)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	return nil
}

func (r *commonWriteRepository) UpdateTimeProperty(ctx context.Context, tenant, nodeLabel, entityId, property string, value *time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWriteRepository.UpdateTimeProperty")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

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

func (r *commonWriteRepository) UpdateInt64Property(ctx context.Context, tenant, nodeLabel, entityId, property string, value int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWriteRepository.UpdateInt64Property")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.SetTag(tracing.SpanTagEntityId, entityId)

	span.LogFields(log.String("property", string(property)), log.String("nodeLabel", nodeLabel), log.Object("value", value))

	cypher := fmt.Sprintf(`MATCH (n:%s:%s_%s {id: $entityId}) SET n.%s = $value, n.updatedAt=datetime()`, nodeLabel, nodeLabel, tenant, property)
	params := map[string]any{
		"entityId": entityId,
		"value":    value,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *commonWriteRepository) UpdateBoolProperty(ctx context.Context, tenant, nodeLabel, entityId, property string, value bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWriteRepository.UpdateBoolProperty")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, entityId)
	span.LogFields(log.String("property", property), log.String("nodeLabel", nodeLabel), log.Object("value", value))

	cypher := fmt.Sprintf(`MATCH (n:%s:%s_%s {id: $entityId}) SET n.%s = $value, n.updatedAt=datetime()`, nodeLabel, nodeLabel, tenant, property)
	params := map[string]any{
		"entityId": entityId,
		"value":    value,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *commonWriteRepository) UpdateStringProperty(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, nodeLabel, entityId, property string, value string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWriteRepository.UpdateStringProperty")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, entityId)
	span.LogFields(log.String("property", property), log.String("nodeLabel", nodeLabel), log.Object("value", value))

	cypher := fmt.Sprintf(`MATCH (n:%s:%s_%s {id: $entityId}) SET n.%s = $value, n.updatedAt=datetime()`, nodeLabel, nodeLabel, tenant, property)
	params := map[string]any{
		"entityId": entityId,
		"value":    value,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *commonWriteRepository) IncrementProperty(ctx context.Context, tenant, nodeLabel, entityId, property string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWriteRepository.IncrementProperty")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.SetTag(tracing.SpanTagEntityId, entityId)

	span.LogFields(log.String("property", property), log.String("nodeLabel", nodeLabel))

	cypher := fmt.Sprintf(`MATCH (n:%s_%s {id: $entityId}) 
			SET n.%s = case WHEN n.%s IS NULL THEN 1 ELSE n.%s+1 END,
				n.updatedAt=datetime()`, nodeLabel, tenant, property, property, property)
	params := map[string]any{
		"entityId": entityId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *commonWriteRepository) RemoveProperty(ctx context.Context, tenant, nodeLabel, entityId, property string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ContactWriteRepository.RemoveProperty")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.SetTag(tracing.SpanTagEntityId, entityId)

	span.LogFields(log.String("property", property), log.String("nodeLabel", nodeLabel))

	cypher := fmt.Sprintf(`MATCH (n:%s_%s {id: $entityId}) REMOVE n.%s SET n.updatedAt=datetime()`, nodeLabel, tenant, property)
	params := map[string]any{
		"entityId": entityId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
