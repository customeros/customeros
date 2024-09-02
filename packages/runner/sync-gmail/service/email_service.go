package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

type emailService struct {
	cfg          *config.Config
	repositories *repository.Repositories
	services     *Services
}

type EmailService interface {
	FindEmailsForUser(tenant, userId string) ([]*entity.EmailEntity, error)

	SyncEmailsForUser(tenant, userSource string)

	SyncEmailByEmailRawId(tenant string, emailId uuid.UUID) (entity.RawState, *string, error)
	SyncEmailByMessageId(tenant, usernameSource, messageId string) (entity.RawState, *string, error)
}

func (s *emailService) FindEmailsForUser(tenant, userId string) ([]*entity.EmailEntity, error) {
	ctx := context.Background()

	emailNodes, err := s.repositories.EmailRepository.FindEmailsByUserId(ctx, tenant, userId)
	if err != nil {
		logrus.Errorf("failed to find user by email: %v", err)
		return nil, err
	}

	emails := make([]*entity.EmailEntity, len(emailNodes))
	for i, node := range emailNodes {
		emails[i] = s.mapDbNodeToEmailEntity(*node)
	}

	return emails, nil
}

func (s *emailService) SyncEmailsForUser(tenant string, userSource string) {
	emailsIdsForSync, err := s.repositories.RawEmailRepository.GetEmailsIdsForUserForSync(tenant, userSource)
	if err != nil {
		logrus.Errorf("failed to get emails for sync: %v", err)
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
		logrus.Errorf("failed to create user source as email node: %v", err)
		return
	}

	for _, externalSystem := range distinctExternalSystems {
		err = s.repositories.Neo4jRepositories.ExternalSystemWriteRepository.CreateIfNotExists(context.Background(), tenant, externalSystem, externalSystem)
		if err != nil {
			return
		}
	}

	s.syncEmails(tenant, emailsIdsForSync)
}

func (s *emailService) SyncEmailByEmailRawId(tenant string, emailId uuid.UUID) (entity.RawState, *string, error) {
	return s.syncEmail(tenant, emailId)
}

func (s *emailService) SyncEmailByMessageId(tenant, usernameSource, messageId string) (entity.RawState, *string, error) {
	rawEmail, err := s.repositories.RawEmailRepository.GetEmailForSyncByMessageId(tenant, usernameSource, messageId)
	if err != nil {
		logrus.Errorf("failed to get emails for sync: %v", err)
		return entity.ERROR, nil, err
	}

	if rawEmail == nil {
		return entity.ERROR, nil, fmt.Errorf("email with message id %v not found", messageId)
	}

	return s.syncEmail(tenant, rawEmail.ID)
}

func (s *emailService) syncEmails(tenant string, emails []entity.RawEmail) {
	for _, email := range emails {
		state, reason, err := s.syncEmail(tenant, email.ID)

		var errMessage *string
		if err != nil {
			s2 := err.Error()
			errMessage = &s2
		}

		err = s.repositories.RawEmailRepository.MarkSentToEventStore(email.ID, state, reason, errMessage)
		if err != nil {
			logrus.Errorf("unable to mark email as sent to event store: %v", err)
		}

		fmt.Println("raw email processed: " + email.ID.String())
	}
}

