package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExternalSystemReadRepository.GetFirstExternalIdForLinkedEntity")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

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

func (r *externalSystemReadRepository) GetFor(ctx context.Context, tenant string, ids []string, label string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RxternalSystemReadRepository.GetFor")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("label", label))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem)<-[rel:IS_LINKED_WITH]-(n:%s)
			WHERE n.id IN $ids RETURN e, rel, n.id order by e.id, rel.syncDate`, label)
	params := map[string]any{
		"tenant": tenant,
		"ids":    ids,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

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
