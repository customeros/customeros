package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/DusanKasan/parsemail"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"net/mail"
)

type mailService struct {
	services *Services
}

type MailService interface {
	SaveMail(ctx context.Context, request dto.MailRequest, email *parsemail.Email, tenant, user, customerOSInternalIdentifier string) (*string, error)
	SendMail(ctx context.Context, request dto.MailRequest, username *string) (*parsemail.Email, error)
}

func (s *mailService) SaveMail(ctx context.Context, request dto.MailRequest, email *parsemail.Email, tenant, user, customerOSInternalIdentifier string) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailService.SaveMail")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)

	threadId := email.Header.Get("Thread-Id")
	span.LogFields(tracingLog.String("threadId", threadId), tracingLog.String("user", user))

	interactionSessionId := ""

	interactionSessionNode, err := s.services.Neo4jRepositories.InteractionSessionReadRepository.GetByIdentifierAndChannel(ctx, tenant, threadId, "EMAIL")
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to get interaction session: %v", err)
	}
	if interactionSessionNode != nil {
		interactionSessionId = neo4jmapper.MapDbNodeToInteractionSessionEntity(interactionSessionNode).Id
	}

	session := utils.NewNeo4jWriteSession(ctx, *s.services.Neo4jRepositories.Neo4jDriver)
	defer session.Close(ctx)

	interactionEventId, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {

		if interactionSessionId == "" {
			sessionIdCreated, err := s.services.InteractionSessionService.CreateInTx(ctx, tx, &neo4jentity.InteractionSessionEntity{
				Status:     "ACTIVE",
				Type:       "THREAD",
				Channel:    "EMAIL",
				Identifier: threadId,
				Name:       utils.StringPtrFirstNonEmpty(request.Subject),

				AppSource: "user-admin-api",
			})
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}

			if sessionIdCreated == nil {
				tracing.TraceErr(span, errors.New("session id is empty"))
				return nil, fmt.Errorf("session id is empty")
			} else {
				interactionSessionId = *sessionIdCreated
			}
		}

		sentBy := make([]InteractionEventParticipantData, 0)
		sentTo := make([]InteractionEventParticipantData, 0)
		sentCc := make([]InteractionEventParticipantData, 0)
		sentBcc := make([]InteractionEventParticipantData, 0)

		for _, from := range email.From {
			sentBy = append(sentBy, InteractionEventParticipantData{
				Email: &from.Address,
			})
		}
		for _, to := range email.To {
			sentTo = append(sentTo, InteractionEventParticipantData{
				Email: &to.Address,
			})

		}
		for _, cc := range email.Cc {
			sentCc = append(sentCc, InteractionEventParticipantData{
				Email: &cc.Address,
			})
		}
		for _, bcc := range email.Bcc {
			sentBcc = append(sentBcc, InteractionEventParticipantData{
				Email: &bcc.Address,
			})
		}

		emailChannelData, err := dto.BuildEmailChannelData(email.Header.Get("Message-Id"), email.Header.Get("Thread-Id"), email.Subject, email.InReplyTo, email.References)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		interactionEventId, err := s.services.InteractionEventService.CreateInTx(ctx, tx, &InteractionEventCreateData{
			InteractionEventEntity: &neo4jentity.InteractionEventEntity{
				Content:                      utils.FirstNotEmptyString(email.HTMLBody, email.TextBody),
				ContentType:                  email.ContentType,
				Channel:                      "EMAIL",
				ChannelData:                  *emailChannelData,
				Identifier:                   utils.EnsureEmailRfcId(email.MessageID),
				CustomerOSInternalIdentifier: customerOSInternalIdentifier,
				Hide:                         false,
				Source:                       "openline", //TODO
				SourceOfTruth:                "openline", //TODO
				AppSource:                    "TODO",     //TODO
			},
			SentBy:            sentBy,
			SentTo:            sentTo,
			SentCc:            sentCc,
			SentBcc:           sentBcc,
			RepliesTo:         request.ReplyTo,
			SessionIdentifier: &interactionSessionId,
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, fmt.Errorf("failed to create interaction event: %v", err)
		}

		err = s.services.Neo4jRepositories.CommonWriteRepository.LinkEntityWithEntityInTx(ctx, tx, tenant, *interactionEventId, commonModel.INTERACTION_EVENT, commonModel.PART_OF, nil, interactionSessionId, commonModel.INTERACTION_SESSION)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, fmt.Errorf("failed to link interaction event with interaction session: %v", err)
		}

		return interactionEventId, nil
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return interactionEventId.(*string), nil
}

func (s *mailService) SendMail(ctx context.Context, request dto.MailRequest, username *string) (*parsemail.Email, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailService.SendMail")
	defer span.Finish()

	span.LogFields(tracingLog.Object("request", request))
	span.LogFields(tracingLog.String("username", *username))

	tenant, err := s.services.TenantService.GetTenantForUserEmail(ctx, *username)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("unable to retrieve tenant for %s", *username)
	}

	oauthToken, err := s.services.PostgresRepositories.OAuthTokenRepository.GetByEmail(ctx, tenant.Name, request.FromProvider, request.From)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("unable to retrieve oauth token for %s: %v", request.From, err)
	}

	if oauthToken == nil {
		mailbox, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, tenant.Name, request.From)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if mailbox == nil {
			return nil, fmt.Errorf("mailbox not found for %s", request.From)
		}

		return s.services.OpenSrsService.Reply(ctx, tenant.Name, request)
	} else {
		if oauthToken.NeedsManualRefresh {
			tracing.TraceErr(span, errors.New("oauth token needs manual refresh"))
			return nil, fmt.Errorf("oauth token needs manual refresh: %v", err)
		}

		if oauthToken.Provider == "google" {
			return s.services.GoogleService.SendEmail(ctx, tenant.Name, request)
		} else if oauthToken.Provider == "azure-ad" {
			return s.services.AzureService.SendEmail(ctx, tenant.Name, request)
		} else {
			return nil, fmt.Errorf("provider %s not supported", oauthToken.Provider)
		}
	}
}

func toParticipantInputArr(from []*mail.Address, participantType *string) []string {
	//var to []cosModel.InteractionEventParticipantInput
	//for _, a := range from {
	//	participantInput := cosModel.InteractionEventParticipantInput{
	//		Email: &a.Address,
	//		Type:  participantType,
	//	}
	//	to = append(to, participantInput)
	//}
	//return to
	return nil
}

func NewMailService(services *Services) MailService {
	return &mailService{
		services: services,
	}
}
