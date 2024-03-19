package neo4j

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jtest "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/test"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform-subscribers/graph_db/entity"
)

func CreateSocial(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, social neo4jentity.SocialEntity) string {
	socialId := utils.NewUUIDIfEmpty(social.Id)
	query := fmt.Sprintf(`MERGE (s:Social:Social_%s {id: $id})
				SET s.url=$url,
					s.platformName=$platformName
				`, tenant)

	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":           socialId,
		"url":          social.Url,
		"platformName": social.PlatformName,
	})
	return socialId
}

func CreateContact(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, contact entity.ContactEntity) string {
	contactId := contact.Id
	if contactId == "" {
		contactId = uuid.New().String()
	}
	query := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})
			  MERGE (t)<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$id})
				SET c:Contact_%s,
					c.firstName=$firstName,
					c.lastName=$lastName
				`, tenant)

	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":    tenant,
		"id":        contactId,
		"firstName": contact.FirstName,
		"lastName":  contact.LastName,
	})
	return contactId
}

func CreateJobRole(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, jobRole entity.JobRoleEntity) string {
	jobRoleId := jobRole.Id
	if jobRoleId == "" {
		jobRoleId = uuid.New().String()
	}
	query := fmt.Sprintf(`CREATE (jobRole:JobRole:JobRole_%s {id:$jobRoleId})`, tenant)

	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"jobRoleId": jobRoleId,
	})
	return jobRoleId

}

func CreateIssue(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, issue entity.IssueEntity) string {
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

func CreateInteractionEvent(ctx context.Context, driver *neo4j.DriverWithContext, tenant string, interactionEvent entity.InteractionEventEntity) string {
	interactionEventId := utils.NewUUIDIfEmpty(interactionEvent.Id)
	query := fmt.Sprintf(`MERGE (i:InteractionEvent {id:$id})
				SET i:InteractionEvent_%s,
					i:TimelineEvent,
					i:TimelineEvent_%s,
					i.content=$content,
					i.contentType=$contentType,
					i.channel=$channel,
					i.channelData=$channelData,
					i.identifier=$identifier,
					i.eventType=$eventType,
					i.source=$source,
					i.sourceOfTruth=$sourceOfTruth
				`, tenant, tenant)

	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"id":            interactionEventId,
		"content":       interactionEvent.Content,
		"contentType":   interactionEvent.ContentType,
		"channel":       interactionEvent.Channel,
		"channelData":   interactionEvent.ChannelData,
		"identifier":    interactionEvent.Identifier,
		"eventType":     interactionEvent.EventType,
		"source":        interactionEvent.Source,
		"sourceOfTruth": interactionEvent.SourceOfTruth,
	})
	return interactionEventId
}

func LinkTag(ctx context.Context, driver *neo4j.DriverWithContext, tagId, entityId string) {

	query := `MATCH (e {id:$entityId})
				MATCH (t:Tag {id:$tagId})
				MERGE (e)-[rel:TAGGED]->(t)
				SET rel.taggedAt=$now`

	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tagId":    tagId,
		"entityId": entityId,
		"now":      utils.Now(),
	})
}

func LinkSocial(ctx context.Context, driver *neo4j.DriverWithContext, socialId, entityId string) {
	query := `MATCH (e {id:$entityId})
				MATCH (s:Social {id:$socialId})
				MERGE (e)-[:HAS]->(s)`

	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"socialId": socialId,
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

func CreateWorkspace(ctx context.Context, driver *neo4j.DriverWithContext, workspace string, provider string, tenant string) {
	query := `MATCH (t:Tenant {name: $tenant})
			  MERGE (t)-[:HAS_WORKSPACE]->(w:Workspace {name:$workspace, provider:$provider})`

	neo4jtest.ExecuteWriteQuery(ctx, driver, query, map[string]any{
		"tenant":    tenant,
		"provider":  provider,
		"workspace": workspace,
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
