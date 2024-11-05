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
	Merge(ctx context.Context, tx *neo4j.ManagedTransaction, entity *entity.FlowActionEntity) (*dbtype.Node, error)
	DeleteForFlow(ctx context.Context, tx *neo4j.ManagedTransaction, id string) error
	Delete(ctx context.Context, id string) error
}

type flowActionWriteRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowActionWriteRepository(driver *neo4j.DriverWithContext, database string) FlowActionWriteRepository {
	return &flowActionWriteRepositoryImpl{driver: driver, database: database}
}

func (r *flowActionWriteRepositoryImpl) Merge(ctx context.Context, tx *neo4j.ManagedTransaction, input *entity.FlowActionEntity) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowActionWriteRepository.Merge")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	onCreate := `ON CREATE SET
				fa.createdAt = $createdAt,
				fa.updatedAt = $updatedAt,
				fa.externalId = $externalId,
				fa.json = $json,
				fa.type = $type,
				fa.waitBefore = $waitBefore,
				fa.action = $action`
	onMatch := `ON MATCH SET
				fa.updatedAt = $updatedAt,
				fa.externalId = $externalId,
				fa.json = $json,
				fa.type = $type,
				fa.waitBefore = $waitBefore,
				fa.action = $action`

	params := map[string]any{
		"tenant":     common.GetTenantFromContext(ctx),
		"id":         input.Id,
		"externalId": input.ExternalId,
		"createdAt":  utils.NowIfZero(input.CreatedAt),
		"updatedAt":  utils.NowIfZero(input.UpdatedAt),
		"json":       input.Json,
		"type":       input.Type,
		"waitBefore": input.Data.WaitBefore,
		"action":     input.Data.Action,
	}

	if input.Data.Action == entity.FlowActionTypeFlowStart {
		onCreate += `,
				fa.data_triggerType = $data_triggerType,
				fa.data_entity = $data_entity`
		onMatch += `,
				fa.data_triggerType = $data_triggerType,
				fa.data_entity = $data_entity`
		params["data_entity"] = input.Data.Entity
		params["data_triggerType"] = input.Data.TriggerType
	}

	if input.Data.Action == entity.FlowActionTypeEmailNew {
		onCreate += `,
				fa.data_subject = $data_subject,
				fa.data_bodyTemplate = $data_bodyTemplate`
		onMatch += `,
				fa.data_subject = $data_subject,
				fa.data_bodyTemplate = $data_bodyTemplate`
		params["data_subject"] = input.Data.Subject
		params["data_bodyTemplate"] = input.Data.BodyTemplate
	}
	if input.Data.Action == entity.FlowActionTypeEmailReply {
		onCreate += `,
				fa.data_bodyTemplate = $data_bodyTemplate`
		onMatch += `,
				fa.data_bodyTemplate = $data_bodyTemplate`
		params["data_bodyTemplate"] = input.Data.BodyTemplate
	}
	if input.Data.Action == entity.FlowActionTypeLinkedinConnectionRequest {
		onCreate += `,
				fa.data_messageTemplate = $data_messageTemplate`
		onMatch += `,
				fa.data_messageTemplate = $data_messageTemplate`
		params["data_messageTemplate"] = input.Data.MessageTemplate
	}
	if input.Data.Action == entity.FlowActionTypeLinkedinMessage {
		onCreate += `,
				fa.data_messageTemplate = $data_messageTemplate`
		onMatch += `,
				fa.data_messageTemplate = $data_messageTemplate`
		params["data_messageTemplate"] = input.Data.MessageTemplate
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

func (r *flowActionWriteRepositoryImpl) DeleteForFlow(ctx context.Context, tx *neo4j.ManagedTransaction, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonWriteRepository.DeleteForFlow")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})<-[:BELONGS_TO_TENANT]-(f:Flow:Flow_%s { id: $id })-[r:HAS]-(fa:FlowAction_%s) detach delete fa`, tenant, tenant)

	params := map[string]any{
		"tenant": tenant,
		"id":     id,
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
		return err
	}

	return nil
}

func (r *flowActionWriteRepositoryImpl) Delete(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "CommonWriteRepository.Delete")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name: $tenant})<-[r:BELONGS_TO_TENANT]-(fa:FlowAction_%s {id:$id}) delete r, fa`, tenant)

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
