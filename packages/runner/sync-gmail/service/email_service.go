package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/sirupsen/logrus"
	"net/mail"
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
	FindEmailForUser(tenant, userId string) (*entity.EmailEntity, error)
	SyncEmails(externalSystemId, tenant string) error
	SyncEmail(externalSystemId, tenant string, emailId uuid.UUID) error
}

func (s *emailService) FindEmailForUser(tenant, userId string) (*entity.EmailEntity, error) {
	ctx := context.Background()

	email, err := s.repositories.EmailRepository.FindUserByEmail(ctx, tenant, userId)
	if err != nil {
		logrus.Errorf("failed to find user by email: %v", err)
		return nil, err
	}
	if email == nil {
		return nil, nil
	}

	return s.mapDbNodeToEmailEntity(*email), nil
}

func (s *emailService) SyncEmails(externalSystemId, tenant string) error {
	emailsIdsForSync, err := s.repositories.RawEmailRepository.GetEmailsIdsForSync(externalSystemId, tenant)
	if err != nil {
		logrus.Errorf("failed to get emails for sync: %v", err)
		return err
	}

	for _, emailForSync := range emailsIdsForSync {

		err := s.SyncEmail(externalSystemId, tenant, emailForSync.ID)

		var errMessage string
		if err != nil {
			errMessage = err.Error()
		}

		err = s.repositories.RawEmailRepository.MarkSentToEventStore(emailForSync.ID, err == nil, errMessage)

		if err != nil {
			logrus.Errorf("unable to mark email as sent to event store: %v", err)
			return err
		}

		fmt.Println("raw email processed: " + emailForSync.ID.String())
	}

	return nil
}

