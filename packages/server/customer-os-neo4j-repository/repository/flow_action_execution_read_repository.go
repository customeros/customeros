package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type FlowActionExecutionReadRepository interface {
	GetById(ctx context.Context, id string) (*dbtype.Node, error)
	GetExecution(ctx context.Context, flowId, actionId, entityId string, entityType model.EntityType) (*dbtype.Node, error)
	GetForEntity(ctx context.Context, tx *neo4j.ManagedTransaction, flowId, entityId string, entityType model.EntityType) ([]*dbtype.Node, error)
	GetFirstSlotForMailbox(ctx context.Context, tx *neo4j.ManagedTransaction, mailbox string) (*time.Time, error)
	GetLastScheduledForMailbox(ctx context.Context, tx *neo4j.ManagedTransaction, mailbox string) (*dbtype.Node, error)
	GetScheduledBefore(ctx context.Context, before time.Time) ([]*dbtype.Node, error)
	GetByMailboxAndTimeInterval(ctx context.Context, tx *neo4j.ManagedTransaction, mailbox string, startTime, endTime time.Time) (*dbtype.Node, error)
	GetForEntityWithActionType(ctx context.Context, entityId string, entityType model.EntityType, actionType neo4jentity.FlowActionType) ([]*dbtype.Node, error)

	CountEmailsPerMailboxPerDay(ctx context.Context, tx *neo4j.ManagedTransaction, mailbox string, startDate, endDate time.Time) (int64, error)
}

type flowActionExecutionReadRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowActionExecutionReadRepository(driver *neo4j.DriverWithContext, database string) FlowActionExecutionReadRepository {
	return &flowActionExecutionReadRepositoryImpl{driver: driver, database: database}
}

func (r *flowActionExecutionReadRepositoryImpl) GetById(ctx context.Context, id string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionReadRepository.GetExecution")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(fae:FlowActionExecution_%s {id:$id}) RETURN fae`, tenant)
	params := map[string]any{
		"tenant": tenant,
		"id":     id,
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

	if err != nil && err.Error() == "Result contains no more records" {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.Bool("result.found", true))

	return result.(*dbtype.Node), nil
}

func (r *flowActionExecutionReadRepositoryImpl) GetExecution(ctx context.Context, flowId, actionId, entityId string, entityType model.EntityType) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionReadRepository.GetExecution")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s {id: $flowId})-[:HAS]->(fa:FlowAction_%s {id: $actionId})-[:HAS_EXECUTION]->(fae:FlowActionExecution_%s {entityId: $entityId, entityType: $entityType}) RETURN fae order by fae.executedAt`, tenant, tenant, tenant)
	params := map[string]any{
		"tenant":     tenant,
		"flowId":     flowId,
		"actionId":   actionId,
		"entityId":   entityId,
		"entityType": entityType.String(),
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
	return result.(*dbtype.Node), nil
}

func (r *flowActionExecutionReadRepositoryImpl) GetForEntity(ctx context.Context, tx *neo4j.ManagedTransaction, flowId, entityId string, entityType model.EntityType) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionReadRepository.GetForEntity")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s {id: $flowId})-[:HAS]->(fa:FlowAction_%s)-[:HAS_EXECUTION]->(fae:FlowActionExecution_%s {entityId: $entityId, entityType: $entityType}) RETURN fae order by fae.executedAt, fae.scheduledAt`, tenant, tenant, tenant)
	params := map[string]any{
		"tenant":     tenant,
		"flowId":     flowId,
		"entityId":   entityId,
		"entityType": entityType.String(),
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})

	if err != nil {
		return nil, nil
	}
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return result.([]*dbtype.Node), nil
}

