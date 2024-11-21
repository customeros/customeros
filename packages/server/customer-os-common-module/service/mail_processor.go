package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	neo4jentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

const AppSource = constants.AppSourceSyncEmail

func (s *mailService) GetEmailsForProcessingForUser(ctx context.Context, tenant, userEmailAddress string) {
	span, ctx := s.initializeTracing(ctx, "MailService.SyncEmailsForUser")
	defer span.Finish()
	span.LogKV("userEmailAddress", userEmailAddress)

	rawEmailsIdsForProcess, err := s.services.PostgresRepositories.RawEmailRepository.GetEmailsIdsForUserForSync(tenant, userEmailAddress)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get emails for sync"))
		return
	}

	if len(rawEmailsIdsForProcess) == 0 {
		err = fmt.Errorf("no emails found for sync")
		tracing.TraceErr(span, err)
		return
	}

	distinctExternalSystems := make([]string, 0)
	for _, email := range rawEmailsIdsForProcess {
		if !utils.Contains(distinctExternalSystems, email.ExternalSystem) {
			distinctExternalSystems = append(distinctExternalSystems, email.ExternalSystem)
		}
	}

	externalSystemStr := ""
	if len(distinctExternalSystems) > 0 {
		externalSystemStr = distinctExternalSystems[0]
	}

	// Create email node in neo4j
	err = s.createUserEmailAddressAsNode(ctx, tenant, userEmailAddress, externalSystemStr, span)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create user source as email node"))
		return
	}

	for _, externalSystem := range distinctExternalSystems {
		// TODO alexb add caching for each tenant of external systems
		err = s.services.Neo4jRepositories.ExternalSystemWriteRepository.CreateIfNotExists(
			ctx, tenant, externalSystem, externalSystem,
		)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to merge external system"))
			return
		}
	}

	s.processRawEmails(ctx, tenant, rawEmailsIdsForProcess, span)
}

func (s *mailService) ProcessEmail(ctx context.Context, tenant string, rawEmailId uuid.UUID) (postgresentity.RawState, *string, error) {
	span, ctx := s.initializeTracing(ctx, "MailService.ProcessEmail")
	defer span.Finish()
	span.LogFields(log.String("rawEmailId", rawEmailId.String()))

	rawEmail, err := s.services.PostgresRepositories.RawEmailRepository.GetEmailForProcess(rawEmailId)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get email for process"))
		return postgresentity.ERROR, nil, err
	}

	emailMessageData, err := s.LoadEmail(ctx, rawEmail)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to load email"))
		return postgresentity.ERROR, nil, err
	}

	if emailMessageData.Identifiers.MessageId == "" {
		return postgresentity.ERROR, nil, fmt.Errorf("email message ID is empty")
	}

	if len(emailMessageData.Participants.AllEmails) == 0 {
		reason := "no email address belongs to a workspace domain"
		return postgresentity.SKIPPED, &reason, nil
	}

	check := s.ProcessEmailCheck(ctx, tenant, &emailMessageData)
	if !check.ProcessEmail {
		return postgresentity.SKIPPED, &check.SkipReason, nil
	}

	sentAt, err := convertToUTC(emailMessageData.Content.SentDate)
	emailMessageData.CreatedAt = sentAt
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to convert email sent date to UTC"))
		return postgresentity.ERROR, nil, err
	}

	interactionEventId, err := s.services.Neo4jRepositories.InteractionEventRepository.GetInteractionEventIdByExternalId(
		ctx, tenant, rawEmail.ExternalSystem, rawEmail.MessageId,
	)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to check if interaction event exists"))
		return postgresentity.ERROR, nil, err
	}

	if interactionEventId != "" {
		tracing.TraceErr(span, errors.Wrap(err, "interaction event already exists"))
		reason := "interaction event already exists"
		return postgresentity.SKIPPED, &reason, nil
	}

	chanErr := s.buildChannelData(&emailMessageData, span)
	if chanErr != nil {
		tracing.TraceErr(span, errors.Wrap(chanErr, "failed to build channel data"))
		return postgresentity.ERROR, nil, chanErr
	}

	return s.processInboundEmail(ctx, tenant, &emailMessageData, rawEmail, utils.Now(), span)
}

func (s *mailService) ProcessEmailByMessageId(ctx context.Context, tenant, usernameSource, messageId string) (postgresentity.RawState, *string, error) {
	span, ctx := s.initializeTracing(ctx, "MailService.ProcessEmailByMessageId")
	defer span.Finish()
	span.LogFields(
		log.String("userSource", usernameSource),
		log.String("messageId", messageId))

	rawEmail, err := s.services.PostgresRepositories.RawEmailRepository.GetEmailForSyncByMessageId(tenant, usernameSource, messageId)
	if err != nil {
		err = fmt.Errorf("failed to get emails for sync: %v", err)
		tracing.TraceErr(span, err)
		return postgresentity.ERROR, nil, err
	}

	if rawEmail == nil {
		return postgresentity.ERROR, nil, fmt.Errorf("email with message id %v not found", messageId)
	}

	return s.ProcessEmail(ctx, tenant, rawEmail.ID)
}

