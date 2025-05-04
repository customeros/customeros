package neo4j_repository

import (
	"context"
	"fmt"
	commonmodel "github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/constants"
	neo4jentity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/model"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

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
	MergeByExternalSystem(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, syncDate time.Time, message commonmodel.SaveEmailMessage, source, appSource string) (string, error)
	LinkInteractionEventToSession(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, interactionEventId, interactionSessionId string) error
	InteractionEventSentByEmail(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, interactionEventId, emailId string) error
	InteractionEventSentToEmails(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, interactionEventId, sentType string, emailIds []string) error
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
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionEventWriteRepository.CreateInTx")
	defer spans.Finish()

	spans.LogKV("interactionEventId", interactionEventId)
	spans.LogObjectAsJson("data", data)

	cypher := fmt.Sprintf(`MERGE (i:InteractionEvent:InteractionEvent_%s {id:$interactionEventId}) 
							ON CREATE SET 
								i:TimelineEvent,
								i:TimelineEvent_%s,
								i.createdAt=$createdAt,
								i.updatedAt=datetime(),
								i.source=$source,
								i.appSource=$appSource,
								i.content=$content,
								i.contentType=$contentType,
								i.channel=$channel,
								i.channelData=$channelData,
								i.identifier=$identifier,
								i.eventType=$eventType,
								i.hide=$hide
							ON MATCH SET 	
								i.content = CASE WHEN $overwrite=true OR i.content is null OR i.content = '' THEN $content ELSE i.content END,
								i.contentType = CASE WHEN $overwrite=true OR i.contentType is null OR i.contentType = '' THEN $contentType ELSE i.contentType END,
								i.channel = CASE WHEN $overwrite=true OR i.channel is null OR i.channel = '' THEN $channel ELSE i.channel END,
								i.channelData = CASE WHEN $overwrite=true OR i.channelData is null OR i.channelData = '' THEN $channelData ELSE i.channelData END,
								i.identifier = CASE WHEN $overwrite=true OR i.identifier is null OR i.identifier = '' THEN $identifier ELSE i.identifier END,
								i.eventType = CASE WHEN $overwrite=true OR i.eventType is null OR i.eventType = '' THEN $eventType ELSE i.eventType END,
								i.hide = CASE WHEN $overwrite=true THEN $hide ELSE i.hide END,
								i.updatedAt = datetime()
							`, tenant, tenant)
	params := map[string]any{
		"tenant":             tenant,
		"interactionEventId": interactionEventId,
		"createdAt":          utils.NowIfZero(data.CreatedAt),
		"source":             data.Source,
		"appSource":          data.AppSource,
		"content":            data.Content,
		"contentType":        data.ContentType,
		"channel":            data.Channel.String(),
		"channelData":        data.ChannelData,
		"identifier":         data.Identifier,
		"eventType":          data.EventType,
		"hide":               data.Hide,
		"overwrite":          data.Source == constants.SourceOpenline,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := tx.Run(ctx, cypher, params)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	return nil
}

func (r *interactionEventWriteRepository) Update(ctx context.Context, tenant, interactionEventId string, data InteractionEventUpdateFields) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionEventWriteRepository.Update")
	defer spans.Finish()

	spans.LogKV("interactionEventId", interactionEventId)
	spans.LogObjectAsJson("data", data)

	cypher := fmt.Sprintf(`MATCH (i:InteractionEvent:InteractionEvent_%s {id:$interactionEventId})
		 	SET	
				i.content= CASE WHEN $overwrite=true OR i.content is null OR i.content = '' THEN $content ELSE i.content END,
				i.contentType= CASE WHEN $overwrite=true OR i.contentType is null OR i.contentType = '' THEN $contentType ELSE i.contentType END,
				i.channel= CASE WHEN $overwrite=true OR i.channel is null OR i.channel = '' THEN $channel ELSE i.channel END,
				i.channelData= CASE WHEN  $overwrite=true OR i.channelData is null OR i.channelData = '' THEN $channelData ELSE i.channelData END,	
				i.identifier= CASE WHEN  $overwrite=true OR i.identifier is null OR i.identifier = '' THEN $identifier ELSE i.identifier END,
				i.eventType= CASE WHEN  $overwrite=true OR i.eventType is null OR i.eventType = '' THEN $eventType ELSE i.eventType END,
				i.hide= CASE WHEN  $overwrite=true THEN $hide ELSE i.hide END,
				i.updatedAt = datetime()`, tenant)
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
		"overwrite":          data.Source == constants.SourceOpenline,
	}
	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	err := utils.ExecuteWriteQuery(ctx, *r.driver, cypher, params)
	if err != nil {
		spans.TraceError(err)
	}
	return err
}

