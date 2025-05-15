package organization

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/enums"

	"github.com/customeros/mailsherpa/mailvalidate"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	nats_common "github.com/customeros/customeros/packages/server/customer-os-common-module/nats"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/proto/pb"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

const (
	SERVICE = "organization-service"
)

func (s *organizationService) Start(ctx context.Context) error {
	if s.natsConns == nil {
		return fmt.Errorf("NATS connection is nil")
	}

	webtrackerVisitorIdentifiedConsumer, err := s.setupWebtrackerVisitorIdentifiedConsumer()
	if err != nil {
		return fmt.Errorf("failed to setup webtracker consumer: %w", err)
	}
	err = webtrackerVisitorIdentifiedConsumer.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start webtracker visitor identified consumer: %w", err)
	}

	tenantCreatedConsumer, err := s.setupTenantCreatedConsumer()
	if err != nil {
		return fmt.Errorf("failed to setup tenant created consumer: %w", err)
	}
	err = tenantCreatedConsumer.Start(ctx)
	if err != nil {
		return fmt.Errorf("failed to start tenant created consumer: %w", err)
	}

	return nil
}

func (s *organizationService) setupWebtrackerVisitorIdentifiedConsumer() (*nats_common.AsyncEventsConsumer, error) {
	config := &nats_common.AsyncConsumerConfig{
		StreamName:        enums.StreamWebtracker,
		ServiceName:       SERVICE,
		SubscribedSubject: string(enums.EventWebtrackerVisitorIdentified),
	}

	consumer, err := nats_common.NewAsyncEventsConsumer(s.natsConns, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create async consumer: %w", err)
	}

	consumer.RegisterHandler(string(enums.EventWebtrackerVisitorIdentified), s.handleWebtrackerVisitorIdentifiedMessage)

	return consumer, nil
}

func (s *organizationService) setupTenantCreatedConsumer() (*nats_common.AsyncEventsConsumer, error) {
	config := &nats_common.AsyncConsumerConfig{
		StreamName:        enums.StreamTenant,
		ServiceName:       SERVICE,
		SubscribedSubject: string(enums.EventTenantCreated),
	}

	consumer, err := nats_common.NewAsyncEventsConsumer(s.natsConns, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create async consumer: %w", err)
	}

	consumer.RegisterHandler(string(enums.EventTenantCreated), s.handleTenantCreatedMessage)

	return consumer, nil
}

// Stop gracefully shuts down the service
func (s *organizationService) Stop() {
	if s.natsConns != nil {
		s.natsConns.Close()
	}
}

