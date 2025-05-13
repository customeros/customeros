package organization

import (
	"context"
	"fmt"
	"log"
	"time"

	nats_common "github.com/customeros/customeros/packages/server/customer-os-common-module/nats"

	"github.com/customeros/mailsherpa/mailvalidate"

	core_crm_pb "github.com/customeros/customeros/packages/server/core-crm/proto/pb"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	leads_enum "github.com/customeros/leads/enum"
	leads_pb "github.com/customeros/leads/proto/pb"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

var ORGANIZATION_SUBJECT = "core.organization.>"
var WEBTRACKER_VISITOR_IDENTIFIED_SUBJECT = string(leads_enum.EventWebtrackerVisitorIdentified)

const (
	// queue groups
	QUEUE_GROUP = "organization-service"

	// consumer configs
	ORGANIZATION_CONSUMER_NAME = "organization-service-consumer"
	WEBTRACKER_CONSUMER_NAME   = "organization-service-webtracker-consumer"
	ACK_WAIT                   = 30 * time.Second
	MAX_DELIVERY_ATTEMPTS      = 5
	MAX_ACK_PENDING            = 100
	FETCH_BATCH_SIZE           = 50
	MAX_FETCH_WAIT             = 500 * time.Millisecond
	ERR_BACKOFF                = 100 * time.Millisecond
)

func (s *organizationService) Start(ctx context.Context) error {
	if s.natsConn == nil {
		return fmt.Errorf("NATS connection is nil")
	}

	// Create durable consumers for both subjects
	_, err := s.natsConn.JS.AddConsumer(nats_common.CORE_STREAM, &nats.ConsumerConfig{
		Durable:       ORGANIZATION_CONSUMER_NAME,
		DeliverGroup:  QUEUE_GROUP,
		AckPolicy:     nats.AckExplicitPolicy,
		AckWait:       ACK_WAIT,
		MaxDeliver:    MAX_DELIVERY_ATTEMPTS,
		FilterSubject: ORGANIZATION_SUBJECT,
		MaxAckPending: MAX_ACK_PENDING,
		DeliverPolicy: nats.DeliverAllPolicy,
	})
	if err != nil {
		return fmt.Errorf("failed to create organization consumer: %w", err)
	}

	_, err = s.natsConn.JS.AddConsumer(nats_common.LEADS_STREAM, &nats.ConsumerConfig{
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

	// Create pull subscriptions
	orgSub, err := s.natsConn.JS.PullSubscribe(
		ORGANIZATION_SUBJECT,
		ORGANIZATION_CONSUMER_NAME,
		nats.Bind("CUSTOMER_OS_STREAM", ORGANIZATION_CONSUMER_NAME),
	)
	if err != nil {
		return fmt.Errorf("failed to create organization subscription: %w", err)
	}

	webSub, err := s.natsConn.JS.PullSubscribe(
		WEBTRACKER_VISITOR_IDENTIFIED_SUBJECT,
		WEBTRACKER_CONSUMER_NAME,
		nats.Bind("CUSTOMER_OS_STREAM", WEBTRACKER_CONSUMER_NAME),
	)
	if err != nil {
		return fmt.Errorf("failed to create webtracker subscription: %w", err)
	}

	// Keep track of subscriptions for cleanup
	s.subscriptions = []*nats.Subscription{orgSub, webSub}

	// Start processing for both subscriptions
	go s.processOrganizationEvents(ctx, orgSub)
	go s.processWebtrackerEvents(ctx, webSub)

	return nil
}

func (s *organizationService) processOrganizationEvents(ctx context.Context, sub *nats.Subscription) {
	log.Println("Organization event processor started")
	for {
		select {
		case <-ctx.Done():
			log.Println("Organization event processor shutting down")
			return
		default:
			s.processOrganizationBatch(ctx, sub)
		}
	}
}

func (s *organizationService) processWebtrackerEvents(ctx context.Context, sub *nats.Subscription) {
	log.Println("Webtracker event processor started")
	for {
		select {
		case <-ctx.Done():
			log.Println("Webtracker event processor shutting down")
			return
		default:
			s.processWebtrackerBatch(ctx, sub)
		}
	}
}

func (s *organizationService) processOrganizationBatch(ctx context.Context, sub *nats.Subscription) {
	msgs, err := sub.Fetch(FETCH_BATCH_SIZE, nats.MaxWait(MAX_FETCH_WAIT))
	if err != nil {
		s.handleFetchError(err)
		return
	}

	for _, msg := range msgs {
		msgCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		s.handleOrganizationMessage(msgCtx, msg)
		cancel()
	}
}

func (s *organizationService) processWebtrackerBatch(ctx context.Context, sub *nats.Subscription) {
	msgs, err := sub.Fetch(FETCH_BATCH_SIZE, nats.MaxWait(MAX_FETCH_WAIT))
	if err != nil {
		s.handleFetchError(err)
		return
	}

	for _, msg := range msgs {
		msgCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		s.handleWebtrackerMessage(msgCtx, msg)
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

func (s *organizationService) handleOrganizationMessage(ctx context.Context, msg *nats.Msg) {
	ctx = common.WithCustomContextFromNats(ctx, msg)
	spans, ctx := telemetry.StartListenerSpan(ctx, "organizationService.handleOrganizationMessage")
	defer spans.Finish()

	if msg == nil {
		spans.TraceError(errors.New("nil nats message"))
		return
	}
	spans.TagString("nats.subject", msg.Subject)
	spans.TagString("nats.reply", msg.Reply)

	resp := &core_crm_pb.OrganizationSaveResponse{}

	request := &core_crm_pb.OrganizationSaveRequest{}
	err := proto.Unmarshal(msg.Data, request)
	if err != nil {
		s.handleProcessingError(ctx, msg, err)
		return
	}

	if resp == nil {
		s.handleProcessingError(ctx, msg, errors.New("empty response"))
		return
	}

	s.sendOrganizationResponse(ctx, msg, resp)
	msg.Ack()
}

func (s *organizationService) handleWebtrackerMessage(ctx context.Context, msg *nats.Msg) {
	ctx = common.WithCustomContextFromNats(ctx, msg)
	span, ctx := telemetry.StartListenerSpan(ctx, "organizationService.handleWebtrackerMessage")
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
		s.handleProcessingError(ctx, msg, errors.New("tenant not set in nats message"))
		return
	}

	request := &leads_pb.WebtrackerVisitorIdentified{}
	err := proto.Unmarshal(msg.Data, request)
	if err != nil {
		s.handleProcessingError(ctx, msg, errors.Wrap(err, "failed to unmarshal webtracker visitor identified request"))
		return
	}
	span.LogObjectAsJson("msg.webtracker.visitor.identified", request)

	// if domain is populated, find global organization by primary domain
	if request.Domain != "" {
		// find global organization by primary domain
		globalOrganization, err := s.postgres.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, request.Domain)
		if err != nil {
			span.TraceError(errors.Wrap(err, "failed to find global organization by primary domain"))
			return
		}
		if globalOrganization != nil {
			// if found, create a new tenant organization
			orgId, err := s.CreateFromGlobalOrganization(ctx, nil, globalOrganization.ID, data_fields.OrganizationFields{})
			if err != nil {
				span.TraceError(errors.Wrap(err, "failed to create organization from global organization"))
				return
			}
			span.LogKV("result.orgId", orgId)
			return
		} else {
			// if not found, create a new organization with the domain
			orgId, err := s.Save(ctx, nil, nil, data_fields.OrganizationFields{
				Domains: []string{request.Domain},
			})
			if err != nil {
				span.TraceError(errors.Wrap(err, "failed to create organization from domain"))
				return
			}
			span.LogKV("result.orgId", orgId)
			return
		}
	} else if request.Email != "" {
		emailValidation := mailvalidate.ValidateEmailSyntax(request.Email)
		if !emailValidation.IsValid {
			span.LogKV("result", "invalid email")
			return
		}
		isPersonalEmailProvider := s.emailService.IsPersonalEmailProvider(ctx, request.Email)
		if isPersonalEmailProvider {
			span.LogKV("result", "personal email provider")
			return
		}
		// if not a personal email provider, create a new organization with the email domain
		domain := utils.ExtractDomain(request.Email)
		orgId, err := s.Save(ctx, nil, nil, data_fields.OrganizationFields{
			Domains: []string{domain},
		})
		if err != nil {
			span.TraceError(errors.Wrap(err, "failed to create organization from email"))
			return
		}
		span.LogKV("result.orgId", orgId)
		return
	}

	msg.Ack()
}

func (s *organizationService) handleProcessingError(ctx context.Context, msg *nats.Msg, err error) {
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

func (s *organizationService) sendOrganizationResponse(ctx context.Context, req *nats.Msg, resp *core_crm_pb.OrganizationSaveResponse) {
	spans, _ := telemetry.StartServiceSpan(ctx, "organizationService.sendOrganizationResponse")
	defer spans.Finish()

	respMessage, err := proto.Marshal(resp)
	if err != nil {
		spans.TraceError(err)
		return
	}
	req.Respond(respMessage)
}
