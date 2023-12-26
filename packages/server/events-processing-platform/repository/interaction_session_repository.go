package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_session/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

type InteractionSessionRepository interface {
	Create(ctx context.Context, tenant, interactionSessionId string, evt event.InteractionSessionCreateEvent) error
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

func (r *interactionSessionRepository) Create(ctx context.Context, tenant, interactionSessionId string, evt event.InteractionSessionCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionSessionRepository.Create")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.SetTag(tracing.SpanTagEntityId, interactionSessionId)
	tracing.LogObjectAsJson(span, "event", evt)

	cypher := fmt.Sprintf(`MERGE (i:InteractionSession:InteractionSession_%s {id:$interactionSessionId}) 
							ON CREATE SET 
								i.createdAt=$createdAt,
								i.updatedAt=$updatedAt,
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
		"createdAt":            evt.CreatedAt,
		"updatedAt":            evt.UpdatedAt,
		"source":               helper.GetSource(evt.Source),
		"sourceOfTruth":        helper.GetSourceOfTruth(evt.Source),
		"appSource":            helper.GetAppSource(evt.AppSource),
		"channel":              evt.Channel,
		"channelData":          evt.ChannelData,
		"identifier":           evt.Identifier,
		"type":                 evt.Type,
		"status":               evt.Status,
		"name":                 evt.Name,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	return r.executeWriteQuery(ctx, cypher, params)
}

func (r *interactionSessionRepository) executeWriteQuery(ctx context.Context, query string, params map[string]any) error {
	return utils.ExecuteWriteQueryOnDb(ctx, *r.driver, r.database, query, params)
}