func (s *emailService) SyncEmail(externalSystemId, tenant string, emailId uuid.UUID) error {
	ctx := context.Background()

	emailIdString := emailId.String()

	rawEmail, err := s.repositories.RawEmailRepository.GetEmailForSync(emailId)
	if err != nil {
		logrus.Errorf("failed to get emails for sync: %v", err)
		return err
	}

	rawEmailData := EmailRawData{}
	err = json.Unmarshal([]byte(rawEmail.Data), &rawEmailData)
	if err != nil {
		logrus.Errorf("failed to unmarshal raw email data: %v", err)
		return err
	}

	interactionEventId, err := s.repositories.InteractionEventRepository.GetInteractionEventIdByExternalId(ctx, tenant, rawEmail.MessageId)
	if err != nil {
		logrus.Errorf("failed to check if interaction event exists for external id %v for tenant %v :%v", rawEmail.MessageId, tenant, err)
		return err
	}

	if interactionEventId == "" {

		now := time.Now().UTC()

		emailSentDate, err := convertToUTC(rawEmailData.Sent)
		if err != nil {
			logrus.Errorf("failed to convert email sent date to UTC for email with id %v :%v", emailIdString, err)
			return err
		}

		from := extractEmailAddresses(rawEmailData.From)[0]
		to := extractEmailAddresses(rawEmailData.To)
		cc := extractEmailAddresses(rawEmailData.Cc)
		bcc := extractEmailAddresses(rawEmailData.Bcc)

		references := extractLines(rawEmailData.Reference)
		inReplyTo := extractLines(rawEmailData.InReplyTo)

		channelData, err := buildEmailChannelData(rawEmailData.Subject, references, inReplyTo)
		if err != nil {
			logrus.Errorf("failed to build email channel data for email with id %v: %v", emailIdString, err)
			return err
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
			return err
		}

		sessionIdentifier := ""
		if references != nil && len(references) > 0 {
			sessionIdentifier = references[0]
		} else {
			sessionIdentifier = rawEmailData.MessageId
		}

		sessionId, err := s.repositories.InteractionEventRepository.MergeInteractionSession(ctx, tx, tenant, sessionIdentifier, now, emailForCustomerOS)
		if err != nil {
			logrus.Errorf("failed merge interaction session for raw email id %v :%v", emailIdString, err)
			return err
		}

		interactionEventId, err = s.repositories.InteractionEventRepository.MergeEmailInteractionEvent(ctx, tx, tenant, now, emailForCustomerOS)
		if err != nil {
			logrus.Errorf("failed merge interaction event for raw email id %v :%v", emailIdString, err)
			return err
		}

		err = s.repositories.InteractionEventRepository.LinkInteractionEventToSession(ctx, tx, tenant, interactionEventId, sessionId)
		if err != nil {
			logrus.Errorf("failed to associate interaction event to session for raw email id %v :%v", emailIdString, err)
			return err
		}

		logrus.Println("fetching email classification")
		emailsClassification, err := s.services.OpenAiService.FetchEmailsClassification(tenant, rawEmail.MessageId, from, to, cc, bcc)
		if err != nil {
			logrus.Errorf("unable to fetch email classification: %v", err)
			return err
		}

		//from
		//check if domain exists for tenant by email. if so, link the email to the user otherwise create a contact and link the email to the contact
		fromEmailId, err := s.getEmailIdForEmail(ctx, tx, tenant, emailsClassification, from, now)
		if err != nil {
			logrus.Errorf("unable to retrieve email id for tenant: %v", err)
			return err
		}
		if fromEmailId == "" {
			logrus.Errorf("unable to retrieve email id for tenant %s and email %s", tenant, from)
			return err
		}

		err = s.repositories.InteractionEventRepository.InteractionEventSentByEmail(ctx, tx, tenant, interactionEventId, fromEmailId)
		if err != nil {
			logrus.Errorf("unable to link email to interaction event: %v", err)
			return err
		}

		//to
		for _, toEmail := range to {
			toEmailId, err := s.getEmailIdForEmail(ctx, tx, tenant, emailsClassification, toEmail, now)
			if err != nil {
				logrus.Errorf("unable to retrieve email id for tenant: %v", err)
				return err
			}
			if toEmailId == "" {
				logrus.Errorf("unable to retrieve email id for tenant %s and email %s", tenant, toEmail)
				return err
			}

			err = s.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tx, tenant, interactionEventId, "TO", []string{toEmailId})
			if err != nil {
				logrus.Errorf("unable to link email to interaction event: %v", err)
				return err
			}
		}

		//cc
		for _, ccEmail := range cc {
			ccEmailId, err := s.getEmailIdForEmail(ctx, tx, tenant, emailsClassification, ccEmail, now)
			if err != nil {
				logrus.Errorf("unable to retrieve email id for tenant: %v", err)
				return err
			}
			if ccEmailId == "" {
				logrus.Errorf("unable to retrieve email id for tenant %s and email %s", tenant, ccEmail)
				return err
			}

			err = s.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tx, tenant, interactionEventId, "CC", []string{ccEmailId})
			if err != nil {
				logrus.Errorf("unable to link email to interaction event: %v", err)
				return err
			}
		}

		//bcc
		for _, bccEmail := range bcc {
			bccEmailId, err := s.getEmailIdForEmail(ctx, tx, tenant, emailsClassification, bccEmail, now)
			if err != nil {
				logrus.Errorf("unable to retrieve email id for tenant: %v", err)
				return err
			}
			if bccEmailId == "" {
				logrus.Errorf("unable to retrieve email id for tenant %s and email %s", tenant, bccEmail)
				return err
			}

			err = s.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tx, tenant, interactionEventId, "BCC", []string{bccEmailId})
			if err != nil {
				logrus.Errorf("unable to link email to interaction event: %v", err)
				return err
			}
		}

		err = tx.Commit(ctx)
		if err != nil {
			logrus.Errorf("failed to commit transaction: %v", err)
			return err
		}

	} else {
		logrus.Infof("interaction event already exists for raw email id %v", emailIdString)
	}

	//TODO HERE PUSH COMMANDS TO EVENT STORE TO GENERATE SUMMARY AND ACTION ITEMS

	//send to event store interactionEventId + html / text
	//promp + promp version

	////get summary for email using claude-2 model
	//summaryExists, err := s.repositories.AnalysisRepository.SummaryExistsForInteractionEvent(ctx, tenant, interactionEventId)
	//if err != nil {
	//	logrus.Errorf("unable to check if summary exists for interaction event: %v", err)
	//	return err
	//}
	//if !summaryExists {
	//	logrus.Println("fetching anthropic summary for email")
	//	summary := s.services.AnthropicService.FetchSummary(utils.StringFirstNonEmpty(rawEmailData.Html, rawEmailData.Text))
	//	_, err = s.repositories.AnalysisRepository.CreateSummaryForEmail(ctx, tenant, interactionEventId, summary, externalSystemId, "sync-gmail", now)
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
	//		_, err = s.repositories.ActionItemRepository.CreateActionItemForEmail(ctx, tenant, interactionEventId, actionItem, externalSystemId, "sync-gmail", now)
	//		if err != nil {
	//			logrus.Errorf("unable to create action item for email: %v", err)
	//			return err
	//		}
	//	}
	//}

	return nil
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

