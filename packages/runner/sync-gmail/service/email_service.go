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
	commonEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

const GmailSource = "gmail"

type emailService struct {
	cfg          *config.Config
	repositories *repository.Repositories
	services     *Services
}

type EmailService interface {
	FindEmailForUser(tenant, userId string) (*entity.EmailEntity, error)

	SyncEmailsForUser(externalSystemId, tenant string, userSource string, personalEmailProviderList []commonEntity.PersonalEmailProvider, organizationAllowedForImport []commonEntity.WhitelistDomain)

	SyncEmailByEmailRawId(externalSystemId, tenant string, emailId uuid.UUID) (entity.RawState, *string, error)
	SyncEmailByMessageId(externalSystemId, tenant, usernameSource, messageId string) (entity.RawState, *string, error)
}

func (s *emailService) FindEmailForUser(tenant, userId string) (*entity.EmailEntity, error) {
	ctx := context.Background()

	email, err := s.repositories.EmailRepository.FindEmailByUserId(ctx, tenant, userId)
	if err != nil {
		logrus.Errorf("failed to find user by email: %v", err)
		return nil, err
	}
	if email == nil {
		return nil, nil
	}

	return s.mapDbNodeToEmailEntity(*email), nil
}

func (s *emailService) SyncEmailsForUser(externalSystemId, tenant string, userSource string, personalEmailProviderList []commonEntity.PersonalEmailProvider, organizationAllowedForImport []commonEntity.WhitelistDomain) {
	emailsIdsForSync, err := s.repositories.RawEmailRepository.GetEmailsIdsForUserForSync(externalSystemId, tenant, userSource)
	if err != nil {
		logrus.Errorf("failed to get emails for sync: %v", err)
	}

	s.syncEmails(externalSystemId, tenant, emailsIdsForSync, personalEmailProviderList, organizationAllowedForImport)
}

func (s *emailService) SyncEmailByEmailRawId(externalSystemId, tenant string, emailId uuid.UUID) (entity.RawState, *string, error) {
	personalEmailProviderList, err := s.services.Repositories.CommonRepositories.PersonalEmailProviderRepository.GetPersonalEmailProviders()
	if err != nil {
		logrus.Errorf("failed to get personal email provider list: %v", err)
		panic(err) //todo handle error
	}

	organizationAllowedForImport, err := s.services.Repositories.CommonRepositories.WhitelistDomainRepository.GetWhitelistDomains(tenant)
	if err != nil {
		logrus.Errorf("failed to check if organization is allowed for import: %v", err)
		return entity.ERROR, nil, err
	}

	return s.syncEmail(externalSystemId, tenant, emailId, personalEmailProviderList, organizationAllowedForImport)
}

func (s *emailService) SyncEmailByMessageId(externalSystemId, tenant, usernameSource, messageId string) (entity.RawState, *string, error) {
	rawEmail, err := s.repositories.RawEmailRepository.GetEmailForSyncByMessageId(externalSystemId, tenant, usernameSource, messageId)
	if err != nil {
		logrus.Errorf("failed to get emails for sync: %v", err)
		return entity.ERROR, nil, err
	}

	if rawEmail == nil {
		return entity.ERROR, nil, fmt.Errorf("email with message id %v not found", messageId)
	}

	personalEmailProviderList, err := s.services.Repositories.CommonRepositories.PersonalEmailProviderRepository.GetPersonalEmailProviders()
	if err != nil {
		logrus.Errorf("failed to get personal email provider list: %v", err)
		panic(err) //todo handle error
	}

	organizationAllowedForImport, err := s.services.Repositories.CommonRepositories.WhitelistDomainRepository.GetWhitelistDomains(tenant)
	if err != nil {
		logrus.Errorf("failed to check if organization is allowed for import: %v", err)
		return entity.ERROR, nil, err
	}

	return s.syncEmail(externalSystemId, tenant, rawEmail.ID, personalEmailProviderList, organizationAllowedForImport)
}

