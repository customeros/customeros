package interfaces

import (
	"context"
	"reflect"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type EventPublisher interface {
	PublishEvent(ctx context.Context, entityId string, entityType model.EntityType, message interface{}) error
	PublishEventOnExchange(ctx context.Context, entityId string, entityType model.EntityType, message interface{}, exchange, routingKey string) error
	PublishEventCompleted(ctx context.Context, tenant string, entityId string, entityType model.EntityType, details *utils.EventCompletedDetails)
	PublishEventCompletedBulk(ctx context.Context, tenant string, entityIds []string, entityType model.EntityType, details *utils.EventCompletedDetails)
	PublishWebhookEvent(ctx context.Context, event dto.WebhookEvent) error
	PublishFlowAgentEvent(ctx context.Context, event dto.FlowAgentEvent) error
	PublishFlowAgentEventResult(ctx context.Context, event dto.FlowAgentExecutionResultEvent) error
	Close() error
}

type EventSubscriber interface {
	RegisterHandler(eventType interface{}, handler EventHandler)
	ListenQueue(queueName string) error
	ListenQueueExclusive(queueName string) error
	Close() error
}

type EventHandler struct {
	HandlerFunc func(ctx context.Context, event any) error
	EventType   string
	DataType    reflect.Type
}
