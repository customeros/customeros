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
)

type IssueRepository interface {
	Create(ctx context.Context, tenant, issueId string, evt event.IssueCreateEvent) error
	Update(ctx context.Context, tenant, logEntryId string, eventData event.IssueUpdateEvent) error
}

type issueRepository struct {
	driver *neo4j.DriverWithContext
}

func NewIssueRepository(driver *neo4j.DriverWithContext) IssueRepository {
	return &issueRepository{
		driver: driver,
	}
}

func (r *issueRepository) Create(ctx context.Context, tenant, issueId string, evt event.IssueCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.Create")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("issueId", issueId), log.String("issueCreateEvent", fmt.Sprintf("%+v", evt)))

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
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
							`, tenant, tenant)
	params := map[string]any{
		"tenant":                   tenant,
		"issueId":                  issueId,
		"createdAt":                evt.CreatedAt,
		"updatedAt":                evt.UpdatedAt,
		"source":                   helper.GetSource(evt.Source),
		"sourceOfTruth":            helper.GetSourceOfTruth(evt.Source),
		"appSource":                helper.GetAppSource(evt.AppSource),
		"subject":                  evt.Subject,
		"description":              evt.Description,
		"status":                   evt.Status,
		"priority":                 evt.Priority,
		"reportedByOrganizationId": evt.ReportedByOrganizationId,
		"overwrite":                helper.GetSourceOfTruth(evt.Source) == constants.SourceOpenline,
	}
	span.LogFields(log.String("query", query), log.Object("params", params))

	return utils.ExecuteWriteQuery(ctx, *r.driver, query, params)
}

func (r *issueRepository) Update(ctx context.Context, tenant, issueId string, evt event.IssueUpdateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("issueId", issueId), log.Object("event", fmt.Sprintf("%+v", evt)))

	query := `MATCH (t:Tenant {name:$tenant})<-[:ISSUE_BELONGS_TO_TENANT]-(i:Issue {id:$issueId})
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
	span.LogFields(log.String("query", query), log.Object("params", params))

	return utils.ExecuteWriteQuery(ctx, *r.driver, query, params)
}