func (r *flowActionExecutionReadRepositoryImpl) GetScheduledBefore(ctx context.Context, before time.Time) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionReadRepository.GetScheduledBefore")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (f:Flow {status: 'ACTIVE'})-[:HAS]->(:FlowAction)-[:HAS_EXECUTION]->(fae:FlowActionExecution) where fae.status = 'SCHEDULED' and fae.scheduledAt < $before RETURN fae order by fae.scheduledAt limit 100`)
	params := map[string]any{
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

func (r *flowActionExecutionReadRepositoryImpl) GetByMailboxAndTimeInterval(ctx context.Context, tx *neo4j.ManagedTransaction, mailbox string, startTime, endTime time.Time) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionReadRepository.GetByMailboxAndTimeInterval")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := `
		MATCH (f:FlowActionExecution)
		WHERE f.mailbox = $mailbox AND f.scheduledAt >= $startTime AND f.scheduledAt < $endTime
		RETURN f
		ORDER BY f.scheduledAt desc
		LIMIT 1`
	params := map[string]interface{}{
		"mailbox":   mailbox,
		"startTime": startTime,
		"endTime":   endTime,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})

	if err != nil && err.Error() == "Result contains no more records" {
		span.LogFields(log.Bool("result.found", false))
		return nil, nil
	}
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.Bool("result.found", result != nil))

	return result.(*dbtype.Node), nil
}

func (r *flowActionExecutionReadRepositoryImpl) GetForEntityWithActionType(ctx context.Context, entityId string, entityType model.EntityType, actionType neo4jentity.FlowActionType) ([]*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionReadRepository.GetForEntityWithActionType")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("entityId", entityId), log.Object("entityType", entityType), log.Object("actionType", actionType))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`
		MATCH (f:FlowActionExecution_%s)-[:HAS_EXECUTION]->(fa:FlowAction_%s)
		WHERE f.entityId = $entityId AND f.entityType = $entityType AND fa.action = $actionType and f.status in ['SCHEDULED', 'IN_PROGRESS']
		RETURN f`, tenant, tenant)
	params := map[string]interface{}{
		"entityId":   entityId,
		"entityType": entityType.String(),
		"actionType": actionType,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	result, err := utils.ExecuteReadInTransaction(ctx, r.driver, r.database, nil, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractAllRecordsFirstValueAsDbNodePtrs(ctx, queryResult, err)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	span.LogFields(log.Int("result.count", len(result.([]*dbtype.Node))))

	return result.([]*dbtype.Node), nil
}

func (r *flowActionExecutionReadRepositoryImpl) GetFirstSlotForMailbox(ctx context.Context, tx *neo4j.ManagedTransaction, mailbox string) (*time.Time, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionReadRepository.GetFirstSlotForMailbox")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.Object("mailbox", mailbox))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(fae:FlowActionExecution_%s) where fae.status = 'SCHEDULED' and fae.mailbox = $mailbox RETURN fae.scheduledAt order by fae.scheduledAt desc limit 1`, tenant)
	params := map[string]any{
		"tenant":  tenant,
		"mailbox": mailbox,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	if tx == nil {
		session := utils.NewNeo4jReadSession(ctx, *r.driver)
		defer session.Close(ctx)

		queryResult, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			qr, err := tx.Run(ctx, cypher, params)
			if err != nil {
				return nil, err
			}
			return utils.ExtractSingleRecordFirstValueAsType[time.Time](ctx, qr, err)
		})
		if err != nil && err.Error() == "Result contains no more records" {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		return queryResult.(*time.Time), nil
	} else {
		queryResult, err := (*tx).Run(ctx, cypher, params)

		if err != nil {
			return nil, err
		}
		t, err := utils.ExtractSingleRecordFirstValueAsType[time.Time](ctx, queryResult, err)
		if err != nil && err.Error() == "Result contains no more records" {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		return &t, err
	}
}

func (r *flowActionExecutionReadRepositoryImpl) GetLastScheduledForMailbox(ctx context.Context, tx *neo4j.ManagedTransaction, mailbox string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionReadRepository.GetLastScheduledForMailbox")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("mailbox", mailbox))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(fae:FlowActionExecution_%s) where fae.status = 'SCHEDULED' and fae.mailbox = $mailbox RETURN fae order by fae.scheduledAt desc limit 1`, tenant)
	params := map[string]any{
		"tenant":  tenant,
		"mailbox": mailbox,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	if tx == nil {
		session := utils.NewNeo4jReadSession(ctx, *r.driver)
		defer session.Close(ctx)

		queryResult, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			qr, err := tx.Run(ctx, cypher, params)
			if err != nil {
				return nil, err
			}
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, qr, err)
		})
		if err != nil && err.Error() == "Result contains no more records" {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		return queryResult.(*neo4j.Node), nil
	} else {
		queryResult, err := (*tx).Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}

		result, err := utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		if err != nil && err.Error() == "Result contains no more records" {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		return result, err
	}
}

func (r *flowActionExecutionReadRepositoryImpl) CountEmailsPerMailboxPerDay(ctx context.Context, tx *neo4j.ManagedTransaction, mailbox string, startDate, endDate time.Time) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionReadRepository.CountEmailsPerMailboxPerDay")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("mailbox", mailbox), log.Object("startDate", startDate), log.Object("endDate", endDate))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(fae:FlowActionExecution_%s) where fae.scheduledAt >= $startDate and fae.scheduledAt <= $endDate and fae.mailbox = $mailbox RETURN count(fae)`, tenant)
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

	if tx == nil {
		session := utils.NewNeo4jReadSession(ctx, *r.driver)
		defer session.Close(ctx)

		queryResult, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			queryResult, err := tx.Run(ctx, cypher, params)
			if err != nil {
				return nil, err
			}
			return queryResult.Single(ctx)
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return 0, err
		}

		count := queryResult.(*db.Record).Values[0].(int64)
		span.LogFields(log.Int64("result", count))
		return count, nil
	} else {
		queryResult, err := (*tx).Run(ctx, cypher, params)
		if err != nil {
			tracing.TraceErr(span, err)
			return 0, err
		}
		if err != nil {
			tracing.TraceErr(span, err)
			return 0, err
		}
		single, err := queryResult.Single(ctx)
		if err != nil {
			tracing.TraceErr(span, err)
			return 0, err
		}
		count := single.Values[0].(int64)
		span.LogFields(log.Int64("result", count))
		return count, nil
	}
}
