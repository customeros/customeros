package organization

import (
	"fmt"

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
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
)

var ORGANIZATION_SUBJECT = "core.organization.>"
var WEBTRACKER_VISITOR_IDENTIFIED_SUBJECT = string(leads_enum.EventWebtrackerVisitorIdentified)

const QUEUE_GROUP = "organization-service"

func (s *organizationService) Start(ctx context.Context) error {
	if s.natsConn == nil {
		return fmt.Errorf("NATS connection is nil")
	}

	// Create subscriptions for handling requests
	subs := []*nats.Subscription{}

	// Subscribe to organization subject
	orgSub, err := s.natsConn.Conn.QueueSubscribe(ORGANIZATION_SUBJECT, QUEUE_GROUP, func(msg *nats.Msg) {
		s.handleOrganizationMessage(ctx, msg)
	})
	if err != nil {
		return fmt.Errorf("failed to create organization subscription: %w", err)
	}
	subs = append(subs, orgSub)

	// Subscribe to webtracker subject
	webSub, err := s.natsConn.Conn.QueueSubscribe(WEBTRACKER_VISITOR_IDENTIFIED_SUBJECT, QUEUE_GROUP, func(msg *nats.Msg) {
		s.handleWebtrackerMessage(ctx, msg)
	})
	if err != nil {
		return fmt.Errorf("failed to create webtracker subscription: %w", err)
	}
	subs = append(subs, webSub)

	// Keep track of subscriptions for cleanup
	s.subscriptions = subs

	// Listen for context cancellation to clean up
	go func() {
		<-ctx.Done()
		for _, sub := range s.subscriptions {
			sub.Unsubscribe()
		}
	}()

	return nil
}

// Close gracefully shuts down the service
func (s *organizationService) Stop() {
	if s.natsConn != nil {
		s.natsConn.Close()
	}
	return
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
		s.sendOrganizationResponse(ctx, msg, resp)
		spans.TraceError(err)
		return
	}

	if resp == nil {
		spans.TraceError(errors.New("empty response"))
		return
	}

	s.sendOrganizationResponse(ctx, msg, resp)
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
		span.TraceError(errors.New("tenant not set in nats message"))
		return
	}

	request := &leads_pb.WebtrackerVisitorIdentified{}
	err := proto.Unmarshal(msg.Data, request)
	if err != nil {
		span.TraceError(errors.Wrap(err, "failed to unmarshal webtracker visitor identified request"))
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
