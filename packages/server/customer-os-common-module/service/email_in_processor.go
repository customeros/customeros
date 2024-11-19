package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
)

const AppSource = "sync-email"

func (p *emailService) syncEmail(tenant string, emailId uuid.UUID) (entity.RawState, *string, error) {
	ctx := context.Background()
	span, ctx := tracing.StartTracerSpan(ctx, "EmailService.syncEmail")
	defer span.Finish()
	span.LogFields(log.String("emailId", emailId.String()))

	emailIdString := emailId.String()

	rawEmail, err := p.services.PostgresRepositories.RawEmailRepository.GetEmailForSync(emailId)
	if err != nil {
		logrus.Errorf("failed to get raw email for sync: %v", err)
		return entity.ERROR, nil, err
	}

	email, err := p.LoadEmail(rawEmail)
	if err != nil {
		logrus.Errorf("failed to load email for sync: %v", err)
		return entity.ERROR, nil, err
	}

	if email.Identifiers.MessageId == "" {
		return entity.ERROR, nil, fmt.Errorf("email message ID is empty")
	}

	interactionEventId, err := p.repositories.InteractionEventRepository.GetInteractionEventIdByExternalId(ctx, tenant, rawEmail.ExternalSystem, rawEmail.MessageId)
	if err != nil {
		logrus.Errorf("failed to check if interaction event exists for external id %v for tenant %v :%v", rawEmail.MessageId, tenant, err)
		return entity.ERROR, nil, err
	}

	now := time.Now().UTC()

	sentAt, err := convertToUTC(email.Content.SentDate)
	email.CreatedAt, err = p.validateDates(sentAt)
	if err != nil {
		logrus.Errorf("%v :%v", err, emailId.String())
		return entity.ERROR, nil, err
	}

	if p.warmingEmailCheck(email) {
		reason := "warming email"
		return entity.SKIPPED, &reason, nil
	}

	if interactionEventId != "" {
		logrus.Infof("interaction event already exists for raw email id %v", emailIdString)
		reason := "interaction event already exists"
		return entity.SKIPPED, &reason, nil
	}

	if len(email.Participants.AllEmails) == 0 {
		reason := "no email address belongs to a workspace domain"
		return entity.SKIPPED, &reason, nil
	}

	chanErr := p.buildChannelData(&email)
	if chanErr != nil {
		logrus.Errorf("failed to build email channel data for email with id %v: %v", emailIdString, chanErr)
		return entity.ERROR, nil, chanErr
	}

	return p.processInboundEmail(ctx, tenant, &email, rawEmail, now)

}

func (p *emailService) processInboundEmail(ctx context.Context, tenant string, email *EmailMessageData, rawEmail *entity.RawEmail, ts time.Time) (entity.RawState, *string, error) {
	session := utils.NewNeo4jWriteSession(ctx, *p.services.Neo4jRepositories.Neo4jDriver)
	defer session.Close(ctx)

	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		return entity.ERROR, nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Close(ctx)

	// Process session and events
	if err := p.processSessionAndEvents(ctx, tx, tenant, email, rawEmail, ts); err != nil {
		return entity.ERROR, nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return entity.ERROR, nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return entity.PROCESSED, nil, nil

}

func (p *emailService) processSessionAndEvents(ctx context.Context, tx neo4j.ExplicitTransaction, tenant string, email *EmailMessageData, rawEmail *entity.RawEmail, ts time.Time) error {

	// get EmailForCustomerOS
	cosEmail := p.buildEmailForCustomerOS(email, rawEmail)

	// Create session
	sessionId, err := p.repositories.InteractionEventRepository.MergeInteractionSession(ctx, tx, tenant, email.Identifiers.EmailThreadId, ts, cosEmail, rawEmail.ExternalSystem, AppSource)
	if err != nil {
		return fmt.Errorf("failed merge interaction session: %v", err)
	}

	// Create event
	eventId, err := p.repositories.InteractionEventRepository.MergeEmailInteractionEvent(ctx, tx, tenant, ts, email, rawEmail.ExternalSystem, AppSource)
	if err != nil {
		return fmt.Errorf("failed merge interaction event: %v", err)
	}

	// Link event to session
	if err := p.repositories.InteractionEventRepository.LinkInteractionEventToSession(ctx, tx, tenant, eventId, sessionId); err != nil {
		return fmt.Errorf("failed to link event to session: %v", err)
	}

	// Process participants
	if err := p.linkParticipants(ctx, tx, tenant, eventId, &email.Participants, ts, rawEmail.ExternalSystem); err != nil {
		return err
	}

	return nil

}

func (p *emailService) buildEmailForCustomerOS(email *EmailMessageData, rawEmail *entity.RawEmail) SaveEmailMessage {
	save := SaveEmailMessage{
		Html:           email.Content.Html,
		Text:           email.Content.Text,
		Subject:        email.Content.Subject,
		CreatedAt:      email.CreatedAt,
		ExternalSystem: rawEmail.ExternalSystem,
		ExternalId:     email.Identifiers.ExternalId,
		EmailThreadId:  email.Identifiers.EmailThreadId,
		Channel:        "EMAIL",
		ChannelData:    email.ChannelData,
	}
	return save
}

func (p *emailService) linkParticipants(ctx context.Context, tx neo4j.ExplicitTransaction, tenant string, eventId string, participants *EmailParticipants, now time.Time, externalSystem string) error {
	emailIds := make(map[string]string)

	// Link From participant
	fromId, err := p.getOrCreateEmailId(ctx, tx, tenant, participants.From.Email, now, externalSystem, emailIds)
	if err != nil {
		return err
	}
	if err := p.repositories.InteractionEventRepository.InteractionEventSentByEtail(ctx, tx, tenant, eventId, fromId); err != nil {
		return err
	}

	// Link To participants
	if err := p.linkEmailGroup(ctx, tx, tenant, eventId, "TO", participants.To, now, externalSystem, emailIds); err != nil {
		return err
	}

	// Link CC participants
	if err := p.linkEmailGroup(ctx, tx, tenant, eventId, "CC", participants.Cc, now, externalSystem, emailIds); err != nil {
		return err
	}

	// Link BCC participants
	if err := p.linkEmailGroup(ctx, tx, tenant, eventId, "BCC", participants.Bcc, now, externalSystem, emailIds); err != nil {
		return err
	}

	return nil
}

func (p *emailService) linkEmailGroup(ctx context.Context, tx neo4j.Transaction, tenant string, eventId string, groupType string, emails []string, now time.Time, externalSystem string, emailIds map[string]string) error {
	var groupEmailIds []string

	for _, email := range emails {
		if email == "" {
			continue
		}

		emailId, err := p.getOrCreateEmailId(ctx, tx, tenant, email, now, externalSystem, emailIds)
		if err != nil {
			return err
		}

		if !utils.Contains(groupEmailIds, emailId) {
			groupEmailIds = append(groupEmailIds, emailId)
		}
	}

	if len(groupEmailIds) > 0 {
		return p.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tx, tenant, eventId, groupType, groupEmailIds)
	}

	return nil
}

