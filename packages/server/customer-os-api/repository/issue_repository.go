package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type IssueRepository interface {
	GetIssueCountByStatusForOrganization(ctx context.Context, tenant, organizationId string) (map[string]int64, error)
	GetById(ctx context.Context, tenant, issueId string) (*dbtype.Node, error)
	GetAllForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
}

type issueRepository struct {
	driver *neo4j.DriverWithContext
}

func NewIssueRepository(driver *neo4j.DriverWithContext) IssueRepository {
	return &issueRepository{
		driver: driver,
	}
}

func (r *issueRepository) GetIssueCountByStatusForOrganization(ctx context.Context, tenant, organizationId string) (map[string]int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.GetIssueCountByStatusForOrganization")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, `
			MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(:Organization {id:$organizationId})<-[:REPORTED_BY]-(i:Issue)
			WITH DISTINCT i
			RETURN i.status AS status, COUNT(i) AS count`,
			map[string]any{
				"tenant":         tenant,
				"organizationId": organizationId,
			}); err != nil {
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

func (r *issueRepository) GetById(ctx context.Context, tenant, issueId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (i:Issue_%s {id:$issueId}) RETURN i`

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]any{
				"issueId": issueId,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Single(ctx)
	})
	return utils.NodePtr(dbRecord.(*db.Record).Values[0].(dbtype.Node)), err
}

func (r *issueRepository) GetAllForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (e:InteractionEvent_%s)-[:PART_OF]->(i:Issue) 
		 WHERE e.id IN $ids 
		 RETURN i, e.id`, tenant)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"tenant": tenant,
				"ids":    ids,
			}); err != nil {
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
