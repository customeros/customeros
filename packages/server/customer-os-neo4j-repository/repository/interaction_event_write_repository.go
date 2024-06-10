package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type InteractionEventCreateFields struct {
	SourceFields       model.Source `json:"sourceFields"`
	CreatedAt          time.Time    `json:"createdAt"`
	Content            string       `json:"content"`
	ContentType        string       `json:"contentType"`
	Channel            string       `json:"channel"`
	ChannelData        string       `json:"channelData"`
	Identifier         string       `json:"identifier"`
	EventType          string       `json:"eventType"`
	BelongsToIssueId   string       `json:"belongsToIssueId"`
	BelongsToSessionId string       `json:"belongsToSessionId"`
	Hide               bool         `json:"hide"`
}

type InteractionEventUpdateFields struct {
	Content     string `json:"content"`
	ContentType string `json:"contentType"`
	Channel     string `json:"channel"`
	ChannelData string `json:"channelData"`
	Identifier  string `json:"identifier"`
	EventType   string `json:"eventType"`
	Hide        bool   `json:"hide"`
	Source      string `json:"source"`
}

type InteractionEventWriteRepository interface {
	Create(ctx context.Context, tenant, interactionEventId string, data InteractionEventCreateFields) error
	Update(ctx context.Context, tenant, interactionEventId string, data InteractionEventUpdateFields) error
	SetAnalysisForInteractionEvent(ctx context.Context, tenant, interactionEventId, content, contentType, analysisType, source, appSource string) error
	RemoveAllActionItemsForInteractionEvent(ctx context.Context, tenant, interactionEventId string) error
	AddActionItemForInteractionEvent(ctx context.Context, tenant, interactionEventId, content, source, appSource string) error
	LinkInteractionEventWithSenderById(ctx context.Context, tenant, interactionEventId, entityId, label, relationType string) error
	LinkInteractionEventWithReceiverById(ctx context.Context, tenant, interactionEventId, entityId, label, relationType string) error

	LinkInteractionEventToSession(ctx context.Context, tenant, interactionEventId, interactionSessionId string) error
	InteractionEventSentByEmail(ctx context.Context, tenant, interactionEventId, emailId string) error
	InteractionEventSentToEmails(ctx context.Context, tenant, interactionEventId, sentType string, emailsId []string) error
}

type interactionEventWriteRepository struct {
	driver   *neo4j.DriverWithContext
	database string
}

func NewInteractionEventWriteRepository(driver *neo4j.DriverWithContext, database string) InteractionEventWriteRepository {
	return &interactionEventWriteRepository{
		driver:   driver,
		database: database,
	}
}

