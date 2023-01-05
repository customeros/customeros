package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store/repository/helper"
	"gorm.io/gorm"
)

type ConversationEventRepository interface {
	Save(conversationEvent *entity.ConversationEvent) helper.QueryResult
}

type conversationEventRepository struct {
	db *gorm.DB
}

func NewConversationEventRepository(db *gorm.DB) ConversationEventRepository {
	return &conversationEventRepository{
		db: db,
	}
}

func (r *conversationEventRepository) Save(conversationEvent *entity.ConversationEvent) helper.QueryResult {
	result := r.db.Create(&conversationEvent)

	if result.Error != nil {
		return helper.QueryResult{Error: result.Error}
	}

	return helper.QueryResult{Result: &conversationEvent}
}
