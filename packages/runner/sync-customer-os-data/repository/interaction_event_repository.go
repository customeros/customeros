package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/constants"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type InteractionEventRepository interface {
	GetMatchedInteractionEvent(ctx context.Context, tenant string, event entity.InteractionEventData) (string, error)
	MergeInteractionEvent(ctx context.Context, tenant string, syncDate time.Time, event entity.InteractionEventData) error
	MergeInteractionSessionForEvent(ctx context.Context, tenant, eventId, externalSource string, syncDate time.Time, session entity.InteractionSession) error

	MergeEmailInteractionSession(ctx context.Context, tenant string, date time.Time, message entity.EmailMessageData) (string, error)
	MergeEmailInteractionEvent(ctx context.Context, tenant, externalSystemId string, date time.Time, message entity.EmailMessageData) (string, error)
	LinkInteractionEventToSession(ctx context.Context, tenant, interactionEventId, interactionSessionId string) error

	InteractionEventSentByEmail(ctx context.Context, tenant, interactionEventId, emailId string) error
	InteractionEventSentToEmails(ctx context.Context, tenant, interactionEventId, sentType string, emails []string) error
	LinkInteractionEventAsPartOf(ctx context.Context, tenant, interactionEventId, id, label string) error

	LinkInteractionEventWithSenderById(ctx context.Context, tenant, eventId, id, label string) error
	LinkInteractionEventWithRecipientById(ctx context.Context, tenant, eventId, id, label, relationType string) error

	GetInteractionSessionIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error)
}

type interactionEventRepository struct {
	driver *neo4j.DriverWithContext
}

func NewInteractionEventRepository(driver *neo4j.DriverWithContext) InteractionEventRepository {
	return &interactionEventRepository{
		driver: driver,
	}
}

func (r *interactionEventRepository) GetMatchedInteractionEvent(ctx context.Context, tenant string, event entity.InteractionEventData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.GetMatchedInteractionEvent")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystem})
				OPTIONAL MATCH (ext)<-[:IS_LINKED_WITH {externalId:$externalId, externalSource:$externalSource}]-(ie:InteractionEvent_%s)
				WITH ie WHERE NOT ie IS NULL
				RETURN ie.id limit 1`

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"externalSystem": event.ExternalSystem,
				"externalId":     event.ExternalId,
				"externalSource": event.ExternalSourceEntity,
			})
		if err != nil {
			return nil, err
		}
		return queryResult.Collect(ctx)
	})
	if err != nil {
		return "", err
	}
	issueIDs := dbRecords.([]*db.Record)
	if len(issueIDs) > 0 {
		return issueIDs[0].Values[0].(string), nil
	}
	return "", nil
}

func (r *interactionEventRepository) MergeInteractionEvent(ctx context.Context, tenant string, syncDate time.Time, event entity.InteractionEventData) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.MergeInteractionEvent")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId}) 
		 MERGE (ie:InteractionEvent_%s {id:$id})-[rel:IS_LINKED_WITH {externalId:$externalId}]->(e) 
		 ON CREATE SET 
		  	ie:InteractionEvent,
		  	ie:TimelineEvent, 
		  	ie:TimelineEvent_%s, 
		  	rel.syncDate=$syncDate, 
			rel.externalUrl=$externalUrl,
			rel.externalSource=$externalSource,
		  	ie.createdAt=$createdAt,
			ie.channel=$channel,
			ie.eventType=$type, 
		  	ie.identifier=$identifier, 
		  	ie.content=$content, 
		  	ie.contentType=$contentType,
			ie.hide=$hide,
		  	ie.source=$source, 
		  	ie.sourceOfTruth=$sourceOfTruth,
		  	ie.appSource=$appSource
		 ON MATCH SET 
			ie.content=$content, 
		  	ie.contentType=$contentType
		 RETURN ie.id`

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		params := map[string]interface{}{
			"tenant":           tenant,
			"id":               event.Id,
			"content":          event.Content,
			"contentType":      event.ContentType,
			"externalSystemId": event.ExternalSystem,
			"createdAt":        utils.TimePtrFirstNonNilNillableAsAny(event.CreatedAt),
			"type":             event.Type,
			"identifier":       utils.FirstNotEmpty(event.Identifier, event.ExternalId),
			"source":           event.ExternalSystem,
			"sourceOfTruth":    event.ExternalSystem,
			"appSource":        constants.AppSourceSyncCustomerOsData,
			"externalId":       event.ExternalId,
			"externalUrl":      event.ExternalUrl,
			"externalSource":   event.ExternalSourceEntity,
			"syncDate":         syncDate,
			"channel":          event.Channel,
			"hide":             event.Hide,
		}

		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
			params)
		if err != nil {
			return nil, err
		}
		_, err = queryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *interactionEventRepository) MergeEmailInteractionSession(ctx context.Context, tenant string, syncDate time.Time, message entity.EmailMessageData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.MergeEmailInteractionSession")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MERGE (is:InteractionSession_%s {identifier:$identifier, channel:$channel}) " +
		" ON CREATE SET " +
		"  is:InteractionSession, " +
		"  is.id=randomUUID(), " +
		"  is.syncDate=$syncDate, " +
		"  is.externalUrl=$externalUrl, " +
		"  is.createdAt=$createdAt, " +
		"  is.name=$name, " +
		"  is.status=$status," +
		"  is.type=$type," +
		"  is.sourceOfTruth=$sourceOfTruth, " +
		"  is.appSource=$appSource " +
		" WITH is " +
		" RETURN is.id"

	dbRecord, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"source":        message.ExternalSystem,
				"sourceOfTruth": message.ExternalSystem,
				"appSource":     constants.AppSourceSyncCustomerOsData,
				"identifier":    message.EmailThreadId,
				"name":          message.Subject,
				"syncDate":      syncDate,
				"externalUrl":   message.ExternalUrl,
				"createdAt":     utils.TimePtrFirstNonNilNillableAsAny(message.CreatedAt),
				"status":        "ACTIVE",
				"type":          "THREAD",
				"channel":       "EMAIL",
			})
		if err != nil {
			return nil, err
		}
		record, err := queryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		return record, nil
	})
	if err != nil {
		return "", err
	}
	return dbRecord.(*db.Record).Values[0].(string), nil
}

