package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	tracingLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

func (s *mailService) SendMail(ctx context.Context, emailMessage *postgresentity.EmailMessage) error {
	span, ctx := s.initializeTracing(ctx, "MailService.SendMail")
	span.LogFields(tracingLog.Object("emailMessage", emailMessage))
	defer span.Finish()

	oauthToken, err := s.getOAuthToken(ctx, span, emailMessage)
	if err != nil {
		return err
	}

	if err := s.prepareEmailMessage(ctx, span, emailMessage); err != nil {
		return err
	}

	if err := s.sendEmailBasedOnProvider(ctx, span, emailMessage, oauthToken); err != nil {
		return err
	}

	return s.storeEmailMessage(ctx, span, emailMessage)
}

func (s *mailService) ProcessSentEmail(ctx context.Context, tx *neo4j.ManagedTransaction, emailMessage *postgresentity.EmailMessage) (*string, error) {
	span, ctx := s.initializeTracing(ctx, "MailService.ProcessSentEmail")
	span.LogFields(tracingLog.Object("emailMessage", emailMessage))
	defer span.Finish()

	id, err := utils.ExecuteWriteInTransaction(ctx, s.services.Neo4jRepositories.Neo4jDriver, s.services.Neo4jRepositories.Database, tx, func(tx neo4j.ManagedTransaction) (any, error) {
		return s.saveEmailInTx(ctx, tx, emailMessage, span)
	})
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	emailMessage.Status = postgresentity.EmailMessageStatusProcessed
	err = s.services.PostgresRepositories.EmailMessageRepository.Store(ctx, emailMessage.Tenant, emailMessage)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("failed to store email message: %v", err)
	}

	return id.(*string), nil
}

func (s *mailService) getOAuthToken(ctx context.Context, span opentracing.Span, emailMessage *postgresentity.EmailMessage) (*postgresentity.OAuthTokenEntity, error) {
	oauthToken, err := s.services.PostgresRepositories.OAuthTokenRepository.GetByEmail(
		ctx,
		emailMessage.Tenant,
		emailMessage.FromProvider,
		emailMessage.From,
	)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, fmt.Errorf("unable to retrieve oauth token for %s: %v", emailMessage.From, err)
	}
	return oauthToken, nil
}

func (s *mailService) prepareEmailMessage(ctx context.Context, span opentracing.Span, emailMessage *postgresentity.EmailMessage) error {
	uniqueInternalIdentifier := utils.GenerateRandomString(64)
	emailMessage.UniqueInternalIdentifier = &uniqueInternalIdentifier

	if emailMessage.ReplyTo == nil {
		emailMessage.Subject = emailMessage.Subject
		return nil
	}

	return s.handleReplyToEmail(ctx, span, emailMessage)
}

func (s *mailService) handleReplyToEmail(ctx context.Context, span opentracing.Span, emailMessage *postgresentity.EmailMessage) error {
	interactionEventNode, err := s.services.Neo4jRepositories.CommonReadRepository.GetById(
		ctx,
		emailMessage.Tenant,
		*emailMessage.ReplyTo,
		commonModel.NodeLabelInteractionEvent,
	)
	if err != nil {
		err = fmt.Errorf("unable to handle reply to email: %v", err)
		tracing.TraceErr(span, err)
		return err
	}

	interactionEvent := neo4jmapper.MapDbNodeToInteractionEventEntity(interactionEventNode)
	emailChannelData := neo4jentity.EmailChannelData{}
	if err := json.Unmarshal([]byte(interactionEvent.ChannelData), &emailChannelData); err != nil {
		err = fmt.Errorf("unable to parse email channel data for %s", *emailMessage.ReplyTo)
		tracing.TraceErr(span, err)
		return err
	}

	s.setReplySubject(emailMessage, emailChannelData)
	s.setReplyReferences(emailMessage, emailChannelData)
	return nil
}

func (s *mailService) setReplySubject(emailMessage *postgresentity.EmailMessage, emailChannelData neo4jentity.EmailChannelData) {
	subject := emailChannelData.Subject
	if len(subject) < 3 || subject[:3] != "Re:" {
		subject = "Re: " + subject
	}
	emailMessage.Subject = subject
}

func (s *mailService) setReplyReferences(emailMessage *postgresentity.EmailMessage, emailChannelData neo4jentity.EmailChannelData) {
	if emailChannelData.Reference != "" {
		emailMessage.ProviderReferences = emailChannelData.Reference + " " + emailChannelData.ProviderMessageId
	} else {
		emailMessage.ProviderReferences = emailChannelData.ProviderMessageId
	}
	emailMessage.ProviderInReplyTo = emailChannelData.ProviderMessageId
}

