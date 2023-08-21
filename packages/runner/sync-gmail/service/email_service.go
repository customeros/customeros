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
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/opentracing/opentracing-go/log"
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
	SyncEmailByEmailRawId(externalSystemId, tenant, usernameSource string, emailId uuid.UUID) error
	SyncEmailByMessageId(externalSystemId, tenant string, fromUsername, messageId string) error
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

		if emailForSync.MessageId == "" {
			err = s.repositories.RawEmailRepository.MarkSentToEventStore(emailForSync.ID, err == nil, "message id is empty")
			if err != nil {
				logrus.Errorf("unable to mark email as sent to event store: %v", err)
				return err
			}
		}

		err := s.syncEmail(externalSystemId, tenant, emailForSync.UsernameSource, emailForSync.ID)

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

func (s *emailService) SyncEmailByEmailRawId(externalSystemId, tenant, usernameSource string, emailId uuid.UUID) error {
	return s.syncEmail(externalSystemId, tenant, usernameSource, emailId)
}

func (s *emailService) SyncEmailByMessageId(externalSystemId, tenant, usernameSource, messageId string) error {
	rawEmail, err := s.repositories.RawEmailRepository.GetEmailForSyncByMessageId(externalSystemId, tenant, usernameSource, messageId)
	if err != nil {
		logrus.Errorf("failed to get emails for sync: %v", err)
		return err
	}

	if rawEmail == nil {
		return fmt.Errorf("email with message id %v not found", messageId)
	}

	if rawEmail.MessageId == "" {
		err = s.repositories.RawEmailRepository.MarkSentToEventStore(rawEmail.ID, err == nil, "message id is empty")
		if err != nil {
			logrus.Errorf("unable to mark email as sent to event store: %v", err)
			return err
		}
	}

	return s.syncEmail(externalSystemId, tenant, rawEmail.UsernameSource, rawEmail.ID)
}

