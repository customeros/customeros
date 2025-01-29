package events

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type SubscriberConfig struct {
	MaxRetries          int
	ReconnectBackoff    time.Duration
	MaxReconnectBackoff time.Duration
}

type RabbitMQSubscriber struct {
	connection      *amqp091.Connection
	connectionMutex sync.Mutex
	url             string
	logger          logger.Logger
	config          SubscriberConfig
	listeners       map[string]EventListener
	listenerMutex   sync.RWMutex
}

func NewRabbitMQSubscriber(rabbitmqURL string, logger logger.Logger, config *SubscriberConfig) (*RabbitMQSubscriber, error) {
	if config == nil {
		config = &SubscriberConfig{
			MaxRetries:          5,
			ReconnectBackoff:    time.Second,
			MaxReconnectBackoff: time.Second * 30,
		}
	}

	subscriber := &RabbitMQSubscriber{
		url:       rabbitmqURL,
		logger:    logger,
		config:    *config,
		listeners: make(map[string]EventListener),
	}

	err := subscriber.connect()
	if err != nil {
		return nil, err
	}

	return subscriber, nil
}

func (r *RabbitMQSubscriber) RegisterListener(listener EventListener) {
	r.listenerMutex.Lock()
	defer r.listenerMutex.Unlock()

	eventType := reflect.TypeOf(listener.GetEventType()).Name()
	r.listeners[eventType] = listener
	r.logger.Infof("Registered listener for event type: %s on queue: %s",
		eventType, listener.GetQueueName())
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
			channel, err := r.connection.Channel()
			if err != nil {
				r.logger.Errorf("Failed to open channel for queue %s: %v. Retrying...", queueName, err)
				time.Sleep(5 * time.Second)
				continue
			}
			defer channel.Close()

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
				r.handleMessage(d, queueName)
			}

			r.logger.Warnf("Connection lost for queue %s. Reconnecting...", queueName)
			time.Sleep(5 * time.Second)
		}
	}()

	return nil
}

func (r *RabbitMQSubscriber) handleMessage(d amqp091.Delivery, queueName string) {
	defer tracing.RecoverAndLogToJaeger(r.logger)

	err := r.processMessage(d, queueName)
	if err != nil {
		r.logger.Errorf("Failed to process message on queue %s: %v", queueName, err)
		r.retryAckNack(d, false)
	} else {
		r.retryAckNack(d, true)
	}
}

func (r *RabbitMQSubscriber) processMessage(d amqp091.Delivery, queueName string) error {
	ctx := context.Background()

	var event dto.Event
	if err := json.Unmarshal(d.Body, &event); err != nil {
		return errors.Wrap(err, "failed to unmarshal message")
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

	// Find the appropriate listener
	r.listenerMutex.RLock()
	listener, exists := r.listeners[event.Event.EventType]
	r.listenerMutex.RUnlock()

	if !exists {
		r.logger.Infof("No listener found for event type: %s on queue: %s", event.Event.EventType, queueName)
		return nil // No listener found, acknowledge the message
	}

	// Verify this listener is for this queue
	if listener.GetQueueName() != queueName {
		r.logger.Warnf("Event type %s received on wrong queue. Expected %s, got %s",
			event.Event.EventType, listener.GetQueueName(), queueName)
		return nil // Wrong queue, acknowledge the message
	}

	return listener.Handle(ctx, event.Event.Data)
}

func (r *RabbitMQSubscriber) connect() error {
	r.connectionMutex.Lock()
	defer r.connectionMutex.Unlock()

	var err error
	r.connection, err = amqp091.Dial(r.url)
	if err != nil {
		return errors.Wrap(err, "Failed to connect to RabbitMQ")
	}

	go func() {
		notifyClose := r.connection.NotifyClose(make(chan *amqp091.Error))
		<-notifyClose
		r.logger.Warn("RabbitMQ connection closed, attempting to reconnect")
		_ = r.connect()
	}()

	return nil
}

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
			return
		}

		time.Sleep(retryDelay)
	}

	r.logger.Errorf("Failed to %s message after %d attempts",
		map[bool]string{true: "acknowledge", false: "negative acknowledge"}[ack],
		maxRetries)
}

func (r *RabbitMQSubscriber) Close() error {
	r.connectionMutex.Lock()
	defer r.connectionMutex.Unlock()

	if r.connection != nil && !r.connection.IsClosed() {
		return r.connection.Close()
	}
	return nil
}
