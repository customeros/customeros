package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/log_entry/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LogEntryRepository interface {
	Create(ctx context.Context, tenant, logEntryId string, evt event.LogEntryCreateEvent) error
	Update(ctx context.Context, tenant, logEntryId string, evt event.LogEntryUpdateEvent) error
}

type logEntryRepository struct {
	driver *neo4j.DriverWithContext
}

func NewLogEntryRepository(driver *neo4j.DriverWithContext) LogEntryRepository {
	return &logEntryRepository{
		driver: driver,
	}
}

func (r *logEntryRepository) Create(ctx context.Context, tenant, logEntryId string, evt event.LogEntryCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryRepository.Create")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("logEntryId", logEntryId), log.Object("event", evt))

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$orgId})
							MERGE (l:LogEntry {id:$logEntryId})<-[:LOGGED]-(o)
							ON CREATE SET 
								l:LogEntry_%s,
								l:TimelineEvent,
								l:TimelineEvent_%s,
								l.createdAt=$createdAt,
								l.updatedAt=$updatedAt,
								l.startedAt=$startedAt,
								l.source=$source,
								l.sourceOfTruth=$sourceOfTruth,
								l.appSource=$appSource,
								l.content=$content,
								l.contentType=$contentType
							WITH l, t
							OPTIONAL MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$authorUserId}) 
							WHERE $authorUserId <> ""
							FOREACH (ignore IN CASE WHEN u IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (l)-[:CREATED_BY]->(u))
							`, tenant, tenant)
	span.LogFields(log.String("query", query))

	return utils.ExecuteWriteQuery(ctx, *r.driver, query, map[string]any{
		"tenant":        tenant,
		"logEntryId":    logEntryId,
		"orgId":         evt.LoggedOrganizationId,
		"createdAt":     evt.CreatedAt,
		"updatedAt":     evt.UpdatedAt,
		"startedAt":     evt.StartedAt,
		"source":        helper.GetSource(evt.Source),
		"sourceOfTruth": helper.GetSourceOfTruth(evt.SourceOfTruth),
		"appSource":     helper.GetAppSource(evt.AppSource),
		"content":       evt.Content,
		"contentType":   evt.ContentType,
		"authorUserId":  evt.AuthorUserId,
	})
}

func (r *logEntryRepository) Update(ctx context.Context, tenant, logEntryId string, evt event.LogEntryUpdateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("logEntryId", logEntryId), log.Object("event", evt))

	query := fmt.Sprintf(`MATCH (l:LogEntry_%s {id:$logEntryId})
								SET 
								l.updatedAt=$updatedAt,
								l.startedAt=$startedAt,
								l.sourceOfTruth=$sourceOfTruth,
								l.content=$content,
								l.contentType=$contentType
								WITH l
							OPTIONAL MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(org:Organization {id:$orgId}) 
							WHERE $orgId <> ""
							FOREACH (ignore IN CASE WHEN org IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (l)<-[:LOGGED]-(org))`, tenant)
	params := map[string]any{
		"tenant":        tenant,
		"logEntryId":    logEntryId,
		"updatedAt":     evt.UpdatedAt,
		"startedAt":     evt.StartedAt,
		"sourceOfTruth": helper.GetSourceOfTruth(evt.SourceOfTruth),
		"content":       evt.Content,
		"contentType":   evt.ContentType,
		"orgId":         evt.LoggedOrganizationId,
	}
	span.LogFields(log.String("query", query), log.Object("params", params))

	return utils.ExecuteWriteQuery(ctx, *r.driver, query, params)
}
