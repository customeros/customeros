package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type FlowActionExecutionReadRepository interface {
	GetByContact(ctx context.Context, flowId, contactId string) ([]*dbtype.Node, error)
	GetFastestMailboxAvailable(ctx context.Context, availableMailboxes []string) (*string, error)
	GetLastScheduledForMailbox(ctx context.Context, mailbox string) (*dbtype.Node, error)
	GetScheduledBefore(ctx context.Context, before time.Time) ([]*dbtype.Node, error)

	CountEmailsPerMailboxPerDay(ctx context.Context, mailbox string, startDate, endDate time.Time) (int64, error)
}

type flowActionExecutionReadRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowActionExecutionReadRepository(driver *neo4j.DriverWithContext, database string) FlowActionExecutionReadRepository {
	return &flowActionExecutionReadRepositoryImpl{driver: driver, database: database}
}

func (r flowActionExecutionReadRepositoryImpl) GetByContact(ctx context.Context, flowId, contactId string) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionReadRepository.GetByContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s {id: $flowId})-[:HAS]->(fae:FlowActionExecution_%s) where fae.contactId = $contactId RETURN f order by fae.executedAt`, tenant, tenant)
	params := map[string]any{
		"tenant":    tenant,
		"flowId":    flowId,
		"contactId": contactId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*dbtype.Node), nil
}

func (r flowActionExecutionReadRepositoryImpl) GetScheduledBefore(ctx context.Context, before time.Time) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionReadRepository.GetScheduledBefore")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	//TODO check that FlowAction is active
	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-[:HAS]->(f:Flow_%s {status: 'ACTIVE'})-[:HAS]->(fae:FlowActionExecution_%s) where fae.status = 'SCHEDULED' and fae.scheduledAt < $before RETURN fae limit 100 order by fae.scheduledAt`, tenant)
	params := map[string]any{
		"tenant": tenant,
		"before": before,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*dbtype.Node), nil
}

func (r flowActionExecutionReadRepositoryImpl) GetFastestMailboxAvailable(ctx context.Context, availableMailboxes []string) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionReadRepository.GetFastestMailboxAvailable")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.Object("availableMailboxes", availableMailboxes))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-[:HAS]->(fae:FlowActionExecution_%s) where fae.status = 'SCHEDULED' and fae.mailbox in $availableMailboxes RETURN fae.mailbox limit 1 order by fae.scheduledAt desc`, tenant)
	params := map[string]any{
		"tenant":             tenant,
		"availableMailboxes": availableMailboxes,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsString(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.(*string), nil
}

func (r flowActionExecutionReadRepositoryImpl) GetLastScheduledForMailbox(ctx context.Context, mailbox string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionReadRepository.GetLastScheduledForMailbox")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("mailbox", mailbox))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-[:HAS]->(fae:FlowActionExecution_%s) where fae.status = 'SCHEDULED' and fae.mailbox = $mailbox RETURN fae limit 1 order by fae.scheduledAt desc`, tenant)
	params := map[string]any{
		"tenant":  tenant,
		"mailbox": mailbox,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.(*neo4j.Node), nil
}

func (r flowActionExecutionReadRepositoryImpl) CountEmailsPerMailboxPerDay(ctx context.Context, mailbox string, startDate, endDate time.Time) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionReadRepository.CountEmailsPerMailboxPerDay")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("mailbox", mailbox), log.Object("startDate", startDate), log.Object("endDate", endDate))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-[:HAS]->(fae:FlowActionExecution_%s) where fae.scheduledAt >= $startDate and fae.scheduledAt <= $endDate and fae.mailbox = $mailbox RETURN count(fae)`, tenant)
	params := map[string]any{
		"tenant":    tenant,
		"mailbox":   mailbox,
		"startDate": startDate,
		"endDate":   endDate,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return queryResult.Single(ctx)
		}
	})
	if err != nil {
		return 0, err
	}
	count := result.(*db.Record).Values[0].(int64)
	span.LogFields(log.Int64("result", count))
	return count, nil
}