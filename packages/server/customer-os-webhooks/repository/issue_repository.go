package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type IssueRepository interface {
	GetById(ctx context.Context, tenant, issueId string) (*dbtype.Node, error)
	GetMatchedIssueId(ctx context.Context, tenant, externalSystem, externalId string) (string, error)
	GetIssueIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error)
}

type issueRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewIssueRepository(driver *neo4j.DriverWithContext, database string) IssueRepository {
	return &issueRepository{
		driver:   driver,
		database: database,
	}
}

func (r *issueRepository) GetById(parentCtx context.Context, tenant, issueId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "IssueRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("issueId", issueId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId}) RETURN i`
	params := map[string]any{
		"tenant":  tenant,
		"issueId": issueId,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return dbRecord.(*dbtype.Node), err
}

func (r *issueRepository) GetMatchedIssueId(ctx context.Context, tenant, externalSystem, externalId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.GetMatchedIssueId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("externalSystem", externalSystem), log.String("externalId", externalId))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				MATCH (t)<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue)-[:IS_LINKED_WITH {externalId:$issueExternalId}]->(e)
				RETURN i.id LIMIT 1`
	params := map[string]interface{}{
		"tenant":          tenant,
		"externalSystem":  externalSystem,
		"issueExternalId": externalId,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	issueIDs := dbRecords.([]*db.Record)
	if len(issueIDs) == 1 {
		return issueIDs[0].Values[0].(string), nil
	}
	return "", nil
}

func (r *issueRepository) GetIssueIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.GetIssueIdByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})
					MATCH (t)<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue_%s)-[:IS_LINKED_WITH {externalId:$externalId}]->(e)
					RETURN i.id ORDER BY is.createdAt`, tenant)
	params := map[string]any{
		"tenant":           tenant,
		"externalId":       externalId,
		"externalSystemId": externalSystemId,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}
