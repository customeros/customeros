package neo4j_repository

import (
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"

	"golang.org/x/net/context"
)

type ExternalSystemReadRepository interface {
	GetFirstExternalIdForLinkedEntity(ctx context.Context, tenant, externalSystemId, entityId, entityLabel string) (string, error)
	GetAllExternalIdsForLinkedEntity(ctx context.Context, tenant, externalSystemId, entityId, entityLabel string) ([]string, error)
	GetAllForTenant(ctx context.Context, tenant string) ([]*dbtype.Node, error)
	GetFor(ctx context.Context, tenant string, ids []string, label string) ([]*utils.DbNodeWithRelationAndId, error)
}

type externalSystemReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewExternalSystemReadRepository(driver *neo4j.DriverWithContext, database string) ExternalSystemReadRepository {
	return &externalSystemReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *externalSystemReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *externalSystemReadRepository) GetFirstExternalIdForLinkedEntity(ctx context.Context, tenant, externalSystemId, entityId, entityLabel string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ExternalSystemReadRepository.GetFirstExternalIdForLinkedEntity")
	defer spans.Finish()

	spans.LogKV("externalSystemId", externalSystemId)
	spans.LogKV("entityId", entityId)
	spans.LogKV("entityLabel", entityLabel)

	cypher := fmt.Sprintf(`
		MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystemId})
		MATCH (entity:%s {id:$entityId})-[rel:IS_LINKED_WITH]->(ext)
		RETURN rel.externalId ORDER BY rel.syncDate`, entityLabel)
	params := map[string]any{
		"tenant":           tenant,
		"externalSystemId": externalSystemId,
		"entityId":         entityId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return "", err
		} else {
			return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return "", err
	}
	spans.LogKV("result.count", len(result.([]string)))
	if len(result.([]string)) > 0 {
		spans.LogKV("result.externalId", result.([]string)[0])
		return result.([]string)[0], nil
	} else {
		return "", nil
	}
}

func (r *externalSystemReadRepository) GetAllExternalIdsForLinkedEntity(ctx context.Context, tenant, externalSystemId, entityId, entityLabel string) ([]string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ExternalSystemReadRepository.GetAllExternalIdsForLinkedEntity")
	defer spans.Finish()

	spans.LogKV("externalSystemId", externalSystemId)
	spans.LogKV("entityId", entityId)
	spans.LogKV("entityLabel", entityLabel)

	cypher := fmt.Sprintf(`
		MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystemId})
		MATCH (entity:%s {id:$entityId})-[rel:IS_LINKED_WITH]->(ext)
		RETURN rel.externalId`, entityLabel)
	params := map[string]any{
		"tenant":           tenant,
		"entityId":         entityId,
		"externalSystemId": externalSystemId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return "", err
		} else {
			return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]string)))
	return result.([]string), nil
}

func (r *externalSystemReadRepository) GetAllForTenant(ctx context.Context, tenant string) ([]*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ExternalSystemReadRepository.GetAllForTenant")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem) RETURN ext`
	params := map[string]any{
		"tenant": tenant,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("result.count", len(result.([]*dbtype.Node)))
	return result.([]*dbtype.Node), err
}

func (r *externalSystemReadRepository) GetFor(ctx context.Context, tenant string, ids []string, label string) ([]*utils.DbNodeWithRelationAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "ExternalSystemReadRepository.GetFor")
	defer spans.Finish()

	spans.LogKV("label", label)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem)<-[rel:IS_LINKED_WITH]-(n:%s)
			WHERE n.id IN $ids RETURN e, rel, n.id order by e.id, rel.syncDate`, label)
	params := map[string]any{
		"tenant": tenant,
		"ids":    ids,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeWithRelationAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}
