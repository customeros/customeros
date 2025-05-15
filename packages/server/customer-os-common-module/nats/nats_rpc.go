package nats_common

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/customeros/customeros/packages/server/enums"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type RpcConsumerConfig struct {
	StreamName          enums.NatsStream
	ServiceName         string
	SubscribedSubject   string
	MaxResponseSize     *int
	MaxDeliveryAttempts *int
}

const (
	DEFAULT_MAX_RESPONSE_SIZE = 1 * 1024 * 1024
)

type RpcEventsConsumer struct {
	NatsConn     *NATSConnections
	Config       *RpcConsumerConfig
	Handlers     map[string]EventHandler
	DLQPublisher func(ctx context.Context, msg *nats.Msg, err error)
	subscription *nats.Subscription
}

func NewRpcEventsConsumer(natsConn *NATSConnections, config *RpcConsumerConfig) (*RpcEventsConsumer, error) {
	validatedConfig, err := validateRpcConsumerConfig(config)
	if err != nil {
		return nil, err
	}

	return &RpcEventsConsumer{
		NatsConn: natsConn,
		Config:   validatedConfig,
		Handlers: make(map[string]EventHandler),
	}, nil
}

func (s *RpcEventsConsumer) RegisterHandler(subject string, handler EventHandler) {
	s.Handlers[subject] = handler
}

func (s *RpcEventsConsumer) Start(ctx context.Context) error {
	queueGroup := fmt.Sprintf("%s-queue-group", strings.ToLower(s.Config.ServiceName))

	stream, err := s.NatsConn.GetNatsConnection(s.Config.StreamName)
	if err != nil {
		return fmt.Errorf("failed to get nats connection: %w", err)
	}

	// Create a queue subscription for handling synchronous requests
	sub, err := stream.Conn.QueueSubscribe(
		s.Config.SubscribedSubject,
		queueGroup,
		func(msg *nats.Msg) {
			reqCtx := telemetry.ExtractTraceContextFromNatsMsg(context.Background(), msg)
			reqCtx = utils.WithCustomContextFromNats(reqCtx, msg)
			s.routeMessage(reqCtx, msg)
		})
	if err != nil {
		return fmt.Errorf("failed to create queue subscription: %w", err)
	}

	// Set subscription options
	sub.SetPendingLimits(-1, -1) // No limits on pending messages
	s.subscription = sub

	// Listen for context cancellation to clean up
	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	log.Printf("%s Service started and listening for requests (queue group: %s)", s.Config.ServiceName, queueGroup)
	return nil
}

func (s *RpcEventsConsumer) routeMessage(ctx context.Context, msg *nats.Msg) {
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
func (s *RpcEventsConsumer) Stop() {
	if s.subscription != nil {
		s.subscription.Unsubscribe()
	}
}

func validateRpcConsumerConfig(config *RpcConsumerConfig) (*RpcConsumerConfig, error) {
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
		err := errors.New("No subscribed subject provided")
		return nil, err
	}

	if config.MaxResponseSize == nil {
		size := DEFAULT_MAX_RESPONSE_SIZE
		config.MaxResponseSize = &size
	}

	if config.MaxDeliveryAttempts == nil {
		attempts := DEFAULT_MAX_DELIVERY_ATTEMPTS
		config.MaxDeliveryAttempts = &attempts
	}

	return config, nil
}
