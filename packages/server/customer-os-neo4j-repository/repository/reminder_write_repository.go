package repository

import (
	"context"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
)

type ReminderWriteRepository interface {
	CreateReminder(ctx context.Context, tenant, id, userId, orgId, content, source, appSource string, createdAt, dueDate time.Time) error
	UpdateReminder(ctx context.Context, tenant, id string, content *string, dueDate *time.Time, dismissed *bool) error
	DeleteReminder(ctx context.Context, tenant, id string) error
}

type reminderWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewReminderWriteRepository(driver *neo4j.DriverWithContext, database string) ReminderWriteRepository {
	return &reminderWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *reminderWriteRepository) CreateReminder(ctx context.Context, tenant, id, userId, orgId, content, source, appSource string, createdAt, dueDate time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderWriteRepository.CreateReminder")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, id)

	cypher := `MATCH (t:Tenant {name:$tenant})
				MERGE (t)<-[:REMINDER_BELONGS_TO_TENANT]-(r:Reminder {id:$id})
				SET r += {
					createdAt: datetime($createdAt),
					updatedAt: datetime($createdAt),
					source: $source,
					sourceOfTruth: $source,
					appSource: $appSource,
					content: $content,
					dueDate: datetime($dueDate),
					dismissed: $dismissed
				}
				MERGE (u:User {id:$userId})-[:REMINDER_BELONGS_TO_USER]->(r)
				MERGE (o:Organization {id:$orgId})-[:REMINDER_BELONGS_TO_ORGANIZATION]->(r)`
	params := map[string]interface{}{
		"tenant":    tenant,
		"id":        id,
		"userId":    userId,
		"orgId":     orgId,
		"content":   content,
		"source":    source,
		"appSource": appSource,
		"createdAt": createdAt,
		"dueDate":   dueDate,
		"dismissed": false,
	}
	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *reminderWriteRepository) UpdateReminder(ctx context.Context, tenant, id string, content *string, dueDate *time.Time, dismissed *bool) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderWriteRepository.UpdateReminder")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, id)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:REMINDER_BELONGS_TO_TENANT]-(r:Reminder {id:$id})
				SET r.updatedAt = datetime($updatedAt)`
	params := map[string]interface{}{
		"tenant":    tenant,
		"id":        id,
		"updatedAt": time.Now(),
	}
	if content != nil {
		cypher += ", r.content = $content"
		params["content"] = *content
	}
	if dueDate != nil {
		cypher += ", r.dueDate = datetime($dueDate)"
		params["dueDate"] = *dueDate
	}
	if dismissed != nil {
		cypher += ", r.dismissed = $dismissed"
		params["dismissed"] = *dismissed
	}

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *reminderWriteRepository) DeleteReminder(ctx context.Context, tenant, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderWriteRepository.DeleteReminder")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, id)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:REMINDER_BELONGS_TO_TENANT]-(r:Reminder {id:$id})
				DETACH DELETE r`
	params := map[string]interface{}{
		"tenant": tenant,
		"id":     id,
	}

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
