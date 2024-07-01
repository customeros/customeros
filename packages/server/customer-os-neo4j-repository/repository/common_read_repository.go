package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type CommonReadRepository interface {
	ExistsById(ctx context.Context, tenant, id, label string) (bool, error)
	ExistsByIdLinkedTo(ctx context.Context, tenant, id, label, linkedToId, linkedToLabel, linkRelationship string) (bool, error)
	ExistsByIdLinkedFrom(ctx context.Context, tenant, id, label, linkedFromId, linkedFromLabel, linkRelationship string) (bool, error)
}

type commonReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewCommonReadRepository(driver *neo4j.DriverWithContext, database string) CommonReadRepository {
	return &commonReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *commonReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func (r *commonReadRepository) ExistsById(ctx context.Context, tenant, id, label string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonReadRepository.ExistsById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("id", id), log.String("label", label))

	cypher := fmt.Sprintf(`MATCH (n:%s {id:$id}) WHERE n:%s_%s RETURN n.id LIMIT 1`, label, label, tenant)
	params := map[string]any{
		"id": id,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return false, err
		} else {
			return queryResult.Next(ctx), nil
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}
	span.LogFields(log.Bool("result.exists", result.(bool)))
	return result.(bool), err
}

func (r *commonReadRepository) ExistsByIdLinkedTo(ctx context.Context, tenant, id, label, linkedToId, linkedToLabel, linkRelationship string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonReadRepository.ExistsById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("id", id), log.String("label", label), log.String("linkedToId", linkedToId), log.String("linkedToLabel", linkedToLabel), log.String("linkRelationship", linkRelationship))

	cypher := fmt.Sprintf(`MATCH (n:%s {id:$id})-`, label)
	if linkRelationship != "" {
		cypher += fmt.Sprintf(`[:%s]`, linkRelationship)
	}
	cypher += fmt.Sprintf(`->(m:%s {id:$linkedToId}) WHERE n:%s_%s RETURN n.id LIMIT 1`, linkedToLabel, label, tenant)
	params := map[string]any{
		"id":         id,
		"linkedToId": linkedToId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return false, err
		} else {
			return queryResult.Next(ctx), nil
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}
	span.LogFields(log.Bool("result.exists", result.(bool)))
	return result.(bool), err
}

func (r *commonReadRepository) ExistsByIdLinkedFrom(ctx context.Context, tenant, id, label, linkedFromId, linkedFromLabel, linkRelationship string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonReadRepository.ExistsById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("id", id), log.String("label", label), log.String("linkedFromId", linkedFromId), log.String("linkedFromLabel", linkedFromLabel), log.String("linkRelationship", linkRelationship))

	cypher := fmt.Sprintf(`MATCH (n:%s {id:$id})<-`, label)
	if linkRelationship != "" {
		cypher += fmt.Sprintf(`[:%s]`, linkRelationship)
	}
	cypher += fmt.Sprintf(`-(m:%s {id:$linkedFromId}) WHERE n:%s_%s RETURN n.id LIMIT 1`, linkedFromLabel, label, tenant)
	params := map[string]any{
		"id":           id,
		"linkedFromId": linkedFromId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return false, err
		} else {
			return queryResult.Next(ctx), nil
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}
	span.LogFields(log.Bool("result.exists", result.(bool)))
	return result.(bool), err
}

func (r *commonReadRepository) ExecuteIntegrityCheckerQuery(ctx context.Context, name, cypherQuery string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Neo4jRepository.ExecuteIntegrityCheckerQuery")
	defer span.Finish()
	span.SetTag("checker-name", name)
	span.LogFields(log.String("cypherQuery", cypherQuery))

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	countFoundRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypherQuery, map[string]any{})
		return utils.ExtractSingleRecordFirstValueAsType[int64](ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
	}
	span.LogFields(log.Int64("output - records", countFoundRecords.(int64)))
	return countFoundRecords.(int64), err
}
