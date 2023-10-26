package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/comment/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type CommentRepository interface {
	Create(ctx context.Context, tenant, commentId string, evt event.CommentCreateEvent) error
	Update(ctx context.Context, tenant, commentId string, evt event.CommentUpdateEvent) error
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

func (r *commentRepository) Create(ctx context.Context, tenant, commentId string, evt event.CommentCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentRepository.Create")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("commentId", commentId), log.String("commentCreateEvent", fmt.Sprintf("%+v", evt)))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
							OPTIONAL MATCH (t)<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$commentedIssueId})
							WHERE $commentedIssueId <> ""
							WITH coalesce(i) as commentedNode, t
							WHERE commentedNode IS NOT NULL
							MERGE (c:Comment {id:$commentId})-[:COMMENTED]->(commentedNode)
							ON CREATE SET 
								c:Comment_%s,
								c.createdAt=$createdAt,
								c.updatedAt=$updatedAt,
								c.source=$source,
								c.sourceOfTruth=$sourceOfTruth,
								c.appSource=$appSource,
								c.content=$content,
								c.contentType=$contentType	
							ON MATCH SET
								c.content = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.content is null OR c.content = '' THEN $content ELSE c.content END,
								c.contentType = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.contentType is null OR c.contentType = '' THEN $contentType ELSE c.contentType END,
								c.updatedAt = $updatedAt,
								c.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE c.sourceOfTruth END
							WITH c, t
							OPTIONAL MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$authorUserId}) 
							WHERE $authorUserId <> ""
							FOREACH (ignore IN CASE WHEN u IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (c)-[:CREATED_BY]->(u))
							`, tenant)
	params := map[string]any{
		"tenant":           tenant,
		"commentId":        commentId,
		"createdAt":        evt.CreatedAt,
		"updatedAt":        evt.UpdatedAt,
		"source":           helper.GetSource(evt.Source),
		"sourceOfTruth":    helper.GetSourceOfTruth(evt.Source),
		"appSource":        helper.GetAppSource(evt.AppSource),
		"content":          evt.Content,
		"contentType":      evt.ContentType,
		"commentedIssueId": evt.CommentedIssueId,
		"authorUserId":     evt.AuthorUserId,
		"overwrite":        helper.GetSourceOfTruth(evt.Source) == constants.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return r.executeWriteQuery(ctx, cypher, params)
}

func (r *commentRepository) Update(ctx context.Context, tenant, commentId string, evt event.CommentUpdateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("commentId", commentId), log.Object("event", fmt.Sprintf("%+v", evt)))

	cypher := `MATCH (c:Comment {id:$commentId})
		 	SET	
				c.content = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.content is null OR c.content = '' THEN $content ELSE c.content END,
				c.contentType = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.contentType is null OR c.contentType = '' THEN $contentType ELSE c.contentType END,
				c.updatedAt = $updatedAt,
				c.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE c.sourceOfTruth END`
	params := map[string]any{
		"tenant":        tenant,
		"commentId":     commentId,
		"updatedAt":     evt.UpdatedAt,
		"content":       evt.Content,
		"contentType":   evt.ContentType,
		"sourceOfTruth": helper.GetSourceOfTruth(evt.Source),
		"overwrite":     helper.GetSourceOfTruth(evt.Source) == constants.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return r.executeWriteQuery(ctx, cypher, params)
}

func (r *commentRepository) executeWriteQuery(ctx context.Context, cypher string, params map[string]any) error {
	return utils.ExecuteWriteQueryOnDb(ctx, *r.driver, r.database, cypher, params)
}
