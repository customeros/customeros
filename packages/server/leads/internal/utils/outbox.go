package utils

import (
	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	"github.com/customeros/customeros/packages/server/leads/internal/models"
)

func BuildOutboxEvent(eventType enum.Events, payload []byte) *models.OutboxEvent {
	return &models.OutboxEvent{
		ID:        GenerateNanoIDWithPrefix("event", 21),
		EventType: eventType,
		Payload:   payload,
		Status:    enum.OutboxPending,
		CreatedAt: Now(),
	}
}