func (s *emailService) syncEmail(tenant string, emailId uuid.UUID) (entity.RawState, *string, error) {
	ctx := context.Background()
	span, ctx := tracing.StartTracerSpan(ctx, "EmailService.syncEmail")
	defer span.Finish()
	span.LogFields(log.String("emailId", emailId.String()))

	emailIdString := emailId.String()

	rawEmail, err := s.repositories.RawEmailRepository.GetEmailForSync(emailId)
	if err != nil {
		logrus.Errorf("failed to get emails for sync: %v", err)
		return entity.ERROR, nil, err
	}

	if rawEmail.MessageId == "" {
		return entity.ERROR, nil, fmt.Errorf("message id is empty")
	}

	rawEmailData := EmailRawData{}
	err = json.Unmarshal([]byte(rawEmail.Data), &rawEmailData)
	if err != nil {
		logrus.Errorf("failed to unmarshal raw email data: %v", err)
		return entity.ERROR, nil, err
	}

	emailExclusion := s.services.Cache.GetEmailExclusion(tenant)

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

	interactionEventId, err := s.repositories.InteractionEventRepository.GetInteractionEventIdByExternalId(ctx, tenant, rawEmail.ExternalSystem, rawEmail.MessageId)
	if err != nil {
		logrus.Errorf("failed to check if interaction event exists for external id %v for tenant %v :%v", rawEmail.MessageId, tenant, err)
		return entity.ERROR, nil, err
	}

	if interactionEventId == "" {

		now := time.Now().UTC()

		emailSentDate, err := s.services.SyncService.ConvertToUTC(rawEmailData.Sent)
		if err != nil {
			logrus.Errorf("failed to convert email sent date to UTC for email with id %v :%v", emailIdString, err)
			return entity.ERROR, nil, err
		}

		from := s.extractEmailAddresses(rawEmailData.From)[0]
		to := s.extractEmailAddresses(rawEmailData.To)
		cc := s.extractEmailAddresses(rawEmailData.Cc)
		bcc := s.extractEmailAddresses(rawEmailData.Bcc)

		references := extractLines(rawEmailData.Reference)
		inReplyTo := extractLines(rawEmailData.InReplyTo)

		allEmailsString, err := s.services.SyncService.BuildEmailsListExcludingPersonalEmails(rawEmail.Username, from, to, cc, bcc)
		if err != nil {
			logrus.Errorf("failed to build emails list: %v", err)
			return entity.ERROR, nil, err
		}

		if len(allEmailsString) == 0 {
			reason := "no emails address belongs to a workspace domain"
			return entity.SKIPPED, &reason, nil
		}

		// Create a map to store the domain counts
		domainCount := make(map[string]int)

		// Iterate through the email addresses
		for _, email := range allEmailsString {
			domain := utils.ExtractDomain(email)
			if domain != "" {
				domainCount[domain]++
			}
		}

		if len(domainCount) > 5 {
			reason := "more than 5 domains belongs to a workspace domain"
			return entity.SKIPPED, &reason, nil
		}

		channelData, err := dto.BuildEmailChannelData(rawEmailData.ProviderMessageId, rawEmailData.ThreadId, rawEmailData.Subject, references, inReplyTo)
		if err != nil {
			logrus.Errorf("failed to build email channel data for email with id %v: %v", emailIdString, err)
			return entity.ERROR, nil, err
		}

		emailForCustomerOS := entity.EmailMessageData{
			Html:           rawEmailData.Html,
			Text:           rawEmailData.Text,
			Subject:        rawEmailData.Subject,
			CreatedAt:      emailSentDate,
			ExternalSystem: rawEmail.ExternalSystem,
			ExternalId:     rawEmailData.MessageId,
			EmailThreadId:  rawEmailData.ThreadId,
			Channel:        "EMAIL",
			ChannelData:    channelData,
		}

		session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Neo4jDriver)
		defer session.Close(ctx)

		tx, err := session.BeginTransaction(ctx)
		if err != nil {
			logrus.Errorf("failed to start transaction for email with id %v: %v", emailIdString, err)
			return entity.ERROR, nil, err
		}

		sessionId, err := s.repositories.InteractionEventRepository.MergeInteractionSession(ctx, tx, tenant, emailForCustomerOS.EmailThreadId, now, emailForCustomerOS, rawEmail.ExternalSystem, AppSource)
		if err != nil {
			logrus.Errorf("failed merge interaction session for raw email id %v :%v", emailIdString, err)
			return entity.ERROR, nil, err
		}

		interactionEventId, err = s.repositories.InteractionEventRepository.MergeEmailInteractionEvent(ctx, tx, tenant, now, emailForCustomerOS, rawEmail.ExternalSystem, AppSource)
		if err != nil {
			logrus.Errorf("failed merge interaction event for raw email id %v :%v", emailIdString, err)
			return entity.ERROR, nil, err
		}

		err = s.repositories.InteractionEventRepository.LinkInteractionEventToSession(ctx, tx, tenant, interactionEventId, sessionId)
		if err != nil {
			logrus.Errorf("failed to associate interaction event to session for raw email id %v :%v", emailIdString, err)
			return entity.ERROR, nil, err
		}

		emailidList := []string{}

		//from
		//check if domain exists for tenant by email. if so, link the email to the user otherwise create a contact and link the email to the contact
		fromEmailId, err := s.services.SyncService.GetEmailIdForEmail(ctx, tx, tenant, from, now, rawEmail.ExternalSystem)
		if err != nil {
			logrus.Errorf("unable to retrieve email id for tenant: %v", err)
			return entity.ERROR, nil, err
		}
		if fromEmailId == "" {
			logrus.Errorf("unable to retrieve email id for tenant %s and email %s", tenant, from)
			return entity.ERROR, nil, err
		}

		err = s.repositories.InteractionEventRepository.InteractionEventSentByEmail(ctx, tx, tenant, interactionEventId, fromEmailId)
		if err != nil {
			logrus.Errorf("unable to link email to interaction event: %v", err)
			return entity.ERROR, nil, err
		}
		emailidList = append(emailidList, fromEmailId)

		//to
		for _, toEmail := range to {
			toEmailId, err := s.services.SyncService.GetEmailIdForEmail(ctx, tx, tenant, toEmail, now, rawEmail.ExternalSystem)
			if err != nil {
				logrus.Errorf("unable to retrieve email id for tenant: %v", err)
				return entity.ERROR, nil, err
			}
			if toEmailId == "" {
				logrus.Errorf("unable to retrieve email id for tenant %s and email %s", tenant, toEmail)
				return entity.ERROR, nil, err
			}
			if utils.Contains(emailidList, toEmailId) {
				continue
			}

			err = s.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tx, tenant, interactionEventId, "TO", []string{toEmailId})
			if err != nil {
				logrus.Errorf("unable to link email to interaction event: %v", err)
				return entity.ERROR, nil, err
			}
			emailidList = append(emailidList, toEmailId)
		}

		//cc
		for _, ccEmail := range cc {
			if ccEmail == "" {
				continue
			}
			ccEmailId, err := s.services.SyncService.GetEmailIdForEmail(ctx, tx, tenant, ccEmail, now, rawEmail.ExternalSystem)
			if err != nil {
				logrus.Errorf("unable to retrieve email id for tenant: %v", err)
				return entity.ERROR, nil, err
			}
			if ccEmailId == "" {
				logrus.Errorf("unable to retrieve email id for tenant %s and email %s", tenant, ccEmail)
				return entity.ERROR, nil, err
			}
			if utils.Contains(emailidList, ccEmailId) {
				continue
			}

			err = s.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tx, tenant, interactionEventId, "CC", []string{ccEmailId})
			if err != nil {
				logrus.Errorf("unable to link email to interaction event: %v", err)
				return entity.ERROR, nil, err
			}
			emailidList = append(emailidList, ccEmailId)
		}

		//bcc
		for _, bccEmail := range bcc {
			if bccEmail == "" {
				continue
			}

			bccEmailId, err := s.services.SyncService.GetEmailIdForEmail(ctx, tx, tenant, bccEmail, now, rawEmail.ExternalSystem)
			if err != nil {
				logrus.Errorf("unable to retrieve email id for tenant: %v", err)
				return entity.ERROR, nil, err
			}
			if bccEmailId == "" {
				logrus.Errorf("unable to retrieve email id for tenant %s and email %s", tenant, bccEmail)
				return entity.ERROR, nil, err
			}
			if utils.Contains(emailidList, bccEmailId) {
				continue
			}

			err = s.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tx, tenant, interactionEventId, "BCC", []string{bccEmailId})
			if err != nil {
				logrus.Errorf("unable to link email to interaction event: %v", err)
				return entity.ERROR, nil, err
			}

			emailidList = append(emailidList, bccEmailId)
		}

		err = tx.Commit(ctx)
		if err != nil {
			logrus.Errorf("failed to commit transaction: %v", err)
			return entity.ERROR, nil, err
		}

	} else {
		logrus.Infof("interaction event already exists for raw email id %v", emailIdString)
		reason := "interaction event already exists"
		return entity.SKIPPED, &reason, nil
	}

	return entity.SENT, nil, err
}