func convertToUTC(datetimeStr string) (time.Time, error) {
	var err error

	t1, err := time.Parse(time.RFC3339, datetimeStr)
	if err == nil {
		return t1.UTC(), nil
	}

	t2, err := time.Parse(time.RFC1123Z, datetimeStr)
	if err == nil {
		return t2.UTC(), nil
	}

	layouts := []string{
		"Mon, 2 Jan 2006 15:04:05 -0700 (MST)",
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 02 Jan 2006 15:04:05 +0000 (GMT)",
		"Thu, 29 Jun 2023 03:53:38 -0700 (PDT)",
		"Wed, 29 Sep 2021 13:02:25 GMT",
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
		return time.Time{}, fmt.Errorf("unable to parse datetime string: %v", err)
	}

	return parsedTime.UTC(), nil
}

func getEmailClassificationByEmail(email string, emailClassifications []*EmailClassification) *EmailClassification {
	for _, emailClassification := range emailClassifications {
		if emailClassification.Email == email {
			return emailClassification
		}
	}
	return nil
}

func extractLines(input string) []string {
	lines := strings.Fields(input)
	return lines
}

func extractEmailAddresses(input string) []string {
	// Regular expression pattern to match email addresses between <>
	emailPattern := `<(.*?)>`

	// Check if the input contains angle brackets
	if strings.Contains(input, "<") && strings.Contains(input, ">") {
		// Extract email addresses using the regular expression pattern
		re := regexp.MustCompile(emailPattern)
		matches := re.FindAllStringSubmatch(input, -1)

		// Create a map to store unique email addresses
		emailMap := make(map[string]bool)
		for _, match := range matches {
			email := match[1]
			emailMap[email] = true
		}

		// Convert the map keys to an array of email addresses
		emailAddresses := make([]string, 0, len(emailMap))
		for email := range emailMap {
			if isValidEmailSyntax(email) {
				emailAddresses = append(emailAddresses, email)
			}
		}

		return emailAddresses
	} else if strings.Contains(input, ",") {
		// Split the input string by commas
		emails := make([]string, 0)
		split := strings.Split(input, ",")

		// Trim any whitespace from the email addresses
		for _, email := range split {
			email = strings.TrimSpace(email)
			if isValidEmailSyntax(email) {
				emails = append(emails, email)
			}
		}

		return emails
	}

	if input != "" {
		return []string{input}
	}

	return nil
}

func isValidEmailSyntax(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "" // Invalid email format
	}
	split := strings.Split(parts[1], ".")
	if len(split) < 2 {
		return parts[1]
	}
	return strings.ToLower(split[len(split)-2] + "." + split[len(split)-1])
}