func (r *interactionEventWriteRepository) Create(ctx context.Context, tenant, interactionEventId string, data InteractionEventCreateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventWriteRepository.Create")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, interactionEventId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MERGE (i:InteractionEvent:InteractionEvent_%s {id:$interactionEventId}) 
							ON CREATE SET 
								i:TimelineEvent,
								i:TimelineEvent_%s,
								i.createdAt=$createdAt,
								i.updatedAt=datetime(),
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
								i.updatedAt = datetime(),
								i.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE i.sourceOfTruth END,
								i.syncedWithEventStore = true
							WITH i
							OPTIONAL MATCH (is:Issue:Issue_%s {id:$belongsToIssueId}) 
							WHERE $belongsToIssueId <> ""
							FOREACH (ignore IN CASE WHEN is IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (i)-[:PART_OF]->(is))
							WITH i
							OPTIONAL MATCH (is:InteractionSession:InteractionSession_%s {id:$belongsToSessionId}) 
							WHERE $belongsToSessionId <> ""
							FOREACH (ignore IN CASE WHEN is IS NOT NULL THEN [1] ELSE [] END |
    							MERGE (i)-[:PART_OF]->(is))
							`, tenant, tenant, tenant, tenant)
	params := map[string]any{
		"tenant":             tenant,
		"interactionEventId": interactionEventId,
		"createdAt":          data.CreatedAt,
		"source":             data.SourceFields.Source,
		"sourceOfTruth":      data.SourceFields.Source,
		"appSource":          data.SourceFields.AppSource,
		"content":            data.Content,
		"contentType":        data.ContentType,
		"channel":            data.Channel,
		"channelData":        data.ChannelData,
		"identifier":         data.Identifier,
		"eventType":          data.EventType,
		"belongsToIssueId":   data.BelongsToIssueId,
		"belongsToSessionId": data.BelongsToSessionId,
		"hide":               data.Hide,
		"overwrite":          data.SourceFields.Source == constants.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *interactionEventWriteRepository) Update(ctx context.Context, tenant, interactionEventId string, data InteractionEventUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventWriteRepository.Update")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, interactionEventId)
	tracing.LogObjectAsJson(span, "data", data)

	cypher := fmt.Sprintf(`MATCH (i:InteractionEvent:InteractionEvent_%s {id:$interactionEventId})
		 	SET	
				i.content= CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.content is null OR i.content = '' THEN $content ELSE i.content END,
				i.contentType= CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.contentType is null OR i.contentType = '' THEN $contentType ELSE i.contentType END,
				i.channel= CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.channel is null OR i.channel = '' THEN $channel ELSE i.channel END,
				i.channelData= CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.channelData is null OR i.channelData = '' THEN $channelData ELSE i.channelData END,	
				i.identifier= CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.identifier is null OR i.identifier = '' THEN $identifier ELSE i.identifier END,
				i.eventType= CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true OR i.eventType is null OR i.eventType = '' THEN $eventType ELSE i.eventType END,
				i.hide= CASE WHEN i.sourceOfTruth=$sourceOfTruth OR $overwrite=true THEN $hide ELSE i.hide END,
				i.updatedAt = datetime(),
				i.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE i.sourceOfTruth END,
				i.syncedWithEventStore = true`, tenant)
	params := map[string]any{
		"tenant":             tenant,
		"interactionEventId": interactionEventId,
		"content":            data.Content,
		"contentType":        data.ContentType,
		"channel":            data.Channel,
		"channelData":        data.ChannelData,
		"identifier":         data.Identifier,
		"eventType":          data.EventType,
		"hide":               data.Hide,
		"sourceOfTruth":      data.Source,
		"overwrite":          data.Source == constants.SourceOpenline,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *interactionEventWriteRepository) SetAnalysisForInteractionEvent(ctx context.Context, tenant, interactionEventId, content, contentType, analysisType, source, appSource string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventWriteRepository.SetAnalysisForInteractionEvent")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, interactionEventId)
	span.LogFields(log.String("content", content), log.String("contentType", contentType), log.String("source", source), log.String("appSource", appSource))

	cypher := fmt.Sprintf(`MATCH (i:InteractionEvent_%s{id:$interactionEventId})
							MERGE (i)<-[r:DESCRIBES]-(a:Analysis_%s {analysisType:$analysisType})
							ON CREATE SET 
								a:Analysis,
								a.id=randomUUID(),
								a.createdAt=datetime(),
								a.updatedAt=datetime(),
								a.analysisType=$analysisType,
								a.source=$source,
								a.sourceOfTruth=$sourceOfTruth,
								a.appSource=$appSource,
								a.content=$content,
								a.contentType=$contentType
							ON MATCH SET
								a.content=$content,
								a.contentType=$contentType`, tenant, tenant)
	params := map[string]any{
		"interactionEventId": interactionEventId,
		"source":             source,
		"sourceOfTruth":      source,
		"appSource":          appSource,
		"content":            content,
		"contentType":        contentType,
		"analysisType":       analysisType,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *interactionEventWriteRepository) AddActionItemForInteractionEvent(ctx context.Context, tenant, interactionEventId, content, source, appSource string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventWriteRepository.AddActionItemForInteractionEvent")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, interactionEventId)

	cypher := fmt.Sprintf(`MATCH (i:InteractionEvent{id:$interactionEventId})
							WHERE i:InteractionEvent_%s
							MERGE (i)-[r:DESCRIBES]->(a:ActionItem_%s {id:randomUUID()})
							SET 
								a:ActionItem,
								a.createdAt=datetime(),
								a.updatedAt=datetime(),
								a.source=$source,
								a.sourceOfTruth=$sourceOfTruth,
								a.appSource=$appSource,
								a.content=$content`, tenant, tenant)
	params := map[string]any{
		"interactionEventId": interactionEventId,
		"source":             source,
		"sourceOfTruth":      source,
		"appSource":          appSource,
		"content":            content,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *interactionEventWriteRepository) RemoveAllActionItemsForInteractionEvent(ctx context.Context, tenant, interactionEventId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventWriteRepository.RemoveActionItemsForInteractionEvent")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, interactionEventId)

	cypher := fmt.Sprintf(`MATCH (i:InteractionEvent_%s{id:$interactionEventId})-[:INCLUDES]->(a:ActionItem)
							DETACH DELETE a`, tenant)
	params := map[string]any{
		"interactionEventId": interactionEventId,
	}
	span.LogFields(log.String("query", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *interactionEventWriteRepository) LinkInteractionEventWithSenderById(ctx context.Context, tenant, interactionEventId, entityId, label, relationType string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventWriteRepository.LinkInteractionEventWithSenderById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, interactionEventId)
	span.LogFields(log.String("entityId", entityId), log.String("label", label), log.String("relationType", relationType))

	cypher := fmt.Sprintf(`MATCH (ie:InteractionEvent {id:$interactionEventId}) 
		MATCH (:Tenant {name:$tenant})-[*1..2]-(sender:Contact|User|Organization|JobRole|Email|PhoneNumber {id:$entityId}) 
		WHERE $label IN labels(sender) AND ie:InteractionEvent_%s
		MERGE (ie)-[rel:SENT_BY]->(sender)
		ON CREATE SET rel.type=$relationType`, tenant)
	params := map[string]interface{}{
		"tenant":             tenant,
		"interactionEventId": interactionEventId,
		"entityId":           entityId,
		"relationType":       relationType,
		"label":              label,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *interactionEventWriteRepository) LinkInteractionEventWithReceiverById(ctx context.Context, tenant, interactionEventId, entityId, label, relationType string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventWriteRepository.LinkInteractionEventWithReceiverById")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.SetTag(tracing.SpanTagEntityId, interactionEventId)
	span.LogFields(log.String("entityId", entityId), log.String("label", label), log.String("relationType", relationType))

	cypher := fmt.Sprintf(`MATCH (ie:InteractionEvent_%s {id:$interactionEventId}) 
		MATCH (:Tenant {name:$tenant})--(sender:Contact|User|Organization|JobRole|Email|PhoneNumber {id:$entityId}) 
		WHERE $label IN labels(sender)
		MERGE (ie)-[rel:SENT_TO]->(sender)
		ON CREATE SET rel.type=$relationType`, tenant)
	params := map[string]interface{}{
		"tenant":             tenant,
		"interactionEventId": interactionEventId,
		"entityId":           entityId,
		"label":              label,
		"relationType":       relationType,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *interactionEventWriteRepository) LinkInteractionEventToSession(ctx context.Context, tenant, interactionEventId, interactionSessionId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.LinkInteractionEventToSession")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("interactionEventId", interactionEventId))

	cypher := fmt.Sprintf(`
	MATCH (is:InteractionSession_%s {id:$interactionSessionId})
	MATCH (ie:InteractionEvent {id:$interactionEventId})
    MERGE (ie)-[:PART_OF]->(is)
		`, tenant)

	params := map[string]interface{}{
		"tenant":               tenant,
		"interactionSessionId": interactionSessionId,
		"interactionEventId":   interactionEventId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *interactionEventWriteRepository) InteractionEventSentByEmail(ctx context.Context, tenant, interactionEventId, emailId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.LinkInteractionEventToSession")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("interactionEventId", interactionEventId))

	cypher := fmt.Sprintf(`
	MATCH (ie:InteractionEvent_%s {id:$interactionEventId})
	MATCH (e:Email_%s {id: $emailId})
	MERGE (ie)-[:SENT_BY]->(e)
`, tenant, tenant)

	params := map[string]interface{}{
		"tenant":             tenant,
		"interactionEventId": interactionEventId,
		"emailId":            emailId,
	}

	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}

func (r *interactionEventWriteRepository) InteractionEventSentToEmails(ctx context.Context, tenant, interactionEventId, sentType string, emailsId []string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.LinkInteractionEventToSession")
	defer span.Finish()
	tracing.SetNeo4jRepositorySpanTags(span, tenant)
	span.LogFields(log.String("interactionEventId", interactionEventId))

	cypher := fmt.Sprintf(`
	MATCH (ie:InteractionEvent_%s {id:$interactionEventId})
	MATCH (e:Email_%s) WHERE e.id IN $emailsId
	MERGE (ie)-[:SENT_TO {type: $sentType}]->(e)
`, tenant, tenant)

	params := map[string]interface{}{
		"tenant":             tenant,
		"interactionEventId": interactionEventId,
		"sentType":           sentType,
		"emailsId":           emailsId,
	}
	span.LogFields(log.String("cypher", cypher))
	tracing.LogObjectAsJson(span, "params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		tracing.TraceErr(span, err)
	}
	return err
}