func (r *interactionEventRepository) MergeEmailInteractionEvent(ctx context.Context, tenant, externalSystemId string, syncDate time.Time, message entity.EmailMessageData) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.MergeEmailInteractionEvent")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId}) " +
		" MERGE (ie:InteractionEvent_%s {source:$source, channel:$channel})-[rel:IS_LINKED_WITH {externalId:$externalId}]->(e) " +
		" ON CREATE SET " +
		"  ie:InteractionEvent, " +
		"  ie:TimelineEvent, " +
		"  ie:TimelineEvent_%s, " +
		"  rel.syncDate=$syncDate, " +
		"  rel.externalUrl=$externalUrl, " +
		"  ie.createdAt=$createdAt, " +
		"  ie.id=randomUUID(), " +
		"  ie.identifier=$identifier, " +
		"  ie.content=$content, " +
		"  ie.contentType=$contentType, " +
		"  ie.sourceOfTruth=$sourceOfTruth, " +
		"  ie.appSource=$appSource " +
		" WITH ie " +
		" RETURN ie.id"

	dbRecord, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		params := map[string]interface{}{
			"tenant":           tenant,
			"externalSystemId": externalSystemId,
			"identifier":       message.EmailMessageId,
			"source":           message.ExternalSystem,
			"sourceOfTruth":    message.ExternalSystem,
			"appSource":        constants.AppSourceSyncCustomerOsData,
			"externalId":       message.ExternalId,
			"externalUrl":      message.ExternalUrl,
			"syncDate":         syncDate,
			"createdAt":        utils.TimePtrFirstNonNilNillableAsAny(message.CreatedAt),
			"channel":          "EMAIL",
		}

		if message.Html != "" {
			params["content"] = message.Html
			params["contentType"] = "text/html"
		} else {
			params["content"] = message.Text
			params["contentType"] = "text/plain"
		}

		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
			params)
		if err != nil {
			return nil, err
		}
		record, err := queryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		return record, nil
	})
	if err != nil {
		return "", err
	}

	return dbRecord.(*db.Record).Values[0].(string), nil
}

func (r *interactionEventRepository) LinkInteractionEventToSession(ctx context.Context, tenant, interactionEventId, interactionSessionId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.LinkInteractionEventToSession")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (is:InteractionSession_%s {id:$interactionSessionId}) " +
		" MATCH (ie:InteractionEvent {id:$interactionEventId})" +
		" MERGE (ie)-[:PART_OF]->(is) "
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":               tenant,
				"interactionSessionId": interactionSessionId,
				"interactionEventId":   interactionEventId,
			})
		return nil, err
	})
	return err
}

func (r *interactionEventRepository) InteractionEventSentByEmail(ctx context.Context, tenant, interactionEventId, emailId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.InteractionEventSentByEmail")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (is:InteractionEvent_%s {id:$interactionEventId}) " +
		" MATCH (e:Email_%s {id: $emailId}) " +
		" MERGE (is)-[:SENT_BY]->(e) "
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
			map[string]interface{}{
				"tenant":             tenant,
				"interactionEventId": interactionEventId,
				"emailId":            emailId,
			})
		return nil, err
	})
	return err
}

func (r *interactionEventRepository) InteractionEventSentToEmails(ctx context.Context, tenant, interactionEventId, sentType string, emails []string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.InteractionEventSentToEmails")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (ie:InteractionEvent_%s {id:$interactionEventId}) " +
		" MATCH (e:Email_%s) WHERE e.rawEmail in $emails " +
		" MERGE (ie)-[:SENT_TO {type: $sentType}]->(e) "
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
			map[string]interface{}{
				"tenant":             tenant,
				"interactionEventId": interactionEventId,
				"sentType":           sentType,
				"emails":             emails,
			})
		return nil, err
	})

	return err
}

