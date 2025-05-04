package neo4j_repository

import (
	"context"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "CommentRepository.GetAllForIssues")
	defer spans.Finish()

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue)<-[:COMMENTED]-(c:Comment) 
				WHERE i.id IN $issueIds
				RETURN c, i.id ORDER BY c.createdAt ASC`
	params := map[string]any{
		"tenant":   tenant,
		"issueIds": issueIds,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
		spans.TraceError(err)
		return nil, err
	}
	return result.([]*utils.DbNodeAndId), nil
}