func (r *interactionEventWriteRepository) MergeByExternalSystem(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, syncDate time.Time, message commonmodel.SaveEmailMessage, source, appSource string) (string, error) {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionEventWriteRepository.MergeByExternalSystem")
	defer spans.Finish()

	spans.LogKV("source", source, "appSource", appSource)
	spans.LogObjectAsJson("message", message)

	cypher := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId}) 
		 MERGE (ie:InteractionEvent_%s {source:$source, channel:$channel})-[rel:IS_LINKED_WITH {externalId:$externalId}]->(e) 
		 ON CREATE SET 
		  ie:InteractionEvent, 
		  ie:TimelineEvent, 
		  ie:TimelineEvent_%s, 
		  rel.syncDate=$syncDate, 
		  ie.createdAt=$createdAt, 
		  ie.id=randomUUID(), 
		  ie.identifier=$identifier, 
		  ie.channel=$channel, 
		  ie.channelData=$channelData, 
		  ie.content=$content, 
		  ie.contentType=$contentType, 
		  ie.appSource=$appSource,
		  ie.updatedAt=datetime()
		 WITH ie 
		 RETURN ie.id`, tenant, tenant)

	params := map[string]interface{}{
		"tenant":           tenant,
		"identifier":       message.ExternalId,
		"source":           source,
		"appSource":        appSource,
		"externalId":       message.ExternalId,
		"externalSystemId": message.ExternalSystem,
		"syncDate":         syncDate,
		"createdAt":        message.CreatedAt,
		"channel":          message.Channel,
		"channelData":      message.ChannelData,
	}

	if message.Html != "" {
		params["content"] = message.Html
		params["contentType"] = "text/html"
	} else {
		params["content"] = message.Text
		params["contentType"] = "text/plain"
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	queryResult, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		qr, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return utils.ExtractSingleRecordFirstValueAsString(ctx, qr, err)
	})
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	return queryResult.(string), nil
}

func (r *interactionEventWriteRepository) LinkInteractionEventToSession(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, interactionEventId, interactionSessionId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionEventWriteRepository.LinkInteractionEventToSession")
	defer spans.Finish()

	spans.LogKV("interactionEventId", interactionEventId)
	spans.LogKV("interactionSessionId", interactionSessionId)

	cypher := fmt.Sprintf(`MATCH (ie:InteractionEvent_%s {id:$interactionEventId}) 
		 MATCH (is:InteractionSession_%s {id:$interactionSessionId}) 
		 MERGE (ie)-[:PART_OF]->(is)`, tenant, tenant)

	params := map[string]interface{}{
		"tenant":               tenant,
		"interactionEventId":   interactionEventId,
		"interactionSessionId": interactionSessionId,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
	}

	return err
}

func (r *interactionEventWriteRepository) InteractionEventSentByEmail(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, interactionEventId, emailId string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionEventWriteRepository.InteractionEventSentByEmail")
	defer spans.Finish()

	spans.LogKV("interactionEventId", interactionEventId)

	cypher := fmt.Sprintf(`MATCH (ie:InteractionEvent_%s {id:$interactionEventId})
		 MATCH (e:Email_%s {id: $emailId})
		 MERGE (ie)-[:SENT_BY]->(e)`, tenant, tenant)

	params := map[string]interface{}{
		"tenant":             tenant,
		"interactionEventId": interactionEventId,
		"emailId":            emailId,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
	}

	return err
}

func (r *interactionEventWriteRepository) InteractionEventSentToEmails(ctx context.Context, tx *neo4j.ManagedTransaction, tenant, interactionEventId, sentType string, emailIds []string) error {
	spans, ctx := telemetry.StartNeo4jSpan(ctx, "InteractionEventWriteRepository.InteractionEventSentToEmails")
	defer spans.Finish()

	spans.LogKV("interactionEventId", interactionEventId)

	cypher := fmt.Sprintf(`MATCH (ie:InteractionEvent_%s {id:$interactionEventId})
		 MATCH (e:Email_%s) WHERE e.id in $emailIds
		 MERGE (ie)-[:SENT_TO {type: $sentType}]->(e)`, tenant, tenant)

	params := map[string]interface{}{
		"tenant":             tenant,
		"interactionEventId": interactionEventId,
		"sentType":           sentType,
		"emailIds":           emailIds,
	}

	spans.LogKV("cypher", cypher)
	spans.LogObjectAsJson("params", params)

	_, err := utils.ExecuteWriteInTransaction(ctx, r.driver, r.database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		spans.TraceError(err)
	}

	return err
}