// handleWebtrackerVisitorIdentifiedMessage processes webtracker visitor identified events.
// Returns nil for success (message will be acked) or error for failure (message will be retried).
// Acking/Nacking is handled by the AsyncEventsConsumer based on the returned error.
func (s *organizationService) handleWebtrackerVisitorIdentifiedMessage(ctx context.Context, msg *nats.Msg) error {
	ctx = common.WithCustomContextFromNats(ctx, msg)
	span, ctx := telemetry.StartListenerSpan(ctx, "OrganizationService.handleWebtrackerVisitorIdentifiedMessage", telemetry.WithNewRoot())
	defer span.Finish()

	if msg == nil {
		err := errors.New("nil nats message")
		span.TraceError(err)
		return err
	}
	span.TagString("nats.subject", msg.Subject)

	// validate tenant
	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("tenant not set in nats header")
		span.TraceError(err)
		return err
	}

	request := &pb.WebtrackerVisitorIdentified{}
	err := proto.Unmarshal(msg.Data, request)
	if err != nil {
		span.TraceError(errors.Wrap(err, "failed to unmarshal webtracker visitor identified request"))
		return err
	}
	span.LogObjectAsJson("msg.webtracker.visitor.identified", request)

	// if domain is populated, find global organization by primary domain
	if request.Domain != "" {
		// find global organization by primary domain
		globalOrganization, err := s.postgres.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, request.Domain)
		if err != nil {
			span.TraceError(errors.Wrap(err, "failed to find global organization by primary domain"))
			return err
		}
		if globalOrganization != nil {
			// if found, create a new tenant organization
			orgId, err := s.CreateFromGlobalOrganization(ctx, nil, globalOrganization.ID, data_fields.OrganizationFields{})
			if err != nil {
				span.TraceError(errors.Wrap(err, "failed to create organization from global organization"))
				return err
			}
			span.LogKV("result.orgId", orgId)
		} else {
			// if not found, create a new organization with the domain
			orgId, err := s.Save(ctx, nil, nil, data_fields.OrganizationFields{
				Domains: []string{request.Domain},
			})
			if err != nil {
				span.TraceError(errors.Wrap(err, "failed to create organization from domain"))
				return err
			}
			span.LogKV("result.orgId", orgId)
		}
	} else if request.Email != "" {
		emailValidation := mailvalidate.ValidateEmailSyntax(request.Email)
		if !emailValidation.IsValid {
			span.LogKV("result", "invalid email")
			return nil
		}
		isPersonalEmailProvider := s.emailService.IsPersonalEmailProvider(ctx, request.Email)
		if isPersonalEmailProvider {
			span.LogKV("result", "personal email provider")
			return nil
		}
		// if not a personal email provider, create a new organization with the email domain
		domain := utils.ExtractDomain(request.Email)
		orgId, err := s.Save(ctx, nil, nil, data_fields.OrganizationFields{
			Domains: []string{domain},
		})
		if err != nil {
			span.TraceError(errors.Wrap(err, "failed to create organization from email"))
			return err
		}
		span.LogKV("result.orgId", orgId)
	}

	return nil
}

func (s *organizationService) handleTenantCreatedMessage(ctx context.Context, msg *nats.Msg) error {
	ctx = common.WithCustomContextFromNats(ctx, msg)
	span, ctx := telemetry.StartListenerSpan(ctx, "OrganizationService.handleTenantCreatedMessage", telemetry.WithNewRoot())
	defer span.Finish()

	if msg == nil {
		err := errors.New("nil nats message")
		span.TraceError(err)
		return err
	}
	span.TagString("nats.subject", msg.Subject)

	// validate tenant
	tenant := common.GetTenantFromContext(ctx)
	if tenant == "" {
		err := errors.New("tenant not set in nats header")
		span.TraceError(err)
		return err
	}

	request := &pb.TenantCreated{}
	err := proto.Unmarshal(msg.Data, request)
	if err != nil {
		span.TraceError(errors.Wrap(err, "failed to unmarshal tenant created request"))
		return err
	}
	span.LogObjectAsJson("msg.tenant.created", request)

	if request.Domain == "" {
		span.LogKV("result", "no domain provided")
		return nil
	}

	// if domain is populated, find global organization by primary domain
	if request.Domain != "" {
		// find global organization by primary domain
		globalOrganization, err := s.postgres.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, request.Domain)
		if err != nil {
			span.TraceError(errors.Wrap(err, "failed to find global organization by primary domain"))
			return err
		}
		if globalOrganization != nil {
			// if found, create a new tenant organization
			orgId, err := s.CreateFromGlobalOrganization(ctx, nil, globalOrganization.ID, data_fields.OrganizationFields{
				LeadSource: utils.StringPtr("Tenant Registration"),
			})
			if err != nil {
				span.TraceError(errors.Wrap(err, "failed to create organization from global organization"))
				return err
			}
			span.LogKV("result.orgId", orgId)
		} else {
			// if not found, create a new organization with the domain
			orgId, err := s.Save(ctx, nil, nil, data_fields.OrganizationFields{
				Domains:    []string{request.Domain},
				LeadSource: utils.StringPtr("Tenant Registration"),
			})
			if err != nil {
				span.TraceError(errors.Wrap(err, "failed to create organization from domain"))
				return err
			}
			span.LogKV("result.orgId", orgId)
		}
	}

	return nil
}
