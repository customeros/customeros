package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type FlowReadRepository interface {
	GetList(ctx context.Context) ([]*dbtype.Node, error)
	GetListWithContact(ctx context.Context, contactIds []string) ([]*utils.DbNodeAndId, error)
	GetById(ctx context.Context, id string) (*dbtype.Node, error)
}

type flowReadRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowReadRepository(driver *neo4j.DriverWithContext, database string) FlowReadRepository {
	return &flowReadRepositoryImpl{driver: driver, database: database}
}

func (r flowReadRepositoryImpl) GetList(ctx context.Context) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowReadRepository.GetList")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s) RETURN f`, common.GetTenantFromContext(ctx))
	params := map[string]any{
		"tenant": common.GetTenantFromContext(ctx),
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*dbtype.Node), nil
}

func (r flowReadRepositoryImpl) GetListWithContact(ctx context.Context, contactIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowReadRepository.GetListWithContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.Object("contactIds", contactIds))

	tenant := common.GetTenantFromContext(ctx)

	params := map[string]any{
		"tenant":     tenant,
		"contactIds": contactIds,
	}

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fc:FlowContact_%s)-[:HAS]->(c:Contact_%s) `, tenant, tenant, tenant)
	cypher += "where c.id in $contactIds RETURN fs, c.id"

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndId))))
	if len(result.([]*utils.DbNodeAndId)) == 0 {
		return nil, nil
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r flowReadRepositoryImpl) GetById(ctx context.Context, id string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowReadRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	tracing.LogObjectAsJson(span, "id", id)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s {id: $id}) RETURN f`, common.GetTenantFromContext(ctx))
	params := map[string]any{
		"tenant": common.GetTenantFromContext(ctx),
		"id":     id,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}
