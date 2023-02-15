package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
	"golang.org/x/net/context"
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
	Create(ctx context.Context, session neo4j.SessionWithContext, tenant string, userIds, contactIds []string, entity entity.ConversationEntity) (*dbtype.Node, error)
	Close(ctx context.Context, session neo4j.SessionWithContext, tenant, conversationId, status string, sourceOfTruth entity.DataSource) (*dbtype.Node, error)
	Update(ctx context.Context, session neo4j.SessionWithContext, tenant string, userIds, contactIds []string, skipMessageCountIncrement bool, entity entity.ConversationEntity) (*dbtype.Node, error)
	GetPaginatedConversationsForUser(ctx context.Context, session neo4j.SessionWithContext, tenant, userId string, skip, limit int, sort *utils.CypherSort) (*ConversationDbNodesWithTotalCount, error)
	GetPaginatedConversationsForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string, skip, limit int, sort *utils.CypherSort) (*ConversationDbNodesWithTotalCount, error)
}

type conversationRepository struct {
	driver *neo4j.DriverWithContext
}

func NewConversationRepository(driver *neo4j.DriverWithContext) ConversationRepository {
	return &conversationRepository{
		driver: driver,
	}
}

func (r *conversationRepository) Create(ctx context.Context, session neo4j.SessionWithContext, tenant string, userIds, contactIds []string, entity entity.ConversationEntity) (*dbtype.Node, error) {
	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "MATCH (t:Tenant {name:$tenant}) " +
			" MERGE (o:Conversation {id:$conversationId}) " +
			" ON CREATE SET o:%s, " +
			"				o.startedAt=$startedAt, o.messageCount=0, o.channel=$channel, o.status=$status, " +
			" 				o.source=$source, o.sourceOfTruth=$sourceOfTruth, o.appSource=$appSource " +
			" %s %s " +
			" RETURN DISTINCT o"
		queryLinkWithContacts := ""
		if len(contactIds) > 0 {
			queryLinkWithContacts = " WITH DISTINCT t, o " +
				" OPTIONAL MATCH (c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(t) WHERE c.id in $contactIds " +
				" WITH t, o, COLLECT(c) as participants " +
				" FOREACH (x in participants | MERGE (x)-[:PARTICIPATES]->(o) )"
		}
		queryLinkWithUsers := ""
		if len(userIds) > 0 {
			queryLinkWithUsers = " WITH DISTINCT t, o " +
				" OPTIONAL MATCH (u:User)-[:USER_BELONGS_TO_TENANT]->(t) WHERE u.id in $userIds " +
				" WITH t, o, COLLECT(u) as participants " +
				" FOREACH (x in participants | MERGE (x)-[:PARTICIPATES]->(o) )"
		}
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, "Conversation_"+tenant, queryLinkWithContacts, queryLinkWithUsers),
			map[string]interface{}{
				"tenant":         tenant,
				"status":         entity.Status,
				"startedAt":      entity.StartedAt,
				"channel":        entity.Channel,
				"conversationId": entity.Id,
				"contactIds":     contactIds,
				"userIds":        userIds,
				"source":         entity.Source,
				"sourceOfTruth":  entity.SourceOfTruth,
				"appSource":      entity.AppSource,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *conversationRepository) Update(ctx context.Context, session neo4j.SessionWithContext, tenant string, userIds, contactIds []string, skipMessageCountIncrement bool, entity entity.ConversationEntity) (*dbtype.Node, error) {
	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "MATCH (o:Conversation {id:$conversationId})--(p)--(t:Tenant {name:$tenant}) " +
			" WHERE 'Contact' IN labels(p) OR 'User' IN labels(p) " +
			" %s " +
			" %s " +
			" WITH DISTINCT o " +
			" SET o.sourceOfTruth=$sourceOfTruth %s" +
			" RETURN o"
		queryLinkWithContacts := ""
		if len(contactIds) > 0 {
			queryLinkWithContacts = " WITH DISTINCT t, o " +
				" OPTIONAL MATCH (c:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(t) WHERE c.id in $contactIds " +
				" WITH t, o, COLLECT(c) as participants " +
				" FOREACH (x in participants | MERGE (x)-[:PARTICIPATES]->(o) )"
		}
		queryLinkWithUsers := ""
		if len(userIds) > 0 {
			queryLinkWithUsers = " WITH DISTINCT t, o " +
				" OPTIONAL MATCH (u:User)-[:USER_BELONGS_TO_TENANT]->(t) WHERE u.id in $userIds " +
				" WITH t, o, COLLECT(u) as participants " +
				" FOREACH (x in participants | MERGE (x)-[:PARTICIPATES]->(o) )"
		}

		querySets := ""
		if len(entity.Channel) > 0 {
			querySets += ", o.channel=$channel "
		}
		if len(entity.Status) > 0 {
			querySets += ", o.status=$status "
		}
		if !skipMessageCountIncrement {
			querySets += ", o.messageCount=o.messageCount+1 "
		}
		queryResult, err := tx.Run(ctx, fmt.Sprintf(query, queryLinkWithContacts, queryLinkWithUsers, querySets),
			map[string]interface{}{
				"tenant":         tenant,
				"status":         entity.Status,
				"channel":        entity.Channel,
				"conversationId": entity.Id,
				"contactIds":     contactIds,
				"userIds":        userIds,
				"sourceOfTruth":  entity.SourceOfTruth,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *conversationRepository) Close(ctx context.Context, session neo4j.SessionWithContext, tenant, conversationId, status string, sourceOfTruth entity.DataSource) (*dbtype.Node, error) {
	if result, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "MATCH (o:Conversation {id:$conversationId})--(p)--(t:Tenant {name:$tenant}) " +
			" WHERE 'Contact' IN labels(p) OR 'User' IN labels(p) " +
			" SET o.endedAt=datetime({timezone: 'UTC'}), o.status=$status, o.sourceOfTruth=$sourceOfTruth" +
			" RETURN DISTINCT o"
		queryResult, err := tx.Run(ctx, query,
			map[string]interface{}{
				"tenant":         tenant,
				"conversationId": conversationId,
				"status":         status,
				"sourceOfTruth":  sourceOfTruth,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(ctx, queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}

func (r *conversationRepository) GetPaginatedConversationsForUser(ctx context.Context, session neo4j.SessionWithContext, tenant, userId string, skip, limit int, sort *utils.CypherSort) (*ConversationDbNodesWithTotalCount, error) {
	result := new(ConversationDbNodesWithTotalCount)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `MATCH (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), 
											(u)-[:PARTICIPATES]->(o:Conversation)
											RETURN count(o) as count`,
			map[string]any{
				"tenant": tenant,
				"userId": userId,
			})
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		result.Count = count.Values[0].(int64)

		queryResult, err = tx.Run(ctx, fmt.Sprintf(
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
		return queryResult.Collect(ctx)
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

func (r *conversationRepository) GetPaginatedConversationsForContact(ctx context.Context, session neo4j.SessionWithContext, tenant, contactId string, skip, limit int, sort *utils.CypherSort) (*ConversationDbNodesWithTotalCount, error) {
	result := new(ConversationDbNodesWithTotalCount)

	dbRecords, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		queryResult, err := tx.Run(ctx, `MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}), 
											(c)-[:PARTICIPATES]->(o:Conversation)
											RETURN count(o) as count`,
			map[string]any{
				"tenant":    tenant,
				"contactId": contactId,
			})
		if err != nil {
			return nil, err
		}
		count, _ := queryResult.Single(ctx)
		result.Count = count.Values[0].(int64)

		queryResult, err = tx.Run(ctx, fmt.Sprintf(
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
		return queryResult.Collect(ctx)
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
