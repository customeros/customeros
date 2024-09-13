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

type FlowActionWriteRepository interface {
	Merge(ctx context.Context, entity *entity.FlowActionEntity) (*dbtype.Node, error)
}

type flowActionWriteRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowActionWriteRepository(driver *neo4j.DriverWithContext, database string) FlowActionWriteRepository {
	return &flowActionWriteRepositoryImpl{driver: driver, database: database}
}

func (r *flowActionWriteRepositoryImpl) Merge(ctx context.Context, input *entity.FlowActionEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionWriteRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	onCreate := `ON CREATE SET
				fa.createdAt = $createdAt,
				fa.updatedAt = $updatedAt,
				fa.index = $index,
				fa.name = $name,
				fa.status = $status,
				fa.actionType = $actionType,`
	onMatch := `ON MATCH SET
				fa.updatedAt = $updatedAt,
				fa.index = $index,
				fa.name = $name,
				fa.status = $status,
				fa.actionType = $actionType,`

	params := map[string]any{
		"tenant":     common.GetTenantFromContext(ctx),
		"id":         input.Id,
		"createdAt":  utils.TimeOrNow(input.CreatedAt),
		"updatedAt":  utils.TimeOrNow(input.UpdatedAt),
		"index":      input.Index,
		"name":       input.Name,
		"status":     input.Status,
		"actionType": input.ActionType,
	}

	if input.ActionType == entity.FlowActionTypeWait {
		onCreate += `
				fa.actionData_minutes = $actionData_waitTime`
		onMatch += `
				fa.actionData_minutes = $actionData_waitTime`
		params["actionData_waitTime"] = input.ActionData.Wait.Minutes
	}
	if input.ActionType == entity.FlowActionTypeEmailNew {
		onCreate += `
				fa.actionData_subject = $actionData_subject,
				fa.actionData_bodyTemplate = $actionData_bodyTemplate`
		onMatch += `
				fa.actionData_subject = $actionData_subject,
				fa.actionData_bodyTemplate = $actionData_bodyTemplate`
		params["actionData_subject"] = input.ActionData.EmailNew.Subject
		params["actionData_bodyTemplate"] = input.ActionData.EmailNew.BodyTemplate
	}
	if input.ActionType == entity.FlowActionTypeEmailReply {
		onCreate += `
				fa.actionData_replyToId = $actionData_replyToId,
				fa.actionData_subject = $actionData_subject,
				fa.actionData_bodyTemplate = $actionData_bodyTemplate`
		onMatch += `
				fa.actionData_replyToId = $actionData_replyToId,
				fa.actionData_subject = $actionData_subject,
				fa.actionData_bodyTemplate = $actionData_bodyTemplate`
		params["actionData_replyToId"] = input.ActionData.EmailReply.ReplyToId
		params["actionData_subject"] = input.ActionData.EmailReply.Subject
		params["actionData_bodyTemplate"] = input.ActionData.EmailReply.BodyTemplate
	}
	if input.ActionType == entity.FlowActionTypeLinkedinConnectionRequest {
		onCreate += `
				fa.actionData_messageTemplate = $actionData_messageTemplate`
		onMatch += `
				fa.actionData_messageTemplate = $actionData_messageTemplate`
		params["actionData_messageTemplate"] = input.ActionData.LinkedinConnectionRequest.MessageTemplate
	}
	if input.ActionType == entity.FlowActionTypeLinkedinMessage {
		onCreate += `
				fa.actionData_messageTemplate = $actionData_messageTemplate`
		onMatch += `
				fa.actionData_messageTemplate = $actionData_messageTemplate`
		params["actionData_messageTemplate"] = input.ActionData.LinkedinMessage.MessageTemplate
	}

	cypher := fmt.Sprintf(`
			MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:BELONGS_TO_TENANT]-(fa:FlowAction:FlowAction_%s {id: $id})
			%s
			%s
			RETURN fa`, common.GetTenantFromContext(ctx), onCreate, onMatch)

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
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