func (s *mailService) sendEmailBasedOnProvider(
	ctx context.Context,
	span opentracing.Span,
	emailMessage *postgresentity.EmailMessage,
	oauthToken *entity.OAuthTokenEntity,
) error {
	if oauthToken == nil {
		return s.sendEmailViaOpenSrs(ctx, span, emailMessage)
	}
	return s.sendEmailViaOAuth(ctx, span, emailMessage, oauthToken)
}

func (s *mailService) sendEmailViaOpenSrs(ctx context.Context, span opentracing.Span, emailMessage *postgresentity.EmailMessage) error {
	mailbox, err := s.services.PostgresRepositories.TenantSettingsMailboxRepository.GetByMailbox(ctx, emailMessage.From)
	if err != nil {
		err = fmt.Errorf("failed to get mailbox: %v", err)
		tracing.TraceErr(span, err)
		return err
	}
	if mailbox == nil {
		err = fmt.Errorf("mailbox not found for %s", emailMessage.From)
		tracing.TraceErr(span, err)
		return err
	}
	return s.services.OpenSrsService.SendEmail(ctx, emailMessage)
}

func (s *mailService) sendEmailViaOAuth(
	ctx context.Context,
	span opentracing.Span,
	emailMessage *entity.EmailMessage,
	oauthToken *postgresentity.OAuthTokenEntity,
) error {
	if oauthToken.NeedsManualRefresh {
		err := errors.New("oauth token needs manual refresh")
		tracing.TraceErr(span, err)
		return fmt.Errorf("oauth token needs manual refresh: %v", err)
	}

	switch oauthToken.Provider {
	case "google":
		return s.services.GoogleService.SendEmail(ctx, emailMessage)
	case "azure-ad":
		return s.services.AzureService.SendEmail(ctx, emailMessage)
	default:
		return fmt.Errorf("provider %s not supported", oauthToken.Provider)
	}
}

func (s *mailService) storeEmailMessage(ctx context.Context, span opentracing.Span, emailMessage *postgresentity.EmailMessage) error {
	emailMessage.Status = entity.EmailMessageStatusSent
	emailMessage.SentAt = utils.TimePtr(utils.Now())

	if err := s.services.PostgresRepositories.EmailMessageRepository.Store(ctx, emailMessage.Tenant, emailMessage); err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("failed to store email message: %v", err)
	}
	return nil
}

func (s *mailService) saveEmailInTx(ctx context.Context, tx neo4j.ManagedTransaction, emailMessage *postgresentity.EmailMessage, span opentracing.Span) (any, error) {
	tenant := common.GetTenantFromContext(ctx)

	sessionID, err := s.getOrCreateInteractionSession(ctx, tx, span, tenant, emailMessage)
	if err != nil {
		err = fmt.Errorf("failed to get or create interaction session: %v", err)
		tracing.TraceErr(span, err)
		return nil, err
	}

	participants := s.buildParticipantLists(emailMessage)

	eventID, err := s.createInteractionEvent(ctx, tx, span, emailMessage, sessionID, participants)
	if err != nil {
		err = fmt.Errorf("failed to create interaction event: %v", err)
		tracing.TraceErr(span, err)
		return nil, err
	}

	if err := s.linkEventToSession(ctx, tx, span, tenant, eventID, sessionID); err != nil {
		err = fmt.Errorf("failed to link event to session: %v", err)
		tracing.TraceErr(span, err)
		return nil, err
	}

	return eventID, nil
}

func (s *mailService) getOrCreateInteractionSession(
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	span opentracing.Span,
	tenant string,
	emailMessage *postgresentity.EmailMessage,
) (string, error) {
	// Try to get existing session
	sessionNode, err := s.services.Neo4jRepositories.InteractionSessionReadRepository.GetByIdentifierAndChannel(
		ctx, tenant, emailMessage.ProviderThreadId, "EMAIL",
	)
	if err != nil {
		err = fmt.Errorf("failed to get interaction session: %v", err)
		tracing.TraceErr(span, err)
		return "", err
	}

	if sessionNode != nil {
		return neo4jmapper.MapDbNodeToInteractionSessionEntity(sessionNode).Id, nil
	}

	// Create new session if none exists
	sessionID, err := s.services.InteractionSessionService.CreateInTx(ctx, tx, &neo4jentity.InteractionSessionEntity{
		Status:     "ACTIVE",
		Type:       "THREAD",
		Channel:    "EMAIL",
		Identifier: emailMessage.ProviderThreadId,
		Name:       emailMessage.Subject,
	})
	if err != nil {
		err = fmt.Errorf("failed to create transaction: %v", err)
		tracing.TraceErr(span, err)
		return "", err
	}

	if sessionID == nil {
		err := errors.New("session id is empty")
		tracing.TraceErr(span, err)
		return "", err
	}

	return *sessionID, nil
}

