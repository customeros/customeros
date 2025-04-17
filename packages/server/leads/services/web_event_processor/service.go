package web_event_processor

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/customeros/customeros/packages/server/leads/interfaces"
	nats_internal "github.com/customeros/customeros/packages/server/leads/internal/nats"
	"github.com/customeros/customeros/packages/server/leads/internal/repository"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
)

type WebEventProcessor struct {
	natsConn     *nats_internal.NATSConnections
	repositories *repository.Repositories
}

func NewWebEventProcessor(
	natsConn *nats_internal.NATSConnections,
	repositories *repository.Repositories,
) interfaces.NatsService {
	return &WebEventProcessor{
		natsConn:     natsConn,
		repositories: repositories,
	}
}

var SUBSCRIBED_SUBJECT = "eventstream.*.event.>"

const (
	// queue group
	QUEUE_GROUP = "web-event-processor"

	// consumer config
	CONSUMER_NAME         = "web-event-consumer"
	ACK_WAIT              = 30 * time.Second
	MAX_DELIVERY_ATTEMPTS = 5
	MAX_ACK_PENDING       = 1000
	FETCH_BATCH_SIZE      = 50
	MAX_FETCH_WAIT        = 500 * time.Millisecond
	ERR_BACKOFF           = 100 * time.Millisecond
)

// Start begins listening for raw email events and processing them
func (s *WebEventProcessor) Start(ctx context.Context) error {
	// Create durable consumer for processing emails
	_, err := s.natsConn.JS.AddConsumer(nats_internal.EVENTSTREAM_STREAM, &nats.ConsumerConfig{
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
		nats.Bind(nats_internal.EVENTSTREAM_STREAM, CONSUMER_NAME),
	)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	// Start processing
	go s.processRawEmailEvents(ctx, sub)

	return nil
}

// processRawEmailEvents continuously processes raw email events
func (s *WebEventProcessor) processRawEmailEvents(ctx context.Context, sub *nats.Subscription) {
	log.Println("Web event processor started")
	for {
		select {
		case <-ctx.Done():
			log.Println("Web event processor shutting down")
			return
		default:
			s.processBatch(ctx, sub)
		}
	}
}

// processBatch fetches and processes a batch of messages
func (s *WebEventProcessor) processBatch(ctx context.Context, sub *nats.Subscription) {
	// Fetch messages batch
	msgs, err := sub.Fetch(FETCH_BATCH_SIZE, nats.MaxWait(MAX_FETCH_WAIT))
	if err != nil {
		s.handleFetchError(err)
		return
	}

	for _, msg := range msgs {
		msgCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		s.processMessage(msgCtx, msg)
		cancel()
	}
}

// handleFetchError handles errors that occur during message fetching
func (s *WebEventProcessor) handleFetchError(err error) {
	if err == nats.ErrTimeout {
		// No messages available, this is normal
		return
	}
	log.Printf("Fetch error: %v", err)
	time.Sleep(ERR_BACKOFF) // Small backoff on error
}

// processMessage processes a single email message
func (s *WebEventProcessor) processMessage(ctx context.Context, msg *nats.Msg) {
	ctx = utils.WithCustomContextFromNats(ctx, msg)
	spans, ctx := telemetry.StartServiceSpan(ctx, "emailStorageService.processMessage")
	defer spans.Finish()

	subject := msg.Subject
	switch {
	case strings.HasSuffix(subject, "event.page_view"):
		s.handlePageView(ctx, msg)
	case strings.HasSuffix(subject, "event.page_exit"):
		s.handlePageExit(ctx, msg)
	case strings.HasSuffix(subject, "event.click"):
		s.handleClick(ctx, msg)
	case strings.HasSuffix(subject, "event.identify"):
		s.handleIdentify(ctx, msg)
	default:
		err := errors.New("Unidentified message")
		spans.TraceError(err)
	}
	return
}

// handleProcessingError deals with errors during email processing
func (s *WebEventProcessor) handleProcessingError(ctx context.Context, msg *nats.Msg, err error) {
	metadata, _ := msg.Metadata()

	// Check if we should retry
	if metadata.NumDelivered <= uint64(MAX_DELIVERY_ATTEMPTS) {
		// Negative acknowledgment triggers redelivery
		msg.Nak()
	} else {
		// Max retries reached, acknowledge but publish to dead letter
		msg.Ack()
		s.publishError(ctx, msg, err)
	}
}

// Stop gracefully shuts down the service
func (s *WebEventProcessor) Stop() error {
	if s.natsConn != nil {
		s.natsConn.Close()
	}
	return nil
}