func (s *emailService) getEmailIdForEmail(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, aiClassification []*EmailClassification, email string, now time.Time) (string, error) {

	fromEmailId, err := s.repositories.EmailRepository.GetEmailId(ctx, tenant, email)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve email id for tenant: %v", err)
	}
	if fromEmailId != "" {
		return fromEmailId, nil
	}

	domainNode, err := s.repositories.DomainRepository.GetDomain(ctx, extractDomain(email))
	if err != nil {
		return "", fmt.Errorf("unable to retrieve domain for tenant: %v", err)
	}

	emailClassification := getEmailClassificationByEmail(email, aiClassification)

	//domain exists already and the AI hasn't been called for it
	if emailClassification == nil && domainNode != nil {
		emailClassification = &EmailClassification{
			Email:               email,
			IsOrganizationEmail: true,
		}
	}

	if emailClassification == nil && domainNode == nil {
		logrus.Errorf("unable to classify email: %v", email)
		return "", fmt.Errorf("unable to classify email: %v", email)
	}

	if emailClassification.IsPersonalEmail {
		emailId, err := s.repositories.EmailRepository.CreateContactWithEmail(ctx, tx, tenant, emailClassification.Email, emailClassification.PersonalFirstName, emailClassification.PersonalLastName, "syng-gmail")
		if err != nil {
			return "", fmt.Errorf("unable to create contact with email: %v", err)
		}
		return emailId, nil
	} else if emailClassification.IsOrganizationEmail {

		if domainNode != nil {
			organizationNode, err := s.repositories.OrganizationRepository.GetOrganizationWithDomain(ctx, tenant, utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*domainNode), "domain"))
			if err != nil {
				return "", fmt.Errorf("unable to retrieve organization for tenant: %v", err)
			}

			if organizationNode == nil {
				organizationNode, err = s.repositories.OrganizationRepository.CreateOrganization(ctx, tx, tenant, emailClassification.OrganizationName, "gmail", "openline", "sync-gmail", now)
				if err != nil {
					return "", fmt.Errorf("unable to create organization for tenant: %v", err)
				}
			}

			emailId, err := s.repositories.EmailRepository.CreateEmailLinkedToOrganization(ctx, tx, tenant, emailClassification.Email, "gmail", "openline", "sync-gmail", utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*organizationNode), "id"), now)
			if err != nil {
				return "", fmt.Errorf("unable to create email linked to organization: %v", err)
			}

			return emailId, nil
		} else {
			domainNode, err := s.repositories.DomainRepository.CreateDomain(ctx, extractDomain(emailClassification.Email), "gmail", "sync-gmail", now)
			if err != nil {
				return "", fmt.Errorf("unable to create domain: %v", err)
			}
			organizationNode, err := s.repositories.OrganizationRepository.CreateOrganization(ctx, tx, tenant, emailClassification.OrganizationName, "gmail", "openline", "sync-gmail", now)
			if err != nil {
				return "", fmt.Errorf("unable to create organization for tenant: %v", err)
			}
			domainName := utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*domainNode), "domain")
			organizationId := utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*organizationNode), "id")
			err = s.repositories.OrganizationRepository.LinkDomainToOrganization(ctx, tx, tenant, domainName, organizationId)
			if err != nil {
				return "", fmt.Errorf("unable to link domain to organization: %v", err)
			}

			emailId, err := s.repositories.EmailRepository.CreateEmailLinkedToOrganization(ctx, tx, tenant, emailClassification.Email, "gmail", "openline", "sync-gmail", utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*organizationNode), "id"), now)
			if err != nil {
				return "", fmt.Errorf("unable to create email linked to organization: %v", err)
			}

			return emailId, nil
		}

	} else {
		return "", fmt.Errorf("unable to determine email type: %v", err)
	}
}

type EmailRawData struct {
	MessageId string            `json:"MessageId"`
	Sent      string            `json:"Sent"`
	Subject   string            `json:"Subject"`
	From      string            `json:"From"`
	To        string            `json:"To"`
	Cc        string            `json:"Cc"`
	Bcc       string            `json:"Bcc"`
	Html      string            `json:"Html"`
	Text      string            `json:"Text"`
	ThreadId  string            `json:"ThreadId"`
	InReplyTo string            `json:"InReplyTo"`
	Reference string            `json:"Reference"`
	Headers   map[string]string `json:"Headers"`
}

func (s *emailService) mapDbNodeToEmailEntity(node dbtype.Node) *entity.EmailEntity {
	props := utils.GetPropsFromNode(node)
	result := entity.EmailEntity{
		Id:       utils.GetStringPropOrEmpty(props, "id"),
		Email:    utils.GetStringPropOrEmpty(props, "email"),
		RawEmail: utils.GetStringPropOrEmpty(props, "rawEmail"),
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
