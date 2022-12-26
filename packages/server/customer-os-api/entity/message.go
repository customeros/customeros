package entity

import (
	"fmt"
	"time"
)

type MessageEntity struct {
	Id             string
	StartedAt      time.Time `neo4jDb:"property:startedAt;lookupName:STARTED_AT;supportCaseSensitive:false"`
	ConversationId string
	Channel        string
}

func (message MessageEntity) ToString() string {
	return fmt.Sprintf("id: %s", message.Id)
}

type MessageEntities []MessageEntity

func (MessageEntity) Action() {
}

func (MessageEntity) ActionName() string {
	return NodeLabel_Message
}

func (message MessageEntity) Labels() []string {
	return []string{"Message", "Action"}
}
