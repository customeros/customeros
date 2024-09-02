package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

type SendDirection string

const (
	SENT_TO SendDirection = "SENT_TO"
	SENT_BY SendDirection = "SENT_BY"
)

type PartOfType string

const (
	PART_OF_INTERACTION_SESSION PartOfType = "InteractionSession"
	PART_OF_MEETING             PartOfType = "Meeting"
)

type InteractionEventRepository interface {
	GetSentByParticipantsForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeWithRelationAndId, error)
	GetSentToParticipantsForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeWithRelationAndId, error)
	GetReplyToInteractionEventsForInteractionEvents(ctx context.Context, tenant string, ids []string, returnContent bool) ([]*utils.DbPropsAndId, error)

	// Deprecated, use events-platform
	LinkWithPartOfXXInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, interactionEventId string, partOfId string, partOfType PartOfType) error
	// Deprecated, use events-platform
	LinkWithRepliesToInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, interactionEventId, repliesToEventId string) error
	// Deprecated, use events-platform
	LinkWithSentXXParticipantInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entityType model.EntityType, interactionEventId, participantId string, sentType *string, direction SendDirection) error
	// Deprecated, use events-platform
	LinkWithSentXXEmailInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, interactionEventId, email string, sentType *string, direction SendDirection) error
	// Deprecated, use events-platform
	LinkWithSentXXPhoneNumberInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, interactionEventId, e164 string, sentType *string, direction SendDirection) error

	Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, newInteractionEvent entity.InteractionEventEntity, source, sourceOfTruth neo4jentity.DataSource) (*dbtype.Node, error)
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

func (r *interactionEventRepository) LinkWithSentXXEmailInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, interactionEventId, email string, sentType *string, direction SendDirection) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.LinkWithSentXXEmailInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (e:Email_%s) `, tenant)
	query += fmt.Sprintf(`MATCH (ie:InteractionEvent_%s {id:$eventId}) `, tenant)
	query += `WHERE e.email = $email OR e.rawEmail = $email `

	if direction == SENT_TO {
		if sentType != nil {
			query += fmt.Sprintf(`MERGE (ie)-[r:SENT_TO {type:$sentType}]->(e) RETURN r`)
		} else {
			query += fmt.Sprintf(`MERGE (ie)-[r:SENT_TO]->(e) RETURN r`)
		}
	} else {
		if sentType != nil {
			query += fmt.Sprintf(`MERGE (ie)-[r:SENT_BY {type:$sentType}]->(e) RETURN r`)
		} else {
			query += fmt.Sprintf(`MERGE (ie)-[r:SENT_BY]->(e) RETURN r`)
		}
	}
	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"email":    email,
			"eventId":  interactionEventId,
			"sentType": sentType,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *interactionEventRepository) LinkWithSentXXPhoneNumberInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, interactionEventId, e164 string, sentType *string, direction SendDirection) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.LinkWithSentXXPhoneNumberInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := fmt.Sprintf(`MATCH (p:PhoneNumber_%s) `, tenant)
	query += fmt.Sprintf(`MATCH (ie:InteractionEvent_%s {id:$eventId}) `, tenant)
	query += `WHERE p.e164 = $e164 OR p.rawPhoneNumber = $e164 `

	if direction == SENT_TO {
		if sentType != nil {
			query += fmt.Sprintf(`MERGE (ie)-[r:SENT_TO {type:$sentType}]->(p) RETURN r`)
		} else {
			query += fmt.Sprintf(`MERGE (ie)-[r:SENT_TO]->(p) RETURN r`)
		}
	} else {
		if sentType != nil {
			query += fmt.Sprintf(`MERGE (ie)-[r:SENT_BY {type:$sentType}]->(p) RETURN r`)
		} else {
			query += fmt.Sprintf(`MERGE (ie)-[r:SENT_BY]->(p) RETURN r`)
		}
	}
	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"e164":    e164,
			"eventId": interactionEventId,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *interactionEventRepository) LinkWithSentXXParticipantInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, entityType model.EntityType, interactionEventId, participantId string, sentType *string, direction SendDirection) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.LinkWithSentXXParticipantInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	query := ""
	switch entityType {
	case model.CONTACT:
		query = fmt.Sprintf(`MATCH (p:Contact_%s {id:$participantId}) `, tenant)
	case model.USER:
		query = fmt.Sprintf(`MATCH (p:User_%s {id:$participantId}) `, tenant)
	}
	query += fmt.Sprintf(`MATCH (ie:InteractionEvent_%s {id:$eventId}) `, tenant)

	if direction == SENT_TO {
		if sentType != nil {
			query += fmt.Sprintf(`MERGE (ie)-[r:SENT_TO {type:$sentType}]->(p) RETURN r`)
		} else {
			query += fmt.Sprintf(`MERGE (ie)-[r:SENT_TO]->(p) RETURN r`)
		}
	} else {
		if sentType != nil {
			query += fmt.Sprintf(`MERGE (ie)-[r:SENT_BY {type:$sentType}]->(p) RETURN r`)
		} else {
			query += fmt.Sprintf(`MERGE (ie)-[r:SENT_BY]->(p) RETURN r`)
		}
	}
	queryResult, err := tx.Run(ctx, query,
		map[string]any{
			"participantId": participantId,
			"eventId":       interactionEventId,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *interactionEventRepository) LinkWithRepliesToInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant, interactionEventId, repliesToEventId string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.LinkWithRepliesToInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	queryResult, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (ie:InteractionEvent_%s {id:$eventId})
			MATCH (rie:InteractionEvent_%s {id:$repliesToEventId})
			MERGE (ie)-[r:REPLIES_TO]->(rie)
			RETURN r`, tenant, tenant),
		map[string]any{
			"eventId":          interactionEventId,
			"repliesToEventId": repliesToEventId,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *interactionEventRepository) LinkWithPartOfXXInTx(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, interactionEventId string, partOfId string, partOfType PartOfType) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.LinkWithPartOfXXInTx")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	queryResult, err := tx.Run(ctx, fmt.Sprintf(`
			MATCH (ie:InteractionEvent_%s {id:$eventId})
			MATCH (is:%s_%s {id:$interactionSessionId})
			MERGE (ie)-[r:PART_OF]->(is)
			RETURN r`, tenant, partOfType, tenant),
		map[string]any{
			"eventId":              interactionEventId,
			"interactionSessionId": partOfId,
		})
	if err != nil {
		return err
	}
	_, err = queryResult.Single(ctx)
	return err
}

func (r *interactionEventRepository) GetSentByParticipantsForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.GetSentByParticipantsForInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (ie:InteractionEvent)-[rel:SENT_BY]->(p:Email|PhoneNumber|User|Contact|Organization|JobRole) 
		WHERE ie.id IN $ids AND ie:InteractionEvent_%s
		RETURN p, rel, ie.id`, tenant)
	params := map[string]any{
		"ids": ids,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeWithRelationAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}