func (s *mailService) buildParticipantLists(emailMessage *postgresentity.EmailMessage) *struct {
	sentBy  []InteractionEventParticipantData
	sentTo  []InteractionEventParticipantData
	sentCc  []InteractionEventParticipantData
	sentBcc []InteractionEventParticipantData
} {
	result := &struct {
		sentBy  []InteractionEventParticipantData
		sentTo  []InteractionEventParticipantData
		sentCc  []InteractionEventParticipantData
		sentBcc []InteractionEventParticipantData
	}{
		sentBy:  make([]InteractionEventParticipantData, 0),
		sentTo:  make([]InteractionEventParticipantData, 0),
		sentCc:  make([]InteractionEventParticipantData, 0),
		sentBcc: make([]InteractionEventParticipantData, 0),
	}

	result.sentBy = append(result.sentBy, InteractionEventParticipantData{
		Email: &emailMessage.From,
	})

	for _, to := range emailMessage.To {
		result.sentTo = append(result.sentTo, InteractionEventParticipantData{
			Email: &to,
		})
	}
	for _, cc := range emailMessage.Cc {
		result.sentCc = append(result.sentCc, InteractionEventParticipantData{
			Email: &cc,
		})
	}
	for _, bcc := range emailMessage.Bcc {
		result.sentBcc = append(result.sentBcc, InteractionEventParticipantData{
			Email: &bcc,
		})
	}

	return result
}

func (s *mailService) createInteractionEvent(
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	span opentracing.Span,
	emailMessage *postgresentity.EmailMessage,
	sessionID string,
	participants *struct {
		sentBy  []InteractionEventParticipantData
		sentTo  []InteractionEventParticipantData
		sentCc  []InteractionEventParticipantData
		sentBcc []InteractionEventParticipantData
	},
) (*string, error) {
	emailChannelData, err := neo4jentity.BuildEmailChannelData(
		emailMessage.ProviderMessageId,
		emailMessage.ProviderThreadId,
		emailMessage.Subject,
		emailMessage.ProviderInReplyTo,
		emailMessage.ProviderReferences,
	)
	if err != nil {
		err = fmt.Errorf("failed to build email channel data: %v", err)
		tracing.TraceErr(span, err)
		return nil, err
	}

	eventID, err := s.services.InteractionEventService.CreateInTx(ctx, tx, &InteractionEventCreateData{
		InteractionEventEntity: &neo4jentity.InteractionEventEntity{
			Content:                      emailMessage.Content,
			ContentType:                  "text/html",
			Channel:                      "EMAIL",
			ChannelData:                  *emailChannelData,
			Identifier:                   emailMessage.ProviderMessageId,
			CustomerOSInternalIdentifier: *emailMessage.UniqueInternalIdentifier,
			Hide:                         false,
			Source:                       "openline", // TODO
			SourceOfTruth:                "openline", // TODO
			AppSource:                    "TODO",     // TODO
		},
		SentBy:            participants.sentBy,
		SentTo:            participants.sentTo,
		SentCc:            participants.sentCc,
		SentBcc:           participants.sentBcc,
		RepliesTo:         emailMessage.ReplyTo,
		SessionIdentifier: &sessionID,
	})
	if err != nil {
		err = fmt.Errorf("failed to create interaction event: %v", err)
		tracing.TraceErr(span, err)
		return nil, err
	}

	return eventID, nil
}

func (s *mailService) linkEventToSession(
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	span opentracing.Span,
	tenant string,
	eventID *string,
	sessionID string,
) error {
	err := s.services.Neo4jRepositories.CommonWriteRepository.Link(ctx, &tx, tenant, repository.LinkDetails{
		FromEntityId:           *eventID,
		FromEntityType:         commonModel.INTERACTION_EVENT,
		Relationship:           commonModel.PART_OF,
		RelationshipProperties: nil,
		ToEntityId:             sessionID,
		ToEntityType:           commonModel.INTERACTION_SESSION,
	})
	if err != nil {
		err = fmt.Errorf("failed to link interaction event with interaction session: %v", err)
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}
