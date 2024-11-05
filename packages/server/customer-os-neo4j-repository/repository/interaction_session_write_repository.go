package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type InteractionSessionWriteRepository interface {
	CreateInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, interactionSessionId string, data entity.InteractionSessionEntity) error
}

type interactionSessionWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewInteractionSessionWriteRepository(driver *neo4j.DriverWithContext, database string) InteractionSessionWriteRepository {
	return &interactionSessionWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *interactionSessionWriteRepository) CreateInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, interactionSessionId string, data entity.InteractionSessionEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionWriteRepository.CreateInTx")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, interactionSessionId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MERGE (i:InteractionSession:InteractionSession_%s {id:$interactionSessionId}) 
							ON CREATE SET 
								i.createdAt=$createdAt,
								i.updatedAt=datetime(),
								i.source=$source,
								i.sourceOfTruth=$sourceOfTruth,
								i.appSource=$appSource,
								i.status=$status,
								i.channel=$channel,
								i.channelData=$channelData,
								i.identifier=$identifier,
								i.type=$type,
								i.name=$name
							`, tenant)
	params := map[string]any{
		"tenant":               tenant,
		"interactionSessionId": interactionSessionId,
		"createdAt":            utils.NowIfZero(data.CreatedAt),
		"source":               data.Source,
		"sourceOfTruth":        data.Source,
		"appSource":            data.AppSource,
		"channel":              data.Channel,
		"channelData":          data.ChannelData,
		"identifier":           data.Identifier,
		"type":                 data.Type,
		"status":               data.Status,
		"name":                 data.Name,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	_, err := tx.Run(ctx, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (r *interactionSessionWriteRepository) executeWriteQuery(ctx context.Context, query string, params map[string]any) error {
	return utils.ExecuteWriteQueryOnDb(ctx, *r.driver, r.database, query, params)
}
