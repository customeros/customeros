package service

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

const (
	NotificationsExchangeName = "notifications"
	NotificationRoutingKey    = "notification"

	EventsExchangeName                      = "customeros"
	EventsRoutingKey                        = "event"
	EventsQueueName                         = "events"
	EventsFlowParticipantScheduleQueueName  = "events-flow-participant-schedule"
	EventsFlowParticipantScheduleRoutingKey = "flow-participant-schedule"
	EventsOpensearchQueueName               = "events-opensearch"
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
	connection      *amqp091.Connection
	connectionMutex sync.Mutex

	publishChannel *amqp091.Channel
	publishMutex   sync.Mutex

	url string

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

	// TODO sigterm
	// defer rabbitMQService.Close()

	return rabbitMQService
}

func (r *RabbitMQService) connect() {
	r.connectionMutex.Lock()
	defer r.connectionMutex.Unlock()

	var err error

	// Attempt to connect to RabbitMQ server
	r.connection, err = amqp091.Dial(r.url)
	failOnError(err, "Failed to connect to RabbitMQ")

	// Open a channel
	r.publishChannel, err = r.connection.Channel()
	failOnError(err, "Failed to open the publish channel")

	log.Println("Connected to RabbitMQ")

	// Set up a notification for connection closures to handle reconnections
	go func() {
		notifyClose := r.connection.NotifyClose(make(chan *amqp091.Error)) // Capture connection close
		err := <-notifyClose                                               // Wait for the error
		if err != nil {
			log.Printf("RabbitMQ connection closed: %v", err)
			r.reconnect(context.Background())
		}
	}()
}

// reconnect attempts to reconnect to RabbitMQ in case of connection loss
func (r *RabbitMQService) reconnect(ctx context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RabbitMQService.reconnect")
	defer span.Finish()

	retryCount := 0
	for {
		// Try reconnecting every 1 seconds
		time.Sleep(250 * time.Millisecond)

		r.connect()
		if r.connection != nil && r.connection.IsClosed() == false {
			return
		}

		retryCount++

		if retryCount > 5 {
			tracing.TraceErr(span, errors.New("Failed to reconnect to RabbitMQ"))
			return
		}
	}
}

// PublishEvent publishes a message to the given exchange and routing key with automatic reconnection
func (r *RabbitMQService) PublishEvent(ctx context.Context, entityId string, entityType model.EntityType, message interface{}) error {
	return r.PublishEventOnExchange(ctx, entityId, entityType, message, EventsExchangeName, EventsRoutingKey)
}

func (r *RabbitMQService) PublishWebhookEvent(ctx context.Context, event dto.WebhookEvent) error {
	return r.PublishEvent(ctx, "", model.WEBHOOK_EVENT, event)
}

func (r *RabbitMQService) PublishFlowActionEvent(ctx context.Context, event dto.FlowActionEvent) error {
	return r.PublishEvent(ctx, "", model.FLOW_ACTION_EVENT, event)
}

func (r *RabbitMQService) PublishEventOnExchange(ctx context.Context, entityId string, entityType model.EntityType, message interface{}, exchange, routingKey string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RabbitMQService.PublishEventOnExchange")
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

	return r.PublishMessageOnExchange(ctx, eventMessage, exchange, routingKey)
}

func (r *RabbitMQService) PublishMessageOnExchange(ctx context.Context, message interface{}, exchange, routingKey string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RabbitMQService.PublishMessageOnExchange")
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
		r.reconnect(ctx)
	}

	if r.connection.IsClosed() {
		tracing.TraceErr(span, errors.New("RabbitMQ connection is closed"))
		return nil
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

	err := r.PublishMessageOnExchange(ctx, event, NotificationsExchangeName, NotificationRoutingKey)
	if err != nil {
		tracing.TraceErr(span, err)
	}
}

// RegisterHandler allows you to register a handler for a specific event type
func (r *RabbitMQService) RegisterHandler(eventType interface{}, handler func(ctx context.Context, services *Services, event any) error) {
	typeOf := reflect.TypeOf(eventType)

	r.handlerRegistry[typeOf.Name()] = EventHandler{
		HandlerFunc: handler,
		EventType:   typeOf.Name(),
		DataType:    typeOf,
	}
}

