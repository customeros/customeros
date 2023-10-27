package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-webhooks/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type InteractionSessionRepository interface {
	GetInteractionSessionIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error)
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionRepository.GetInteractionSessionIdByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})
					MATCH (is:InteractionSession_%s)-[:IS_LINKED_WITH {externalId:$externalId}]->(e)
					RETURN is.id ORDER BY is.createdAt`, tenant)
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
