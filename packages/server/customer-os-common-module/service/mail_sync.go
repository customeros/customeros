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
	"github.com/opentracing/opentracing-go/log"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

const AppSource = "sync-email"

func (s *mailService) SyncEmail(tenant string, emailId uuid.UUID) (postgresentity.RawState, *string, error) {
	ctx := context.Background()
	span, ctx := s.initializeTracing(ctx, "MailService.SyncEmail")
	defer span.Finish()
	span.LogFields(
		log.String("emailId", emailId.String()),
		log.String("tenant", tenant))

	var reason string

	rawEmail, err := s.services.PostgresRepositories.RawEmailRepository.GetEmailForSync(emailId)
	if err != nil {
		err = fmt.Errorf("failed to get raw email for sync: %w", err)
		tracing.TraceErr(span, err)
		return postgresentity.ERROR, nil, err
	}

	email, err := s.LoadEmail(rawEmail)
	if err != nil {
		err = fmt.Errorf("failed to load email for sync: %w", err)
		tracing.TraceErr(span, err)
		return postgresentity.ERROR, nil, err
	}

	if email.Identifiers.MessageId == "" {
		return postgresentity.ERROR, nil, fmt.Errorf("email message ID is empty")
	}

	check := s.ProcessEmailCheck(&email)
	if !check.ProcessEmail {
		if check.IsBounce {
			reason = "email bounced"
		}
		if check.IsAutoResponder {
			reason = "email autoresponder"
		}
		if check.IsBulkMail {
			reason = "bulk email"
		}
		return postgresentity.SKIPPED, &reason, nil
	}

	if len(email.Participants.AllEmails) == 0 {
		reason := "no email address belongs to a workspace domain"
		return postgresentity.SKIPPED, &reason, nil
	}

	if s.warmingEmailCheck(tenant, email) {
		reason := "warming email"
		return postgresentity.SKIPPED, &reason, nil
	}

	now := utils.Now()

	sentAt, err := convertToUTC(email.Content.SentDate)
	email.CreatedAt = sentAt
	if err != nil {
		err = fmt.Errorf("%v :%v", err, emailId.String())
		tracing.TraceErr(span, err)
		return postgresentity.ERROR, nil, err
	}

	interactionEventId, err := s.services.Neo4jRepositories.InteractionEventRepository.GetInteractionEventIdByExternalId(ctx, tenant, rawEmail.ExternalSystem, rawEmail.MessageId)
	if err != nil {
		err = fmt.Errorf("failed to check if interaction event exists for external id %v for tenant %v :%v", rawEmail.MessageId, tenant, err)
		tracing.TraceErr(span, err)
		return postgresentity.ERROR, nil, err
	}

	if interactionEventId != "" {
		err = fmt.Errorf("interaction event already exists for raw email id %v", emailId.String())
		tracing.TraceErr(span, err)
		reason := "interaction event already exists"
		return postgresentity.SKIPPED, &reason, nil
	}

	chanErr := s.buildChannelData(&email)
	if chanErr != nil {
		err = fmt.Errorf("failed to build email channel data for email with id %v: %v", emailId.String(), chanErr)
		tracing.TraceErr(span, err)
		return postgresentity.ERROR, nil, chanErr
	}

	return s.processInboundEmail(ctx, tenant, &email, rawEmail, now)
}

func (s *mailService) SyncEmailsForUser(tenant string, userSource string) {
	ctx := context.Background()
	span, ctx := s.initializeTracing(ctx, "MailService.SyncEmailsForUser")
	defer span.Finish()
	span.LogFields(
		log.String("tenant", tenant),
		log.String("userSource", userSource))

	emailsIdsForSync, err := s.services.PostgresRepositories.RawEmailRepository.GetEmailsIdsForUserForSync(tenant, userSource)
	if err != nil {
		err = fmt.Errorf("failed to get emails for sync: %v", err)
		tracing.TraceErr(span, err)
	}

	if len(emailsIdsForSync) == 0 {
		return
	}

	distinctExternalSystems := make([]string, 0)
	for _, email := range emailsIdsForSync {
		if !utils.Contains(distinctExternalSystems, email.ExternalSystem) {
			distinctExternalSystems = append(distinctExternalSystems, email.ExternalSystem)
		}
	}

	externalSystemStr := ""
	if len(distinctExternalSystems) > 0 {
		externalSystemStr = distinctExternalSystems[0]
	}

	_, err = s.createUserSourceAsEmailNode(tenant, userSource, externalSystemStr)
	if err != nil {
		err = fmt.Errorf("failed to create user source as email node: %v", err)
		tracing.TraceErr(span, err)
		return
	}

	for _, externalSystem := range distinctExternalSystems {
		err = p.services.Neo4jRepositories.ExternalSystemWriteRepository.CreateIfNotExists(context.Background(), tenant, externalSystem, externalSystem)
		if err != nil {
			return
		}
	}

	s.syncEmails(tenant, emailsIdsForSync)
}

