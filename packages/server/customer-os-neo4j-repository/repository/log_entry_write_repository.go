package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type LogEntryCreateFields struct {
	Content              string       `json:"content"`
	ContentType          string       `json:"contentType"`
	StartedAt            time.Time    `json:"startedAt" `
	AuthorUserId         string       `json:"authorUserId"`
	LoggedOrganizationId string       `json:"loggedOrganizationId"`
	SourceFields         model.Source `json:"sourceFields"`
	CreatedAt            time.Time    `json:"createdAt"`
	UpdatedAt            time.Time    `json:"updatedAt"`
}

type LogEntryUpdateFields struct {
	Content              string    `json:"content"`
	ContentType          string    `json:"contentType"`
	StartedAt            time.Time `json:"startedAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	Source               string    `json:"source"`
	LoggedOrganizationId string    `json:"loggedOrganizationId"`
}

type LogEntryWriteRepository interface {
	Create(ctx context.Context, tenant, logEntryId string, data LogEntryCreateFields) error
	Update(ctx context.Context, tenant, logEntryId string, data LogEntryUpdateFields) error
}

type logEntryWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewLogEntryWriteRepository(driver *neo4j.DriverWithContext, database string) LogEntryWriteRepository {
	return &logEntryWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *logEntryWriteRepository) Create(ctx context.Context, tenant, logEntryId string, data LogEntryCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryWriteRepository.Create")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, logEntryId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$orgId})
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
	params := map[string]any{
		"tenant":        tenant,
		"logEntryId":    logEntryId,
		"orgId":         data.LoggedOrganizationId,
		"createdAt":     data.CreatedAt,
		"updatedAt":     data.UpdatedAt,
		"startedAt":     data.StartedAt,
		"source":        data.SourceFields.Source,
		"sourceOfTruth": data.SourceFields.SourceOfTruth,
		"appSource":     data.SourceFields.AppSource,
		"content":       data.Content,
		"contentType":   data.ContentType,
		"authorUserId":  data.AuthorUserId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *logEntryWriteRepository) Update(ctx context.Context, tenant, logEntryId string, data LogEntryUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryWriteRepository.Create")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, logEntryId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (l:LogEntry_%s {id:$logEntryId})
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
		"updatedAt":     data.UpdatedAt,
		"startedAt":     data.StartedAt,
		"sourceOfTruth": data.Source,
		"content":       data.Content,
		"contentType":   data.ContentType,
		"orgId":         data.LoggedOrganizationId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
