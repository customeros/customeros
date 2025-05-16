package nats_common

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/customeros/customeros/packages/server/enums"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

type EventHandler func(ctx context.Context, msg *nats.Msg) error

type AsyncConsumerConfig struct {
	StreamName          enums.NatsStream
	ServiceName         enums.Services
	SubscribedSubject   string
	AckWait             *time.Duration
	MaxDeliveryAttempts *int
	MaxAckPending       *int
	FetchBatchSize      *int
	MaxFetchWait        *time.Duration
	ErrBackoff          *time.Duration
}

const (
	DEFAULT_ACK_WAIT              = 30 * time.Second
	DEFAULT_MAX_DELIVERY_ATTEMPTS = 3
	DEFAULT_MAX_ACK_PENDING       = 1000
	DEFAULT_FETCH_BATCH_SIZE      = 50
	DEFAULT_MAX_FETCH_WAIT        = 500 * time.Millisecond
	DEFAULT_ERR_BACKOFF           = 100 * time.Millisecond
)

type AsyncEventsConsumer struct {
	NatsConn     *NATSConnections
	Config       *AsyncConsumerConfig
	Handlers     map[string]EventHandler
	DLQPublisher func(ctx context.Context, msg *nats.Msg, err error)
}

func NewAsyncEventsConsumer(natsConn *NATSConnections, config *AsyncConsumerConfig) (*AsyncEventsConsumer, error) {
	validatedConfig, err := validateConsumerConfig(config)
	if err != nil {
		return nil, err
	}

	return &AsyncEventsConsumer{
		NatsConn: natsConn,
		Config:   validatedConfig,
		Handlers: make(map[string]EventHandler),
	}, nil
}

func (s *AsyncEventsConsumer) RegisterHandler(subject string, handler EventHandler) {
	s.Handlers[subject] = handler
}

func (s *AsyncEventsConsumer) Start(ctx context.Context) error {
	// Create durable consumer for processing messages
	consumer := fmt.Sprintf("%s-consumer", s.Config.ServiceName.String())
	queueGroup := fmt.Sprintf("%s", s.Config.ServiceName.String())

	stream, err := s.NatsConn.GetNatsConnection(s.Config.StreamName)
	if err != nil {
		return fmt.Errorf("failed to get nats connection: %w", err)
	}

	_, err = stream.JS.AddConsumer(s.Config.StreamName.String(), &nats.ConsumerConfig{
		Durable:       consumer,
		DeliverGroup:  queueGroup,
		AckPolicy:     nats.AckExplicitPolicy,
		AckWait:       *s.Config.AckWait,
		MaxDeliver:    *s.Config.MaxDeliveryAttempts,
		FilterSubject: s.Config.SubscribedSubject, // Note: For multi-subject, use more advanced setup
		MaxAckPending: *s.Config.MaxAckPending,
		DeliverPolicy: nats.DeliverAllPolicy,
	})
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	// Create pull subscription
	sub, err := stream.JS.PullSubscribe(
		s.Config.SubscribedSubject,
		consumer,
		nats.Bind(s.Config.StreamName.String(), consumer),
	)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	// Start processing
	go s.processEvents(ctx, sub)

	return nil
}

func (s *AsyncEventsConsumer) processEvents(ctx context.Context, sub *nats.Subscription) {
	log.Printf("%s service started", s.Config.ServiceName)
	for {
		select {
		case <-ctx.Done():
			log.Printf("%s service shutting down", s.Config.ServiceName)
			return
		default:
			s.processBatch(ctx, sub)
		}
	}
}

func (s *AsyncEventsConsumer) processBatch(ctx context.Context, sub *nats.Subscription) {
	msgs, err := sub.Fetch(*s.Config.FetchBatchSize, nats.MaxWait(*s.Config.MaxFetchWait))
	if err != nil {
		s.handleFetchError(err)
		return
	}

	for _, msg := range msgs {
		msgCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		s.routeMessage(msgCtx, msg)
		cancel()
	}
}

func (s *AsyncEventsConsumer) handleFetchError(err error) {
	if errors.Is(err, nats.ErrTimeout) {
		// No messages available, this is normal
		return
	}
	log.Printf("Fetch error: %v", err)
	time.Sleep(*s.Config.ErrBackoff)
}

func (s *AsyncEventsConsumer) routeMessage(ctx context.Context, msg *nats.Msg) {
	if msg == nil {
		log.Println("Error: received nil NATS message")
		return
	}

	handler, exists := s.Handlers[msg.Subject]
	if !exists {
		err := fmt.Errorf("no handler for subject: %s", msg.Subject)
		log.Println(err)
		if s.DLQPublisher != nil {
			s.DLQPublisher(ctx, msg, err)
		}
		msg.Ack()
		return
	}

	err := handler(ctx, msg)
	if err != nil {
		log.Printf("Error handling message: %v", err)
		metadata, _ := msg.Metadata()
		if metadata.NumDelivered <= uint64(*s.Config.MaxDeliveryAttempts) {
			msg.Nak()
		} else {
			msg.Ack()
			if s.DLQPublisher != nil {
				s.DLQPublisher(ctx, msg, err)
			}
		}
		return
	}

	msg.Ack()
}

// Stop gracefully shuts down the service
func (s *AsyncEventsConsumer) Stop() {
	// Any cleanup needed
}

func validateConsumerConfig(config *AsyncConsumerConfig) (*AsyncConsumerConfig, error) {
	if config == nil {
		err := errors.New("Consumer config is empty")
		return nil, err
	}

	if config.ServiceName == "" {
		err := errors.New("Service name is empty")
		return nil, err
	}

	if config.StreamName == "" {
		err := errors.New("Stream name is empty")
		return nil, err
	}

	if config.SubscribedSubject == "" {
		err := errors.New("No subscribed subjects provided")
		return nil, err
	}

	if config.AckWait == nil {
		wait := 30 * time.Second
		config.AckWait = &wait
	}

	if config.MaxDeliveryAttempts == nil {
		attempts := 3
		config.MaxDeliveryAttempts = &attempts
	}

	if config.MaxAckPending == nil {
		pending := 1000
		config.MaxAckPending = &pending
	}

	if config.FetchBatchSize == nil {
		batch := 50
		config.FetchBatchSize = &batch
	}

	if config.MaxFetchWait == nil {
		wait := 500 * time.Millisecond
		config.MaxFetchWait = &wait
	}

	if config.ErrBackoff == nil {
		backoff := 100 * time.Millisecond
		config.ErrBackoff = &backoff
	}

	return config, nil
}
