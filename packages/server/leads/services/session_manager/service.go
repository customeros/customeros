package session_manager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/customeros/customeros/packages/server/leads/interfaces"
	"github.com/customeros/customeros/packages/server/leads/internal/enum"
	nats_internal "github.com/customeros/customeros/packages/server/leads/internal/nats"
	"github.com/customeros/customeros/packages/server/leads/internal/repository"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
)

type sessionManager struct {
	natsConn     *nats_internal.NATSConnections
	repositories *repository.Repositories
}

func NewSessionManager(
	natsConn *nats_internal.NATSConnections,
	repository *repository.Repositories,
) interfaces.NatsService {
	return &sessionManager{
		natsConn:     natsConn,
		repositories: repository,
	}
}

var SUBSCRIBED_SUBJECT = "webtracker.session.>"

const (
	// queue group
	QUEUE_GROUP = "session-manager"

	// consumer config
	CONSUMER_NAME         = "session-manager-consumer"
	ACK_WAIT              = 30 * time.Second
	MAX_DELIVERY_ATTEMPTS = 5
	MAX_ACK_PENDING       = 100
	FETCH_BATCH_SIZE      = 50
	MAX_FETCH_WAIT        = 500 * time.Millisecond
	ERR_BACKOFF           = 100 * time.Millisecond
)

// Start begins listening for raw email events and processing them
func (s *sessionManager) Start(ctx context.Context) error {
	// Create durable consumer for processing emails
	_, err := s.natsConn.JS.AddConsumer(nats_internal.LEADS_STREAM, &nats.ConsumerConfig{
		Durable:       CONSUMER_NAME,
		DeliverGroup:  QUEUE_GROUP,
		AckPolicy:     nats.AckExplicitPolicy,
		AckWait:       ACK_WAIT,
		MaxDeliver:    MAX_DELIVERY_ATTEMPTS,
		FilterSubject: SUBSCRIBED_SUBJECT,
		MaxAckPending: MAX_ACK_PENDING,
		DeliverPolicy: nats.DeliverAllPolicy,
	})
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	// Create pull subscription
	sub, err := s.natsConn.JS.PullSubscribe(
		SUBSCRIBED_SUBJECT,
		CONSUMER_NAME,
		nats.Bind(nats_internal.LEADS_STREAM, CONSUMER_NAME),
	)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	// Start processing
	go s.processRawEvents(ctx, sub)

	return nil
}

// processRawEmailEvents continuously processes raw email events
func (s *sessionManager) processRawEvents(ctx context.Context, sub *nats.Subscription) {
	log.Println("Session Manager started")
	for {
		select {
		case <-ctx.Done():
			log.Println("Session Manager shutting down")
			return
		default:
			s.processBatch(ctx, sub)
		}
	}
}

// processBatch fetches and processes a batch of messages
func (s *sessionManager) processBatch(ctx context.Context, sub *nats.Subscription) {
	// Fetch messages batch
	msgs, err := sub.Fetch(FETCH_BATCH_SIZE, nats.MaxWait(MAX_FETCH_WAIT))
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

// handleFetchError handles errors that occur during message fetching
func (s *sessionManager) handleFetchError(err error) {
	if errors.Is(err, nats.ErrTimeout) {
		// No messages available, this is normal
		return
	}
	log.Printf("Fetch error: %v", err)
	time.Sleep(ERR_BACKOFF) // Small backoff on error
}

func (s *sessionManager) routeMessage(ctx context.Context, msg *nats.Msg) {
	ctx = utils.WithCustomContextFromNats(ctx, msg)
	spans, ctx := telemetry.StartServiceSpan(ctx, "sessionManager.processMessage")
	defer spans.Finish()

	if msg == nil {
		spans.TraceError(errors.New("nil nats message"))
		return
	}
	spans.TagString("nats.subject", msg.Subject)
	spans.TagString("nats.reply", msg.Reply)

	switch {
	case msg.Subject == enum.EventWebtrackerSessionCreated.String():
		s.NewSession(ctx, msg)

	case msg.Subject == enum.EventWebtrackerSessionClosed.String():
		s.CloseSession(ctx, msg)

	default:
		spans.TraceError(errors.New("Unsupported event type"))
	}

	msg.Ack()
	return
}

// handleProcessingError deals with errors during email processing
func (s *sessionManager) handleProcessingError(ctx context.Context, msg *nats.Msg, err error) {
	metadata, _ := msg.Metadata()

	// Check if we should retry
	if metadata.NumDelivered <= uint64(MAX_DELIVERY_ATTEMPTS) {
		// Negative acknowledgment triggers redelivery
		msg.Nak()
	} else {
		// Max retries reached, acknowledge but publish to dead letter
		msg.Ack()
		// s.publishError(ctx, msg, err)
	}
}

// Close gracefully shuts down the service
func (s *sessionManager) Stop() {
	if s.natsConn != nil {
		s.natsConn.Close()
	}
	return
}
