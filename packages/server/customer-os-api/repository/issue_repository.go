package repository

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/opentracing/opentracing-go/log"
)

type IssueRepository interface {
	GetAllForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
	GetSubmitterParticipantsForIssues(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
	GetReporterParticipantsForIssues(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
	GetAssigneeParticipantsForIssues(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
	GetFollowerParticipantsForIssues(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error)
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

func (r *issueRepository) GetAllForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "IssueRepository.GetAllForInteractionEvents")
	defer spans.Finish()

	query := fmt.Sprintf(`MATCH (e:InteractionEvent_%s)-[:PART_OF]->(i:Issue) 
		 WHERE e.id IN $ids 
		 RETURN i, e.id`, tenant)
	spans.LogKV(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
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

func (r *issueRepository) GetSubmitterParticipantsForIssues(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "IssueRepository.GetSubmitterParticipantsForIssues")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue)-[:SUBMITTED_BY]->(p:User|Contact|Organization)
			WHERE i.id IN $ids
			RETURN p, i.id`
	params := map[string]any{
		"tenant": tenant,
		"ids":    ids,
	}
	spans.LogKV(log.String("cypher", cypher), log.Object("params", params))

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
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *issueRepository) GetReporterParticipantsForIssues(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "IssueRepository.GetReporterParticipantsForIssues")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue)-[:REPORTED_BY]->(p:User|Contact|Organization)
			WHERE i.id IN $ids
			RETURN p, i.id`
	params := map[string]any{
		"tenant": tenant,
		"ids":    ids,
	}
	spans.LogKV(log.String("cypher", cypher), log.Object("params", params))

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
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *issueRepository) GetAssigneeParticipantsForIssues(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "IssueRepository.GetAssigneeParticipantsForIssues")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue)-[:ASSIGNED_TO]->(p:User|Contact|Organization)
			WHERE i.id IN $ids
			RETURN p, i.id`
	params := map[string]any{
		"tenant": tenant,
		"ids":    ids,
	}
	spans.LogKV(log.String("cypher", cypher), log.Object("params", params))

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
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r *issueRepository) GetFollowerParticipantsForIssues(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeAndId, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "IssueRepository.GetFollowerParticipantsForIssues")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue)-[:FOLLOWED_BY]->(p:User|Contact|Organization)
			WHERE i.id IN $ids
			RETURN p, i.id`
	params := map[string]any{
		"tenant": tenant,
		"ids":    ids,
	}
	spans.LogKV(log.String("cypher", cypher), log.Object("params", params))

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
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), err
}
