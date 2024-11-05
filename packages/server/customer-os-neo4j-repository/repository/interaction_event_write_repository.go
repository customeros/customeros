package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/constants"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/model"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type InteractionEventCreateFields struct {
	SourceFields       model.SourceFields `json:"sourceFields"`
	CreatedAt          time.Time          `json:"createdAt"`
	Content            string             `json:"content"`
	ContentType        string             `json:"contentType"`
	Channel            string             `json:"channel"`
	ChannelData        string             `json:"channelData"`
	Identifier         string             `json:"identifier"`
	EventType          string             `json:"eventType"`
	BelongsToIssueId   string             `json:"belongsToIssueId"`
	BelongsToSessionId string             `json:"belongsToSessionId"`
	Hide               bool               `json:"hide"`
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
	CreateInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, interactionEventId string, data neo4jentity.InteractionEventEntity) error
	Update(ctx context.Context, tenant, interactionEventId string, data InteractionEventUpdateFields) error
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

func (r *interactionEventWriteRepository) CreateInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, interactionEventId string, data neo4jentity.InteractionEventEntity) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventWriteRepository.CreateInTx")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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
								i.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE i.sourceOfTruth END
							`, tenant, tenant)
	params := map[string]any{
		"tenant":             tenant,
		"interactionEventId": interactionEventId,
		"createdAt":          utils.NowIfZero(data.CreatedAt),
		"updatedAt":          utils.Now(),
		"source":             data.Source,
		"sourceOfTruth":      data.Source,
		"appSource":          data.AppSource,
		"content":            data.Content,
		"contentType":        data.ContentType,
		"channel":            data.Channel,
		"channelData":        data.ChannelData,
		"identifier":         data.Identifier,
		"eventType":          data.EventType,
		"hide":               data.Hide,
		"overwrite":          data.Source == constants.SourceOpenline,
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

func (r *interactionEventWriteRepository) Update(ctx context.Context, tenant, interactionEventId string, data InteractionEventUpdateFields) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventWriteRepository.Update")
	defer span.Finish()
	tracing.TagComponentNeo4jRepository(span)
	tracing.TagTenant(span, tenant)
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
				i.sourceOfTruth = case WHEN $overwrite=true THEN $sourceOfTruth ELSE i.sourceOfTruth END`, tenant)
	params := map[string]any{
		"tenant":             tenant,
		"interactionEventId": interactionEventId,
		"updatedAt":          time.Now(),
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