func (s *emailService) syncEmails(externalSystemId, tenant string, emails []entity.RawEmail, personalEmailProviderList []commonEntity.PersonalEmailProvider, organizationAllowedForImport []commonEntity.WhitelistDomain) {
	for _, email := range emails {
		state, reason, err := s.syncEmail(externalSystemId, tenant, email.ID, personalEmailProviderList, organizationAllowedForImport)

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

func (s *emailService) syncEmail(externalSystemId, tenant string, emailId uuid.UUID, personalEmailProviderList []commonEntity.PersonalEmailProvider, whitelistDomainList []commonEntity.WhitelistDomain) (entity.RawState, *string, error) {
	ctx := context.Background()

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

	if strings.HasSuffix(rawEmailData.Subject, "• lemwarmup") || strings.HasSuffix(rawEmailData.Subject, "• lemwarm") {
		reason := "warmer email"
		return entity.SKIPPED, &reason, nil
	}

	interactionEventId, err := s.repositories.InteractionEventRepository.GetInteractionEventIdByExternalId(ctx, tenant, externalSystemId, rawEmail.MessageId)
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

		allEmailsString, err := s.services.SyncService.BuildEmailsListExcludingPersonalEmails(personalEmailProviderList, rawEmail.UsernameSource, from, to, cc, bcc)
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

		channelData, err := buildEmailChannelData(rawEmailData.Subject, references, inReplyTo)
		if err != nil {
			logrus.Errorf("failed to build email channel data for email with id %v: %v", emailIdString, err)
			return entity.ERROR, nil, err
		}

		emailForCustomerOS := entity.EmailMessageData{
			Html:           rawEmailData.Html,
			Text:           rawEmailData.Text,
			Subject:        rawEmailData.Subject,
			CreatedAt:      emailSentDate,
			ExternalSystem: externalSystemId,
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

		sessionIdentifier := ""
		if references != nil && len(references) > 0 {
			sessionIdentifier = references[0]
		} else {
			sessionIdentifier = rawEmailData.MessageId
		}

		sessionId, err := s.repositories.InteractionEventRepository.MergeInteractionSession(ctx, tx, tenant, sessionIdentifier, now, emailForCustomerOS, GmailSource, AppSource)
		if err != nil {
			logrus.Errorf("failed merge interaction session for raw email id %v :%v", emailIdString, err)
			return entity.ERROR, nil, err
		}

		interactionEventId, err = s.repositories.InteractionEventRepository.MergeEmailInteractionEvent(ctx, tx, tenant, now, emailForCustomerOS, GmailSource, AppSource)
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
		fromEmailId, err := s.services.SyncService.GetEmailIdForEmail(ctx, tx, tenant, interactionEventId, from, s.services.SyncService.GetWhitelistedDomain(utils.ExtractDomain(from), whitelistDomainList), personalEmailProviderList, now, GmailSource)
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
			toEmailId, err := s.services.SyncService.GetEmailIdForEmail(ctx, tx, tenant, interactionEventId, toEmail, s.services.SyncService.GetWhitelistedDomain(utils.ExtractDomain(toEmail), whitelistDomainList), personalEmailProviderList, now, GmailSource)
			if err != nil {
				logrus.Errorf("unable to retrieve email id for tenant: %v", err)
				return entity.ERROR, nil, err
			}
			if toEmailId == "" {
				logrus.Errorf("unable to retrieve email id for tenant %s and email %s", tenant, toEmail)
				return entity.ERROR, nil, err
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
			ccEmailId, err := s.services.SyncService.GetEmailIdForEmail(ctx, tx, tenant, interactionEventId, ccEmail, s.services.SyncService.GetWhitelistedDomain(utils.ExtractDomain(ccEmail), whitelistDomainList), personalEmailProviderList, now, GmailSource)
			if err != nil {
				logrus.Errorf("unable to retrieve email id for tenant: %v", err)
				return entity.ERROR, nil, err
			}
			if ccEmailId == "" {
				logrus.Errorf("unable to retrieve email id for tenant %s and email %s", tenant, ccEmail)
				return entity.ERROR, nil, err
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

			bccEmailId, err := s.services.SyncService.GetEmailIdForEmail(ctx, tx, tenant, interactionEventId, bccEmail, s.services.SyncService.GetWhitelistedDomain(utils.ExtractDomain(bccEmail), whitelistDomainList), personalEmailProviderList, now, GmailSource)
			if err != nil {
				logrus.Errorf("unable to retrieve email id for tenant: %v", err)
				return entity.ERROR, nil, err
			}
			if bccEmailId == "" {
				logrus.Errorf("unable to retrieve email id for tenant %s and email %s", tenant, bccEmail)
				return entity.ERROR, nil, err
			}

			err = s.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tx, tenant, interactionEventId, "BCC", []string{bccEmailId})
			if err != nil {
				logrus.Errorf("unable to link email to interaction event: %v", err)
				return entity.ERROR, nil, err
			}

			emailidList = append(emailidList, bccEmailId)
		}

		organizationIdList, err := s.repositories.OrganizationRepository.GetOrganizationsLinkedToEmailsInTx(ctx, tx, tenant, emailidList)
		if err != nil {
			logrus.Errorf("unable to retrieve organization id list for tenant: %v", err)
			return entity.ERROR, nil, err
		}

		for _, organizationId := range organizationIdList {
			lastTouchpointAt, lastTouchpointId, err := s.repositories.TimelineEventRepository.CalculateAndGetLastTouchpointInTx(ctx, tx, tenant, organizationId)
			if err != nil {
				logrus.Errorf("unable to calculate last touchpoint for organization: %v", err)
				return entity.ERROR, nil, err
			}

			if lastTouchpointAt != nil && lastTouchpointId != "" {
				err := s.repositories.OrganizationRepository.UpdateLastTouchpointInTx(ctx, tx, tenant, organizationId, *lastTouchpointAt, lastTouchpointId)
				if err != nil {
					logrus.Errorf("unable to update last touchpoint for organization: %v", err)
					return entity.ERROR, nil, err
				}
			}
		}

		err = tx.Commit(ctx)
		if err != nil {
			logrus.Errorf("failed to commit transaction: %v", err)
			return entity.ERROR, nil, err
		}

	} else {
		logrus.Infof("interaction event already exists for raw email id %v", emailIdString)
	}

	//TODO HERE PUSH COMMANDS TO EVENT STORE TO GENERATE SUMMARY AND ACTION ITEMS

	//get summary for email using claude-2 model
	//summaryExists, err := s.repositories.AnalysisRepository.SummaryExistsForInteractionEvent(ctx, tenant, interactionEventId)
	//if err != nil {
	//	logrus.Errorf("unable to check if summary exists for interaction event: %v", err)
	//	return err
	//}
	//if !summaryExists {
	//	logrus.Println("fetching anthropic summary for email")
	//	summary := s.services.AnthropicService.FetchSummary(utils.StringFirstNonEmpty(rawEmailData.Html, rawEmailData.Text))
	//	_, err = s.repositories.AnalysisRepository.CreateSummaryForEmail(ctx, tenant, interactionEventId, summary, externalSystemId, "sync-gmail", time.Now().UTC())
	//	if err != nil {
	//		logrus.Errorf("unable to create summary for email: %v", err)
	//		return err
	//	}
	//}
	//
	////get the action items for the email using claude-2 model
	//actionItemsExists, err := s.repositories.ActionItemRepository.ActionsItemsExistsForInteractionEvent(ctx, tenant, interactionEventId)
	//if err != nil {
	//	logrus.Errorf("unable to check if action items exists for interaction event: %v", err)
	//	return err
	//}
	//if !actionItemsExists {
	//	logrus.Println("fetching anthropic action items for email")
	//	actionItems := s.services.AnthropicService.FetchActionItems(utils.StringFirstNonEmpty(rawEmailData.Html, rawEmailData.Text))
	//
	//	//TODO insert should be done in a single transaction to follow ActionsItemsExistsForInteractionEvent logic
	//	for _, actionItem := range actionItems {
	//		_, err = s.repositories.ActionItemRepository.CreateActionItemForEmail(ctx, tenant, interactionEventId, actionItem, externalSystemId, AppSource, time.Now().UTC())
	//		if err != nil {
	//			logrus.Errorf("unable to create action item for email: %v", err)
	//			return err
	//		}
	//	}
	//}

	return entity.SENT, nil, err
}

type EmailChannelData struct {
	Subject   string   `json:"Subject"`
	InReplyTo []string `json:"InReplyTo"`
	Reference []string `json:"Reference"`
}

func buildEmailChannelData(subject string, references, inReplyTo []string) (*string, error) {
	emailContent := EmailChannelData{
		Subject:   subject,
		InReplyTo: utils.EnsureEmailRfcIds(inReplyTo),
		Reference: utils.EnsureEmailRfcIds(references),
	}
	jsonContent, err := json.Marshal(emailContent)
	if err != nil {
		return nil, err
	}
	jsonContentString := string(jsonContent)

	return &jsonContentString, nil
}
func extractLines(input string) []string {
	lines := strings.Fields(input)
	return lines
}

func (s *emailService) extractEmailAddresses(input string) []string {
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

	if input != "" {
		return []string{input}
	}

	return nil
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
