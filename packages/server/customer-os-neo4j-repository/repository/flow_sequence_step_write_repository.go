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

type FlowSequenceStepWriteRepository interface {
	Merge(ctx context.Context, entity *entity.FlowSequenceStepEntity) (*dbtype.Node, error)
}

type flowSequenceStepWriteRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowSequenceStepWriteRepository(driver *neo4j.DriverWithContext, database string) FlowSequenceStepWriteRepository {
	return &flowSequenceStepWriteRepositoryImpl{driver: driver, database: database}
}

func (r *flowSequenceStepWriteRepositoryImpl) Merge(ctx context.Context, input *entity.FlowSequenceStepEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowSequenceStepWriteRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	onCreate := `ON CREATE SET
				fss.createdAt = $createdAt,
				fss.updatedAt = $updatedAt,
				fss.index = $index,
				fss.name = $name,
				fss.status = $status,
				fss.action = $action,`
	onMatch := `ON MATCH SET
				fss.updatedAt = $updatedAt,
				fss.index = $index,
				fss.name = $name,
				fss.status = $status,
				fss.action = $action,`

	params := map[string]any{
		"tenant":    common.GetTenantFromContext(ctx),
		"id":        input.Id,
		"createdAt": utils.TimeOrNow(input.CreatedAt),
		"updatedAt": utils.TimeOrNow(input.UpdatedAt),
		"index":     input.Index,
		"name":      input.Name,
		"status":    input.Status,
		"action":    input.Action,
	}

	if input.Action == entity.FlowSequenceStepActionWait {
		onCreate += `
				fss.actionData_minutes = $actionData_waitTime`
		onMatch += `
				fss.actionData_minutes = $actionData_waitTime`
		params["actionData_waitTime"] = input.ActionData.Wait.Minutes
	}
	if input.Action == entity.FlowSequenceStepActionEmailNew {
		onCreate += `
				fss.actionData_subject = $actionData_subject,
				fss.actionData_bodyTemplate = $actionData_bodyTemplate`
		onMatch += `
				fss.actionData_subject = $actionData_subject,
				fss.actionData_bodyTemplate = $actionData_bodyTemplate`
		params["actionData_subject"] = input.ActionData.EmailNew.Subject
		params["actionData_bodyTemplate"] = input.ActionData.EmailNew.BodyTemplate
	}
	if input.Action == entity.FlowSequenceStepActionEmailReply {
		onCreate += `
				fss.actionData_stepId = $actionData_stepId,
				fss.actionData_subject = $actionData_subject,
				fss.actionData_bodyTemplate = $actionData_bodyTemplate`
		onMatch += `
				fss.actionData_stepId = $actionData_stepId,
				fss.actionData_subject = $actionData_subject,
				fss.actionData_bodyTemplate = $actionData_bodyTemplate`
		params["actionData_stepId"] = input.ActionData.EmailReply.StepID
		params["actionData_subject"] = input.ActionData.EmailReply.Subject
		params["actionData_bodyTemplate"] = input.ActionData.EmailReply.BodyTemplate
	}
	if input.Action == entity.FlowSequenceStepActionLinkedinConnectionRequest {
		onCreate += `
				fss.actionData_messageTemplate = $actionData_messageTemplate`
		onMatch += `
				fss.actionData_messageTemplate = $actionData_messageTemplate`
		params["actionData_messageTemplate"] = input.ActionData.LinkedinConnectionRequest.MessageTemplate
	}
	if input.Action == entity.FlowSequenceStepActionLinkedinMessage {
		onCreate += `
				fss.actionData_messageTemplate = $actionData_messageTemplate`
		onMatch += `
				fss.actionData_messageTemplate = $actionData_messageTemplate`
		params["actionData_messageTemplate"] = input.ActionData.LinkedinMessage.MessageTemplate
	}

	cypher := fmt.Sprintf(`
			MATCH (t:Tenant {name:$tenant})
			MERGE (t)<-[:BELONGS_TO_TENANT]-(fss:FlowSequenceStep:FlowSequenceStep_%s {id: $id})
			%s
			%s
			RETURN fss`, common.GetTenantFromContext(ctx), onCreate, onMatch)

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
