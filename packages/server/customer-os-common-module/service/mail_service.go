package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/DusanKasan/parsemail"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
)

type mailService struct {
	services *Services
}

type MailService interface {
	SendMail(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, request dto.MailRequest) (*string, error)
}

func (s *mailService) SendMail(ctx context.Context, tx *neo4j.ManagedTransaction, tenant string, request dto.MailRequest) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailService.SendMail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(tracingLog.Object("request", request))

	var parseEmail *parsemail.Email
	var err error

	oauthToken, err := s.services.PostgresRepositories.OAuthTokenRepository.GetByEmail(ctx, tenant, request.FromProvider, request.From)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("unable to retrieve oauth token for %s: %v", request.From, err)
	}

	uniqueInternalIdentifier := utils.GenerateRandomString(64)
	request.UniqueInternalIdentifier = &uniqueInternalIdentifier

	footer := `
					<div>
						<div style="font-size: 12px; font-weight: normal; font-family: Barlow, sans-serif; color: rgb(102, 112, 133); line-height: 32px;">
							<img width="16px" src="https://customer-os.imgix.net/website/favicon.png" alt="CustomerOS" style="vertical-align: middle; margin-right: 5px; margin-bottom: 2px;" />
							Sent from <a href="https://customeros.ai/?utm_content=timeline_email&utm_medium=email" style="text-decoration: underline; color: rgb(102, 112, 133);">CustomerOS</a>
						</div>
					</div>
					`
	request.Content += footer

	// Append an image tag pointing to the spy endpoint to the request content
	imgTag := "<img id=\"customer-os-email-track-open\" height=1 width=1 src=\"" + s.services.GlobalConfig.InternalServices.UserAdminApiPublicPath + "/mail/" + uniqueInternalIdentifier + "/track\" />"
	request.Content += imgTag

	if oauthToken == nil {
		mailbox, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, tenant, request.From)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if mailbox == nil {
			return nil, fmt.Errorf("mailbox not found for %s", request.From)
		}

		parseEmail, err = s.services.OpenSrsService.SendEmail(ctx, tenant, request)
	} else {
		if oauthToken.NeedsManualRefresh {
			tracing.TraceErr(span, errors.New("oauth token needs manual refresh"))
			return nil, fmt.Errorf("oauth token needs manual refresh: %v", err)
		}

		if oauthToken.Provider == "google" {
			parseEmail, err = s.services.GoogleService.SendEmail(ctx, tenant, request)
		} else if oauthToken.Provider == "azure-ad" {
			parseEmail, err = s.services.AzureService.SendEmail(ctx, tenant, request)
		} else {
			return nil, fmt.Errorf("provider %s not supported", oauthToken.Provider)
		}
	}

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to send email: %v", err)
	}

	//store email
	interactionEventId, err := s.saveMail(ctx, tx, request, parseEmail, uniqueInternalIdentifier)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return interactionEventId, nil
}

func (s *mailService) saveMail(ctx context.Context, tx *neo4j.ManagedTransaction, request dto.MailRequest, email *parsemail.Email, customerOSInternalIdentifier string) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailService.SaveMail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	interactionEventId := ""

	if tx == nil {
		session := utils.NewNeo4jWriteSession(ctx, *s.services.Neo4jRepositories.Neo4jDriver)
		defer session.Close(ctx)

		id, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			return s.saveEmailInTx(ctx, tx, request, email, customerOSInternalIdentifier)
		})
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if id != nil {
			interactionEventId = id.(string)
		}
	} else {
		id, err := s.saveEmailInTx(ctx, *tx, request, email, customerOSInternalIdentifier)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if id != nil {
			interactionEventId = *id.(*string)
		}
	}

	return &interactionEventId, nil
}

func (s *mailService) saveEmailInTx(ctx context.Context, tx neo4j.ManagedTransaction, request dto.MailRequest, email *parsemail.Email, customerOSInternalIdentifier string) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailService.saveEmailInTx")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	interactionSessionId := ""
	threadId := email.Header.Get("Thread-Id")
	span.LogFields(tracingLog.String("threadId", threadId))

	interactionSessionNode, err := s.services.Neo4jRepositories.InteractionSessionReadRepository.GetByIdentifierAndChannel(ctx, tenant, threadId, "EMAIL")
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to get interaction session: %v", err)
	}
	if interactionSessionNode != nil {
		interactionSessionId = neo4jmapper.MapDbNodeToInteractionSessionEntity(interactionSessionNode).Id
	}

	if interactionSessionId == "" {
		sessionIdCreated, err := s.services.InteractionSessionService.CreateInTx(ctx, tx, &neo4jentity.InteractionSessionEntity{
			Status:     "ACTIVE",
			Type:       "THREAD",
			Channel:    "EMAIL",
			Identifier: threadId,
			Name:       utils.StringPtrFirstNonEmpty(request.Subject),
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

	err = s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, &tx, tenant, repository.LinkDetails{
		FromEntityId:           *interactionEventId,
		FromEntityType:         commonModel.INTERACTION_EVENT,
		Relationship:           commonModel.PART_OF,
		RelationshipProperties: nil,
		ToEntityId:             interactionSessionId,
		ToEntityType:           commonModel.INTERACTION_SESSION,
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to link interaction event with interaction session: %v", err)
	}

	return interactionEventId, nil
}

func NewMailService(services *Services) MailService {
	return &mailService{
		services: services,
	}
}