func (s *mailService) SyncEmailByMessageId(tenant, usernameSource, messageId string) (postgresentity.RawState, *string, error) {
	ctx := context.Background()
	span, ctx := s.initializeTracing(ctx, "MailService.SyncEmailByMessageId")
	defer span.Finish()
	span.LogFields(
		log.String("tenant", tenant),
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

	return s.SyncEmail(tenant, rawEmail.ID)
}

func (s *mailService) SyncEmailByEmailRawId(tenant string, emailId uuid.UUID) (postgresentity.RawState, *string, error) {
	return s.SyncEmail(tenant, emailId)
}

func (s *mailService) syncEmails(tenant string, emails []postgresentity.RawEmail) {
	ctx := context.Background()
	span, ctx := s.initializeTracing(ctx, "MailService.syncEmails")
	defer span.Finish()
	span.LogFields(
		log.String("tenant", tenant))

	for _, email := range emails {
		// TODO here is control to call new service !!!
		// state, reason, err := s.syncEmail(tenant, email.ID)
		state, reason, err := s.SyncEmail(tenant, email.ID)

		var errMessage *string
		if err != nil {
			s2 := err.Error()
			errMessage = &s2
		}

		err = s.services.PostgresRepositories.RawEmailRepository.MarkSentToEventStore(email.ID, state, reason, errMessage)
		if err != nil {
			fmt.Errorf("unable to mark email as sent to event store: %v", err)
			tracing.TraceErr(span, err)
		}

		fmt.Println("raw email processed: " + email.ID.String())
	}
}

func (s *mailService) processInboundEmail(ctx context.Context, tenant string, email *EmailMessageData, rawEmail *postgresentity.RawEmail, ts time.Time) (postgresentity.RawState, *string, error) {
	session := utils.NewNeo4jWriteSession(ctx, *s.services.Neo4jRepositories.Neo4jDriver)
	defer session.Close(ctx)

	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		err = fmt.Errorf("failed to start transaction: %v", err)
		return postgresentity.ERROR, nil, err
	}
	defer tx.Close(ctx)

	// Process session and events
	if err := s.processSessionAndEvents(ctx, tx, tenant, email, rawEmail, ts); err != nil {
		return postgresentity.ERROR, nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return postgresentity.ERROR, nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return postgresentity.PROCESSED, nil, nil
}

func (s *mailService) processSessionAndEvents(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, email *EmailMessageData, rawEmail *postgresentity.RawEmail, ts time.Time) error {
	// get EmailForCustomerOS
	cosEmail := s.buildEmailForCustomerOS(email, rawEmail.ExternalSystem)

	// Create session
	sessionId, err := s.services.Neo4jRepositories.InteractionEventRepository.MergeInteractionSession(ctx, tx, tenant, email.Identifiers.EmailThreadId, ts, cosEmail, rawEmail.ExternalSystem, AppSource)
	if err != nil {
		return fmt.Errorf("failed merge interaction session: %v", err)
	}

	// Create event
	eventId, err := s.services.Neo4jRepositories.InteractionEventRepository.MergeEmailInteractionEvent(ctx, tx, tenant, ts, cosEmail, rawEmail.ExternalSystem, AppSource)
	if err != nil {
		return fmt.Errorf("failed merge interaction event: %v", err)
	}

	// Link event to session
	if err := s.services.Neo4jRepositories.InteractionEventRepository.LinkInteractionEventToSession(ctx, tx, tenant, eventId, sessionId); err != nil {
		return fmt.Errorf("failed to link event to session: %v", err)
	}

	// Process participants
	if err := s.linkParticipants(ctx, tx, tenant, eventId, &email.Participants, ts, rawEmail.ExternalSystem); err != nil {
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
		ExternalId:     email.Identifiers.ExternalId,
		EmailThreadId:  email.Identifiers.EmailThreadId,
		Channel:        "EMAIL",
		ChannelData:    email.ChannelData,
	}
	return save
}

func (s *mailService) linkParticipants(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, eventId string, participants *EmailParticipants, now time.Time, externalSystem string) error {
	emailIds := make(map[string]string)

	// Link From participant
	fromId, err := s.getOrCreateEmailId(ctx, tx, tenant, participants.From.Email, now, externalSystem, emailIds)
	if err != nil {
		return err
	}
	if err := s.services.Neo4jRepositories.InteractionEventRepository.InteractionEventSentByEmail(ctx, tx, tenant, eventId, fromId); err != nil {
		return err
	}

	// Link To participants
	if err := s.linkEmailGroup(ctx, tx, tenant, eventId, "TO", participants.GetToEmailAddresses(), now, externalSystem, emailIds); err != nil {
		return err
	}

	// Link CC participants
	if err := s.linkEmailGroup(ctx, tx, tenant, eventId, "CC", participants.GetCcEmailAddresses(), now, externalSystem, emailIds); err != nil {
		return err
	}

	// Link BCC participants
	if err := s.linkEmailGroup(ctx, tx, tenant, eventId, "BCC", participants.GetBccEmailAddresses(), now, externalSystem, emailIds); err != nil {
		return err
	}

	return nil
}

func (s *mailService) linkEmailGroup(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, eventId string, groupType string, emails []string, now time.Time, externalSystem string, emailIds map[string]string) error {
	var groupEmailIds []string

	for _, email := range emails {
		if email == "" {
			continue
		}

		emailId, err := s.getOrCreateEmailId(ctx, tx, tenant, email, now, externalSystem, emailIds)
		if err != nil {
			return err
		}

		if !utils.Contains(groupEmailIds, emailId) {
			groupEmailIds = append(groupEmailIds, emailId)
		}
	}

	if len(groupEmailIds) > 0 {
		return s.services.Neo4jRepositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tx, tenant, eventId, groupType, groupEmailIds)
	}

	return nil
}

