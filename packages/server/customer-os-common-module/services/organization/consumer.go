package organization

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/customeros/customeros/packages/server/enums"

	"github.com/customeros/mailsherpa/mailvalidate"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/proto/pb"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

var ORGANIZATION_SUBJECT = "core.organization.>"
var WEBTRACKER_VISITOR_IDENTIFIED_SUBJECT = string(enums.EventWebtrackerVisitorIdentified)

const (
	// queue groups
	QUEUE_GROUP = "organization-service"

	// consumer configs
	WEBTRACKER_CONSUMER_NAME = "organization-service-webtracker-consumer"
	ACK_WAIT                 = 30 * time.Second
	MAX_DELIVERY_ATTEMPTS    = 5
	MAX_ACK_PENDING          = 100
	FETCH_BATCH_SIZE         = 50
	MAX_FETCH_WAIT           = 500 * time.Millisecond
	ERR_BACKOFF              = 100 * time.Millisecond
)

func (s *organizationService) Start(ctx context.Context) error {
	if s.natsConns == nil {
		return fmt.Errorf("NATS connection is nil")
	}
	webtrackerNatsConn, err := s.natsConns.GetNatsConnection(enums.StreamWebtracker)
	if err != nil {
		return fmt.Errorf("failed to get webtracker NATS connection: %w", err)
	}

	_, err = webtrackerNatsConn.JS.AddConsumer(enums.StreamWebtracker.String(), &nats.ConsumerConfig{
		Durable:       WEBTRACKER_CONSUMER_NAME,
		DeliverGroup:  QUEUE_GROUP,
		AckPolicy:     nats.AckExplicitPolicy,
		AckWait:       ACK_WAIT,
		MaxDeliver:    MAX_DELIVERY_ATTEMPTS,
		FilterSubject: WEBTRACKER_VISITOR_IDENTIFIED_SUBJECT,
		MaxAckPending: MAX_ACK_PENDING,
		DeliverPolicy: nats.DeliverAllPolicy,
	})
	if err != nil {
		return fmt.Errorf("failed to create webtracker consumer: %w", err)
	}

	webSub, err := webtrackerNatsConn.JS.PullSubscribe(
		WEBTRACKER_VISITOR_IDENTIFIED_SUBJECT,
		WEBTRACKER_CONSUMER_NAME,
		nats.Bind(enums.StreamWebtracker.String(), WEBTRACKER_CONSUMER_NAME),
	)
	if err != nil {
		return fmt.Errorf("failed to create webtracker subscription: %w", err)
	}

	go s.processWebtrackerVisitorIdentifiedEvents(ctx, webSub)

	return nil
}

// Stop gracefully shuts down the service
func (s *organizationService) Stop() {
	if s.natsConns != nil {
		s.natsConns.Close()
	}
}

func (s *organizationService) processWebtrackerVisitorIdentifiedEvents(ctx context.Context, sub *nats.Subscription) {
	log.Println("Webtracker visitor identified event processor started")
	for {
		select {
		case <-ctx.Done():
			log.Println("Webtracker visitor identified event processor shutting down")
			return
		default:
			s.processWebtrackerVisitorIdentifiedEventsBatch(sub)
		}
	}
}

func (s *organizationService) processWebtrackerVisitorIdentifiedEventsBatch(sub *nats.Subscription) {
	msgs, err := sub.Fetch(FETCH_BATCH_SIZE, nats.MaxWait(MAX_FETCH_WAIT))
	if err != nil {
		s.handleFetchError(err)
		return
	}

	for _, msg := range msgs {
		msgCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		s.handleWebtrackerVisitorIdentifiedMessage(msgCtx, msg)
		cancel()
	}
}

func (s *organizationService) handleFetchError(err error) {
	if errors.Is(err, nats.ErrTimeout) {
		// No messages available, this is normal
		return
	}
	log.Printf("Fetch error: %v", err)
	time.Sleep(ERR_BACKOFF)
}

func (s *organizationService) handleWebtrackerVisitorIdentifiedMessage(ctx context.Context, msg *nats.Msg) {
	ctx = common.WithCustomContextFromNats(ctx, msg)
	span, ctx := telemetry.StartListenerSpan(ctx, "OrganizationService.handleWebtrackerVisitorIdentifiedMessage", telemetry.WithNewRoot())
	defer span.Finish()

	if msg == nil {
		span.TraceError(errors.New("nil nats message"))
		return
	}
	span.TagString("nats.subject", msg.Subject)
	span.TagString("nats.reply", msg.Reply)

	// validate tenant
	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("tenant not set in nats header")
		span.TraceError(err)
		s.handleProcessingError(msg)
		return
	}

	request := &pb.WebtrackerVisitorIdentified{}
	err := proto.Unmarshal(msg.Data, request)
	if err != nil {
		span.TraceError(errors.Wrap(err, "failed to unmarshal webtracker visitor identified request"))
		s.handleProcessingError(msg)
		return
	}
	span.LogObjectAsJson("msg.webtracker.visitor.identified", request)

	// if domain is populated, find global organization by primary domain
	if request.Domain != "" {
		// find global organization by primary domain
		globalOrganization, err := s.postgres.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, request.Domain)
		if err != nil {
			span.TraceError(errors.Wrap(err, "failed to find global organization by primary domain"))
			s.handleProcessingError(msg)
			return
		}
		if globalOrganization != nil {
			// if found, create a new tenant organization
			orgId, err := s.CreateFromGlobalOrganization(ctx, nil, globalOrganization.ID, data_fields.OrganizationFields{})
			if err != nil {
				span.TraceError(errors.Wrap(err, "failed to create organization from global organization"))
				s.handleProcessingError(msg)
				return
			}
			span.LogKV("result.orgId", orgId)
		} else {
			// if not found, create a new organization with the domain
			orgId, err := s.Save(ctx, nil, nil, data_fields.OrganizationFields{
				Domains: []string{request.Domain},
			})
			if err != nil {
				span.TraceError(errors.Wrap(err, "failed to create organization from domain"))
				s.handleProcessingError(msg)
				return
			}
			span.LogKV("result.orgId", orgId)
		}
	} else if request.Email != "" {
		emailValidation := mailvalidate.ValidateEmailSyntax(request.Email)
		if !emailValidation.IsValid {
			span.LogKV("result", "invalid email")
			msg.Ack()
			return
		}
		isPersonalEmailProvider := s.emailService.IsPersonalEmailProvider(ctx, request.Email)
		if isPersonalEmailProvider {
			span.LogKV("result", "personal email provider")
			msg.Ack()
			return
		}
		// if not a personal email provider, create a new organization with the email domain
		domain := utils.ExtractDomain(request.Email)
		orgId, err := s.Save(ctx, nil, nil, data_fields.OrganizationFields{
			Domains: []string{domain},
		})
		if err != nil {
			span.TraceError(errors.Wrap(err, "failed to create organization from email"))
			s.handleProcessingError(msg)
			return
		}
		span.LogKV("result.orgId", orgId)
	}

	msg.Ack()
}

func (s *organizationService) handleProcessingError(msg *nats.Msg) {
	metadata, _ := msg.Metadata()

	// Check if we should retry
	if metadata.NumDelivered <= uint64(MAX_DELIVERY_ATTEMPTS) {
		// Negative acknowledgment triggers redelivery
		msg.Nak()
	} else {
		// Max retries reached, acknowledge but could publish to dead letter
		msg.Ack()
	}
}
