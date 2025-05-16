package repository

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

type InteractionSessionRepository interface {
	// Deprecated
	GetInteractionSessionIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error)
	// Deprecated
	GetById(ctx context.Context, tenant, interactionEventId string) (*dbtype.Node, error)
}

type interactionSessionRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewInteractionSessionRepository(driver *neo4j.DriverWithContext, database string) InteractionSessionRepository {
	return &interactionSessionRepository{
		driver:   driver,
		database: database,
	}
}

func (r *interactionSessionRepository) GetInteractionSessionIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionSessionRepository.GetInteractionSessionIdByExternalId")
	defer spans.Finish()
	spans.LogKV("externalSystemId", externalSystemId)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})
					MATCH (is:InteractionSession_%s)-[:IS_LINKED_WITH {externalId:$externalId}]->(e)
					RETURN is.id ORDER BY is.createdAt`, tenant)
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

func (r *interactionSessionRepository) GetById(parentCtx context.Context, tenant, interactionSessionId string) (*dbtype.Node, error) {
	spans, ctx := telemetry.StartNeo4jSpan(parentCtx, "InteractionSessionRepository.GetById")
	defer spans.Finish()
	spans.LogKV("interactionSessionId", interactionSessionId)

	cypher := fmt.Sprintf(`MATCH (i:InteractionSession {id:$interactionSessionId}) WHERE i:InteractionSession_%s RETURN i`, tenant)
	params := map[string]any{
		"interactionSessionId": interactionSessionId,
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
