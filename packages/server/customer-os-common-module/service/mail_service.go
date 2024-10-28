package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

type mailService struct {
	services *Services
}

type MailService interface {
	SendMail(ctx context.Context, emailMessage *entity.EmailMessage) error
	ProcessSentEmail(ctx context.Context, tx *neo4j.ManagedTransaction, emailMessage *entity.EmailMessage) (*string, error)
}

func (s *mailService) SendMail(ctx context.Context, emailMessage *entity.EmailMessage) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailService.SendMail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	span.LogFields(tracingLog.Object("emailMessage", emailMessage))

	tenant := emailMessage.Tenant

	var err error

	oauthToken, err := s.services.PostgresRepositories.OAuthTokenRepository.GetByEmail(ctx, tenant, emailMessage.FromProvider, emailMessage.From)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("unable to retrieve oauth token for %s: %v", emailMessage.From, err)
	}

	uniqueInternalIdentifier := utils.GenerateRandomString(64)
	emailMessage.UniqueInternalIdentifier = &uniqueInternalIdentifier

	//footer := `
	//				<div>
	//					<div style="font-size: 12px; font-weight: normal; font-family: Barlow, sans-serif; color: rgb(102, 112, 133); line-height: 32px;">
	//						<img width="16px" src="https://customer-os.imgix.net/website/favicon.png" alt="CustomerOS" style="vertical-align: middle; margin-right: 5px; margin-bottom: 2px;" />
	//						Sent from <a href="https://customeros.ai/?utm_content=timeline_email&utm_medium=email" style="text-decoration: underline; color: rgb(102, 112, 133);">CustomerOS</a>
	//					</div>
	//				</div>
	//				`
	//emailMessage.Content += footer

	// Append an image tag pointing to the spy endpoint to the request content
	//imgTag := "<img id=\"customer-os-email-track-open\" height=1 width=1 src=\"" + s.services.GlobalConfig.InternalServices.UserAdminApiPublicPath + "/mail/" + uniqueInternalIdentifier + "/track\" />"
	//emailMessage.Content += imgTag

	subject := ""
	inReplyTo := emailMessage.ProviderInReplyTo
	references := emailMessage.ProviderReferences

	if emailMessage.ReplyTo != nil {
		interactionEventNode, err := s.services.Neo4jRepositories.CommonReadRepository.GetById(ctx, tenant, *emailMessage.ReplyTo, commonModel.NodeLabelInteractionEvent)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		interactionEvent := neo4jmapper.MapDbNodeToInteractionEventEntity(interactionEventNode)

		emailChannelData := neo4jentity.EmailChannelData{}
		err = json.Unmarshal([]byte(interactionEvent.ChannelData), &emailChannelData)
		if err != nil {
			tracing.TraceErr(span, err)
			return fmt.Errorf("unable to parse email channel data for %s", *emailMessage.ReplyTo)
		}

		subject = emailChannelData.Subject
		if len(emailChannelData.Subject) < 3 || emailChannelData.Subject[:3] != "Re:" {
			subject = "Re: " + emailChannelData.Subject
		}

		if emailChannelData.Reference != "" {
			references = emailChannelData.Reference + " " + emailChannelData.ProviderMessageId
		} else {
			references = emailChannelData.ProviderMessageId
		}

		inReplyTo = emailChannelData.ProviderMessageId
	} else {
		subject = emailMessage.Subject
	}

	emailMessage.Subject = subject
	emailMessage.ProviderInReplyTo = inReplyTo
	emailMessage.ProviderReferences = references

	userNode, err := s.services.Neo4jRepositories.UserReadRepository.GetFirstUserByEmail(ctx, tenant, emailMessage.From)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get first user by email"))
		return err
	}

	if userNode == nil {
		tracing.TraceErr(span, errors.New("user not found"))
		return err
	}

	user := neo4jmapper.MapDbNodeToUserEntity(userNode)

	emailMessage.FromName = user.FirstName + " " + user.LastName

	if oauthToken == nil {
		mailbox, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, tenant, emailMessage.From)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if mailbox == nil {
			return fmt.Errorf("mailbox not found for %s", emailMessage.From)
		}

		err = s.services.OpenSrsService.SendEmail(ctx, tenant, emailMessage)
	} else {
		if oauthToken.NeedsManualRefresh {
			tracing.TraceErr(span, errors.New("oauth token needs manual refresh"))
			return fmt.Errorf("oauth token needs manual refresh: %v", err)
		}

		if oauthToken.Provider == "google" {
			err = s.services.GoogleService.SendEmail(ctx, tenant, emailMessage)
		} else if oauthToken.Provider == "azure-ad" {
			err = s.services.AzureService.SendEmail(ctx, tenant, emailMessage)
		} else {
			return fmt.Errorf("provider %s not supported", oauthToken.Provider)
		}
	}

	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to send email: %v", err)
	}

	emailMessage.Status = entity.EmailMessageStatusSent
	emailMessage.SentAt = utils.TimePtr(utils.Now())

	err = s.services.PostgresRepositories.EmailMessageRepository.Store(ctx, tenant, emailMessage)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to store email message: %v", err)
	}

	return nil
}

