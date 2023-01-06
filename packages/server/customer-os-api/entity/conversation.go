package entity

import (
	"fmt"
	"time"
)

type ConversationEntity struct {
	Id           string
	StartedAt    time.Time  `neo4jDb:"property:startedAt;lookupName:STARTED_AT;supportCaseSensitive:false"`
	EndedAt      *time.Time `neo4jDb:"property:endedAt;lookupName:ENDED_AT;supportCaseSensitive:false"`
	Status       string
	Channel      string
	MessageCount int64
}

func (conversation ConversationEntity) ToString() string {
	return fmt.Sprintf("id: %s", conversation.Id)
}

type ConversationEntities []ConversationEntity

func (conversation ConversationEntity) Labels() []string {
	return []string{"Conversation"}
}
