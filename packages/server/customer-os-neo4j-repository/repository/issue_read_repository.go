package neo4j_repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type IssueReadRepository interface {
	GetById(ctx context.Context, tenant, issueId string) (*dbtype.Node, error)
	GetMatchedIssueId(ctx context.Context, tenant, externalSystem, externalId string) (string, error)
	GetIssueIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error)
	GetIssueCountByStatusForOrganization(ctx context.Context, tenant, organizationId string) (map[string]int64, error)
}

type issueReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewIssueReadRepository(driver *neo4j.DriverWithContext, database string) IssueReadRepository {
	return &issueReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *issueReadRepository) GetById(ctx context.Context, tenant, issueId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "IssueRepository.GetById")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId}) RETURN i`
	params := map[string]any{
		"tenant":  tenant,
		"issueId": issueId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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

func (r *issueReadRepository) GetMatchedIssueId(ctx context.Context, tenant, externalSystem, externalId string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "IssueRepository.GetMatchedIssueId")
	defer spans.Finish()

	spans.LogKV("externalSystem", externalSystem)
	spans.LogKV("externalId", externalId)

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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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

func (r *issueReadRepository) GetIssueIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "IssueRepository.GetIssueIdByExternalId")
	defer spans.Finish()

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})
					MATCH (t)<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue_%s)-[:IS_LINKED_WITH {externalId:$externalId}]->(e)
					RETURN i.id ORDER BY i.createdAt`, tenant)
	params := map[string]any{
		"tenant":           tenant,
		"externalId":       externalId,
		"externalSystemId": externalSystemId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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

func (r *issueReadRepository) GetIssueCountByStatusForOrganization(ctx context.Context, tenant, organizationId string) (map[string]int64, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "IssueReadRepository.GetIssueCountByStatusForOrganization")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(:Organization {id:$organizationId})<-[:REPORTED_BY]-(i:Issue)
			WITH DISTINCT i
			RETURN i.status AS status, COUNT(i) AS count`
	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return queryResult.Collect(ctx)
		}
	})
	if err != nil {
		return nil, err
	}
	output := make(map[string]int64)
	for _, v := range result.([]*neo4j.Record) {
		status := ""
		if v.Values[0] != nil {
			status = v.Values[0].(string)
		}
		output[status] = v.Values[1].(int64)
	}
	return output, err
}
