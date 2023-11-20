package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/net/context"
)

type OpportunityRepository interface {
	GetById(ctx context.Context, tenant, opportunityId string) (*dbtype.Node, error)
	GetForContracts(ctx context.Context, tenant string, contractIds []string) ([]*utils.DbNodeAndId, error)
}

type opportunityRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewOpportunityRepository(driver *neo4j.DriverWithContext, database string) OpportunityRepository {
	return &opportunityRepository{
		driver:   driver,
		database: database,
	}
}

func (r *opportunityRepository) GetById(ctx context.Context, tenant, opportunityId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (o:Opportunity {id:$opportunityId}) 
		WHERE o:Opportunity_%s
		RETURN o`, tenant)
	params := map[string]any{
		"opportunityId": opportunityId,
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

func (r *opportunityRepository) GetForContracts(ctx context.Context, tenant string, contractIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "OpportunityRepository.GetForContracts")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:CONTRACT_BELONGS_TO_TENANT]-(c:Contract)-[:HAS_OPPORTUNITY]->(o:Opportunity)
			WHERE c.id IN $contractIds
			RETURN o, c.id ORDER BY o.createdAt DESC`
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