func (r *RabbitMQService) ListenQueue(queueName string) {
	go func() {
		for {
			// Open a new channel for each queue (ensure it's properly closed after use)
			channel, err := r.connection.Channel()
			if err != nil {
				log.Printf("Failed to open channel for queue %s: %v. Retrying...", queueName, err)
				r.reconnect(context.Background()) // Reconnect in case of connection failure
				continue
			}
			defer channel.Close() // Ensure the channel is closed after use

			msgs, err := channel.Consume(
				queueName, // queue
				"",        // consumer tag
				false,     // auto-ack
				false,     // exclusive
				false,     // no-local
				false,     // no-wait
				nil,       // args
			)
			if err != nil {
				log.Printf("Failed to register consumer on queue %s: %v. Reconnecting...", queueName, err)
				r.reconnect(context.Background())
				continue
			}

			done := make(chan bool)

			log.Printf("Listening for messages on queue %s", queueName)

			go func() {
				for d := range msgs {
					err := r.ProcessMessage(d)
					if err != nil {
						log.Printf("Failed to process message: %v", err)
						retryAckNack(d, false)
					} else {
						retryAckNack(d, true)
					}
				}
				// Signal disconnection when msgs channel is closed
				done <- true
			}()

			// Wait for channel closure to trigger reconnection in connectAndConsume
			<-done
			log.Printf("Connection lost for queue %s. Reconnecting...", queueName)
			r.reconnect(context.Background())
		}
	}()
}

func (r *RabbitMQService) ListenQueueExclusive(queueName string) {
	go func() {
		for {
			// Open a new channel for each queue (ensure it's properly closed after use)
			channel, err := r.connection.Channel()
			if err != nil {
				log.Printf("Failed to open channel for exclusive queue %s: %v. Retrying...", queueName, err)
				r.reconnect(context.Background()) // Reconnect in case of connection failure
				continue
			}
			defer channel.Close() // Ensure the channel is closed after use

			msgs, err := channel.Consume(
				queueName, // queue
				"",        // consumer tag
				false,     // auto-ack
				true,      // exclusive
				false,     // no-local
				false,     // no-wait
				nil,       // args
			)
			if err != nil {
				if strings.Contains(err.Error(), "ACCESS_REFUSED") && strings.Contains(err.Error(), "exclusive") {
					log.Printf("Exclusive consumer conflict for queue %s. Only one instance can consume exclusively.", queueName)
					time.Sleep(10 * time.Second)
					continue
				} else {
					log.Printf("Failed to register exclusive consumer on queue %s: %v. Reconnecting...", queueName, err)
					r.reconnect(context.Background())
					continue
				}
			}

			log.Printf("Listening for messages on exclusive queue %s", queueName)

			done := make(chan bool)

			go func() {
				for d := range msgs {
					err := r.ProcessMessage(d)
					if err != nil {
						log.Printf("Failed to process message: %v", err)
						retryAckNack(d, false)
					} else {
						retryAckNack(d, true)
					}
				}
				// Signal disconnection when msgs channel is closed
				done <- true
			}()

			// Wait for channel closure to trigger reconnection
			<-done
			log.Printf("Connection lost for queue %s. Reconnecting...", queueName)
			r.reconnect(context.Background())
		}
	}()
}

func retryAckNack(d amqp091.Delivery, ack bool) {
	maxRetries := 5
	retryDelay := 100 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		var err error
		if ack {
			err = d.Ack(false)
		} else {
			err = d.Nack(false, false)
		}

		if err == nil {
			return // Successfully Acked/Nacked, exit the retry loop
		}

		time.Sleep(retryDelay) // Wait before retrying
	}
}

func (r *RabbitMQService) ProcessMessage(d amqp091.Delivery) error {
	ctx := context.Background()

	var event dto.Event
	if err := json.Unmarshal(d.Body, &event); err != nil {
		log.Printf("Failed to unmarshal message: %s", err)
		return err
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
		return errors.New("Data not found in message")
	}

	// Invoke the appropriate handler based on the event type
	eventHandler, found := r.handlerRegistry[event.Event.EventType]
	if !found {
		return errors.New("Handler not found for event type: " + event.Event.EventType)
	}

	if data == nil {
		return errors.New("Data not found in message")
	}

	eventDataPtr := reflect.New(eventHandler.DataType).Interface()
	if err := mapstructure.Decode(data, eventDataPtr); err != nil {
		return err
	}

	event.Event.Data = eventDataPtr

	err := eventHandler.HandlerFunc(ctx, r.services, &event) // Pass the entire delivery struct to the handler
	if err != nil {
		log.Printf("Failed to handle message: %s", err)
		return err
	}

	return nil
}

// Close closes the RabbitMQ connection and channel
func (r *RabbitMQService) Close() {
	r.connectionMutex.Lock()
	defer r.connectionMutex.Unlock()

	if r.connection != nil {
		r.connection.Close()
	}
}
