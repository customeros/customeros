package event_logger

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/customeros/mailstack/proto/pb"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	leads_errors "github.com/customeros/customeros/packages/server/leads/errors"
	"github.com/customeros/customeros/packages/server/leads/interfaces"
	"github.com/customeros/customeros/packages/server/leads/internal/models"
	nats_internal "github.com/customeros/customeros/packages/server/leads/internal/nats"
	"github.com/customeros/customeros/packages/server/leads/internal/repository"
	"github.com/customeros/customeros/packages/server/leads/internal/telemetry"
	"github.com/customeros/customeros/packages/server/leads/internal/utils"
)

type LeadEventLoggerService struct {
	natsConn     *nats_internal.NATSConnections
	repositories *repository.Repositories
}

func NewLeadEventLoggerService(
	natsConn *nats_internal.NATSConnections,
	repos *repository.Repositories,
) interfaces.NatsService {
	return &LeadEventLoggerService{
		natsConn:     natsConn,
		repositories: repos,
	}
}

var SUBSCRIBED_SUBJECT = "leads.>"

func (s *LeadEventLoggerService) NewLeadEventRecord(ctx context.Context) *models.LeadEvent {
	spans, ctx := telemetry.StartServiceSpan(ctx, "LeadEventLoggerService.NewLeadEventRecord")
	defer spans.Finish()

	// validate context
	var err error
	var errorMessage string

	tenant := utils.GetTenantFromContext(ctx)
	userID := utils.GetUserIdFromContext(ctx)

	if tenant == "" {
		err = leads_errors.ErrTenantMissing
		spans.TraceError(err)
		errorMessage = err.Error()
	}

	if userID == "" {
		err = leads_errors.ErrUserIdMissing
		spans.TraceError(err)
		errorMessage = err.Error()
	}

	return &models.LeadEvent{
		ErrorMessage: errorMessage,
	}
}

// Start begins listening for raw email events and processing them
func (s *LeadEventLoggerService) Start(ctx context.Context) error {
	// Subscribe to all standard request/reply messages
	_, err := s.natsConn.Conn.Subscribe(SUBSCRIBED_SUBJECT, func(msg *nats.Msg) {
		spans, ctx := telemetry.StartServiceSpan(ctx, "LeadEventLoggerService.setupNonPersistedSubscriptions")
		defer spans.Finish()

		s.processMessage(ctx, msg)
	})
	if err != nil {
		return fmt.Errorf("failed to create standard subscription: %w", err)
	}

	log.Println("Email Logger Service started standard NATS subscriptions")
	return nil
}

// processMessage processes a single email message
func (s *LeadEventLoggerService) processMessage(ctx context.Context, msg *nats.Msg) {
	ctx = utils.WithCustomContextFromNats(ctx, msg)
	spans, ctx := telemetry.StartServiceSpan(ctx, "LeadEventLoggerService.processMessage")
	defer spans.Finish()

	// Extract the subject to determine message type
	subject := msg.Subject

	// Check if it's an error message
	if strings.HasPrefix(subject, "emails.errors.") {
		s.processErrorMessage(ctx, msg)
		return
	}

	// Dispatch based on subject pattern
	switch {
	// case strings.HasPrefix(subject, enum.EventEmailInboundClassifiedSkip.String()):
	//     s.processClassifiedSkipMessage(ctx, msg)
	//
	// case strings.HasPrefix(subject, enum.EventEmailInboundClassifiedBounce.String()):
	//     s.processClassifiedBounceMessage(ctx, msg)
	//
	// case strings.HasPrefix(subject, enum.EventEmailInboundClassifiedAutoresponder.String()):
	//     s.processClassifiedAutoresponderMessage(ctx, msg)

	default:
		err := errors.New("Unidentified message")
		spans.TraceError(err)
	}

	return
}

// Close gracefully shuts down the service
func (s *LeadEventLoggerService) Stop() {
	if s.natsConn != nil {
		s.natsConn.Close()
	}
	return
}

func (s *LeadEventLoggerService) publishError(ctx context.Context, msg *nats.Msg, err error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "LeadEventLoggerService.publishError")
	defer spans.Finish()

	errorEvent := &pb.ErrorEvent{
		Timestamp:    timestamppb.Now(),
		Subject:      msg.Subject,
		ErrorMessage: err.Error(),
		RawData:      msg.Data,
		// Publisher:    pb_mappers.MailstackServiceToServiceName(enum.MailstackEventLoggerService),
	}

	data, err := proto.Marshal(errorEvent)
	if err != nil {
		spans.TraceError(err)
		log.Printf("Failed to marshal error event: %v", err)
		return
	}

	_, pubErr := s.natsConn.JS.Publish("", data)
	if pubErr != nil {
		spans.TraceError(pubErr)
		log.Printf("Failed to publish error event: %v", pubErr)
	}
}
