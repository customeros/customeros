package events

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

const (
	NotificationsExchangeName = "notifications"
	NotificationRoutingKey    = "notification"

	EventsExchangeName                      = "customeros"
	EventsDirectExchangeName                = "customeros-direct"
	EventsRoutingKey                        = "event"
	EventsQueueName                         = "events"
	EventsFlowParticipantScheduleQueueName  = "events-flow-participant-schedule"
	EventsFlowParticipantScheduleRoutingKey = "flow-participant-schedule"
	EventsOpensearchQueueName               = "events-opensearch"
)

// RabbitMQPublisher implements the EventPublisher interface for RabbitMQ
type RabbitMQPublisher struct {
	connection      *amqp091.Connection
	connectionMutex sync.Mutex
	publishChannel  *amqp091.Channel
	publishMutex    sync.Mutex
	url             string
	logger          logger.Logger
}

// NewRabbitMQPublisher creates a new RabbitMQ publisher
func NewRabbitMQPublisher(rabbitmqURL string, logger logger.Logger) (*RabbitMQPublisher, error) {
	publisher := &RabbitMQPublisher{
		url:    rabbitmqURL,
		logger: logger,
	}

	err := publisher.connect()
	if err != nil {
		return nil, err
	}

	return publisher, nil
}

func (r *RabbitMQPublisher) connect() error {
	r.connectionMutex.Lock()
	defer r.connectionMutex.Unlock()

	var err error
	r.connection, err = amqp091.Dial(r.url)
	if err != nil {
		return errors.Wrap(err, "Failed to connect to RabbitMQ")
	}

	r.publishChannel, err = r.connection.Channel()
	if err != nil {
		return errors.Wrap(err, "Failed to open publish channel")
	}

	// Set up reconnection mechanism
	go func() {
		notifyClose := r.connection.NotifyClose(make(chan *amqp091.Error))
		<-notifyClose
		r.logger.Warn("RabbitMQ connection closed, attempting to reconnect")
		_ = r.connect()
	}()

	return nil
}

// PublishEvent publishes an event to the default exchange and routing key
func (r *RabbitMQPublisher) PublishEvent(ctx context.Context, entityId string, entityType model.EntityType, message interface{}) error {
	return r.PublishEventOnExchange(ctx, entityId, entityType, message, EventsExchangeName, EventsRoutingKey)
}

// PublishWebhookEvent publishes a webhook event
func (r *RabbitMQPublisher) PublishWebhookEvent(ctx context.Context, event dto.WebhookEvent) error {
	return r.PublishEvent(ctx, "", model.WEBHOOK_EVENT, event)
}

// PublishFlowAgentEvent publishes a flow agent event
func (r *RabbitMQPublisher) PublishFlowAgentEvent(ctx context.Context, event dto.FlowAgentEvent) error {
	return r.PublishEvent(ctx, "", model.FLOW_ACTION_EVENT, event)
}

// PublishFlowAgentEventResult publishes a flow agent event result
func (r *RabbitMQPublisher) PublishFlowAgentEventResult(ctx context.Context, event dto.FlowAgentExecutionResultEvent) error {
	return r.PublishEvent(ctx, "", model.FLOW_ACTION_RESULT_EVENT, event)
}

// PublishEventOnExchange publishes an event to a specific exchange and routing key
func (r *RabbitMQPublisher) PublishEventOnExchange(ctx context.Context, entityId string, entityType model.EntityType, message interface{}, exchange, routingKey string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RabbitMQPublisher.PublishEventOnExchange")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tracingData := tracing.ExtractTextMapCarrier((span).Context())

	eventMessage := dto.Event{
		Event: dto.EventDetails{
			Id:         utils.GenerateRandomString(32),
			EntityId:   entityId,
			EntityType: entityType,
			Tenant:     common.GetTenantFromContext(ctx),
			EventType:  reflect.TypeOf(message).Name(),
			Data:       message,
		},
		Metadata: dto.EventMetadata{
			UberTraceId: tracingData["uber-trace-id"],
			AppSource:   common.GetAppSourceFromContext(ctx),
			UserId:      common.GetUserIdFromContext(ctx),
			UserEmail:   common.GetUserEmailFromContext(ctx),
			Timestamp:   utils.Now().String(),
		},
	}

	return r.publishMessageOnExchange(ctx, eventMessage, exchange, routingKey)
}

func (r *RabbitMQPublisher) publishMessageOnExchange(ctx context.Context, message interface{}, exchange, routingKey string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RabbitMQPublisher.PublishMessageOnExchange")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tracing.LogObjectAsJson(span, "message", message)

	r.publishMutex.Lock()
	defer r.publishMutex.Unlock()

	// Convert the message to JSON
	jsonBody, err := json.Marshal(message)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Ensure the connection is not nil
	if r.connection == nil || r.connection.IsClosed() || r.publishChannel == nil || r.publishChannel.IsClosed() {
		if err := r.connect(); err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	// Try publishing the message
	err = r.publishChannel.Publish(
		exchange,   // Exchange name
		routingKey, // Routing key
		false,      // Mandatory
		false,      // Immediate
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "application/json",
			Body:         jsonBody,
		})
	if err != nil {
		tracing.TraceErr(span, err)
	}

	return err
}

// PublishEventCompleted publishes an event completion notification
func (r *RabbitMQPublisher) PublishEventCompleted(ctx context.Context, tenant string, entityId string, entityType model.EntityType, details *utils.EventCompletedDetails) {
	r.PublishEventCompletedBulk(ctx, tenant, []string{entityId}, entityType, details)
}

func (r *RabbitMQPublisher) PublishEventCompletedBulk(ctx context.Context, tenant string, entityIds []string, entityType model.EntityType, details *utils.EventCompletedDetails) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RabbitMQPublisher.PublishEventCompletedBulk")
	defer span.Finish()
	span.LogKV("tenant", tenant, "entityType", entityType, "entityIds", entityIds)

	event := dto.EventCompleted{
		Tenant:     tenant,
		EntityType: entityType,
		EntityIds:  entityIds,
		Create:     false,
		Update:     false,
		Delete:     false,
	}

	if details != nil {
		event.Create = details.Create
		event.Update = details.Update
		event.Delete = details.Delete
	}

	err := r.publishMessageOnExchange(ctx, event, NotificationsExchangeName, NotificationRoutingKey)
	if err != nil {
		tracing.TraceErr(span, err)
	}
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQPublisher) Close() error {
	r.connectionMutex.Lock()
	defer r.connectionMutex.Unlock()

	if r.connection != nil {
		return r.connection.Close()
	}
	return nil
}
