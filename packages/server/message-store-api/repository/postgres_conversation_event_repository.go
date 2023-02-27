package repository

import (
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/message-store-api/repository/helper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ConversationEventRepository interface {
	GetEventById(eventId string) helper.QueryResult
	GetEventsForConversation(conversationId string, t *time.Time, pageSize int) helper.QueryResult
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

func (r *conversationEventRepository) GetEventById(eventId string) helper.QueryResult {
	var event entity.ConversationEvent

	err := r.db.Where("id = ?", eventId).
		Find(&entity.ConversationEvent{}).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true}).
		First(&event).
		Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &event}
}

func (r *conversationEventRepository) GetEventsForConversation(conversationId string, t *time.Time, pageSize int) helper.QueryResult {
	var rows []entity.ConversationEvent

	intermediate := r.db.Where("conversation_id = ?", conversationId).
		Find(&entity.ConversationEvent{}).
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true})

	if t != nil {
		intermediate = intermediate.Where(clause.Lt{Column: clause.Column{Name: "created_at"}, Value: t})
	}

	err := intermediate.Limit(pageSize).
		Find(&rows).
		Error

	if err != nil {
		return helper.QueryResult{Error: err}
	}

	return helper.QueryResult{Result: &rows}
}

func (r *conversationEventRepository) Save(conversationEvent *entity.ConversationEvent) helper.QueryResult {
	result := r.db.Create(&conversationEvent)

	if result.Error != nil {
		return helper.QueryResult{Error: result.Error}
	}

	return helper.QueryResult{Result: &conversationEvent}
}