func (s *mailService) ProcessSentEmail(ctx context.Context, tx *neo4j.ManagedTransaction, emailMessage *entity.EmailMessage) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailService.ProcessSentEmail")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	id, err := utils.ExecuteWriteInTransaction(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		return s.saveEmailInTx(ctx, tx, emailMessage)
	})

	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	emailMessage.Status = entity.EmailMessageStatusProcessed
	err = s.services.PostgresRepositories.EmailMessageRepository.Store(ctx, emailMessage.Tenant, emailMessage)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to store email message: %v", err)
	}

	return id.(*string), nil
}

func (s *mailService) saveEmailInTx(ctx context.Context, tx neo4j.ManagedTransaction, emailMessage *entity.EmailMessage) (any, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "MailService.saveEmailInTx")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	tenant := common.GetTenantFromContext(ctx)

	interactionSessionId := ""
	threadId := emailMessage.ProviderThreadId
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
			Name:       emailMessage.Subject,
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

	sentBy = append(sentBy, InteractionEventParticipantData{
		Email: &emailMessage.From,
	})

	for _, to := range emailMessage.To {
		sentTo = append(sentTo, InteractionEventParticipantData{
			Email: &to,
		})

	}
	for _, cc := range emailMessage.Cc {
		sentCc = append(sentCc, InteractionEventParticipantData{
			Email: &cc,
		})
	}
	for _, bcc := range emailMessage.Bcc {
		sentBcc = append(sentBcc, InteractionEventParticipantData{
			Email: &bcc,
		})
	}

	emailChannelData, err := neo4jentity.BuildEmailChannelData(emailMessage.ProviderMessageId, emailMessage.ProviderThreadId, emailMessage.Subject, emailMessage.ProviderInReplyTo, emailMessage.ProviderReferences)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	interactionEventId, err := s.services.InteractionEventService.CreateInTx(ctx, tx, &InteractionEventCreateData{
		InteractionEventEntity: &neo4jentity.InteractionEventEntity{
			Content:                      emailMessage.Content,
			ContentType:                  "text/html",
			Channel:                      "EMAIL",
			ChannelData:                  *emailChannelData,
			Identifier:                   emailMessage.ProviderMessageId,
			CustomerOSInternalIdentifier: *emailMessage.UniqueInternalIdentifier,
			Hide:                         false,
			Source:                       "openline", //TODO
			SourceOfTruth:                "openline", //TODO
			AppSource:                    "TODO",     //TODO
		},
		SentBy:            sentBy,
		SentTo:            sentTo,
		SentCc:            sentCc,
		SentBcc:           sentBcc,
		RepliesTo:         emailMessage.ReplyTo,
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
