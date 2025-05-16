package repository

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionEventRepository.GetMatchedInteractionEventId")
	defer spans.Finish()
	spans.LogKV("externalSystem", externalSystem)
	spans.LogKV("externalId", externalId)
	spans.LogKV("externalSourceEntity", externalSourceEntity)

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
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionEventRepository.GetInteractionEventIdByExternalId")
	defer spans.Finish()
	spans.LogKV("externalSystemId", externalSystemId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})
					MATCH (ie:InteractionEvent_%s)-[:IS_LINKED_WITH {externalId:$externalId}]->(e)
					RETURN ie.id ORDER BY ie.createdAt`, tenant)
	params := map[string]any{
		"tenant":           tenant,
		"externalId":       externalId,
		"externalSystemId": externalSystemId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

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
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "InteractionEventRepository.GetById")
	defer spans.Finish()
	spans.LogKV("interactionEventId", interactionEventId)

	cypher := fmt.Sprintf(`MATCH (i:InteractionEvent_%s {id:$interactionEventId}) RETURN i`, tenant)
	params := map[string]any{
		"interactionEventId": interactionEventId,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	dbRecord, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, cypher, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	})
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	return dbRecord.(*dbtype.Node), err
}