func extractLines(input string) []string {
	lines := strings.Fields(input)
	return lines
}

func (s *emailService) extractEmailAddresses(input string) []string {
	if input == "" {
		return []string{""}
	}
	// Regular expression pattern to match email addresses between <>
	emailPattern := `<(.*?)>`

	emails := make([]string, 0)
	emailAddresses := make([]string, 0)

	if strings.Contains(input, ",") {
		split := strings.Split(input, ",")

		for _, email := range split {
			email = strings.TrimSpace(email)
			email = strings.ToLower(email)
			emails = append(emails, email)
		}
	} else {
		emails = append(emails, input)
	}

	for _, email := range emails {
		email = strings.TrimSpace(email)
		email = strings.ToLower(email)
		if strings.Contains(email, "<") && strings.Contains(email, ">") {
			// Extract email addresses using the regular expression pattern
			re := regexp.MustCompile(emailPattern)
			matches := re.FindAllStringSubmatch(email, -1)

			// Create a map to store unique email addresses
			emailMap := make(map[string]bool)
			for _, match := range matches {
				email := match[1]
				emailMap[email] = true
			}

			// Convert the map keys to an array of email addresses
			for email := range emailMap {
				if s.services.SyncService.IsValidEmailSyntax(email) {
					emailAddresses = append(emailAddresses, email)
				}
			}

		} else if s.services.SyncService.IsValidEmailSyntax(email) {
			emailAddresses = append(emailAddresses, email)
		}
	}

	if len(emailAddresses) > 0 {
		return emailAddresses
	}

	return []string{input}
}