func (s *mailService) getOrCreateEmailId(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, email string, now time.Time, externalSystem string, emailIds map[string]string) (string, error) {
	if id, exists := emailIds[email]; exists {
		return id, nil
	}

	id, err := s.services.SyncService.GetEmailIdForEmail(ctx, tx, tenant, email, now, externalSystem)
	if err != nil {
		return "", fmt.Errorf("failed to get email ID for %s: %v", email, err)
	}
	if id == "" {
		return "", fmt.Errorf("no email ID found for %s", email)
	}

	emailIds[email] = id
	return id, nil
}

func (s *mailService) buildChannelData(email *EmailMessageData) error {
	channelData, err := neo4jentity.BuildEmailChannelData(email.Identifiers.ProviderMessageId, email.Identifiers.EmailThreadId, email.Content.Subject, strings.Join(email.Participants.GetInReplyToEmailAddresses(), " "), strings.Join(email.Identifiers.References, " "))
	if err != nil {
		return err
	}
	email.Channel = "EMAIL"
	email.ChannelData = channelData

	return nil
}

func (s *mailService) warmingEmailCheck(tenant string, email EmailMessageData) bool {
	emailExclusion := s.services.Cache.GetEmailExclusion(tenant)

	for _, exclusion := range emailExclusion {
		if exclusion.ExcludeSubject != nil {
			if strings.Contains(email.Content.Subject, *exclusion.ExcludeSubject) {
				return true
			}
		}
		if exclusion.ExcludeBody != nil {
			if strings.Contains(email.Content.Html, *exclusion.ExcludeBody) {
				return true
			}
			if strings.Contains(email.Content.Text, *exclusion.ExcludeBody) {
				return true
			}
		}
	}
	return false
}

func (s *mailService) createUserSourceAsEmailNode(tenant, userSource, externalSystem string) (string, error) {
	ctx := context.Background()

	emailId, err := s.services.EmailService.Merge(ctx, tenant, EmailFields{
		Email:     userSource,
		AppSource: AppSource,
		Source:    neo4jentity.DecodeDataSource(externalSystem),
	}, nil)
	if err != nil {
		return "", fmt.Errorf("unable to create email: %v", err)
	}

	return *emailId, nil
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
