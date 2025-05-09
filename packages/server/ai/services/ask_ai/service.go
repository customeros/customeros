package ai

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/customeros/customeros/packages/server/ai/interfaces"
	"github.com/customeros/customeros/packages/server/ai/internal/enum"
	"github.com/customeros/customeros/packages/server/ai/internal/repository"
	"github.com/customeros/customeros/packages/server/ai/internal/telemetry"
	"github.com/customeros/customeros/packages/server/ai/internal/utils"
	nats_internal "github.com/customeros/customeros/packages/server/ai/nats"
	"github.com/customeros/customeros/packages/server/ai/services/anthropic"
	"github.com/customeros/customeros/packages/server/ai/services/deepseek"
	"github.com/customeros/customeros/packages/server/ai/services/gemini"
	"github.com/customeros/customeros/packages/server/ai/services/groq"
)

type aiService struct {
	natsConn     *nats_internal.NATSConnections
	anthropic    *anthropic.AnthropicService
	deepseek     *deepseek.DeepseekService
	groq         *groq.GroqService
	gemini       *gemini.GeminiService
	repositories *repository.Repositories
}

func NewAIService(
	natsConn *nats_internal.NATSConnections,
	anthropic *anthropic.AnthropicService,
	deepseek *deepseek.DeepseekService,
	groq *groq.GroqService,
	gemini *gemini.GeminiService,
	repositories *repository.Repositories,
) interfaces.NatsService {
	return &aiService{
		natsConn:     natsConn,
		anthropic:    anthropic,
		deepseek:     deepseek,
		groq:         groq,
		gemini:       gemini,
		repositories: repositories,
	}
}

var SUBSCRIBED_SUBJECT = "ai.>"

const (
	// queue group
	QUEUE_GROUP = "ai-service"

	// consumer config
	CONSUMER_NAME         = "ai-consumer"
	ACK_WAIT              = 30 * time.Second
	MAX_DELIVERY_ATTEMPTS = 5
	MAX_ACK_PENDING       = 100
	FETCH_BATCH_SIZE      = 50
	MAX_FETCH_WAIT        = 500 * time.Millisecond
	ERR_BACKOFF           = 100 * time.Millisecond
)

// Start begins listening for raw email events and processing them
func (s *aiService) Start(ctx context.Context) error {
	// Create durable consumer for processing emails
	_, err := s.natsConn.JS.AddConsumer(nats_internal.AI_STREAM, &nats.ConsumerConfig{
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
		nats.Bind(nats_internal.AI_STREAM, CONSUMER_NAME),
	)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	// Start processing
	go s.processRawEvents(ctx, sub)

	return nil
}

// processRawEvents continuously processes raw email events
func (s *aiService) processRawEvents(ctx context.Context, sub *nats.Subscription) {
	log.Println("AI Service started")
	for {
		select {
		case <-ctx.Done():
			log.Println("AI Service shutting down")
			return
		default:
			s.processBatch(ctx, sub)
		}
	}
}

// processBatch fetches and processes a batch of messages
func (s *aiService) processBatch(ctx context.Context, sub *nats.Subscription) {
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
func (s *aiService) handleFetchError(err error) {
	if errors.Is(err, nats.ErrTimeout) {
		// No messages available, this is normal
		return
	}
	log.Printf("Fetch error: %v", err)
	time.Sleep(ERR_BACKOFF) // Small backoff on error
}

// processMessage processes a single email message
func (s *aiService) routeMessage(ctx context.Context, msg *nats.Msg) {
	ctx = utils.WithCustomContextFromNats(ctx, msg)
	spans, ctx := telemetry.StartServiceSpan(ctx, "aiService.processMessage")
	defer spans.Finish()

	if msg == nil {
		spans.TraceError(errors.New("nil nats message"))
		return
	}
	spans.TagString("nats.subject", msg.Subject)

	var err error
	switch msg.Subject {
	case enum.AIRequestGeneric.String():
		err = s.AskAIForString(ctx, msg)
	// TODO handle all events
	default:
		err = errors.ErrUnsupported
	}

	if err != nil {
		if !strings.Contains(err.Error(), "skipping") {
			spans.TraceError(err)
		}
		s.handleProcessingError(ctx, msg, err)
		return
	}

	msg.Ack()
	return
}

// handleProcessingError deals with errors during email processing
func (s *aiService) handleProcessingError(ctx context.Context, msg *nats.Msg, err error) {
	metadata, _ := msg.Metadata()

	// Check if we should retry
	if metadata.NumDelivered <= uint64(MAX_DELIVERY_ATTEMPTS) {
		// Negative acknowledgment triggers redelivery
		msg.Nak()
	} else {
		// Max retries reached, acknowledge but publish to dead letter
		msg.Ack()
		// TODO
		// s.publishError(ctx, msg, err)
	}
}

// Close gracefully shuts down the service
func (s *aiService) Stop() {
	if s.natsConn != nil {
		s.natsConn.Close()
	}
	return
}

//TODO
// func (s *aiService) publishError(ctx context.Context, msg *nats.Msg, err error) {
// 	spans, ctx := telemetry.StartServiceSpan(ctx, "aiService.publishError")
// 	defer spans.Finish()
//
// 	errorEvent := &pb.ErrorEvent{
// 		Timestamp:    timestamppb.Now(),
// 		Subject:      msg.Subject,
// 		ErrorMessage: err.Error(),
// 		RawData:      msg.Data,
// 		Service:      pb.ServiceName_LEADS_PROXY_MANAGER_SERVICE,
// 	}
//
// 	data, err := proto.Marshal(errorEvent)
// 	if err != nil {
// 		spans.TraceError(err)
// 		log.Printf("Failed to marshal error event: %v", err)
// 		return
// 	}
//
// 	// Create message with headers
// 	newMsg := nats.NewMsg(enum.EventLeadError.String())
// 	newMsg.Data = data
// 	newMsg.Header.Set(nats_internal.HEADER_TENANT, utils.GetTenantFromContext(ctx))
// 	newMsg.Header.Set(nats_internal.HEADER_USERID, utils.GetUserIdFromContext(ctx))
//
// 	// Publish to the stored subject
// 	_, err = s.natsConn.JS.PublishMsg(newMsg)
// 	if err != nil {
// 		spans.TraceError(err)
// 		return
// 	}
// }