func (s *emailService) createUserSourceAsEmailNode(tenant, userSource, externalSystem string) (string, error) {

	ctx := context.Background()

	session := utils.NewNeo4jWriteSession(ctx, *s.repositories.Neo4jDriver)
	defer session.Close(ctx)

	tx, err := session.BeginTransaction(ctx)

	emailId, err := s.repositories.EmailRepository.CreateEmail(ctx, tx, tenant, userSource, externalSystem, AppSource)
	if err != nil {
		return "", fmt.Errorf("unable to create email: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		logrus.Errorf("failed to commit transaction: %v", err)
		return "", err
	}

	return emailId, nil
}

type EmailRawData struct {
	ProviderMessageId string            `json:"ProviderMessageId"`
	MessageId         string            `json:"MessageId"`
	Sent              string            `json:"Sent"`
	Subject           string            `json:"Subject"`
	From              string            `json:"From"`
	To                string            `json:"To"`
	Cc                string            `json:"Cc"`
	Bcc               string            `json:"Bcc"`
	Html              string            `json:"Html"`
	Text              string            `json:"Text"`
	ThreadId          string            `json:"ThreadId"`
	InReplyTo         string            `json:"InReplyTo"`
	Reference         string            `json:"Reference"`
	Headers           map[string]string `json:"Headers"`
}

func (s *emailService) mapDbNodeToEmailEntity(node dbtype.Node) *entity.EmailEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.EmailEntity{
		Id:          utils.GetStringPropOrEmpty(props, "id"),
		Email:       utils.GetStringPropOrEmpty(props, "email"),
		RawEmail:    utils.GetStringPropOrEmpty(props, "rawEmail"),
		IsReachable: utils.GetStringPropOrNil(props, "isReachable"),
	}
	return &result
}

func NewEmailService(cfg *config.Config, repositories *repository.Repositories, services *Services) EmailService {
	return &emailService{
		cfg:          cfg,
		repositories: repositories,
		services:     services,
	}
}