func (s *mailService) ProcessEmailByEmailRawId(ctx context.Context, tenant string, emailId uuid.UUID) (postgresentity.RawState, *string, error) {
	return s.ProcessEmail(ctx, tenant, emailId)
}

func (s *mailService) processRawEmails(ctx context.Context, tenant string, rawEmails []postgresentity.RawEmail, span opentracing.Span) {
	for _, rawEmail := range rawEmails {
		state, reason, err := s.ProcessEmail(ctx, tenant, rawEmail.ID)
		var errMessage *string
		if err != nil {
			s2 := err.Error()
			errMessage = &s2
		}

		err = s.services.PostgresRepositories.RawEmailRepository.MarkProcessed(rawEmail.ID, state, reason, errMessage)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to mark raw email as processed in postgres"))
		}
	}
}

func (s *mailService) processInboundEmail(
	ctx context.Context,
	tenant string,
	email *EmailMessageData,
	rawEmail *postgresentity.RawEmail,
	ts time.Time,
	span opentracing.Span,
) (postgresentity.RawState, *string, error) {
	session := utils.NewNeo4jWriteSession(ctx, *s.services.Neo4jRepositories.Neo4jDriver)
	defer session.Close(ctx)

	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		err = fmt.Errorf("failed to start transaction: %v", err)
		tracing.TraceErr(span, err)
		return postgresentity.ERROR, nil, err
	}
	defer tx.Close(ctx)

	// Process session and events
	if err := s.processSessionAndEvents(ctx, tx, tenant, email, rawEmail, ts, span); err != nil {
		err = fmt.Errorf("failed to process session and events: %v", err)
		tracing.TraceErr(span, err)
		return postgresentity.ERROR, nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		err = fmt.Errorf("failed to commit transaction: %v", err)
		tracing.TraceErr(span, err)
		return postgresentity.ERROR, nil, err
	}

	return postgresentity.PROCESSED, nil, nil
}

