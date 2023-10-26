package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/domain/interaction_event/event"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/helper"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type InteractionEventRepository interface {
	Create(ctx context.Context, tenant, interactionEventId string, evt event.InteractionEventCreateEvent) error
	Update(ctx context.Context, tenant, interactionEventId string, eventData event.InteractionEventUpdateEvent) error
	SetAnalysisForInteractionEvent(ctx context.Context, tenant, interactionEventId, content, contentType, analysisType, source, appSource string, updatedAt time.Time) error
	RemoveAllActionItemsForInteractionEvent(ctx context.Context, tenant, interactionEventId string) error
	AddActionItemForInteractionEvent(ctx context.Context, tenant, interactionEventId, content, source, appSource string, updatedAt time.Time) error
	GetInteractionEvent(ctx context.Context, tenant, interactionEventId string) (*dbtype.Node, error)
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

func (r *interactionEventRepository) Create(ctx context.Context, tenant, interactionEventId string, evt event.InteractionEventCreateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IssueRepository.Create")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("interactionEventId", interactionEventId), log.String("createEvent", fmt.Sprintf("%+v", evt)))

	query := fmt.Sprintf(`MERGE (i:InteractionEvent:InteractionEvent_%s {id:$interactionEventId}) 
							ON CREATE SET 
								i:TimelineEvent,
								i:TimelineEvent_%s,
								i.createdAt=$createdAt,
								i.updatedAt=$updatedAt,
								i.source=$source,
								i.sourceOfTruth=$sourceOfTruth,
								i.appSource=$appSource,
								i.content=$content,
								i.contentType=$contentType,
								i.channel=$channel,
								i.channelData=$channelData,
								i.identifier=$identifier,
								i.eventType=$eventType,
								i.hide=$hide
							ON MATCH SET 	
								i.content = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.content is null OR i.content = '' THEN $content ELSE i.content END,
								i.contentType = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.contentType is null OR i.contentType = '' THEN $contentType ELSE i.contentType END,
								i.channel = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.channel is null OR i.channel = '' THEN $channel ELSE i.channel END,
								i.channelData = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.channelData is null OR i.channelData = '' THEN $channelData ELSE i.channelData END,
								i.identifier = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.identifier is null OR i.identifier = '' THEN $identifier ELSE i.identifier END,
								i.eventType = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.eventType is null OR i.eventType = '' THEN $eventType ELSE i.eventType END,
								i.hide = CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $hide ELSE i.hide END,
								i.updatedAt = $updatedAt,
								i.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE i.sourceOfTruth END,
								i.syncedWithEventStore = true
							WITH i
							OPTIONAL MATCH (is:Issue:Issue_%s {id:$partOfIssueId}) 
							WHERE $partOfIssueId <> ""
							FOREACH (ignore IN CASE WHEN is IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (i)-[:PART_OF]->(is))
							WITH i
							OPTIONAL MATCH (is:InteractionSession:InteractionSession_%s {id:$partOfSessionId}) 
							WHERE $partOfSessionId <> ""
							FOREACH (ignore IN CASE WHEN is IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (i)-[:PART_OF]->(is))
							`, tenant, tenant, tenant, tenant)
	params := map[string]any{
		"tenant":             tenant,
		"interactionEventId": interactionEventId,
		"createdAt":          evt.CreatedAt,
		"updatedAt":          evt.UpdatedAt,
		"source":             helper.GetSource(evt.Source),
		"sourceOfTruth":      helper.GetSourceOfTruth(evt.Source),
		"appSource":          helper.GetAppSource(evt.AppSource),
		"content":            evt.Content,
		"contentType":        evt.ContentType,
		"channel":            evt.Channel,
		"channelData":        evt.ChannelData,
		"identifier":         evt.Identifier,
		"eventType":          evt.EventType,
		"partOfIssueId":      evt.PartOfIssueId,
		"partOfSessionId":    evt.PartOfSessionId,
		"hide":               evt.Hide,
		"overwrite":          helper.GetSourceOfTruth(evt.Source) == constants.SourceOpenline,
	}
	span.LogFields(log.String("query", query), log.Object("params", params))

	return r.executeWriteQuery(ctx, query, params)
}

func (r *interactionEventRepository) Update(ctx context.Context, tenant, interactionEventId string, evt event.InteractionEventUpdateEvent) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("interactionEventId", interactionEventId), log.Object("event", fmt.Sprintf("%+v", evt)))

	query := fmt.Sprintf(`MATCH (i:InteractionEvent:InteractionEvent_%s {id:$interactionEventId})
		 	SET	
				i.content= CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.content is null OR i.content = '' THEN $content ELSE i.content END,
				i.contentType= CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.contentType is null OR i.contentType = '' THEN $contentType ELSE i.contentType END,
				i.channel= CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.channel is null OR i.channel = '' THEN $channel ELSE i.channel END,
				i.channelData= CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.channelData is null OR i.channelData = '' THEN $channelData ELSE i.channelData END,	
				i.identifier= CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.identifier is null OR i.identifier = '' THEN $identifier ELSE i.identifier END,
				i.eventType= CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.eventType is null OR i.eventType = '' THEN $eventType ELSE i.eventType END,
				i.hide= CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $hide ELSE i.hide END,
				i.updatedAt = $updatedAt,
				i.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE i.sourceOfTruth END,
				i.syncedWithEventStore = true`, tenant)
	params := map[string]any{
		"tenant":             tenant,
		"interactionEventId": interactionEventId,
		"updatedAt":          evt.UpdatedAt,
		"content":            evt.Content,
		"contentType":        evt.ContentType,
		"channel":            evt.Channel,
		"channelData":        evt.ChannelData,
		"identifier":         evt.Identifier,
		"eventType":          evt.EventType,
		"hide":               evt.Hide,
		"sourceOfTruth":      helper.GetSourceOfTruth(evt.Source),
		"overwrite":          helper.GetSourceOfTruth(evt.Source) == constants.SourceOpenline,
	}
	span.LogFields(log.String("query", query), log.Object("params", params))

	return r.executeWriteQuery(ctx, query, params)
}

