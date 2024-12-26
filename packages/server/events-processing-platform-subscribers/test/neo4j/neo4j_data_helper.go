package neo4j

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
)

func CreateIssue(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, issue neo4jentity.IssueEntity) string {
	issueId := utils.NewUUIDIfEmpty(issue.Id)
	query := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})
			  MERGE (t)<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$id})
				SET i:Issue_%s,
					i:TimelineEvent,
					i:TimelineEvent_%s,
					i.subject=$subject,
					i.status=$status,
					i.priority=$priority,
					i.description=$description,
					i.source=$source,
					i.sourceOfTruth=$sourceOfTruth
				`, tenant, tenant)

	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":        tenant,
		"id":            issueId,
		"subject":       issue.Subject,
		"status":        issue.Status,
		"priority":      issue.Priority,
		"description":   issue.Description,
		"source":        issue.Source,
		"sourceOfTruth": issue.SourceOfTruth,
	})
	return issueId
}

func CreateComment(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, comment neo4jentity.CommentEntity) string {
	commentId := utils.NewUUIDIfEmpty(comment.Id)
	query := fmt.Sprintf(`MERGE (c:Comment:Comment_%s {id:$id})
				SET c.content=$content,
					c.contentType=$contentType,
					c.source=$source,
					c.sourceOfTruth=$sourceOfTruth
				`, tenant)

	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":            commentId,
		"content":       comment.Content,
		"contentType":   comment.ContentType,
		"source":        comment.Source,
		"sourceOfTruth": comment.SourceOfTruth,
	})
	return commentId
}

func LinkIssueReportedBy(ctx context.Context, driver *neo4j.DriverWithContext, issueId, entityId string) {
	query := `MATCH (e {id:$entityId})
				MATCH (i:Issue {id:$issueId})
				MERGE (i)-[:REPORTED_BY]->(e)`

	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"issueId":  issueId,
		"entityId": entityId,
	})
}

func LinkIssueAssignedTo(ctx context.Context, driver *neo4j.DriverWithContext, issueId, entityId string) {
	query := `MATCH (e {id:$entityId})
				MATCH (i:Issue {id:$issueId})
				MERGE (i)-[:ASSIGNED_TO]->(e)`

	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"issueId":  issueId,
		"entityId": entityId,
	})
}

func LinkIssueFollowedBy(ctx context.Context, driver *neo4j.DriverWithContext, issueId, entityId string) {
	query := `MATCH (e {id:$entityId})
				MATCH (i:Issue {id:$issueId})
				MERGE (i)-[:FOLLOWED_BY]->(e)`

	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"issueId":  issueId,
		"entityId": entityId,
	})
}

func CreateExternalSystem(ctx context.Context, driver *neo4j.DriverWithContext, tenant, externalSystem string) {
	query := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})
			  MERGE (t)<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystemId})
				ON CREATE SET ext:ExternalSystem_%s`, tenant)

	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":           tenant,
		"externalSystemId": externalSystem,
	})
}

func CreateCountry(ctx context.Context, driver *neo4j.DriverWithContext, codeA2, codeA3, name, phoneCode string) {
	query := `MERGE (c:Country{codeA3: $codeA3}) 
				ON CREATE SET 
					c.phoneCode = $phoneCode,
					c.codeA2 = $codeA2,
					c.name = $name, 
					c.createdAt = $now, 
					c.updatedAt = $now`
	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"codeA2":    codeA2,
		"codeA3":    codeA3,
		"phoneCode": phoneCode,
		"name":      name,
		"now":       utils.Now(),
	})
}
