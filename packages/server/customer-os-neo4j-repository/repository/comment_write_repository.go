package neo4j_repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type CommentCreateFields struct {
	Content          string             `json:"content"`
	CreatedAt        time.Time          `json:"createdAt"`
	ContentType      string             `json:"contentType"`
	AuthorUserId     string             `json:"authorUserId"`
	CommentedIssueId string             `json:"commentedIssueId"`
	SourceFields     model.SourceFields `json:"sourceFields"`
}

type CommentUpdateFields struct {
	Content     string `json:"content"`
	ContentType string `json:"contentType"`
	Source      string `json:"source"`
}

type CommentWriteRepository interface {
	Create(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, commentId string, data data_fields.CommentFields) error
	Update(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, commentId string, data data_fields.CommentFields) error
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

func (r *commentWriteRepository) Create(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, commentId string, data data_fields.CommentFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentWriteRepository.Create")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	tracing.TagEntity(span, commentId)
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
								c.updatedAt=datetime(),
								c.source=$source,
								c.appSource=$appSource,
								c.content=$content,
								c.contentType=$contentType	
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
		"source":           utils.IfNotNilString(data.Source),
		"appSource":        utils.IfNotNilString(data.AppSource),
		"content":          utils.IfNotNilString(data.Content),
		"contentType":      utils.IfNotNilString(data.ContentType),
		"commentedIssueId": utils.IfNotNilString(data.CommentedIssueId),
		"authorUserId":     utils.IfNotNilString(data.AuthorUserId),
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

func (r *commentWriteRepository) Update(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, commentId string, data data_fields.CommentFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommentWriteRepository.Update")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, commentId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (c:Comment_%s {id:$commentId})
		 	SET updatedAt = datetime()`, tenant)
	params := map[string]any{
		commentId: commentId,
	}
	if data.Content != nil {
		cypher += `, c.content = $content`
		params["content"] = *data.Content
	}
	if data.ContentType != nil {
		cypher += `, c.contentType = $contentType`
		params["contentType"] = *data.ContentType
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}
