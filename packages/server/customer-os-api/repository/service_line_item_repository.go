package repository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type ServiceLineItemRepository interface {
	GetById(ctx context.Context, tenant, serviceLineItemId string) (*dbtype.Node, error)
	GetForContracts(ctx context.Context, tenant string, contractIds []string) ([]*utils.DbNodeAndId, error)
}

type serviceLineItemRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewServiceLineItemRepository(driver *neo4j.DriverWithContext, database string) ServiceLineItemRepository {
	return &serviceLineItemRepository{
		driver:   driver,
		database: database,
	}
}

func (r *serviceLineItemRepository) GetById(ctx context.Context, tenant, serviceLineItemId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:SERVICE_LINE_ITEM_BELONGS_TO_TENANT]-(s:ServiceLineItem {id:$serviceLineItemId}) 
		RETURN s`
	params := map[string]any{
		"tenant":            tenant,
		"serviceLineItemId": serviceLineItemId,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)
	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Single(ctx)
	})
	return utils.NodePtr(dbRecord.(*db.Record).Values[0].(dbtype.Node)), err
}

func (r *serviceLineItemRepository) GetForContracts(ctx context.Context, tenant string, contractIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ServiceLineItemRepository.GetForContracts")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract)-[:HAS_SERVICE_LINE_ITEM]->(serviceLineItem:ServiceLineItem)-[:SERVICE_LINE_ITEM_BELONGS_TO_TENANT]->(t)
			WHERE c.id IN $contractIds
			RETURN serviceLineItem, c.id ORDER BY serviceLineItem.createdAt DESC`
	params := map[string]any{
		"tenant":      tenant,
		"contractIds": contractIds,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}
