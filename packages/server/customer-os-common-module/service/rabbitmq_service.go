package service

import (
	"context"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto/events"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"reflect"
	"sync"
	"time"
)

const (
	NotificationsExchangeName = "notifications"
	NotificationRoutingKey    = "notification"

	EventsExchangeName        = "customeros"
	EventsRoutingKey          = "event"
	EventsQueueName           = "events"
	EventsOpensearchQueueName = "events-opensearch"
)

//{
//	event: {
//		id: "123",
//		entity: "flow",
//		tenant: "tenant",
//		type: "FlowInitialSchedule",
//
//		data: {
//			flowId: "123",
//		}
//	},
//	metadata: {
//		appSource: "customer-os",
//		uber-trace-id: "123",
//		userId: "ABC",
//		userEmail: ""
//		timestamp: "2021-09-01T12:00:00Z"
//	}
//}

type EventHandler struct {
	HandlerFunc func(ctx context.Context, services *Services, event any) error
	EventType   string
	DataType    reflect.Type
}

type RabbitMQService struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	mu      sync.Mutex
	url     string

	handlerRegistry map[string]EventHandler

	services *Services
}

// failOnError logs errors and exits if a critical failure occurs
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// NewRabbitMQService initializes the RabbitMQ service with connection retry logic
func NewRabbitMQService(url string, services *Services) RabbitMQService {
	rabbitMQService := RabbitMQService{
		url:             url,
		services:        services,
		handlerRegistry: make(map[string]EventHandler),
	}
	rabbitMQService.connect()

	//TODO sigterm
	//defer rabbitMQService.Close()

	return rabbitMQService
}

// connect establishes a connection to RabbitMQ and opens a channel
func (r *RabbitMQService) connect() {
	var err error

	// Attempt to connect to RabbitMQ server
	r.conn, err = amqp091.Dial(r.url)
	failOnError(err, "Failed to connect to RabbitMQ")

	// Open a channel
	r.channel, err = r.conn.Channel()
	failOnError(err, "Failed to open a channel")

	log.Println("Connected to RabbitMQ")

	// Set up a notification for connection closures to handle reconnections
	go func() {
		notifyClose := r.conn.NotifyClose(make(chan *amqp091.Error)) // Capture connection close
		err := <-notifyClose                                         // Wait for the error
		if err != nil {
			log.Printf("RabbitMQ connection closed: %v", err)
			r.reconnect()
		}
	}()
}

// reconnect attempts to reconnect to RabbitMQ in case of connection loss
func (r *RabbitMQService) reconnect() {
	for {
		log.Println("Attempting to reconnect to RabbitMQ...")

		// Try reconnecting every 1 seconds
		time.Sleep(1 * time.Second)

		r.connect() // Re-establish connection
		if r.conn != nil && r.conn.IsClosed() == false {
			log.Println("Reconnected to RabbitMQ")
			return
		}
	}
}

// PublishEvent publishes a message to the given exchange and routing key with automatic reconnection
func (r *RabbitMQService) PublishEvent(ctx context.Context, entityId string, entityType model.EntityType, message interface{}) error {
	return r.PublishEventOnQueue(ctx, entityId, entityType, message, EventsExchangeName, EventsRoutingKey)
}

func (r *RabbitMQService) PublishEventOnQueue(ctx context.Context, entityId string, entityType model.EntityType, message interface{}, exchange, routingKey string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RabbitMQService.PublishOnQueue")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if r.conn == nil {
		tracing.TraceErr(span, errors.New("RabbitMQ connection is nil"))
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	tracingData := tracing.ExtractTextMapCarrier((span).Context())

	eventMessage := events.Event{
		Event: events.EventDetails{
			Id:         utils.GenerateRandomString(32),
			EntityId:   entityId,
			EntityType: entityType.String(),
			Tenant:     common.GetTenantFromContext(ctx),
			EventType:  reflect.TypeOf(message).Name(),
			Data:       message,
		},
		Metadata: events.EventMetadata{
			UberTraceId: tracingData["uber-trace-id"],
			AppSource:   common.GetAppSourceFromContext(ctx),
			UserId:      common.GetUserIdFromContext(ctx),
			UserEmail:   common.GetUserEmailFromContext(ctx),
			Timestamp:   utils.Now().String(),
		},
	}

	return r.PublishOnQueue(ctx, eventMessage, exchange, routingKey)
}

