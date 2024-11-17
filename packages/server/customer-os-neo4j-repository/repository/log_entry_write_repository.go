package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/data_fields"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type LogEntryWriteRepository interface {
	CreateInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, logEntryId string, data data_fields.LogEntryFields) error
	UpdateInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, logEntryId string, data data_fields.LogEntryFields) error
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

func (r *logEntryWriteRepository) CreateInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, logEntryId string, data data_fields.LogEntryFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryWriteRepository.CreateInTx")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, logEntryId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$orgId})
							MERGE (l:LogEntry {id:$logEntryId})<-[:LOGGED]-(o)
							ON CREATE SET 
								l:LogEntry_%s,
								l:TimelineEvent,
								l:TimelineEvent_%s,
								l.createdAt=$createdAt,
								l.updatedAt=datetime(),
								l.startedAt=$startedAt,
								l.source=$source,
								l.appSource=$appSource,
								l.content=$content,
								l.contentType=$contentType,
								o.updatedAt=datetime()
							WITH l, t
							OPTIONAL MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$authorUserId}) 
							WHERE $authorUserId <> ""
							FOREACH (ignore IN CASE WHEN u IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (l)-[:CREATED_BY]->(u))
							`, tenant, tenant)
	params := map[string]any{
		"tenant":       tenant,
		"logEntryId":   logEntryId,
		"orgId":        utils.IfNotNilString(data.OrganizationId),
		"createdAt":    utils.IfNotNilTimeWithDefault(data.StartedAt, utils.Now()),
		"startedAt":    utils.IfNotNilTimeWithDefault(data.StartedAt, utils.Now()),
		"source":       utils.IfNotNilString(data.Source),
		"appSource":    utils.IfNotNilString(data.AppSource),
		"content":      utils.IfNotNilString(data.Content),
		"contentType":  utils.IfNotNilString(data.ContentType),
		"authorUserId": utils.IfNotNilString(data.AuthorUserId),
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

func (r *logEntryWriteRepository) UpdateInTx(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, logEntryId string, data data_fields.LogEntryFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "LogEntryWriteRepository.UpdateInTx")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, logEntryId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (l:LogEntry_%s {id:$logEntryId}) SET l.updatedAt=datetime()`, tenant)
	params := map[string]any{
		"tenant":      tenant,
		"logEntryId":  logEntryId,
		"startedAt":   data.StartedAt,
		"content":     data.Content,
		"contentType": data.ContentType,
	}

	if data.Content != nil {
		params["content"] = *data.Content
		cypher += ", l.content=$content"
	}
	if data.ContentType != nil {
		params["contentType"] = *data.ContentType
		cypher += ", l.contentType=$contentType"
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