func (s *mailService) processSessionAndEvents(
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	tenant string,
	emailMessageData *EmailMessageData,
	rawEmail *postgresentity.RawEmail,
	ts time.Time,
	span opentracing.Span,
) error {
	// get EmailForCustomerOS
	cosEmail := s.buildEmailForCustomerOS(emailMessageData, rawEmail.ExternalSystem)

	// Create session
	sessionId, err := s.services.Neo4jRepositories.InteractionEventRepository.MergeInteractionSession(
		ctx, tx, tenant, emailMessageData.Identifiers.EmailThreadId, ts, cosEmail, rawEmail.ExternalSystem, AppSource,
	)
	if err != nil {
		err = fmt.Errorf("failed merge interaction session: %v", err)
		tracing.TraceErr(span, err)
		return err
	}

	// Create event
	eventId, err := s.services.Neo4jRepositories.InteractionEventRepository.MergeEmailInteractionEvent(
		ctx, tx, tenant, ts, cosEmail, rawEmail.ExternalSystem, AppSource,
	)
	if err != nil {
		err = fmt.Errorf("failed merge interaction event: %v", err)
		tracing.TraceErr(span, err)
		return err
	}

	// Link event to session
	if err := s.services.Neo4jRepositories.InteractionEventRepository.LinkInteractionEventToSession(
		ctx, tx, tenant, eventId, sessionId,
	); err != nil {
		err = fmt.Errorf("failed to link event to session: %v", err)
		tracing.TraceErr(span, err)
		return err
	}

	// Process participants
	if err := s.linkParticipants(ctx, tx, tenant, eventId, &emailMessageData.Participants, ts, rawEmail.ExternalSystem, span); err != nil {
		err = fmt.Errorf("failed to link participants: %v", err)
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *mailService) buildEmailForCustomerOS(email *EmailMessageData, externalSystem string) model.SaveEmailMessage {
	save := model.SaveEmailMessage{
		Html:           email.Content.Html,
		Text:           email.Content.Text,
		Subject:        email.Content.Subject,
		CreatedAt:      email.CreatedAt,
		ExternalSystem: externalSystem,
		ExternalId:     email.Identifiers.MessageId,
		EmailThreadId:  email.Identifiers.EmailThreadId,
		Channel:        "EMAIL",
		ChannelData:    email.ChannelData,
	}
	return save
}

func (s *mailService) linkParticipants(
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	tenant string,
	eventId string,
	participants *EmailParticipants,
	now time.Time,
	externalSystem string,
	span opentracing.Span,
) error {
	emailIds := make(map[string]string)

	// Link From participant
	fromId, err := s.getOrCreateEmailId(ctx, tx, tenant, participants.From.Email, now, externalSystem, emailIds, span)
	if err != nil {
		err = fmt.Errorf("failed to get or create email id: %v", err)
		tracing.TraceErr(span, err)
		return err
	}
	if err := s.services.Neo4jRepositories.InteractionEventRepository.InteractionEventSentByEmail(ctx, tx, tenant, eventId, fromId); err != nil {
		err = fmt.Errorf("failed to create interaction event sent by email: %v", err)
		tracing.TraceErr(span, err)
		return err
	}

	// Link To participants
	if err := s.linkEmailGroup(ctx, tx, tenant, eventId, "TO", participants.GetToEmailAddresses(), now, externalSystem, emailIds, span); err != nil {
		err = fmt.Errorf("failed to link email group for TO: %v", err)
		tracing.TraceErr(span, err)
		return err
	}

	// Link CC participants
	if err := s.linkEmailGroup(ctx, tx, tenant, eventId, "CC", participants.GetCcEmailAddresses(), now, externalSystem, emailIds, span); err != nil {
		err = fmt.Errorf("failed to link email group for CC: %v", err)
		tracing.TraceErr(span, err)
		return err
	}

	// Link BCC participants
	if err := s.linkEmailGroup(ctx, tx, tenant, eventId, "BCC", participants.GetBccEmailAddresses(), now, externalSystem, emailIds, span); err != nil {
		err = fmt.Errorf("failed to link email group for BCC: %v", err)
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *mailService) linkEmailGroup(
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	tenant string,
	eventId string,
	groupType string,
	emails []string,
	now time.Time,
	externalSystem string,
	emailIds map[string]string,
	span opentracing.Span,
) error {
	var groupEmailIds []string

	for _, email := range emails {
		if email == "" {
			continue
		}

		emailId, err := s.getOrCreateEmailId(ctx, tx, tenant, email, now, externalSystem, emailIds, span)
		if err != nil {
			err = fmt.Errorf("failed to get or create email ID: %v", err)
			tracing.TraceErr(span, err)
			return err
		}

		if !utils.Contains(groupEmailIds, emailId) {
			groupEmailIds = append(groupEmailIds, emailId)
		}
	}

	if len(groupEmailIds) > 0 {
		return s.services.Neo4jRepositories.InteractionEventRepository.InteractionEventSentToEmails(
			ctx, tx, tenant, eventId, groupType, groupEmailIds,
		)
	}

	return nil
}

func (s *mailService) getOrCreateEmailId(
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	tenant string,
	email string,
	now time.Time,
	externalSystem string,
	emailIds map[string]string,
	span opentracing.Span,
) (string, error) {
	if id, exists := emailIds[email]; exists {
		return id, nil
	}

	id, err := s.services.SyncService.GetEmailIdForEmail(ctx, tx, tenant, email, now, externalSystem)
	if err != nil {
		err = fmt.Errorf("failed to get email ID for %s: %v", email, err)
		tracing.TraceErr(span, err)
		return "", err
	}
	if id == "" {
		err = fmt.Errorf("no email ID found for %s", email)
		tracing.TraceErr(span, err)
		return "", err
	}

	emailIds[email] = id
	return id, nil
}

func (s *mailService) buildChannelData(email *EmailMessageData, span opentracing.Span) error {
	channelData, err := neo4jentity.BuildEmailChannelData(
		email.Identifiers.ProviderMessageId,
		email.Identifiers.EmailThreadId,
		email.Content.Subject,
		strings.Join(email.Participants.GetReplyToEmailAddresses(), " "),
		strings.Join(email.Identifiers.References, " "),
	)
	if err != nil {
		err = fmt.Errorf("building channel data failed: %v", err)
		tracing.TraceErr(span, err)
		return err
	}
	email.Channel = "EMAIL"
	email.ChannelData = channelData

	return nil
}

func (s *mailService) createUserEmailAddressAsNode(ctx context.Context, tenant, userEmailAddress, externalSystem string, span opentracing.Span) error {
	_, err := s.services.EmailService.Merge(ctx, tenant, EmailFields{
		Email:     userEmailAddress,
		AppSource: AppSource,
		Source:    neo4jentity.DecodeDataSource(externalSystem),
	}, nil)
	if err != nil {
		err = fmt.Errorf("unable to create email: %v", err)
		tracing.TraceErr(span, err)
		return err
	}

	return nil
}

func (s *mailService) mapDbNodeToEmailEntity(node dbtype.Node) *neo4jentity.EmailEntity {
	props := utils.GetPropsFromNode(node)
	result := neo4jentity.EmailEntity{
		Id:       utils.GetStringPropOrEmpty(props, "id"),
		Email:    utils.GetStringPropOrEmpty(props, "email"),
		RawEmail: utils.GetStringPropOrEmpty(props, "rawEmail"),
	}
	return &result
}