func (r *RabbitMQService) PublishOnQueue(ctx context.Context, message interface{}, exchange, routingKey string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RabbitMQService.PublishOnQueue")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tracing.LogObjectAsJson(span, "message", message)

	if r.conn == nil {
		tracing.TraceErr(span, errors.New("RabbitMQ connection is nil"))
		return nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Convert the message to JSON
	jsonBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	// Retry logic in case the connection is closed
	for {
		// Ensure the connection and channel are open
		if r.conn.IsClosed() {
			tracing.TraceErr(span, errors.New("RabbitMQ connection is closed"))
			r.reconnect()
		}

		// Try publishing the message
		err = r.channel.Publish(
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
			r.reconnect()
		} else {
			break // Message sent successfully, exit retry loop
		}
	}

	log.Printf(" [x] Sent message to exchange %s with routing key %s", EventsExchangeName, EventsRoutingKey)
	return nil
}

func (r *RabbitMQService) PublishEventCompleted(ctx context.Context, tenant string, entityId string, entityType model.EntityType, details *utils.EventCompletedDetails) {
	r.PublishEventCompletedBulk(ctx, tenant, []string{entityId}, entityType, details)
}

func (r *RabbitMQService) PublishEventCompletedBulk(ctx context.Context, tenant string, entityIds []string, entityType model.EntityType, details *utils.EventCompletedDetails) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RabbitMQService.PublishEventCompletedBulk")
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

	err := r.PublishOnQueue(ctx, event, NotificationsExchangeName, NotificationRoutingKey)

	if err != nil {
		tracing.TraceErr(span, err)
	}
}

// RegisterHandler allows you to register a handler for a specific event type
func (r *RabbitMQService) RegisterHandler(eventType interface{}, handler func(ctx context.Context, services *Services, event any) error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	typeOf := reflect.TypeOf(eventType)

	r.handlerRegistry[typeOf.Name()] = EventHandler{
		HandlerFunc: handler,
		EventType:   typeOf.Name(),
		DataType:    typeOf,
	}
}

// Listen listens for messages from the specified queue and processes them with the provided handler
func (r *RabbitMQService) Listen() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Start consuming messages
	msgs, err := r.channel.Consume(
		EventsQueueName, // queue
		"",              // consumer tag
		false,           // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	failOnError(err, "Failed to register consumer")

	// Handle messages in a separate goroutine
	go func() {
		for d := range msgs {
			r.ProcessMessage(d)
		}
	}()
}

func (r *RabbitMQService) ProcessMessage(d amqp091.Delivery) {
	ctx := context.Background()

	var event events.Event
	if err := json.Unmarshal(d.Body, &event); err != nil {
		log.Printf("Failed to unmarshal message: %s", err)
		return
	}

	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    event.Event.Tenant,
		AppSource: event.Metadata.AppSource,
		UserId:    event.Metadata.UserId,
		UserEmail: event.Metadata.UserEmail,
	})

	ctx, span := tracing.StartRabbitMQMessageTracerSpanWithHeader(ctx, "RabbitMQService.Listen", event.Metadata.UberTraceId)
	defer span.Finish()

	data, ok := event.Event.Data.(map[string]interface{})
	if !ok {
		log.Printf("Data not found in message: %s", d.Body)
		if err := d.Nack(false, false); err != nil {
			log.Printf("Failed to negatively acknowledge message: %s", err)
		}
		return
	}

	// Invoke the appropriate handler based on the event type
	eventHandler, found := r.handlerRegistry[event.Event.EventType]
	if !found {
		if err := d.Ack(false); err != nil {
			log.Printf("Failed to acknowledge message: %s", err.Error())
			//TODO retry nack
		}
		return
	}

	if data == nil {
		tracing.TraceErr(nil, errors.New("Data not found in message"))
		if err := d.Nack(false, false); err != nil {
			log.Printf("Failed to negatively acknowledge message: %s", err)
		}
		return
	}

	eventDataPtr := reflect.New(eventHandler.DataType).Interface()
	if err := mapstructure.Decode(data, eventDataPtr); err != nil {
		log.Printf("Failed to decode data for event type %s: %s", event.Event.EventType, err)
		if err := d.Nack(false, false); err != nil {
			log.Printf("Failed to negatively acknowledge message: %s", err)
		}
		return
	}

	event.Event.Data = eventDataPtr

	err := eventHandler.HandlerFunc(ctx, r.services, &event) // Pass the entire delivery struct to the handler
	if err != nil {
		log.Printf("Failed to handle message: %s", err)
		if err := d.Nack(false, false); err != nil {
			log.Printf("Failed to negatively acknowledge message: %s", err)
		}
		return
	}

	if err := d.Ack(false); err != nil {
		log.Printf("Failed to acknowledge message: %s", err)
		//TODO retry nack
	}
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQService) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