func (p *emailService) getOrCreateEmailId(ctx context.Context, tx neo4j.Transaction, tenant string, email string, now time.Time, externalSystem string, emailIds map[string]string) (string, error) {
	if id, exists := emailIds[email]; exists {
		return id, nil
	}

	id, err := p.services.SyncService.GetEmailIdForEmail(ctx, tx, tenant, email, now, externalSystem)
	if err != nil {
		return "", fmt.Errorf("failed to get email ID for %s: %v", email, err)
	}
	if id == "" {
		return "", fmt.Errorf("no email ID found for %s", email)
	}

	emailIds[email] = id
	return id, nil
}

func (p *emailService) buildChannelData(email *EmailMessageData) error {
	channelData, err := neo4jentity.BuildEmailChannelData(email.Identifiers.ProviderMessageId, email.Identifiers.EmailThreadId, email.Content.Subject, strings.Join(email.Participants.InReplyTo.Email, " "), strings.Join(email.Identifiers.References, " "))
	if err != nil {
		return err
	}
	email.Channel = "EMAIL"
	email.ChannelData = channelData

	return nil
}

func (p *emailService) validateDates(sent time.Time) (time.Time, error) {
	emailSentDate, err := p.services.SyncService.ConvertToUTC(sent)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to convert sent date to UTC: %v", err)
	}
	return emailSentDate, nil
}

func (p *emailService) warmingEmailCheck(email EmailMessageData) bool {

	emailExclusion := p.services.Cache.GetEmailExclusion(tenant)

	for _, exclusion := range emailExclusion {
		if exclusion.ExcludeSubject != nil {
			if strings.Contains(rawEmailData.Subject, *exclusion.ExcludeSubject) {
				reason := "excluded by subject"
				return entity.SKIPPED, &reason, nil
			}
		}
		if exclusion.ExcludeBody != nil {
			if strings.Contains(rawEmailData.Html, *exclusion.ExcludeBody) {
				reason := "excluded by html body"
				return entity.SKIPPED, &reason, nil
			}
			if strings.Contains(rawEmailData.Text, *exclusion.ExcludeBody) {
				reason := "excluded by text body"
				return entity.SKIPPED, &reason, nil
			}
		}
	}
}

func convertToUTC(datetimeStr string) (time.Time, error) {
	var err error

	layouts := []string{
		"2006-01-02T15:04:05Z07:00",

		"Mon, 2 Jan 2006 15:04:05 -0700 (MST)",

		"Mon, 2 Jan 2006 15:04:05 MST",

		"Mon, 2 Jan 2006 15:04:05 -0700",

		"Mon, 2 Jan 2006 15:04:05 +0000 (GMT)",

		"Mon, 2 Jan 2006 15:04:05 -0700 (MST)",

		"2 Jan 2006 15:04:05 -0700",
	}
	var parsedTime time.Time

	// Try parsing with each layout until successful
	for _, layout := range layouts {
		parsedTime, err = time.Parse(layout, datetimeStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		return time.Time{}, fmt.Errorf("unable to parse datetime string: %s", datetimeStr)
	}

	return parsedTime.UTC(), nil
}
