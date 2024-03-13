package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type ExternalSystemReadRepository interface {
	GetFirstExternalIdForLinkedEntity(ctx context.Context, tenant, externalSystemId, entityId, entityLabel string) (string, error)
	GetAllExternalIdsForLinkedEntity(ctx context.Context, tenant, externalSystemId, entityId, entityLabel string) ([]string, error)
	GetAllForTenant(ctx context.Context, tenant string) ([]*dbtype.Node, error)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemReadRepository.GetFirstExternalIdForLinkedEntity")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("externalSystemId", externalSystemId), log.String("entityId", entityId), log.String("entityLabel", entityLabel))

	cypher := fmt.Sprintf(`
		MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystemId})
		MATCH (entity:%s {id:$entityId})-[rel:IS_LINKED_WITH]->(ext)
		RETURN rel.externalId ORDER BY rel.syncDate`, entityLabel)
	params := map[string]any{
		"tenant":           tenant,
		"externalSystemId": externalSystemId,
		"entityId":         entityId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
		tracing.TraceErr(span, err)
		return "", err
	}
	span.LogFields(log.Int("result.count", len(result.([]string))))
	if len(result.([]string)) > 0 {
		span.LogFields(log.String("result.externalId", result.([]string)[0]))
		return result.([]string)[0], nil
	} else {
		return "", nil
	}
}

func (r *externalSystemReadRepository) GetAllExternalIdsForLinkedEntity(ctx context.Context, tenant, externalSystemId, entityId, entityLabel string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemReadRepository.GetAllExternalIdsForLinkedEntity")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("externalSystemId", externalSystemId), log.String("entityId", entityId), log.String("entityLabel", entityLabel))

	cypher := fmt.Sprintf(`
		MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystemId})
		MATCH (entity:%s {id:$entityId})-[rel:IS_LINKED_WITH]->(ext)
		RETURN rel.externalId`, entityLabel)
	params := map[string]any{
		"tenant":           tenant,
		"entityId":         entityId,
		"externalSystemId": externalSystemId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]string))))
	return result.([]string), nil
}

func (r *externalSystemReadRepository) GetAllForTenant(ctx context.Context, tenant string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemReadRepository.GetAllForTenant")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem) RETURN ext`
	params := map[string]any{
		"tenant": tenant,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*dbtype.Node))))
	return result.([]*dbtype.Node), err

}
