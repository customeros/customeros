package neo4j_repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

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
	UpdateProperties(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, nodeLabel, entityId string, properties map[string]interface{}) error
	UpdateTimeProperty(ctx context.Context, tenant, nodeLabel, entityId, property string, value *time.Time) error
	UpdateTimePropertyInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, nodeLabel, entityId, property string, value *time.Time) error
	UpdateInt64Property(ctx context.Context, tenant, nodeLabel, entityId, property string, value int64) error
	UpdateBoolProperty(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, nodeLabel, entityId, property string, value bool) error
	UpdateStringProperty(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, nodeLabel, entityId, property string, value string) error
	IncrementProperty(ctx context.Context, tenant, nodeLabel, entityId, property string) error
	RemoveProperty(ctx context.Context, tenant, nodeLabel, entityId, property string) error
	TouchEntity(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, nodeLabel, entityId string) error
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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "CommonWriteRepository.Link")
	defer spans.Finish()

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

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *commonWriteRepository) Unlink(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, details LinkDetails) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "CommonWriteRepository.Unlink")
	defer spans.Finish()

	cypher := fmt.Sprintf(`MATCH (parent:%s_%s {id:$entityId})-[r:%s]-(child:%s_%s {id:$withEntityId}) `, details.FromEntityType.Neo4jLabel(), tenant, details.Relationship.String(), details.ToEntityType.Neo4jLabel(), tenant)
	cypher += `DELETE r`

	params := map[string]any{
		"tenant":       tenant,
		"entityId":     details.FromEntityId,
		"withEntityId": details.ToEntityId,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *commonWriteRepository) Delete(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, id, label string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "CommonWriteRepository.Delete")
	defer spans.Finish()

	spans.LogKV("id", id)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})<-[:BELONGS_TO_TENANT]-(n:%s_%s {id:$id}) delete n`, label, tenant)

	params := map[string]any{
		"tenant": tenant,
		"id":     id,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
			spans.TraceError(err)
			return err
		}
	} else {
		_, err := (*tx).Run(ctx, cypher, params)
		if err != nil {
			spans.TraceError(err)
			return err
		}
	}

	return nil
}

func (r *commonWriteRepository) UpdateProperties(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, nodeLabel, entityId string, properties map[string]interface{}) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "CommonWriteRepository.UpdateProperties")
	defer spans.Finish()

	spans.TagEntity(entityId)
	spans.LogKV("nodeLabel", nodeLabel)
	spans.LogObjectAsJson("properties", properties)

	// Build the dynamic Cypher query
	setClauses := make([]string, 0, len(properties))
	params := map[string]any{
		"entityId": entityId,
	}
	for prop, value := range properties {
		paramKey := fmt.Sprintf("prop_%s", prop) // Unique key for each property
		setClauses = append(setClauses, fmt.Sprintf("n.%s = $%s", prop, paramKey))
		params[paramKey] = value
	}
	setClauses = append(setClauses, "n.updatedAt = datetime()") // Ensure updatedAt is always updated

	cypher := fmt.Sprintf(
		`MATCH (n:%s:%s_%s {id: $entityId}) SET %s`,
		nodeLabel, nodeLabel, tenant, strings.Join(setClauses, ", "),
	)
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *commonWriteRepository) UpdateTimeProperty(ctx context.Context, tenant, nodeLabel, entityId, property string, value *time.Time) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "CommonWriteRepository.UpdateTimeProperty")
	defer spans.Finish()

	spans.TagEntity(entityId)

	spans.LogKV("property", string(property))
	spans.LogKV("nodeLabel", nodeLabel)
	spans.LogObjectAsJson("value", value)

	cypher := fmt.Sprintf(`MATCH (n:%s:%s_%s {id: $entityId}) SET n.%s = $value`, nodeLabel, nodeLabel, tenant, property)
	params := map[string]any{
		"entityId": entityId,
		"value":    utils.TimePtrAsAny(value),
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *commonWriteRepository) UpdateTimePropertyInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, nodeLabel, entityId, property string, value *time.Time) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "CommonWriteRepository.UpdateTimeProperty")
	defer spans.Finish()

	spans.TagEntity(entityId)

	spans.LogKV("property", string(property))
	spans.LogKV("nodeLabel", nodeLabel)
	spans.LogObjectAsJson("value", value)

	cypher := fmt.Sprintf(`MATCH (n:%s:%s_%s {id: $entityId}) SET n.%s = $value`, nodeLabel, nodeLabel, tenant, property)
	params := map[string]any{
		"entityId": entityId,
		"value":    utils.TimePtrAsAny(value),
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		return tx.Run(ctx, cypher, params)
	})
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *commonWriteRepository) UpdateInt64Property(ctx context.Context, tenant, nodeLabel, entityId, property string, value int64) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "CommonWriteRepository.UpdateInt64Property")
	defer spans.Finish()

	spans.TagEntity(entityId)

	spans.LogKV("property", string(property))
	spans.LogKV("nodeLabel", nodeLabel)
	spans.LogObjectAsJson("value", value)

	cypher := fmt.Sprintf(`MATCH (n:%s:%s_%s {id: $entityId}) SET n.%s = $value, n.updatedAt=datetime()`, nodeLabel, nodeLabel, tenant, property)
	params := map[string]any{
		"entityId": entityId,
		"value":    value,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *commonWriteRepository) UpdateBoolProperty(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, nodeLabel, entityId, property string, value bool) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "CommonWriteRepository.UpdateBoolProperty")
	defer spans.Finish()

	spans.TagEntity(entityId)
	spans.LogKV("property", property)
	spans.LogKV("nodeLabel", nodeLabel)
	spans.LogObjectAsJson("value", value)

	cypher := fmt.Sprintf(`MATCH (n:%s:%s_%s {id: $entityId}) SET n.%s = $value, n.updatedAt=datetime()`, nodeLabel, nodeLabel, tenant, property)
	params := map[string]any{
		"entityId": entityId,
		"value":    value,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *commonWriteRepository) UpdateStringProperty(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, nodeLabel, entityId, property string, value string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "CommonWriteRepository.UpdateStringProperty")
	defer spans.Finish()

	spans.TagEntity(entityId)
	spans.LogKV("property", property)
	spans.LogKV("nodeLabel", nodeLabel)
	spans.LogObjectAsJson("value", value)

	cypher := fmt.Sprintf(`MATCH (n:%s:%s_%s {id: $entityId}) SET n.%s = $value, n.updatedAt=datetime()`, nodeLabel, nodeLabel, tenant, property)
	params := map[string]any{
		"entityId": entityId,
		"value":    value,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		return nil, err
	})
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *commonWriteRepository) IncrementProperty(ctx context.Context, tenant, nodeLabel, entityId, property string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "CommonWriteRepository.IncrementProperty")
	defer spans.Finish()

	spans.TagEntity(entityId)
	spans.LogKV("property", property)
	spans.LogKV("nodeLabel", nodeLabel)

	cypher := fmt.Sprintf(`MATCH (n:%s_%s {id: $entityId}) 
			SET n.%s = case WHEN n.%s IS NULL THEN 1 ELSE n.%s+1 END,
				n.updatedAt=datetime()`, nodeLabel, tenant, property, property, property)
	params := map[string]any{
		"entityId": entityId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *commonWriteRepository) RemoveProperty(ctx context.Context, tenant, nodeLabel, entityId, property string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "CommonWriteRepository.RemoveProperty")
	defer spans.Finish()

	spans.TagEntity(entityId)
	spans.LogKV("property", property)
	spans.LogKV("nodeLabel", nodeLabel)

	cypher := fmt.Sprintf(`MATCH (n:%s_%s {id: $entityId}) REMOVE n.%s SET n.updatedAt=datetime()`, nodeLabel, tenant, property)
	params := map[string]any{
		"entityId": entityId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

// update the updatedAt property of the entity to current time
func (r *commonWriteRepository) TouchEntity(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, nodeLabel, entityId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "CommonWriteRepository.TouchEntity")
	defer spans.Finish()

	spans.TagEntity(entityId)
	spans.LogKV("nodeLabel", nodeLabel)

	cypher := fmt.Sprintf(`MATCH (n:%s_%s {id: $entityId}) SET n.updatedAt=datetime()`, nodeLabel, tenant)
	params := map[string]any{
		"entityId": entityId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
	}

	return err
}
