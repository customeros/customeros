package events

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

// EventHandlerRegistration holds the handler information for a specific event type
type EventHandlerRegistration struct {
	Handler   interfaces.EventHandler
	EventType string
	DataType  reflect.Type
}

// RabbitMQSubscriber implements the EventSubscriber interface for RabbitMQ
type RabbitMQSubscriber struct {
	connection      *amqp091.Connection
	connectionMutex sync.Mutex
	url             string
	logger          logger.Logger

	handlerRegistryMutex sync.RWMutex
	handlerRegistry      map[string]EventHandlerRegistration
}

// NewRabbitMQSubscriber creates a new RabbitMQ subscriber
func NewRabbitMQSubscriber(rabbitmqURL string, logger logger.Logger) (*RabbitMQSubscriber, error) {
	subscriber := &RabbitMQSubscriber{
		url:             rabbitmqURL,
		logger:          logger,
		handlerRegistry: make(map[string]EventHandlerRegistration),
	}

	err := subscriber.connect()
	if err != nil {
		return nil, err
	}

	return subscriber, nil
}

// connect establishes a connection to RabbitMQ
func (r *RabbitMQSubscriber) connect() error {
	r.connectionMutex.Lock()
	defer r.connectionMutex.Unlock()

	var err error
	r.connection, err = amqp091.Dial(r.url)
	if err != nil {
		return errors.Wrap(err, "Failed to connect to RabbitMQ")
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

// RegisterHandler registers a handler for a specific event type
func (r *RabbitMQSubscriber) RegisterHandler(eventType interface{}, handler interfaces.EventHandler) {
	r.handlerRegistryMutex.Lock()
	defer r.handlerRegistryMutex.Unlock()

	typeOf := reflect.TypeOf(eventType)
	r.handlerRegistry[typeOf.Name()] = EventHandlerRegistration{
		Handler:   handler,
		EventType: typeOf.Name(),
		DataType:  typeOf,
	}
}

// processMessage handles the processing of a single message
func (r *RabbitMQSubscriber) processMessage(d amqp091.Delivery) error {
	ctx := context.Background()

	var event dto.Event
	if err := json.Unmarshal(d.Body, &event); err != nil {
		r.logger.Errorf("Failed to unmarshal message: %v", err)
		return err
	}

	// Enrich context with event metadata
	ctx = common.WithCustomContext(ctx, &common.CustomContext{
		Tenant:    event.Event.Tenant,
		AppSource: event.Metadata.AppSource,
		UserId:    event.Metadata.UserId,
		UserEmail: event.Metadata.UserEmail,
	})

	ctx, span := tracing.StartRabbitMQMessageTracerSpanWithHeader(ctx, "RabbitMQSubscriber.ProcessMessage", event.Metadata.UberTraceId)
	defer span.Finish()

	// Retrieve the data as a map
	data, ok := event.Event.Data.(map[string]interface{})
	if !ok {
		r.logger.Errorf("Data not found in message: %s", d.Body)
		err := errors.New("data not found in message")
		tracing.TraceErr(span, err)
		return err
	}

	// Find the appropriate handler
	r.handlerRegistryMutex.RLock()
	handlerReg, found := r.handlerRegistry[event.Event.EventType]
	r.handlerRegistryMutex.RUnlock()

	if !found {
		return nil // No handler found, ignore the message
	}

	if data == nil {
		tracing.TraceErr(span, errors.New("data is nil in message"))
		return errors.New("data is nil in message")
	}

	// Decode the data into the specific event type
	eventDataPtr := reflect.New(handlerReg.DataType).Interface()
	if err := mapstructure.Decode(data, eventDataPtr); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "Failed to decode data"))
		return err
	}

	event.Event.Data = eventDataPtr

	// Call the handler
	return handlerReg.Handler.HandlerFunc(ctx, &event)
}

// ListenQueue starts listening to a standard queue
func (r *RabbitMQSubscriber) ListenQueue(queueName string) error {
	return r.listenQueueWithExclusive(queueName, false)
}

// ListenQueueExclusive starts listening to an exclusive queue
func (r *RabbitMQSubscriber) ListenQueueExclusive(queueName string) error {
	return r.listenQueueWithExclusive(queueName, true)
}

// listenQueueWithExclusive is the internal method to listen to a queue with optional exclusivity
func (r *RabbitMQSubscriber) listenQueueWithExclusive(queueName string, exclusive bool) error {
	go func() {
		for {
			// Open a new channel for each queue
			channel, err := r.connection.Channel()
			if err != nil {
				r.logger.Errorf("Failed to open channel for queue %s: %v. Retrying...", queueName, err)
				time.Sleep(5 * time.Second)
				continue
			}
			defer channel.Close()

			// Consume messages
			msgs, err := channel.Consume(
				queueName, // queue
				"",        // consumer tag
				false,     // auto-ack
				exclusive, // exclusive
				false,     // no-local
				false,     // no-wait
				nil,       // args
			)
			if err != nil {
				if exclusive && strings.Contains(err.Error(), "ACCESS_REFUSED") && strings.Contains(err.Error(), "exclusive") {
					r.logger.Warnf("Exclusive consumer conflict for queue %s. Only one instance can consume exclusively.", queueName)
					time.Sleep(10 * time.Second)
					continue
				}
				r.logger.Errorf("Failed to register consumer on queue %s: %v. Retrying...", queueName, err)
				time.Sleep(5 * time.Second)
				continue
			}

			r.logger.Infof("Listening for messages on queue %s", queueName)

			for d := range msgs {
				err := r.processMessage(d)
				if err != nil {
					r.logger.Errorf("Failed to process message: %v", err)
					r.retryAckNack(d, false)
				} else {
					r.retryAckNack(d, true)
				}
			}

			r.logger.Warn("Connection lost for queue %s. Reconnecting...", queueName)
			time.Sleep(5 * time.Second)
		}
	}()

	return nil
}

// retryAckNack attempts to acknowledge or negative acknowledge a message with retries
func (r *RabbitMQSubscriber) retryAckNack(d amqp091.Delivery, ack bool) {
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

	r.logger.Errorf("Failed to %s message after %d attempts", func() string {
		if ack {
			return "acknowledge"
		}
		return "negative acknowledge"
	}(), maxRetries)
}

func (r *RabbitMQSubscriber) Close() error {
	r.connectionMutex.Lock()
	defer r.connectionMutex.Unlock()

	if r.connection != nil && !r.connection.IsClosed() {
		return r.connection.Close()
	}
	return nil
}
