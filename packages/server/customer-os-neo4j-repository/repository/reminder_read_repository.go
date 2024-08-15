package repository

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type ReminderReadRepository interface {
	GetReminderById(ctx context.Context, tenant, id string) (*dbtype.Node, error)
	GetRemindersOrderByDueDateAsc(ctx context.Context, tenant, organizationId string, dismissed *bool) ([]*dbtype.Node, error)
	GetRemindersOrderByDueDateDesc(ctx context.Context, tenant, organizationId string, dismissed *bool) ([]*dbtype.Node, error)
	GetRemindersForOrganization(ctx context.Context, tenant, organizationId string, dismissed *bool) ([]*dbtype.Node, error)
}

type reminderReadRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func (r *reminderReadRepository) prepareReadSession(ctx context.Context) neo4j.SessionWithContext {
	return utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
}

func NewReminderReadRepository(driver *neo4j.DriverWithContext, database string) ReminderReadRepository {
	return &reminderReadRepository{
		driver:   driver,
		database: database,
	}
}

func (r *reminderReadRepository) GetReminderById(ctx context.Context, tenant, id string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderReadRepository.GetReminderById")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, id)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:REMINDER_BELONGS_TO_TENANT]-(r:Reminder {id:$id}) RETURN r`
	params := map[string]any{
		"tenant": tenant,
		"id":     id,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)

	})
	if err != nil {
		span.LogFields(log.Bool("result.found", false))
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Bool("result.found", result != nil))
	return result.(*dbtype.Node), nil
}

// GetReminderOrderByDueDateAsc returns a list of reminders ordered by due date in ascending order
func (r *reminderReadRepository) GetRemindersOrderByDueDateAsc(ctx context.Context, tenant, organizationId string, dismissed *bool) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderReadRepository.GetRemindersOrderByDueDateAsc")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:REMINDER_BELONGS_TO_TENANT]-(r:Reminder)-[:REMINDER_BELONGS_TO_ORGANIZATION]->(o:Organization {id:$organizationId})`
	if dismissed != nil {
		if *dismissed {
			cypher += ` WHERE r.dismissed = true`
		} else {
			cypher += ` WHERE r.dismissed IS NULL OR r.dismissed = false`
		}
	}
	cypher += ` RETURN r ORDER BY r.dueDate ASC`

	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*dbtype.Node))))
	return result.([]*dbtype.Node), nil
}

// GetReminderOrderByDueDateDesc returns a list of reminders ordered by due date in descending order
func (r *reminderReadRepository) GetRemindersOrderByDueDateDesc(ctx context.Context, tenant, organizationId string, dismissed *bool) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderReadRepository.GetRemindersOrderByDueDateDesc")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:REMINDER_BELONGS_TO_TENANT]-(r:Reminder)-[:REMINDER_BELONGS_TO_ORGANIZATION]->(o:Organization {id:$organizationId})`
	if dismissed != nil {
		if *dismissed {
			cypher += ` WHERE r.dismissed = true`
		} else {
			cypher += ` WHERE r.dismissed IS NULL OR r.dismissed = false`
		}
	}
	cypher += ` RETURN r ORDER BY r.dueDate DESC`

	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*dbtype.Node))))
	return result.([]*dbtype.Node), nil
}

// GetRemindersForOrganization returns a list of reminders for an organization
func (r *reminderReadRepository) GetRemindersForOrganization(ctx context.Context, tenant, organizationId string, dismissed *bool) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ReminderReadRepository.GetRemindersForOrganization")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)

	cypher := `MATCH (:Tenant {name:$tenant})<-[:REMINDER_BELONGS_TO_TENANT]-(r:Reminder)-[:REMINDER_BELONGS_TO_ORGANIZATION]->(o:Organization {id:$organizationId})`
	if dismissed != nil {
		if *dismissed {
			cypher += ` WHERE r.dismissed = true`
		} else {
			cypher += ` WHERE r.dismissed IS NULL OR r.dismissed = false`
		}
	}
	cypher += ` RETURN r`

	params := map[string]any{
		"tenant":         tenant,
		"organizationId": organizationId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := r.prepareReadSession(ctx)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*dbtype.Node))))
	return result.([]*dbtype.Node), nil
}
