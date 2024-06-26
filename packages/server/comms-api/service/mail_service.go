package service

import (
	"context"
	"fmt"
	"github.com/DusanKasan/parsemail"
	c "github.com/openline-ai/openline-customer-os/packages/server/comms-api/config"
	"github.com/openline-ai/openline-customer-os/packages/server/comms-api/model"
	cosModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"log"
	"net/mail"
)

type mailService struct {
	services *Services
	config   *c.Config
}

type MailService interface {
	SaveMail(ctx context.Context, email *parsemail.Email, tenant, user, customerOSInternalIdentifier string) (*model.InteractionEventCreateResponse, error)
	SendMail(ctx context.Context, request dto.MailRequest, username *string) (*parsemail.Email, error)
}

func (s *mailService) SaveMail(ctx context.Context, email *parsemail.Email, tenant, user, customerOSInternalIdentifier string) (*model.InteractionEventCreateResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailService.SaveMail")
	defer span.Finish()

	threadId := email.Header.Get("Thread-Id")
	span.LogFields(tracingLog.String("threadId", threadId), tracingLog.String("tenant", tenant), tracingLog.String("user", user))

	sessionId, err := s.services.CustomerOsService.GetInteractionSession(&threadId, &tenant, &user)

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to get interaction session: %v", err)
	}

	channelValue := "EMAIL"
	appSource := "COMMS_API"
	sessionStatus := "ACTIVE"
	sessionType := "THREAD"
	if sessionId == nil {
		sessionOpts := []SessionOption{
			WithSessionIdentifier(&threadId),
			WithSessionChannel(&channelValue),
			WithSessionName(&email.Subject),
			WithSessionAppSource(&appSource),
			WithSessionStatus(&sessionStatus),
			WithSessionTenant(&tenant),
			WithSessionUsername(&user),
			WithSessionType(&sessionType),
		}

		sessionId, err = s.services.CustomerOsService.CreateInteractionSession(sessionOpts...)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, fmt.Errorf("failed to create interaction session: %v", err)
		}
	}

	participantTypeTO, participantTypeCC, participantTypeBCC := "TO", "CC", "BCC"
	participantsTO := toParticipantInputArr(email.To, &participantTypeTO)
	participantsCC := toParticipantInputArr(email.Cc, &participantTypeCC)
	participantsBCC := toParticipantInputArr(email.Bcc, &participantTypeBCC)
	sentTo := append(append(participantsTO, participantsCC...), participantsBCC...)
	sentBy := toParticipantInputArr(email.From, nil)

	emailChannelData, err := dto.BuildEmailChannelData(email.Header.Get("Message-Id"), email.Header.Get("Thread-Id"), email.Subject, email.InReplyTo, email.References)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	eventOpts := []EventOption{
		WithTenant(&tenant),
		WithUsername(&user),
		WithSessionId(sessionId),
		WithEventIdentifier(utils.EnsureEmailRfcId(email.MessageID)),
		WithExternalId(utils.EnsureEmailRfcId(email.MessageID)),
		WithExternalSystemId("gmail"),
		WithCustomerOSInternalIdentifier(customerOSInternalIdentifier),
		WithChannel(&channelValue),
		WithChannelData(emailChannelData),
		WithContent(utils.FirstNotEmpty(email.HTMLBody, email.TextBody)),
		WithContentType(&email.ContentType),
		WithSentBy(sentBy),
		WithSentTo(sentTo),
		WithAppSource(&appSource),
	}
	response, err := s.services.CustomerOsService.CreateInteractionEvent(eventOpts...)

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to create interaction event: %v", err)
	}

	return response, nil
}

func (s *mailService) SendMail(ctx context.Context, request dto.MailRequest, username *string) (*parsemail.Email, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailService.SendMail")
	defer span.Finish()

	span.LogFields(tracingLog.Object("request", request))
	span.LogFields(tracingLog.String("username", *username))

	tenant, err := s.services.CustomerOsService.GetTenant(username)
	if err != nil {
		log.Printf("unable to retrieve tenant for %s", *username)
		return nil, fmt.Errorf("unable to retrieve tenant for %s", *username)
	}

	oauthToken, err := s.services.CommonServices.PostgresRepositories.OAuthTokenRepository.GetByEmail(ctx, tenant.Tenant, request.FromProvider, request.From)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("unable to retrieve oauth token for %s: %v", request.From, err)
	}

	if oauthToken == nil {
		tracing.TraceErr(span, errors.New("unable to retrieve oauth token for new gmail service"))
		return nil, fmt.Errorf("unable to retrieve oauth token for new gmail service: %v", err)
	}

	if oauthToken.NeedsManualRefresh {
		tracing.TraceErr(span, errors.New("oauth token needs manual refresh"))
		return nil, fmt.Errorf("oauth token needs manual refresh: %v", err)
	}

	if oauthToken.Provider == "google" {
		return s.services.CommonServices.GoogleService.SendEmail(ctx, tenant.Tenant, request)
	} else if oauthToken.Provider == "azure-ad" {
		return s.services.CommonServices.AzureService.SendEmail(ctx, tenant.Tenant, request)
	} else {
		return nil, fmt.Errorf("provider %s not supported", oauthToken.Provider)
	}
}

func toParticipantInputArr(from []*mail.Address, participantType *string) []cosModel.InteractionEventParticipantInput {
	var to []cosModel.InteractionEventParticipantInput
	for _, a := range from {
		participantInput := cosModel.InteractionEventParticipantInput{
			Email: &a.Address,
			Type:  participantType,
		}
		to = append(to, participantInput)
	}
	return to
}

func NewMailService(config *c.Config, services *Services) MailService {
	return &mailService{
		config:   config,
		services: services,
	}
}
