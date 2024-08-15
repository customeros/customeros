package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type CommentReadRepository interface {
	GetAllForIssues(ctx context.Context, tenant string, issueIds []string) ([]*utils.DbNodeAndId, error)
}

type commentReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewCommentReadRepository(driver *neo4j.DriverWithContext, database string) CommentReadRepository {
	return &commentReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *commentReadRepository) GetAllForIssues(ctx context.Context, tenant string, issueIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentRepository.GetAllForIssues")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

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
