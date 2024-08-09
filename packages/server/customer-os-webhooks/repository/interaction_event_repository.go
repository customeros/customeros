package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	tracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	localtracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type InteractionEventRepository interface {
	// Deprecated
	GetMatchedInteractionEventId(ctx context.Context, tenant, externalId, externalSystem, externalSourceEntity string) (string, error)
	// Deprecated
	GetInteractionEventIdByExternalId(ctx context.Context, tenant, externalId, externalSystem string) (string, error)
	// Deprecated
	GetById(ctx context.Context, tenant, interactionEventId string) (*dbtype.Node, error)
}

type interactionEventRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewInteractionEventRepository(driver *neo4j.DriverWithContext, database string) InteractionEventRepository {
	return &interactionEventRepository{
		driver:   driver,
		database: database,
	}
}

func (r *interactionEventRepository) GetMatchedInteractionEventId(ctx context.Context, tenant, externalId, externalSystem, externalSourceEntity string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.GetMatchedInteractionEventId")
	defer span.Finish()
	localtracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("externalSystem", externalSystem), log.String("externalId", externalId), log.String("externalSourceEntity", externalSourceEntity))

	filter := ""
	params := map[string]interface{}{
		"tenant":         tenant,
		"externalSystem": externalSystem,
		"externalId":     externalId,
	}
	if externalSourceEntity != "" {
		params["externalSource"] = externalSourceEntity
		filter = " WHERE i.externalSource = $externalSource "
	}
	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem})
				MATCH (i:InteractionEvent_%s)-[:IS_LINKED_WITH {externalId:$externalId}]->(e)
				%s
				RETURN i.id LIMIT 1`, tenant, filter)
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	interactionEventIds := dbRecords.([]*db.Record)
	if len(interactionEventIds) == 1 {
		return interactionEventIds[0].Values[0].(string), nil
	}
	return "", nil
}

func (r *interactionEventRepository) GetInteractionEventIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.GetInteractionEventIdByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})
					MATCH (ie:InteractionEvent_%s)-[:IS_LINKED_WITH {externalId:$externalId}]->(e)
					RETURN ie.id ORDER BY ie.createdAt`, tenant)
	params := map[string]any{
		"tenant":           tenant,
		"externalId":       externalId,
		"externalSystemId": externalSystemId,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractAllRecordsAsString(ctx, queryResult, err)
	})
	if err != nil {
		return "", err
	}
	if len(records.([]string)) == 0 {
		return "", nil
	}
	return records.([]string)[0], nil
}

func (r *interactionEventRepository) GetById(parentCtx context.Context, tenant, interactionEventId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(parentCtx, "InteractionEventRepository.GetById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.String("interactionEventId", interactionEventId))

	cypher := fmt.Sprintf(`MATCH (i:InteractionEvent_%s {id:$interactionEventId}) RETURN i`, tenant)
	params := map[string]any{
		"interactionEventId": interactionEventId,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		return nil, err
	}
	return dbRecord.(*dbtype.Node), err
}
