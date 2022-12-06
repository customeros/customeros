package repository

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/utils"
)

type MessageRepository interface {
	CreateMessage(tenant string, conversationId string, entity *entity.MessageEntity) (*dbtype.Node, error)
}

type messageRepository struct {
	driver *neo4j.Driver
}

func NewMessageRepository(driver *neo4j.Driver) MessageRepository {
	return &messageRepository{
		driver: driver,
	}
}

func (r *messageRepository) CreateMessage(tenant string, conversationId string, entity *entity.MessageEntity) (*dbtype.Node, error) {
	session := utils.NewNeo4jWriteSession(*r.driver)
	defer session.Close()

	if result, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		queryResult, err := tx.Run(`
			MATCH (o:Conversation {id:$conversationId})<-[:PARTICIPATES]-(:Contact)-[:CONTACT_BELONGS_TO_TENANT]->(:Tenant {name:$tenant})
			MERGE (o)-[:CONSISTS_OF]->(m:Message:Action {id:$messageId})
			ON CREATE SET m.channel=$channel, m.startedAt=$startedAt, m.conversationId=$conversationId
			RETURN m`,
			map[string]any{
				"tenant":         tenant,
				"conversationId": conversationId,
				"messageId":      entity.Id,
				"channel":        entity.Channel,
				"startedAt":      entity.StartedAt,
			})
		return utils.ExtractSingleRecordFirstValueAsNode(queryResult, err)
	}); err != nil {
		return nil, err
	} else {
		return result.(*dbtype.Node), nil
	}
}
