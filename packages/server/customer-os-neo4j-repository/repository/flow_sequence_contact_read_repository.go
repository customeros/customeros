package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type FlowSequenceContactReadRepository interface {
	GetList(ctx context.Context, sequenceIds []string) ([]*utils.DbNodeAndId, error)
	Identify(ctx context.Context, sequenceId, contactId, emailId string) (*neo4j.Node, error)
	GetById(ctx context.Context, id string) (*neo4j.Node, error)
}

type flowSequenceContactReadRepositoryImpl struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewFlowSequenceContactReadRepository(driver *neo4j.DriverWithContext, database string) FlowSequenceContactReadRepository {
	return &flowSequenceContactReadRepositoryImpl{driver: driver, database: database}
}

func (r flowSequenceContactReadRepositoryImpl) GetList(ctx context.Context, sequenceIds []string) ([]*utils.DbNodeAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowSequenceContactReadRepository.GetList")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	if sequenceIds != nil && len(sequenceIds) > 0 {
		span.LogFields(log.String("sequenceIds", fmt.Sprintf("%v", sequenceIds)))
	}

	tenant := common.GetTenantFromContext(ctx)

	params := map[string]any{
		"tenant": tenant,
	}

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fs:FlowSequence_%s)-[:HAS]->(fsc:FlowSequenceContact_%s) `, tenant, tenant, tenant)
	if sequenceIds != nil && len(sequenceIds) > 0 {
		cypher += "WHERE fs.id in $sequenceIds "
		params["sequenceIds"] = sequenceIds
	}
	cypher += "RETURN fsc, fs.id"

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

func (r flowSequenceContactReadRepositoryImpl) Identify(ctx context.Context, sequenceId, contactId, emailId string) (*neo4j.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowSequenceContactReadRepository.Identify")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("sequenceId", sequenceId), log.String("contactId", contactId), log.String("emailId", emailId))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`
MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fs:FlowSequence_%s {id: $sequenceId})-[:HAS]->(fsc:FlowSequenceContact_%s) 
WITH fsc
MATCH (e:Email_%s {id: $emailId})<-[:HAS]-(fsc)-[:HAS]->(c:Contact_%s {id: $contactId})

RETURN fsc`, tenant, tenant, tenant, tenant, tenant)
	params := map[string]any{
		"tenant":     tenant,
		"sequenceId": sequenceId,
		"contactId":  contactId,
		"emailId":    emailId,
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

func (r flowSequenceContactReadRepositoryImpl) GetById(ctx context.Context, id string) (*neo4j.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "FlowSequenceContactReadRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	span.LogFields(log.String("id", id))

	tenant := common.GetTenantFromContext(ctx)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:BELONGS_TO_TENANT]-(f:Flow_%s)-[:HAS]->(fs:FlowSequence_%s)-[:HAS]->(fsc:FlowSequenceContact_%s {id: $id}) RETURN fsc`, tenant, tenant)
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
