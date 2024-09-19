package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type FlowActionExecutionWriteRepository interface {
	Merge(ctx context.Context, tx *neo4j.ManagedTransaction, entity *entity.FlowActionExecutionEntity) (*dbtype.Node, error)
	Delete(ctx context.Context, id string) error
}

type flowActionExecutionWriteRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowActionExecutionWriteRepository(driver *neo4j.DriverWithContext, database string) FlowActionExecutionWriteRepository {
	return &flowActionExecutionWriteRepositoryImpl{driver: driver, database: database}
}

func (r *flowActionExecutionWriteRepositoryImpl) Merge(ctx context.Context, tx *neo4j.ManagedTransaction, entity *entity.FlowActionExecutionEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionExecutionWriteRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`
			MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(:Flow_%s {id: $flowId})
			MERGE (t)<-[:BELONGS_TO_TENANT]-(fae:FlowActionExecution:FlowActionExecution_%s {id: $id})<-[:HAS]-(f)
			ON MATCH SET
				fae.updatedAt = $updatedAt,
				fae.error = $error,
				fae.subject = $subject,
				fae.body = $body,
				fae.from = $from,
				fae.to = $to,
				fae.cc = $cc,
				fae.bcc = $bcc
			ON CREATE SET
				fae.createdAt = $createdAt,
				fae.updatedAt = $updatedAt,
				fae.flowId = $flowId,
				fae.contactId = $contactId,
				fae.actionId = $actionId,
				fae.scheduledAt = $scheduledAt,
				fae.status = $status,
				fae.mailbox = $mailbox
			RETURN fae`, tenant, tenant)

	params := map[string]any{
		"tenant":      tenant,
		"id":          entity.Id,
		"createdAt":   utils.TimeOrNow(entity.CreatedAt),
		"updatedAt":   utils.TimeOrNow(entity.UpdatedAt),
		"flowId":      entity.FlowId,
		"contactId":   entity.ContactId,
		"actionId":    entity.ActionId,
		"scheduledAt": entity.ScheduledAt,
		"executedAt":  entity.ExecutedAt,
		"status":      entity.Status,
		"error":       entity.Error,

		//config
		"mailbox": entity.Mailbox,

		//data
		"subject": entity.Subject,
		"body":    entity.Body,
		"from":    entity.From,
		"to":      entity.To,
		"cc":      entity.Cc,
		"bcc":     entity.Bcc,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	if tx == nil {
		session := utils.NewNeo4jWriteSession(ctx, *r.driver)
		defer session.Close(ctx)

		queryResult, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			qr, err := tx.Run(ctx, cypher, params)
			if err != nil {
				return nil, err
			}
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, qr, err)
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		return queryResult.(*neo4j.Node), nil
	} else {
		queryResult, err := (*tx).Run(ctx, cypher, params)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}
}

func (r *flowActionExecutionWriteRepositoryImpl) Delete(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonWriteRepository.Delete")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})<-[r:BELONGS_TO_TENANT]-(fc:FlowContact_%s {id:$id}) delete r, fc`, tenant)

	params := map[string]any{
		"tenant": tenant,
		"id":     id,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}
