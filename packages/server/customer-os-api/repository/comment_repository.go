package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type CommentRepository interface {
	GetAllForIssues(ctx context.Context, tenant string, issueIds []string) ([]*utils.DbNodeAndId, error)
}

type commentRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewCommentRepository(driver *neo4j.DriverWithContext, database string) CommentRepository {
	return &commentRepository{
		driver:   driver,
		database: database,
	}
}

func (r *commentRepository) GetAllForIssues(ctx context.Context, tenant string, issueIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentRepository.GetAllForIssues")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue)<-[:COMMENTED]-(c:Comment) 
				WHERE i.id IN $issueIds
				RETURN c, i.id ORDER BY c.createdAt ASC`
	params := map[string]any{
		"tenant":   tenant,
		"issueIds": issueIds,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

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
	return result.([]*utils.DbNodeAndId), nil
}
