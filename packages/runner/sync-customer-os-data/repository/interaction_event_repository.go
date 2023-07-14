package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"time"
)

type InteractionEventRepository interface {
	GetMatchedInteractionEvent(ctx context.Context, tenant string, event entity.InteractionEventData) (string, error)
	MergeInteractionEvent(ctx context.Context, tenant string, syncDate time.Time, event entity.InteractionEventData) error

	MergeInteractionSession(ctx context.Context, tenant string, date time.Time, message entity.EmailMessageData) (string, error)
	MergeEmailInteractionEvent(ctx context.Context, tenant, externalSystemId string, date time.Time, message entity.EmailMessageData) (string, error)
	LinkInteractionEventToSession(ctx context.Context, tenant, interactionEventId, interactionSessionId string) error

	InteractionEventSentByEmail(ctx context.Context, tenant, interactionEventId, emailId string) error
	InteractionEventSentToEmails(ctx context.Context, tenant, interactionEventId, sentType string, emails []string) error
	LinkInteractionEventAsPartOfByExternalId(ctx context.Context, tenant string, event entity.InteractionEventData) error
	LinkInteractionEventWithSenderByExternalId(ctx context.Context, tenant, eventId, externalSystem string, sender entity.InteractionEventParticipant) error
	LinkInteractionEventWithRecipientByExternalId(ctx context.Context, tenant, eventId, externalSystem string, recipient entity.InteractionEventParticipant) error
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
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (t:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystem})
				OPTIONAL MATCH (ext)<-[:IS_LINKED_WITH {externalId:$issueExternalId}]-(ie:InteractionEvent_%s)
				WITH ie WHERE ie is not null
				return ie.id limit 1`

	dbRecords, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":          tenant,
				"externalSystem":  event.ExternalSystem,
				"issueExternalId": event.ExternalId,
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
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId}) 
		 MERGE (ie:InteractionEvent_%s {id:$id})-[rel:IS_LINKED_WITH {externalId:$externalId}]->(e) 
		 ON CREATE SET 
		  	ie:InteractionEvent,
		  	ie:TimelineEvent, 
		  	ie:TimelineEvent_%s, 
		  	rel.syncDate=$syncDate, 
		  	ie.createdAt=$createdAt,
			ie.channel=$channel,
			ie.eventType=$type, 
		  	ie.identifier=$identifier, 
		  	ie.content=$content, 
		  	ie.contentType=$contentType,
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
			"createdAt":        event.CreatedAt,
			"type":             event.Type,
			"identifier":       event.ExternalId,
			"source":           event.ExternalSystem,
			"sourceOfTruth":    event.ExternalSystem,
			"appSource":        event.ExternalSystem,
			"externalId":       event.ExternalId,
			"syncDate":         syncDate,
			"channel":          event.Channel,
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

func (r *interactionEventRepository) MergeInteractionSession(ctx context.Context, tenant string, syncDate time.Time, message entity.EmailMessageData) (string, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MERGE (is:InteractionSession_%s {identifier:$identifier, source:$source, channel:$channel}) " +
		" ON CREATE SET " +
		"  is:InteractionSession, " +
		"  is.id=randomUUID(), " +
		"  is.syncDate=$syncDate, " +
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
				"appSource":     message.ExternalSystem,
				"identifier":    message.EmailThreadId,
				"name":          message.Subject,
				"syncDate":      syncDate,
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
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystemId}) " +
		" MERGE (ie:InteractionEvent_%s {source:$source, channel:$channel})-[rel:IS_LINKED_WITH {externalId:$externalId}]->(e) " +
		" ON CREATE SET " +
		"  ie:InteractionEvent, " +
		"  ie:TimelineEvent, " +
		"  ie:TimelineEvent_%s, " +
		"  rel.syncDate=$syncDate, " +
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
			"appSource":        message.ExternalSystem,
			"externalId":       message.ExternalId,
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

func (r *interactionEventRepository) LinkInteractionEventAsPartOfByExternalId(ctx context.Context, tenant string, event entity.InteractionEventData) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (ie:InteractionEvent_%s {id:$interactionEventId}) 
		MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystemId})<-[:IS_LINKED_WITH {externalId:$partOfExternalId}]-(n) 
		WHERE n:Issue
		MERGE (ie)-[result:PART_OF]->(n) 
		return result`
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":             tenant,
				"interactionEventId": event.Id,
				"partOfExternalId":   event.PartOfExternalId,
				"externalSystemId":   event.ExternalSystem,
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

func (r *interactionEventRepository) LinkInteractionEventWithSenderByExternalId(ctx context.Context, tenant, eventId, externalSystem string, sender entity.InteractionEventParticipant) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (ie:InteractionEvent_%s {id:$interactionEventId}) 
		MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystemId})<-[:IS_LINKED_WITH {externalId:$sentByExternalId}]-(n) 
		WHERE $nodeLabel in labels(n)
		MERGE (ie)-[result:SENT_BY]->(n) 
		return result`
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":             tenant,
				"interactionEventId": eventId,
				"sentByExternalId":   sender.ExternalId,
				"nodeLabel":          sender.GetNodeLabel(),
				"externalSystemId":   externalSystem,
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

func (r *interactionEventRepository) LinkInteractionEventWithRecipientByExternalId(ctx context.Context, tenant, eventId, externalSystem string, recipient entity.InteractionEventParticipant) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := `MATCH (ie:InteractionEvent_%s {id:$interactionEventId}) 
		MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(ext:ExternalSystem {id:$externalSystemId})<-[:IS_LINKED_WITH {externalId:$sentByExternalId}]-(n) 
		WHERE $nodeLabel in labels(n)
		MERGE (ie)-[result:SENT_TO]->(n)
		ON CREATE SET result.type=$relationType
		return result`
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":             tenant,
				"interactionEventId": eventId,
				"sentByExternalId":   recipient.ExternalId,
				"nodeLabel":          recipient.GetNodeLabel(),
				"relationType":       recipient.RelationType,
				"externalSystemId":   externalSystem,
			})
		return nil, err
	})
	return err
}
