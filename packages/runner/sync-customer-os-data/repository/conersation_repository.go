package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	"time"
)

type ConversationRepository interface {
	MergeEmailConversation(tenant string, date time.Time, message entity.EmailMessageData) (string, int64, error)
	UserInitiateConversation(tenant, conversationId, userExternalI, externalSystem string) error
	ContactInitiateConversation(tenant, conversationId, contactId string) error
	ContactByIdParticipateInConversation(tenant, conversationId, contactId string) error
	ContactsByExternalIdParticipateInConversation(tenant, conversationId, externalSystem string, contactExternalIds []string) error
	UserByExternalIdParticipateInConversation(tenant, conversationId, externalSystem, userExternalId string) error
	UsersByEmailParticipateInConversation(tenant, conversationId string, userEmails []string) error
	IncrementMessageCount(tenant, conversationId string) error
}

type conversationRepository struct {
	driver *neo4j.Driver
}

func NewConversationRepository(driver *neo4j.Driver) ConversationRepository {
	return &conversationRepository{
		driver: driver,
	}
}

func (r *conversationRepository) MergeEmailConversation(tenant string, syncDate time.Time, message entity.EmailMessageData) (string, int64, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MERGE (c:Conversation_%s {threadId:$threadId, source:$source, channel:$channel}) " +
		" ON CREATE SET c:Conversation, " +
		" 				c.syncDate=$syncDate, c.id=randomUUID(), c.startedAt=$createdAt, " +
		"             	c.sourceOfTruth=$sourceOfTruth, c.appSource=$appSource, c.status=$status," +
		"				c.messageCount=0 " +
		" ON MATCH SET 	c.syncDate=$syncDate, c.status=$status " +
		" WITH c " +
		" REMOVE c.endedAt " +
		" RETURN c.id, c.messageCount"

	dbRecord, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":        tenant,
				"source":        message.ExternalSystem,
				"sourceOfTruth": message.ExternalSystem,
				"appSource":     message.ExternalSystem,
				"threadId":      message.EmailThreadId,
				"syncDate":      syncDate,
				"createdAt":     message.CreatedAt,
				"status":        "ACTIVE",
				"channel":       "EMAIL",
			})
		if err != nil {
			return nil, err
		}
		record, err := queryResult.Single()
		if err != nil {
			return nil, err
		}
		return record, nil
	})
	if err != nil {
		return "", 0, err
	}
	return dbRecord.(*db.Record).Values[0].(string), dbRecord.(*db.Record).Values[1].(int64), nil
}

func (r *conversationRepository) UserInitiateConversation(tenant, conversationId, userExternalId, externalSystem string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (o:Conversation_%s {id:$conversationId}) " +
		" MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem}) " +
		" MATCH (u:User)-[:IS_LINKED_WITH {externalId:$userExternalId}]->(e)" +
		" MERGE (u)-[:INITIATED]->(o) "
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"conversationId": conversationId,
				"userExternalId": userExternalId,
				"externalSystem": externalSystem,
			})
		return nil, err
	})
	return err
}

func (r *conversationRepository) ContactInitiateConversation(tenant, conversationId, contactId string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (o:Conversation_%s {id:$conversationId}) " +
		" MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId}) " +
		" MERGE (c)-[:INITIATED]->(o) "
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"conversationId": conversationId,
				"contactId":      contactId,
			})
		return nil, err
	})
	return err
}

func (r *conversationRepository) ContactByIdParticipateInConversation(tenant, conversationId, contactId string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (o:Conversation_%s {id:$conversationId}) " +
		" MATCH (:Tenant {name:$tenant})<-[:CONTACT_BELONGS_TO_TENANT]-(c:Contact {id:$contactId}) " +
		" MERGE (c)-[:PARTICIPATES]->(o) "
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"conversationId": conversationId,
				"contactId":      contactId,
			})
		return nil, err
	})
	return err
}

func (r *conversationRepository) ContactsByExternalIdParticipateInConversation(tenant, conversationId, externalSystem string, contactExternalIds []string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (o:Conversation_%s {id:$conversationId}) " +
		" MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem}) " +
		" MATCH (c:Contact)-[r:IS_LINKED_WITH]->(e) WHERE r.externalId in $contactExternalIds " +
		" MERGE (c)-[:PARTICIPATES]->(o)  "
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, tenant),
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

func (r *conversationRepository) UserByExternalIdParticipateInConversation(tenant, conversationId, externalSystem, userExternalId string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (o:Conversation_%s {id:$conversationId}) " +
		" MATCH (:Tenant {name:$tenant})<-[:EXTERNAL_SYSTEM_BELONGS_TO_TENANT]-(e:ExternalSystem {id:$externalSystem}) " +
		" MATCH (u:User)-[:IS_LINKED_WITH {externalId:$userExternalId}]->(e) " +
		" MERGE (u)-[:PARTICIPATES]->(o)  "
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, tenant),
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

func (r *conversationRepository) UsersByEmailParticipateInConversation(tenant, conversationId string, userEmails []string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (o:Conversation_%s {id:$conversationId}) " +
		" MATCH (:Tenant {name:$tenant})<-[:USER_BELONGS_TO_TENANT]-(u:User) WHERE u.email in $userEmails " +
		" MERGE (u)-[:PARTICIPATES]->(o)  "
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"conversationId": conversationId,
				"userEmails":     userEmails,
			})
		return nil, err
	})
	return err
}

func (r *conversationRepository) IncrementMessageCount(tenant, conversationId string) error {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	query := "MATCH (o:Conversation_%s {id:$conversationId}) " +
		" SET o.messageCount=o.messageCount+1 "
	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		_, err := tx.Run(fmt.Sprintf(query, tenant),
			map[string]interface{}{
				"tenant":         tenant,
				"conversationId": conversationId,
			})
		return nil, err
	})
	return err
}
