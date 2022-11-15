package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/openline-ai/openline-customer-os/customer-os-api/utils"
)

type ConversationRepository interface {
	Create(tenant string, userId string, contactId string, conversationId string) (any, error)
}

type conversationRepository struct {
	driver *neo4j.Driver
	repos  *RepositoryContainer
}

func NewConversationRepository(driver *neo4j.Driver, repos *RepositoryContainer) ConversationRepository {
	return &conversationRepository{
		driver: driver,
		repos:  repos,
	}
}

func (r *conversationRepository) Create(tenant string, userId string, contactId string, conversationId string) (any, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	return session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		txResult, err := tx.Run(`
			MATCH (c:Contact {id:$contactId})-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant}),
				  (u:User {id:$userId})-[:USER_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			MERGE (c)-[:PARTICIPATES]->(o:Conversation {id:$conversationId})<-[:PARTICIPATES]-(u)
            ON CREATE SET o.started=datetime({timezone: 'UTC'})
			RETURN o`,
			map[string]interface{}{
				"tenant":         tenant,
				"contactId":      contactId,
				"userId":         userId,
				"conversationId": conversationId,
			})
		record, err := txResult.Single()
		if err != nil {
			return nil, err
		}
		return record.Values[0], nil
	})
}
