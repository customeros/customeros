package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type ConversationDbNodeWithParticipantIDs struct {
	Node      *dbtype.Node
	UserId    string
	ContactId string
}

type ConversationDbNodesWithTotalCount struct {
	Nodes []*ConversationDbNodeWithParticipantIDs
	Count int64
}

type ConversationRepository interface {
	Create(tenant string, userIds, contactIds []string, entity entity.ConversationEntity) (any, error)
	GetPaginatedConversationsForUser(session neo4j.Session, tenant, userId string, skip, limit int, sort *utils.CypherSort) (*ConversationDbNodesWithTotalCount, error)
	GetPaginatedConversationsForContact(session neo4j.Session, tenant, contactId string, skip, limit int, sort *utils.CypherSort) (*ConversationDbNodesWithTotalCount, error)
}

type conversationRepository struct {
	driver *neo4j.Driver
}

func NewConversationRepository(driver *neo4j.Driver) ConversationRepository {
	return &conversationRepository{
		driver: driver,
	}
}

func (r *conversationRepository) Create(tenant string, userIds, contactIds []string, entity entity.ConversationEntity) (any, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	return session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		query := "MATCH (t:Tenant {name:$tenant}) " +
			" MERGE (o:Conversation {id:$conversationId}) " +
			" ON CREATE SET o.startedAt=$startedAt, o.itemCount=0, o.channel=$channel, o.status=$status, o:%s " +
			" %s %s " +
			" RETURN DISTINCT o"
		queryLinkWithContacts := ""
		if len(contactIds) > 0 {
			queryLinkWithContacts = " WITH DISTINCT t, o " +
				" OPTIONAL MATCH (c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(t) WHERE c.id in $contactIds " +
				" MERGE (c)-[:PARTICIPATES]->(o) "
		}
		queryLinkWithUsers := ""
		if len(userIds) > 0 {
			queryLinkWithUsers = " WITH DISTINCT t, o " +
				" OPTIONAL MATCH (u:User)-[:USER_BELONGS_TO_TENANT]->(t) WHERE u.id in $userIds " +
				" MERGE (u)-[:PARTICIPATES]->(o) "
		}
		txResult, err := tx.Run(fmt.Sprintf(query, "Conversation_"+tenant, queryLinkWithContacts, queryLinkWithUsers),
			map[string]interface{}{
				"tenant":         tenant,
				"status":         entity.Status,
				"startedAt":      entity.StartedAt,
				"channel":        entity.Channel,
				"conversationId": entity.Id,
				"contactIds":     contactIds,
				"userIds":        userIds,
			})
		record, err := txResult.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
}

func (r *conversationRepository) GetPaginatedConversationsForUser(session neo4j.Session, tenant, userId string, skip, limit int, sort *utils.CypherSort) (*ConversationDbNodesWithTotalCount, error) {
	result := new(ConversationDbNodesWithTotalCount)

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), 
											(u)-[:PARTICIPATES]->(o:Conversation)
											RETURN count(o) as count`,
			map[string]any{
				"tenant": tenant,
				"userId": userId,
			})
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single()
		result.Count = count.Values[0].(int64)

		queryResult, err = tx.Run(fmt.Sprintf(
			"MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), "+
				" (u)-[:PARTICIPATES]->(o:Conversation)<-[:PARTICIPATES]-(c:Contact) "+
				" RETURN o, u.id, c.id "+
				" %s "+
				" SKIP $skip LIMIT $limit", sort.SortingCypherFragment("o")),
			map[string]any{
				"tenant": tenant,
				"userId": userId,
				"skip":   skip,
				"limit":  limit,
			})
		return queryResult.Collect()
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		conversationWithParticipantIDs := new(ConversationDbNodeWithParticipantIDs)
		conversationWithParticipantIDs.Node = utils.NodePtr(v.Values[0].(neo4j.Node))
		conversationWithParticipantIDs.UserId = v.Values[1].(string)
		conversationWithParticipantIDs.ContactId = v.Values[2].(string)
		result.Nodes = append(result.Nodes, conversationWithParticipantIDs)
	}
	return result, nil
}

func (r *conversationRepository) GetPaginatedConversationsForContact(session neo4j.Session, tenant, contactId string, skip, limit int, sort *utils.CypherSort) (*ConversationDbNodesWithTotalCount, error) {
	result := new(ConversationDbNodesWithTotalCount)

	dbRecords, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), 
											(c)-[:PARTICIPATES]->(o:Conversation)
											RETURN count(o) as count`,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
			})
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single()
		result.Count = count.Values[0].(int64)

		queryResult, err = tx.Run(fmt.Sprintf(
			"MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), "+
				" (c)-[:PARTICIPATES]->(o:Conversation)<-[:PARTICIPATES]-(u:User) "+
				" RETURN o, u.id, c.id "+
				" %s "+
				" SKIP $skip LIMIT $limit", sort.SortingCypherFragment("o")),
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
				"skip":      skip,
				"limit":     limit,
			})
		return queryResult.Collect()
	})
	if err != nil {
		return nil, err
	}
	for _, v := range dbRecords.([]*neo4j.Record) {
		conversationWithParticipantIDs := new(ConversationDbNodeWithParticipantIDs)
		conversationWithParticipantIDs.Node = utils.NodePtr(v.Values[0].(neo4j.Node))
		conversationWithParticipantIDs.UserId = v.Values[1].(string)
		conversationWithParticipantIDs.ContactId = v.Values[2].(string)
		result.Nodes = append(result.Nodes, conversationWithParticipantIDs)
	}
	return result, nil
}
