package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/events-processing-platform/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type InteractionEventRepository interface {
	SetAnalysisForInteractionEvent(ctx context.Context, tenant, interactionEventId, content, contentType, analysisType, source, appSource string, updatedAt time.Time) error

	RemoveAllActionItemsForInteractionEvent(ctx context.Context, tenant, interactionEventId string) error
	AddActionItemForInteractionEvent(ctx context.Context, tenant, interactionEventId, content, source, appSource string, updatedAt time.Time) error

	GetInteractionEvent(ctx context.Context, tenant, interactionEventId string) (*dbtype.Node, error)
}

type interactionEventRepository struct {
	driver *neo4j.DriverWithContext
}

func NewInteractionEventRepository(driver *neo4j.DriverWithContext) InteractionEventRepository {
	return &interactionEventRepository{
		driver: driver,
	}
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
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
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
				"contentType":        contentType,
				"analysisType":       analysisType,
			})
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

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, query,
			map[string]any{
				"id": interactionEventId,
			}); err != nil {
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

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
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
	span.LogFields(log.String("query", query))

	return utils.ExecuteQuery(ctx, *r.driver, query, map[string]any{
		"interactionEventId": interactionEventId,
	})
}
