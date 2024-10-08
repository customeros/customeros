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
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type FlowParticipantReadRepository interface {
	GetList(ctx context.Context, flowIds []string) ([]*utils.DbNodeAndId, error)
	CountWithStatus(ctx context.Context, flowId string, status entity.FlowParticipantStatus) (int64, error)
	IdentifyForContact(ctx context.Context, flowId, contactId string) (*neo4j.Node, error)
	GetById(ctx context.Context, id string) (*neo4j.Node, error)
}

type flowParticipantReadRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowParticipantReadRepository(driver *neo4j.DriverWithContext, database string) FlowParticipantReadRepository {
	return &flowParticipantReadRepositoryImpl{driver: driver, database: database}
}

func (r flowParticipantReadRepositoryImpl) GetList(ctx context.Context, flowIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowParticipantReadRepository.GetList")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	if flowIds != nil && len(flowIds) > 0 {
		span.LogFields(log.String("flowIds", fmt.Sprintf("%v", flowIds)))
	}

	tenant := common.GetTenantFromContext(ctx)

	params := map[string]any{
		"tenant": tenant,
	}

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fc:FlowParticipant_%s) `, tenant, tenant)
	if flowIds != nil && len(flowIds) > 0 {
		cypher += "WHERE f.id in $flowIds "
		params["flowIds"] = flowIds
	}
	cypher += "RETURN fc, f.id"

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	span.LogFields(log.Int("result.count", len(result.([]*utils.DbNodeAndId))))
	if len(result.([]*utils.DbNodeAndId)) == 0 {
		return nil, nil
	}
	return result.([]*utils.DbNodeAndId), err
}

func (r flowParticipantReadRepositoryImpl) CountWithStatus(ctx context.Context, flowId string, status entity.FlowParticipantStatus) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowParticipantReadRepository.CountWithStatus")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s {id: $flowId})-[:HAS]->(fc:FlowParticipant_%s {status: $status}) return count(fc)`, tenant, tenant)
	params := map[string]any{
		"tenant": tenant,
		"flowId": flowId,
		"status": status,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return queryResult.Single(ctx)
		}
	})
	if err != nil {
		return 0, err
	}
	organizationsCount := dbRecord.(*db.Record).Values[0].(int64)
	span.LogFields(log.Int64("result", organizationsCount))
	return organizationsCount, nil
}

func (r flowParticipantReadRepositoryImpl) IdentifyForContact(ctx context.Context, flowId, contactId string) (*neo4j.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowParticipantReadRepository.IdentifyForContact")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("flowId", flowId), log.String("contactId", contactId))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s { id: $flowId})-[:HAS]->(fc:FlowParticipant_%s)-[:HAS]->(c:Contact_%s {id: $contactId}) RETURN fc`, tenant, tenant, tenant)
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
			return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
		}
	})
	if err != nil && err.Error() == "Result contains no more records" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}

func (r flowParticipantReadRepositoryImpl) GetById(ctx context.Context, id string) (*neo4j.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowParticipantReadRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fc:FlowParticipant_%s {id: $id}) RETURN fc`, tenant, tenant)
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
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return result.(*dbtype.Node), nil
}
