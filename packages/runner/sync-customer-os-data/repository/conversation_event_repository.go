package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/entity"
	"gorm.io/gorm"
)

type ConversationEventRepository interface {
	Save(entity entity.ConversationEvent)
}

type conversationEventRepository struct {
	db *gorm.DB
}

func NewConversationEventRepository(gormDb *gorm.DB) ConversationEventRepository {
	return &conversationEventRepository{
		db: gormDb,
	}
}

func (r *conversationEventRepository) Save(entity entity.ConversationEvent) {
	r.db.Create(&entity)
}
