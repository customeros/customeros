package service

import (
	"context"
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
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

type Event struct {
	Event    EventDetails  `json:"event"`
	Metadata EventMetadata `json:"metadata"`
}

type EventDetails struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type EventMetadata struct {
	UberTraceId string `json:"uber-trace-id"`
	AppSource   string `json:"appSource"`
	Tenant      string `json:"tenant"`
	UserId      string `json:"userId"`
	UserEmail   string `json:"userEmail"`
	Timestamp   string `json:"timestamp"`
}

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
func NewRabbitMQService(url string, services *Services) *RabbitMQService {
	rabbitMQService := &RabbitMQService{
		url:             url,
		services:        services,
		handlerRegistry: make(map[string]EventHandler),
	}
	rabbitMQService.connect()
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

// Publish publishes a message to the given exchange and routing key with automatic reconnection
func (r *RabbitMQService) Publish(ctx context.Context, exchange string, routingKey string, message interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "RabbitMQService.Publish")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tracing.LogObjectAsJson(span, "message", message)

	r.mu.Lock()
	defer r.mu.Unlock()

	tracingData := tracing.ExtractTextMapCarrier((span).Context())

	eventMessage := Event{
		Event: EventDetails{
			Type: reflect.TypeOf(message).Name(),
			Data: message,
		},
		Metadata: EventMetadata{
			UberTraceId: tracingData["uber-trace-id"],
			AppSource:   common.GetAppSourceFromContext(ctx),
			Tenant:      common.GetTenantFromContext(ctx),
			UserId:      common.GetUserIdFromContext(ctx),
			UserEmail:   common.GetUserEmailFromContext(ctx),
			Timestamp:   utils.Now().String(),
		},
	}

	// Convert the message to JSON
	jsonBody, err := json.Marshal(eventMessage)
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
				ContentType: "application/json",
				Body:        jsonBody,
			})

		if err != nil {
			tracing.TraceErr(span, err)
			r.reconnect()
		} else {
			break // Message sent successfully, exit retry loop
		}
	}

	log.Printf(" [x] Sent message to exchange %s with routing key %s", exchange, routingKey)
	return nil
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
		"events", // queue
		"",       // consumer tag
		false,    // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
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

	var event Event
	if err := json.Unmarshal(d.Body, &event); err != nil {
		log.Printf("Failed to unmarshal message: %s", err)
		return
	}

	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    event.Metadata.Tenant,
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
	eventHandler, found := r.handlerRegistry[event.Event.Type]
	if !found {
		log.Printf("No handler registered for event type: %s", event.Event.Type)
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
		log.Printf("Failed to decode data for event type %s: %s", event.Event.Type, err)
		if err := d.Nack(false, false); err != nil {
			log.Printf("Failed to negatively acknowledge message: %s", err)
		}
		return
	}

	err := eventHandler.HandlerFunc(ctx, r.services, eventDataPtr) // Pass the entire delivery struct to the handler
	if err != nil {
		log.Printf("Failed to handle message: %s", err)
		if err := d.Nack(false, false); err != nil {
			log.Printf("Failed to negatively acknowledge message: %s", err)
		}
		//TODO dead letter queue
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
