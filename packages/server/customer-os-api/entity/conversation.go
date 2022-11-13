package entity

import (
	"fmt"
	"time"
)

type ConversationEntity struct {
	Id        string
	StartedAt time.Time
}

func (conversation ConversationEntity) ToString() string {
	return fmt.Sprintf("id: %s", conversation.Id)
}

type ConversationEntities []ConversationEntity
