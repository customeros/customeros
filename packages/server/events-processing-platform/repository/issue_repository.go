package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/issue/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type IssueRepository interface {
	Create(ctx context.Context, tenant, issueId string, evt event.IssueCreateEvent) error
	Update(ctx context.Context, tenant, issueId string, eventData event.IssueUpdateEvent) error
	AddUserAssignee(ctx context.Context, tenant, issueId, userId string, at time.Time) error
	RemoveUserAssignee(ctx context.Context, tenant, issueId, userId string, at time.Time) error
	AddUserFollower(ctx context.Context, tenant, issueId, userId string, at time.Time) error
	RemoveUserFollower(ctx context.Context, tenant, issueId, userId string, at time.Time) error
	ExistsById(ctx context.Context, tenant, issueId string) (bool, error)
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

func (r *issueRepository) Create(ctx context.Context, tenant, issueId string, evt event.IssueCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.Create")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("issueId", issueId), log.String("issueCreateEvent", fmt.Sprintf("%+v", evt)))

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
							MERGE (t)<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId}) 
							ON CREATE SET 
								i:Issue_%s,
								i:TimelineEvent,
								i:TimelineEvent_%s,
								i.createdAt=$createdAt,
								i.updatedAt=$updatedAt,
								i.source=$source,
								i.sourceOfTruth=$sourceOfTruth,
								i.appSource=$appSource,
								i.subject=$subject,
								i.description=$description,	
								i.status=$status,	
								i.priority=$priority
							ON MATCH SET 	
								i.subject = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.subject is null OR i.subject = '' THEN $subject ELSE i.subject END,
								i.description = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.description is null OR i.description = '' THEN $description ELSE i.description END,
								i.status = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.status is null OR i.status = '' THEN $status ELSE i.status END,
								i.priority = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.priority is null OR i.priority = '' THEN $priority ELSE i.priority END,
								i.updatedAt = $updatedAt,
								i.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE i.sourceOfTruth END,
								i.syncedWithEventStore = true
							WITH i, t
							OPTIONAL MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$reportedByOrganizationId}) 
							WHERE $reportedByOrganizationId <> ""
							FOREACH (ignore IN CASE WHEN o IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (i)-[:REPORTED_BY]->(o))
							WITH i, t
							OPTIONAL MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$submittedByOrganizationId}) 
							WHERE $submittedByOrganizationId <> ""
							FOREACH (ignore IN CASE WHEN o IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (i)-[:SUBMITTED_BY]->(o))
							WITH i, t
							OPTIONAL MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$submittedByUserId}) 
							WHERE $submittedByUserId <> ""
							FOREACH (ignore IN CASE WHEN u IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (i)-[:SUBMITTED_BY]->(u))
							`, tenant, tenant)
	params := map[string]any{
		"tenant":                    tenant,
		"issueId":                   issueId,
		"createdAt":                 evt.CreatedAt,
		"updatedAt":                 evt.UpdatedAt,
		"source":                    helper.GetSource(evt.Source),
		"sourceOfTruth":             helper.GetSourceOfTruth(evt.Source),
		"appSource":                 helper.GetAppSource(evt.AppSource),
		"subject":                   evt.Subject,
		"description":               evt.Description,
		"status":                    evt.Status,
		"priority":                  evt.Priority,
		"reportedByOrganizationId":  evt.ReportedByOrganizationId,
		"submittedByOrganizationId": evt.SubmittedByOrganizationId,
		"submittedByUserId":         evt.SubmittedByUserId,
		"overwrite":                 helper.GetSourceOfTruth(evt.Source) == constants.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return r.executeWriteQuery(ctx, cypher, params)
}

func (r *issueRepository) Update(ctx context.Context, tenant, issueId string, evt event.IssueUpdateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("issueId", issueId), log.Object("event", evt))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId})
		 	SET	
				i.subject = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.subject is null OR i.subject = '' THEN $subject ELSE i.subject END,
				i.description = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.description is null OR i.description = '' THEN $description ELSE i.description END,
				i.status = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.status is null OR i.status = '' THEN $status ELSE i.status END,
				i.priority = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.priority is null OR i.priority = '' THEN $priority ELSE i.priority END,
				i.updatedAt = $updatedAt,
				i.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE i.sourceOfTruth END,
				i.syncedWithEventStore = true`
	params := map[string]any{
		"tenant":        tenant,
		"issueId":       issueId,
		"updatedAt":     evt.UpdatedAt,
		"subject":       evt.Subject,
		"description":   evt.Description,
		"status":        evt.Status,
		"priority":      evt.Priority,
		"sourceOfTruth": helper.GetSourceOfTruth(evt.Source),
		"overwrite":     helper.GetSourceOfTruth(evt.Source) == constants.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return r.executeWriteQuery(ctx, cypher, params)
}

func (r *issueRepository) AddUserAssignee(ctx context.Context, tenant, issueId, userId string, at time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.AddUserAssignee")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("issueId", issueId), log.String("userId", userId), log.Object("at", at))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId}),
				(t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
		 	MERGE (i)-[:ASSIGNED_TO]->(u)
				ON CREATE SET i.updatedAt = $updatedAt`
	params := map[string]any{
		"tenant":    tenant,
		"issueId":   issueId,
		"updatedAt": at,
		"userId":    userId,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return r.executeWriteQuery(ctx, cypher, params)
}

func (r *issueRepository) AddUserFollower(ctx context.Context, tenant, issueId, userId string, at time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.AddUserFollower")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("issueId", issueId), log.String("userId", userId), log.Object("at", at))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId}),
				(t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
		 	MERGE (i)-[:FOLLOWED_BY]->(u)
				ON CREATE SET i.updatedAt = $updatedAt`
	params := map[string]any{
		"tenant":    tenant,
		"issueId":   issueId,
		"updatedAt": at,
		"userId":    userId,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return r.executeWriteQuery(ctx, cypher, params)
}

func (r *issueRepository) RemoveUserAssignee(ctx context.Context, tenant, issueId, userId string, at time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.RemoveUserAssignee")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("issueId", issueId), log.String("userId", userId), log.Object("at", at))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId}),
				(t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId}),
				(i)-[r:ASSIGNED_TO]->(u)
				SET i.updatedAt = $updatedAt
		 		DELETE r`
	params := map[string]any{
		"tenant":    tenant,
		"issueId":   issueId,
		"updatedAt": at,
		"userId":    userId,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return r.executeWriteQuery(ctx, cypher, params)
}

func (r *issueRepository) RemoveUserFollower(ctx context.Context, tenant, issueId, userId string, at time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.RemoveUserFollower")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("issueId", issueId), log.String("userId", userId), log.Object("at", at))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId}),
				(t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId}),
				(i)-[r:FOLLOWED_BY]->(u)
				SET i.updatedAt = $updatedAt
		 		DELETE r`
	params := map[string]any{
		"tenant":    tenant,
		"issueId":   issueId,
		"updatedAt": at,
		"userId":    userId,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	return r.executeWriteQuery(ctx, cypher, params)
}

func (r *issueRepository) ExistsById(ctx context.Context, tenant, issueId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.ExistsById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("issueId", issueId))

	cypher := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId}) RETURN i LIMIT 1`
	params := map[string]any{
		"tenant":  tenant,
		"issueId": issueId,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return false, err
		} else {
			return queryResult.Next(ctx), nil
		}
	})
	if err != nil {
		return false, err
	}
	return result.(bool), err
}

func (r *issueRepository) executeWriteQuery(ctx context.Context, query string, params map[string]any) error {
	return utils.ExecuteWriteQueryOnDb(ctx, *r.driver, r.database, query, params)
}
