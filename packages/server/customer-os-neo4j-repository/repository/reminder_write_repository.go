package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
)

type ReminderUpdateFields struct {
	Content         *string
	DueDate         *time.Time
	Dismissed       *bool
	UpdateContent   bool
	UpdateDueDate   bool
	UpdateDismissed bool
}

type ReminderWriteRepository interface {
	CreateReminder(ctx context.Context, tenant, id, userId, organizationId, content, source, appSource string, createdAt, dueDate time.Time) error
	UpdateReminder(ctx context.Context, tenant, id string, data ReminderUpdateFields) error
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

func (r *reminderWriteRepository) CreateReminder(ctx context.Context, tenant, id, userId, organizationId, content, source, appSource string, createdAt, dueDate time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderWriteRepository.CreateReminder")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, id)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})
				MERGE (t)<-[:REMINDER_BELONGS_TO_TENANT]-(r:Reminder {id:$id})
				ON CREATE SET  
					r:Reminder_%s,
					r.createdAt=$createdAt,
					r.updatedAt=datetime(),	
					r.source=$source,
					r.sourceOfTruth=$source,
					r.appSource=$appSource,
					r.content=$content,	
					r.dueDate=$dueDate,
					r.dismissed=$dismissed
					
				WITH t, r	
			
				MATCH (t)<-[:USER_BELONGS_TO_TENANT]-(u:User {id:$userId})
				MERGE (u)<-[:REMINDER_BELONGS_TO_USER]-(r)

				WITH t, r
				MATCH (t)<-[:ORGANIZATION_BELONGS_TO_TENANT]-(o:Organization {id:$organizationId})
				MERGE (o)<-[:REMINDER_BELONGS_TO_ORGANIZATION]-(r),
				SET o.updatedAt=datetime()`, tenant)
	params := map[string]interface{}{
		"tenant":         tenant,
		"id":             id,
		"userId":         userId,
		"organizationId": organizationId,
		"content":        content,
		"source":         source,
		"appSource":      appSource,
		"createdAt":      createdAt,
		"dueDate":        dueDate,
		"dismissed":      false,
	}
	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *reminderWriteRepository) UpdateReminder(ctx context.Context, tenant, id string, data ReminderUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderWriteRepository.UpdateReminder")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, id)

	if data.Content == nil && data.DueDate == nil && data.Dismissed == nil {
		return nil
	}

	cypher := `MATCH (:Tenant {name:$tenant})<-[:REMINDER_BELONGS_TO_TENANT]-(r:Reminder {id:$id})
				SET r.updatedAt = datetime()`
	params := map[string]interface{}{
		"tenant": tenant,
		"id":     id,
	}
	if data.UpdateContent {
		cypher += ", r.content = $content"
		params["content"] = data.Content
	}
	if data.UpdateDueDate {
		cypher += ", r.dueDate = datetime($dueDate)"
		params["dueDate"] = *data.DueDate
	}
	if data.UpdateDismissed {
		cypher += ", r.dismissed = $dismissed"
		params["dismissed"] = *data.Dismissed
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
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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
