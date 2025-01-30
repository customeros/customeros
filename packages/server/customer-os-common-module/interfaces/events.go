package interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type EventPublisher interface {
	PublishFanoutEvent(ctx context.Context, entityId string, entityType model.EntityType, message interface{}) error
	PublishDirectEvent(ctx context.Context, entityId string, entityType model.EntityType, message interface{}) error
	PublishNotification(ctx context.Context, tenant string, entityId string, entityType model.EntityType, details *utils.EventCompletedDetails)
	PublishNotificationBulk(ctx context.Context, tenant string, entityIds []string, entityType model.EntityType, details *utils.EventCompletedDetails)
	Close() error
}

type EventListener interface {
	Handle(ctx context.Context, baseEvent any) error
	GetEventType() string
	GetQueueName() string
}

type EventSubscriber interface {
	RegisterListener(listener EventListener)
	ListenQueue(queueName string) error
	ListenQueueExclusive(queueName string) error
	Close() error
}
