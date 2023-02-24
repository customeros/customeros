package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"golang.org/x/net/context"
	"time"
)

type ConversationRepository interface {
	MergeEmailConversation(ctx context.Context, tenant string, date time.Time, message entity.EmailMessageData) (string, int64, string, error)
	UserInitiateConversation(ctx context.Context, tenant, conversationId string, initiator entity.ConversationInitiator) error
	ContactInitiateConversation(ctx context.Context, tenant, conversationId string, initiator entity.ConversationInitiator) error
	ContactByIdParticipateInConversation(ctx context.Context, tenant, conversationId, contactId string) error
	ContactsByExternalIdParticipateInConversation(ctx context.Context, tenant, conversationId, externalSystem string, contactExternalIds []string) error
	UserByExternalIdParticipateInConversation(ctx context.Context, tenant, conversationId, externalSystem, userExternalId string) error
	UsersByEmailParticipateInConversation(ctx context.Context, tenant, conversationId string, userEmails []string) error
	IncrementMessageCount(ctx context.Context, tenant, conversationId string, updatedAt time.Time) error
}

type conversationRepository struct {
	driver *neo4j.DriverWithContext
}

func NewConversationRepository(driver *neo4j.DriverWithContext) ConversationRepository {
	return &conversationRepository{
		driver: driver,
	}
}

func (r *conversationRepository) MergeEmailConversation(ctx context.Context, tenant string, syncDate time.Time, message entity.EmailMessageData) (string, int64, string, error) {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MERGE (o:Conversation_%s {threadId:$threadId, source:$source, channel:$channel}) " +
		" ON CREATE SET " +
		"  o:Conversation, " +
		"  o.syncDate=$syncDate, " +
		"  o.id=randomUUID(), " +
		"  o.subject=$subject, " +
		"  o.startedAt=$createdAt, " +
		"  o.sourceOfTruth=$sourceOfTruth, " +
		"  o.appSource=$appSource, " +
		"  o.status=$status," +
		"  o.messageCount=0 " +
		" ON MATCH SET 	" +
		"  o.syncDate=$syncDate, " +
		"  o.status=$status " +
		" WITH o " +
		" REMOVE o.endedAt " +
		" RETURN o.id, o.messageCount, coalesce(o.initiatorUsername, $emptyInitiator) "

	dbRecord, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"subject":        message.Subject,
				"source":         message.ExternalSystem,
				"sourceOfTruth":  message.ExternalSystem,
				"appSource":      message.ExternalSystem,
				"threadId":       message.EmailThreadId,
				"syncDate":       syncDate,
				"createdAt":      message.CreatedAt,
				"status":         "ACTIVE",
				"channel":        "EMAIL",
				"emptyInitiator": "",
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
		return "", 0, "", err
	}
	return dbRecord.(*db.Record).Values[0].(string), dbRecord.(*db.Record).Values[1].(int64), dbRecord.(*db.Record).Values[2].(string), nil
}

func (r *conversationRepository) UserInitiateConversation(ctx context.Context, tenant, conversationId string, initiator entity.ConversationInitiator) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (o:Conversation_%s {id:$conversationId}) " +
		" MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem}) " +
		" MATCH (u:User)-[:IS_LINKED_WITH {externalId:$userExternalId}]->(e)" +
		" MERGE (u)-[:INITIATED]->(o) " +
		" SET o.initiatorFirstName=$firstName, o.initiatorLastName=$lastName, o.initiatorUsername=$email, o.initiatorType=$initiatorType "
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"conversationId": conversationId,
				"userExternalId": initiator.ExternalId,
				"externalSystem": initiator.ExternalSystem,
				"firstName":      initiator.FirstName,
				"lastName":       initiator.LastName,
				"email":          initiator.Email,
				"initiatorType":  initiator.InitiatorType,
			})
		return nil, err
	})
	return err
}

func (r *conversationRepository) ContactInitiateConversation(ctx context.Context, tenant, conversationId string, initiator entity.ConversationInitiator) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (o:Conversation_%s {id:$conversationId}) " +
		" MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId}) " +
		" MERGE (c)-[:INITIATED]->(o) " +
		" SET o.initiatorFirstName=$firstName, o.initiatorLastName=$lastName, o.initiatorUsername=$email, o.initiatorType=$initiatorType "
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"conversationId": conversationId,
				"contactId":      initiator.Id,
				"firstName":      initiator.FirstName,
				"lastName":       initiator.LastName,
				"email":          initiator.Email,
				"initiatorType":  initiator.InitiatorType,
			})
		return nil, err
	})
	return err
}

func (r *conversationRepository) ContactByIdParticipateInConversation(ctx context.Context, tenant, conversationId, contactId string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (o:Conversation_%s {id:$conversationId}) " +
		" MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId}) " +
		" MERGE (c)-[:PARTICIPATES]->(o) "
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"conversationId": conversationId,
				"contactId":      contactId,
			})
		return nil, err
	})
	return err
}

func (r *conversationRepository) ContactsByExternalIdParticipateInConversation(ctx context.Context, tenant, conversationId, externalSystem string, contactExternalIds []string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (o:Conversation_%s {id:$conversationId}) " +
		" MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem}) " +
		" MATCH (c:Contact)-[r:IS_LINKED_WITH]->(e) WHERE r.externalId in $contactExternalIds " +
		" MERGE (c)-[:PARTICIPATES]->(o)  "
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":             tenant,
				"conversationId":     conversationId,
				"externalSystem":     externalSystem,
				"contactExternalIds": contactExternalIds,
			})
		return nil, err
	})
	return err
}

func (r *conversationRepository) UserByExternalIdParticipateInConversation(ctx context.Context, tenant, conversationId, externalSystem, userExternalId string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (o:Conversation_%s {id:$conversationId}) " +
		" MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem}) " +
		" MATCH (u:User)-[:IS_LINKED_WITH {externalId:$userExternalId}]->(e) " +
		" MERGE (u)-[:PARTICIPATES]->(o)  "
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"conversationId": conversationId,
				"externalSystem": externalSystem,
				"userExternalId": userExternalId,
			})
		return nil, err
	})
	return err
}

func (r *conversationRepository) UsersByEmailParticipateInConversation(ctx context.Context, tenant, conversationId string, userEmails []string) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (o:Conversation_%s {id:$conversationId}) " +
		" MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User) WHERE u.email in $userEmails " +
		" MERGE (u)-[:PARTICIPATES]->(o)  "
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"conversationId": conversationId,
				"userEmails":     userEmails,
			})
		return nil, err
	})
	return err
}

func (r *conversationRepository) IncrementMessageCount(ctx context.Context, tenant, conversationId string, updatedAt time.Time) error {
	session := utils.NewNeo4jWriteSession(ctx, *r.driver)
	defer session.Close(ctx)

	query := "MATCH (o:Conversation_%s {id:$conversationId}) " +
		" SET o.messageCount=o.messageCount+1, o.updatedAt=$updatedAt "
	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"conversationId": conversationId,
				"updatedAt":      updatedAt,
			})
		return nil, err
	})
	return err
}
