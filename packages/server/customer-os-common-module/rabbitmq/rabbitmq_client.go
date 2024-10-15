package rabbimq

import (
	"context"
	"encoding/json"
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

type RabbitMQService struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	mu      sync.Mutex
	url     string
}

// failOnError logs errors and exits if a critical failure occurs
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// NewRabbitMQService initializes the RabbitMQ service with connection retry logic
func NewRabbitMQService(url string) *RabbitMQService {
	rabbitMQService := &RabbitMQService{
		url: url,
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

	eventMessage := make(map[string]interface{})
	eventMessage["event"] = make(map[string]interface{})
	eventMessage["event"].(map[string]interface{})["type"] = reflect.TypeOf(message).Name()
	eventMessage["event"].(map[string]interface{})["data"] = message
	eventMessage["metadata"] = tracing.ExtractTextMapCarrier((span).Context())
	eventMessage["metadata"].(opentracing.TextMapCarrier)["appSource"] = common.GetAppSourceFromContext(ctx)
	eventMessage["metadata"].(opentracing.TextMapCarrier)["tenant"] = common.GetTenantFromContext(ctx)
	eventMessage["metadata"].(opentracing.TextMapCarrier)["userId"] = common.GetUserIdFromContext(ctx)
	eventMessage["metadata"].(opentracing.TextMapCarrier)["userEmail"] = common.GetUserEmailFromContext(ctx)
	eventMessage["metadata"].(opentracing.TextMapCarrier)["timestamp"] = utils.Now().String()

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
