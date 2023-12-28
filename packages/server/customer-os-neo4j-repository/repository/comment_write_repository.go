package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/customer-os-neo4j-repository/constant"
	"github.com/openline-ai/customer-os-neo4j-repository/model"
	"github.com/openline-ai/customer-os-neo4j-repository/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type CommentCreateFields struct {
	Content          string       `json:"content"`
	CreatedAt        time.Time    `json:"createdAt"`
	UpdatedAt        time.Time    `json:"updatedAt"`
	ContentType      string       `json:"contentType"`
	AuthorUserId     string       `json:"authorUserId"`
	CommentedIssueId string       `json:"commentedIssueId"`
	SourceFields     model.Source `json:"sourceFields"`
}

type CommentUpdateFields struct {
	Content     string    `json:"content"`
	ContentType string    `json:"contentType"`
	Source      string    `json:"source"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CommentWriteRepository interface {
	Create(ctx context.Context, tenant, commentId string, data CommentCreateFields) error
	Update(ctx context.Context, tenant, commentId string, data CommentUpdateFields) error
}

type commentWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewCommentWriteRepository(driver *neo4j.DriverWithContext, database string) CommentWriteRepository {
	return &commentWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *commentWriteRepository) Create(ctx context.Context, tenant, commentId string, data CommentCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentWriteRepository.Create")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, commentId)
	tracing.LogObjectAsJson(span, "data", data)

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
		"createdAt":        data.CreatedAt,
		"updatedAt":        data.UpdatedAt,
		"source":           data.SourceFields.Source,
		"sourceOfTruth":    data.SourceFields.SourceOfTruth,
		"appSource":        data.SourceFields.AppSource,
		"content":          data.Content,
		"contentType":      data.ContentType,
		"commentedIssueId": data.CommentedIssueId,
		"authorUserId":     data.AuthorUserId,
		"overwrite":        data.SourceFields.Source == constant.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *commentWriteRepository) Update(ctx context.Context, tenant, commentId string, data CommentUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentWriteRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, commentId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := `MATCH (c:Comment {id:$commentId})
		 	SET	
				c.content = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.content is null OR c.content = '' THEN $content ELSE c.content END,
				c.contentType = CASE WHEN c.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR c.contentType is null OR c.contentType = '' THEN $contentType ELSE c.contentType END,
				c.updatedAt = $updatedAt,
				c.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE c.sourceOfTruth END`
	params := map[string]any{
		"tenant":        tenant,
		"commentId":     commentId,
		"updatedAt":     data.UpdatedAt,
		"content":       data.Content,
		"contentType":   data.ContentType,
		"sourceOfTruth": data.Source,
		"overwrite":     data.Source == constant.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