func (r *interactionEventRepository) LinkInteractionEventAsPartOf(ctx context.Context, tenant, interactionEventId, id, label string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.LinkInteractionEventAsPartOf")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (ie:InteractionEvent_%s {id:$interactionEventId}) 
		MATCH (linkedEntity:Issue|InteractionSession {id:$id}) 
		WHERE $label IN labels(linkedEntity) AND $label+'_%s' IN labels(linkedEntity)
		MERGE (ie)-[rel:PART_OF]->(linkedEntity)
		RETURN rel`, tenant, tenant)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":             tenant,
				"interactionEventId": interactionEventId,
				"id":                 id,
				"label":              label,
			})
		if err != nil {
			return nil, err
		}
		_, err = queryResult.Single(ctx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}

func (r *interactionEventRepository) MergeInteractionSessionForEvent(ctx context.Context, tenant, eventId, externalSource string, syncDate time.Time, session entity.InteractionSession) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.MergeInteractionSessionForEvent")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId}) 
								MATCH (ie:InteractionEvent_%s {id:$interactionEventId})
		 						MERGE (is:InteractionSession_%s)-[rel:IS_LINKED_WITH {externalId:$externalId}]->(e)
								ON CREATE SET
									is:InteractionSession,
									is.id=randomUUID(),
									rel.syncDate=$syncDate,
									rel.externalUrl=$externalUrl,
									is.createdAt=$createdAt,
									is.updatedAt=$createdAt,
									is.source=$source,
									is.sourceOfTruth=$sourceOfTruth,
									is.appSource=$appSource,
									is.identifier=$identifier,
									is.status=$status,
									is.type=$type,
									is.channel=$channel,
									is.name=$name
		WITH is, ie
		MERGE (ie)-[r:PART_OF]->(is)`, tenant, tenant)

	return utils.ExecuteWriteQuery(ctx, *r.driver, query, map[string]interface{}{
		"tenant":             tenant,
		"interactionEventId": eventId,
		"externalId":         session.ExternalId,
		"identifier":         utils.FirstNotEmpty(session.Identifier, session.ExternalId),
		"externalSystemId":   externalSource,
		"source":             externalSource,
		"sourceOfTruth":      externalSource,
		"appSource":          constants.AppSourceSyncCustomerOsData,
		"syncDate":           syncDate,
		"externalUrl":        session.ExternalUrl,
		"createdAt":          utils.TimePtrFirstNonNilNillableAsAny(session.CreatedAt),
		"channel":            session.Channel,
		"type":               session.Type,
		"status":             session.Status,
		"name":               session.Name,
	})
}

func (r *interactionEventRepository) LinkInteractionEventWithSenderById(ctx context.Context, tenant, eventId, id, label string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.LinkInteractionEventWithSenderById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (ie:InteractionEvent_%s {id:$interactionEventId}) 
		MATCH (:Tenant {name:$tenant})-[*1..2]-(sender:Contact|User|Organization|JobRole|Email|PhoneNumber {id:$id}) 
		WHERE $label IN labels(sender)
		MERGE (ie)-[:SENT_BY]->(sender)`, tenant)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query, map[string]interface{}{
			"tenant":             tenant,
			"interactionEventId": eventId,
			"id":                 id,
			"label":              label,
		})
		return nil, err
	})
	return err
}

func (r *interactionEventRepository) LinkInteractionEventWithRecipientById(ctx context.Context, tenant, eventId, id, label, relationType string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.LinkInteractionEventWithReceiverById")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (ie:InteractionEvent_%s {id:$interactionEventId}) 
		MATCH (:Tenant {name:$tenant})--(sender:Contact|User|Organization|JobRole|Email|PhoneNumber {id:$id}) 
		WHERE $label IN labels(sender)
		MERGE (ie)-[rel:SENT_TO]->(sender)
		ON CREATE SET rel.type=$relationType`, tenant)
	span.LogFields(log.String("query", query))

	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, query, map[string]interface{}{
			"tenant":             tenant,
			"interactionEventId": eventId,
			"id":                 id,
			"label":              label,
			"relationType":       relationType,
		})
		return nil, err
	})
	return err
}

func (r *interactionEventRepository) GetInteractionSessionIdByExternalId(ctx context.Context, tenant, externalId, externalSystemId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.GetInteractionSessionIdByExternalId")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId})
					MATCH (is:InteractionSession_%s)-[:IS_LINKED_WITH {externalId:$externalId}]->(e)
				return is.id order by is.createdAt`, tenant)

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	records, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, query, map[string]any{
			"tenant":           tenant,
			"externalId":       externalId,
			"externalSystemId": externalSystemId,
		})
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