func (r *interactionEventRepository) SetAnalysisForInteractionEvent(ctx context.Context, tenant, interactionEventId, content, contentType, analysisType, source, appSource string, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.SetAnalysisForInteractionEvent")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("interactionEventId", interactionEventId), log.Object("updatedAt", updatedAt),
		log.String("content", content), log.String("contentType", contentType), log.String("source", source), log.String("appSource", appSource))

	query := fmt.Sprintf(`MATCH (i:InteractionEvent_%s{id:$interactionEventId})
							MERGE (i)<-[r:DESCRIBES]-(a:Analysis_%s {analysisType:$analysisType})
							ON CREATE SET 
								a:Analysis,
								a.id=randomUUID(),
								a.createdAt=$createdAt,
								a.updatedAt=$updatedAt,
								a.analysisType=$analysisType,
								a.source=$source,
								a.sourceOfTruth=$sourceOfTruth,
								a.appSource=$appSource,
								a.content=$content,
								a.contentType=$contentType
							ON MATCH SET
								a.content=$content,
								a.contentType=$contentType
							RETURN a`, tenant, tenant)
	params := map[string]any{
		"interactionEventId": interactionEventId,
		"createdAt":          updatedAt,
		"updatedAt":          updatedAt,
		"source":             source,
		"sourceOfTruth":      source,
		"appSource":          appSource,
		"content":            content,
		"contentType":        contentType,
		"analysisType":       analysisType,
	}
	span.LogFields(log.String("query", query), log.Object("params", params))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, params)
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return err
	}
	return nil
}

func (r *interactionEventRepository) GetInteractionEvent(ctx context.Context, tenant, interactionEventId string) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.GetInteractionEvent")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("interactionEventId", interactionEventId))

	query := fmt.Sprintf(`MATCH (i:InteractionEvent_%s {id:$id}) RETURN i`, tenant)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jReadSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	params := map[string]any{
		"id": interactionEventId,
	}
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query, params); err != nil {
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

func (r *interactionEventRepository) AddActionItemForInteractionEvent(ctx context.Context, tenant, interactionEventId, content, source, appSource string, updatedAt time.Time) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.AddActionItemForInteractionEvent")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("interactionEventId", interactionEventId))

	query := fmt.Sprintf(`MATCH (i:InteractionEvent_%s{id:$interactionEventId})
							MERGE (i)-[r:DESCRIBES]->(a:ActionItem_%s {id:randomUUID()})
							SET 
								a:ActionItem,
								a.createdAt=$createdAt,
								a.updatedAt=$updatedAt,
								a.source=$source,
								a.sourceOfTruth=$sourceOfTruth,
								a.appSource=$appSource,
								a.content=$content
							RETURN a`, tenant, tenant)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver, utils.WithDatabaseName(r.database))
	defer session.Close(ctx)

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"interactionEventId": interactionEventId,
				"createdAt":          updatedAt,
				"updatedAt":          updatedAt,
				"source":             source,
				"sourceOfTruth":      source,
				"appSource":          appSource,
				"content":            content,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return err
	}
	return nil
}

func (r *interactionEventRepository) RemoveAllActionItemsForInteractionEvent(ctx context.Context, tenant, interactionEventId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.RemoveActionItemsForInteractionEvent")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(ctx, span, tenant)
	span.LogFields(log.String("interactionEventId", interactionEventId))

	query := fmt.Sprintf(`MATCH (i:InteractionEvent_%s{id:$interactionEventId})-[:INCLUDES]->(a:ActionItem)
							DETACH DELETE a`, tenant)
	params := map[string]any{
		"interactionEventId": interactionEventId,
	}
	span.LogFields(log.String("query", query), log.Object("params", params))

	return r.executeWriteQuery(ctx, query, params)
}

func (r *interactionEventRepository) executeWriteQuery(ctx context.Context, query string, params map[string]any) error {
	return utils.ExecuteWriteQueryOnDb(ctx, *r.driver, r.database, query, params)
}