func (s *emailService) syncEmail(externalSystemId, tenant, usernameSource string, emailId uuid.UUID) error {
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

		//lista de domenii care stim ca sunt personale: @gmail.com, @yahoo.com ...( asta e lista in baza de date ). se scot toate domeniile personale
		//nu creem domenii si organizatii pentru ele si contacte

		//adunam toate adresele de email din to, cc, bcc
		//check daca emailul exista deja in DB si este safe. daca exista, se importa tot automat cu logica de mai jos, fara verificare in reacher

		//ne scoatem pe ownerul emailului
		//verificam pentru fiecare adresa de email daca avem domeniu in system
		//daca nu avem domeniu, atunci il intrebam pe reacher daca emailul este safe

		//daca avem macar un email safe in lista de emailuri, se importa InteractionEvent + all emals + check each each email if it's an org + org creation + domain creation
		//open ai: care e numele companiei in baza acestui domeniu?

		//for personal emails, we don't create contacts, organizations and domains
		allEmailsString, err := s.buildEmailsListExcludingPersonalEmails(usernameSource, from, to, cc, bcc)
		if err != nil {
			logrus.Errorf("failed to build emails list: %v", err)
			return err
		}

		allEmailsNodes, err := s.repositories.EmailRepository.GetEmailsByRawEmail(ctx, tenant, allEmailsString)
		if err != nil {
			logrus.Errorf("failed to get emails by raw email: %v", err)
			return err
		}

		allEmailsEntities := make([]*entity.EmailEntity, len(allEmailsNodes))
		for i, emailNode := range allEmailsNodes {
			allEmailsEntities[i] = s.mapDbNodeToEmailEntity(*emailNode)
		}

		shouldAddInteractionEvent := false

		for _, emailEntity := range allEmailsEntities {
			if emailEntity.IsReachable != nil && *emailEntity.IsReachable == "safe" {
				shouldAddInteractionEvent = true
			}
		}

		if !shouldAddInteractionEvent {
			for _, emailString := range allEmailsString {
				emailValidated, err := s.services.ValidationApiService.ValidateEmail(emailString)
				if err != nil {
					logrus.Errorf("failed to validate email: %v", err)
					return err
				}

				if emailValidated.IsReachable == "safe" {
					shouldAddInteractionEvent = true
					break
				}
			}
		}

		if shouldAddInteractionEvent {
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

			//from
			//check if domain exists for tenant by email. if so, link the email to the user otherwise create a contact and link the email to the contact
			fromEmailId, err := s.getEmailIdForEmail(ctx, tx, tenant, interactionEventId, from, now)
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
				toEmailId, err := s.getEmailIdForEmail(ctx, tx, tenant, interactionEventId, toEmail, now)
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
				ccEmailId, err := s.getEmailIdForEmail(ctx, tx, tenant, interactionEventId, ccEmail, now)
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
				bccEmailId, err := s.getEmailIdForEmail(ctx, tx, tenant, interactionEventId, bccEmail, now)
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
	//		_, err = s.repositories.ActionItemRepository.CreateActionItemForEmail(ctx, tenant, interactionEventId, actionItem, externalSystemId, "sync-gmail", time.Now().UTC())
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

func (s *emailService) buildEmailsListExcludingPersonalEmails(usernameSource, from string, to []string, cc []string, bcc []string) ([]string, error) {
	personalEmailProviderList, err := s.repositories.PersonalEmailProviderRepository.GetPersonalEmailProviderList()
	if err != nil {
		logrus.Errorf("failed to get personal email provider list: %v", err)
		return nil, err
	}

	var allEmails []string

	if from != "" && from != usernameSource && !hasPersonalEmailProvider(personalEmailProviderList, extractDomain(from)) {
		allEmails = append(allEmails, from)
	}
	allEmails = append(allEmails, from)
	for _, email := range [][]string{to, cc, bcc} {
		for _, email := range email {
			if email != "" && email != usernameSource && !hasPersonalEmailProvider(personalEmailProviderList, extractDomain(email)) {
				allEmails = append(allEmails, email)
			}
		}
	}
	return allEmails, nil
}

func hasPersonalEmailProvider(providers []entity.PersonalEmailProvider, domain string) bool {
	for _, provider := range providers {
		if provider.ProviderDomain == domain {
			return true
		}
	}
	return false
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

func (s *emailService) getEmailIdForEmail(ctx context.Context, tx neo4j.ManagedTransaction, tenant string, interactionEventId, email string, now time.Time) (string, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "EmailService.getEmailIdForEmail")
	defer span.Finish()
	span.LogFields(log.String("tenant", tenant))
	span.LogFields(log.String("email", email))

	fromEmailId, err := s.repositories.EmailRepository.GetEmailId(ctx, tenant, email)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve email id for tenant: %v", err)
	}
	if fromEmailId != "" {
		return fromEmailId, nil
	}

	var domainNode *neo4j.Node
	var organizationNode *neo4j.Node
	var organizationId string

	domainNode, err = s.repositories.DomainRepository.GetDomainInTx(ctx, tx, extractDomain(email))
	if err != nil {
		return "", fmt.Errorf("unable to retrieve domain for tenant: %v", err)
	}

	if domainNode == nil {
		domainNode, err = s.repositories.DomainRepository.CreateDomain(ctx, tx, extractDomain(email), "gmail", "sync-gmail", now)
		if err != nil {
			return "", fmt.Errorf("unable to create domain: %v", err)
		}

		organizationName, err := s.services.OpenAiService.AskForOrganizationNameByDomain(tenant, interactionEventId, extractDomain(email))
		if err != nil {
			return "", fmt.Errorf("unable to retrieve organization name for tenant: %v", err)
		}
		if organizationName == "" {
			return "", fmt.Errorf("unable to retrieve organization name for tenant: %v", err)
		}

		organizationNode, err = s.repositories.OrganizationRepository.CreateOrganization(ctx, tx, tenant, organizationName, "gmail", "openline", "sync-gmail", now)
		if err != nil {
			return "", fmt.Errorf("unable to create organization for tenant: %v", err)
		}

		organizationId = utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*organizationNode), "id")
		domainName := utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*domainNode), "domain")
		err = s.repositories.OrganizationRepository.LinkDomainToOrganization(ctx, tx, tenant, domainName, organizationId)
		if err != nil {
			return "", fmt.Errorf("unable to link domain to organization: %v", err)
		}
	} else {

		organizationNode, err = s.repositories.OrganizationRepository.GetOrganizationWithDomain(ctx, tx, tenant, utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*domainNode), "domain"))
		if err != nil {
			return "", fmt.Errorf("unable to retrieve organization for tenant: %v", err)
		}

		if organizationNode == nil {
			organizationName, err := s.services.OpenAiService.AskForOrganizationNameByDomain(tenant, interactionEventId, extractDomain(email))
			if err != nil {
				return "", fmt.Errorf("unable to retrieve organization name for tenant: %v", err)
			}
			if organizationName == "" {
				return "", fmt.Errorf("unable to retrieve organization name for tenant: %v", err)
			}

			organizationNode, err = s.repositories.OrganizationRepository.CreateOrganization(ctx, tx, tenant, organizationName, "gmail", "openline", "sync-gmail", now)
			if err != nil {
				return "", fmt.Errorf("unable to create organization for tenant: %v", err)
			}

			organizationId = utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*organizationNode), "id")
			domainName := utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*domainNode), "domain")
			err = s.repositories.OrganizationRepository.LinkDomainToOrganization(ctx, tx, tenant, domainName, organizationId)
			if err != nil {
				return "", fmt.Errorf("unable to link domain to organization: %v", err)
			}
		}

	}

	emailId, err := s.repositories.EmailRepository.CreateEmailLinkedToOrganization(ctx, tx, tenant, email, "gmail", "openline", "sync-gmail", organizationId, now)
	if err != nil {
		return "", fmt.Errorf("unable to create email linked to organization: %v", err)
	}

	if emailId == "" {
		return "", fmt.Errorf("unable to create email linked to organization: %v", err)
	}

	return emailId, nil
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