func (r *interactionEventRepository) Create(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, newInteractionEvent entity.InteractionEventEntity, source, sourceOfTruth neo4jentity.DataSource) (*dbtype.Node, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.Create")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	var createdAt time.Time
	createdAt = utils.Now()
	if newInteractionEvent.CreatedAt != nil {
		createdAt = *newInteractionEvent.CreatedAt
	}

	query := "MERGE (ie:InteractionEvent_%s {id:randomUUID()}) ON CREATE SET " +
		"  ie:InteractionEvent, " +
		"  ie:TimelineEvent, " +
		"  ie:TimelineEvent_%s, " +
		" ie.source=$source, " +
		" ie.channel=$channel, " +
		" ie.channelData=$channelData, " +
		" ie.createdAt=$createdAt, " +
		" ie.identifier=$identifier, " +
		" ie.customerOSInternalIdentifier=$customerOSInternalIdentifier, " +
		" ie.content=$content, " +
		" ie.contentType=$contentType, " +
		" ie.eventType=$eventType, " +
		" ie.sourceOfTruth=$sourceOfTruth, " +
		" ie.appSource=$appSource " +
		" RETURN ie"

	if queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant, tenant),
		map[string]interface{}{
			"tenant":                       tenant,
			"source":                       source,
			"channel":                      newInteractionEvent.Channel,
			"channelData":                  newInteractionEvent.ChannelData,
			"createdAt":                    createdAt,
			"identifier":                   newInteractionEvent.EventIdentifier,
			"customerOSInternalIdentifier": newInteractionEvent.CustomerOSInternalIdentifier,
			"content":                      newInteractionEvent.Content,
			"contentType":                  newInteractionEvent.ContentType,
			"eventType":                    newInteractionEvent.EventType,
			"sourceOfTruth":                sourceOfTruth,
			"appSource":                    newInteractionEvent.AppSource,
		}); err != nil {
		return nil, err
	} else {
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}
}

func (r *interactionEventRepository) GetReplyToInteractionEventsForInteractionEvents(ctx context.Context, tenant string, ids []string, returnContent bool) ([]*utils.DbPropsAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.GetReplyToInteractionEventsForInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)
	span.LogFields(log.Bool("returnContent", returnContent))

	cypherReturnFragment := "rie {.*}"
	if !returnContent {
		cypherReturnFragment = "rie {.*, content: ''}"
	}

	cypher := fmt.Sprintf(`MATCH (ie:InteractionEvent)-[rel:REPLIES_TO]->(rie:InteractionEvent_%s) 
		 	WHERE ie.id IN $ids AND ie:InteractionEvent_%s 
			RETURN %s, ie.id`, tenant, tenant, cypherReturnFragment)
	params := map[string]any{
		"ids": ids,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbPropsAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbPropsAndId), err
}

func (r *interactionEventRepository) GetSentToParticipantsForInteractionEvents(ctx context.Context, tenant string, ids []string) ([]*utils.DbNodeWithRelationAndId, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "InteractionEventRepository.GetSentToParticipantsForInteractionEvents")
	defer span.Finish()
	tracing.SetDefaultNeo4jRepositorySpanTags(ctx, span)

	cypher := fmt.Sprintf(`MATCH (ie:InteractionEvent)-[rel:SENT_TO]->(p:Email|PhoneNumber|User|Contact|Organization|JobRole) 
		 WHERE ie.id IN $ids AND ie:InteractionEvent_%s RETURN p, rel, ie.id`, tenant)
	params := map[string]any{
		"ids": ids,
	}
	span.LogFields(log.String("cypher", cypher), log.Object("params", params))

	session := utils.NewNeo4jReadSession(ctx, *r.driver)
	defer session.Close(ctx)
	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		if queryResult, err := tx.Run(ctx, cypher, params); err != nil {
			return nil, err
		} else {
			return utils.ExtractAllRecordsAsDbNodeWithRelationAndId(ctx, queryResult, err)
		}
	})
	if err != nil {
		return nil, err
	}
	return result.([]*utils.DbNodeWithRelationAndId), err
}
