package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
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
	ReadNewEmailsForUsername(tenant, username string) error
}

func (s *emailService) ReadNewEmailsForUsername(tenant, username string) error {
	ctx := context.Background()

	err := s.repositories.ExternalSystemRepository.Merge(ctx, tenant, "gmail")
	if err != nil {
		return fmt.Errorf("unable to merge external system: %v", err)
	}

	googleServer, err := s.newGmailService(username, tenant)
	if err != nil {
		return fmt.Errorf("unable to retrieve mail token for new gmail service: %v", err)
	}

	forUsername, err := s.repositories.UserGmailImportPageTokenRepository.GetGmailImportPageTokenForUsername(tenant, username)
	if err != nil {
		return fmt.Errorf("unable to retrieve history id for username: %v", err)
	}

	//empty cursor means all messages have been read already
	if forUsername != nil && *forUsername == "" {
		return nil
	} else if forUsername == nil {
		emptyString := ""
		forUsername = &emptyString
	}

	userMessages, err := googleServer.Users.Messages.List(username).MaxResults(s.cfg.SyncData.BatchSize).PageToken(*forUsername).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve emails for user: %v", err)
	}

	for _, message := range userMessages.Messages {
		email, err := googleServer.Users.Messages.Get(username, message.Id).Format("full").Do()
		if err != nil {
			return fmt.Errorf("unable to retrieve email: %v", err)
		}

		messageId := ""
		emailSubject := ""
		emailHtml := ""
		emailText := ""
		now := time.Now().UTC()

		from := ""
		to := make([]string, 0)
		cc := make([]string, 0)
		bcc := make([]string, 0)

		for i := range email.Payload.Headers {
			if email.Payload.Headers[i].Name == "Message-ID" {
				messageId = email.Payload.Headers[i].Value
			} else if email.Payload.Headers[i].Name == "Subject" {
				emailSubject = email.Payload.Headers[i].Value
			} else if email.Payload.Headers[i].Name == "From" {
				from = extractEmailAddresses(email.Payload.Headers[i].Value)[0]
			} else if email.Payload.Headers[i].Name == "To" {
				to = extractEmailAddresses(email.Payload.Headers[i].Value)
			} else if email.Payload.Headers[i].Name == "Cc" {
				cc = extractEmailAddresses(email.Payload.Headers[i].Value)
			} else if email.Payload.Headers[i].Name == "Bcc" {
				bcc = extractEmailAddresses(email.Payload.Headers[i].Value)
			}
		}

		for i := range email.Payload.Parts {
			if email.Payload.Parts[i].MimeType == "text/html" {
				emailHtmlBytes, _ := base64.URLEncoding.DecodeString(email.Payload.Parts[i].Body.Data)
				emailHtml = fmt.Sprintf("%s", emailHtmlBytes)
			}
		}

		emailForCustomerOS := entity.EmailMessageData{
			Html:           emailHtml,
			Text:           emailText,
			Subject:        emailSubject,
			CreatedAt:      now,
			ExternalId:     messageId,
			ExternalSystem: "gmail",
			EmailThreadId:  email.ThreadId,
			EmailMessageId: messageId,
		}

		interactionEventId, err := s.repositories.InteractionEventRepository.GetInteractionEventIdByExternalId(ctx, tenant, messageId)
		if err != nil {
			logrus.Errorf("failed to check if interaction event exists for external id %v for tenant %v :%v", messageId, tenant, err)
			return err
		}

		//TODO we need to check for each item inserted part of this email, not only for the email itself
		if interactionEventId == "" {

			sessionId, err := s.repositories.InteractionEventRepository.MergeInteractionSession(ctx, tenant, time.Now().UTC(), emailForCustomerOS)
			if err != nil {
				logrus.Errorf("failed merge interaction session with external reference %v for tenant %v :%v", message.Id, tenant, err)
				return err
			}

			interactionEventId, err = s.repositories.InteractionEventRepository.MergeEmailInteractionEvent(ctx, tenant, time.Now().UTC(), emailForCustomerOS)
			if err != nil {
				logrus.Errorf("failed merge interaction event with external reference %v for tenant %v :%v", message.Id, tenant, err)
				return err
			}

			err = s.repositories.InteractionEventRepository.LinkInteractionEventToSession(ctx, tenant, interactionEventId, sessionId)
			if err != nil {
				logrus.Errorf("failed to associate interaction event to session %v for tenant %v :%v", message.Id, tenant, err)
				return err
			}

			logrus.Println("fetching email classification")
			emailsClassification, err := s.services.OpenAiService.FetchEmailsClassification(from, to, cc, bcc)
			if err != nil {
				return fmt.Errorf("unable to fetch email classification: %v", err)
			}

			//from
			//check if domain exists for tenant by email. if so, link the email to the user otherwise create a contact and link the email to the contact
			emailClassification := getEmailClassificationByEmail(from, emailsClassification)
			fromEmailId, err := s.getEmailIdForEmail(ctx, tenant, emailClassification, now)
			if err != nil {
				return fmt.Errorf("unable to retrieve email id for tenant: %v", err)
			}
			if fromEmailId == "" {
				return fmt.Errorf("unable to retrieve email id for tenant %s and email %s", tenant, from)
			}

			err = s.repositories.InteractionEventRepository.InteractionEventSentByEmail(ctx, tenant, interactionEventId, fromEmailId)
			if err != nil {
				return fmt.Errorf("unable to link email to interaction event: %v", err)
			}

			//to
			for _, toEmail := range to {
				emailClassification := getEmailClassificationByEmail(toEmail, emailsClassification)
				toEmailId, err := s.getEmailIdForEmail(ctx, tenant, emailClassification, now)
				if err != nil {
					return fmt.Errorf("unable to retrieve email id for tenant: %v", err)
				}
				if toEmailId == "" {
					return fmt.Errorf("unable to retrieve email id for tenant %s and email %s", tenant, toEmail)
				}

				err = s.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tenant, interactionEventId, "TO", []string{toEmailId})
				if err != nil {
					return fmt.Errorf("unable to link email to interaction event: %v", err)
				}
			}

			//cc
			for _, ccEmail := range cc {
				emailClassification := getEmailClassificationByEmail(ccEmail, emailsClassification)
				ccEmailId, err := s.getEmailIdForEmail(ctx, tenant, emailClassification, now)
				if err != nil {
					return fmt.Errorf("unable to retrieve email id for tenant: %v", err)
				}
				if ccEmailId == "" {
					return fmt.Errorf("unable to retrieve email id for tenant %s and email %s", tenant, ccEmail)
				}

				err = s.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tenant, interactionEventId, "CC", []string{ccEmailId})
				if err != nil {
					return fmt.Errorf("unable to link email to interaction event: %v", err)
				}
			}

			//bcc
			for _, bccEmail := range bcc {
				emailClassification := getEmailClassificationByEmail(bccEmail, emailsClassification)
				bccEmailId, err := s.getEmailIdForEmail(ctx, tenant, emailClassification, now)
				if err != nil {
					return fmt.Errorf("unable to retrieve email id for tenant: %v", err)
				}
				if bccEmailId == "" {
					return fmt.Errorf("unable to retrieve email id for tenant %s and email %s", tenant, bccEmail)
				}

				err = s.repositories.InteractionEventRepository.InteractionEventSentToEmails(ctx, tenant, interactionEventId, "BCC", []string{bccEmailId})
				if err != nil {
					return fmt.Errorf("unable to link email to interaction event: %v", err)
				}
			}

		}

		//get summary for email using claude-2 model
		summaryExists, err := s.repositories.AnalysisRepository.SummaryExistsForInteractionEvent(ctx, tenant, interactionEventId)
		if err != nil {
			return fmt.Errorf("unable to check if summary exists for interaction event: %v", err)
		}
		if !summaryExists {
			logrus.Println("fetching anthropic summary for email")
			summary := s.services.AnthropicService.FetchSummary(emailHtml)
			_, err = s.repositories.AnalysisRepository.CreateSummaryForEmail(ctx, tenant, interactionEventId, summary, "gmail", "sync-gmail", now)
			if err != nil {
				return fmt.Errorf("unable to create summary for email: %v", err)
			}
		}

		//get the action items for the email using claude-2 model
		actionItemsExists, err := s.repositories.ActionItemRepository.ActionsItemsExistsForInteractionEvent(ctx, tenant, interactionEventId)
		if err != nil {
			return fmt.Errorf("unable to check if action items exists for interaction event: %v", err)
		}
		if !actionItemsExists {
			logrus.Println("fetching anthropic action items for email")
			actionItems := s.services.AnthropicService.FetchActionItems(emailHtml)

			//TODO insert should be done in a single transaction to follow ActionsItemsExistsForInteractionEvent logic
			for _, actionItem := range actionItems {
				_, err = s.repositories.ActionItemRepository.CreateActionItemForEmail(ctx, tenant, interactionEventId, actionItem, "gmail", "sync-gmail", now)
				if err != nil {
					return fmt.Errorf("unable to create action item for email: %v", err)
				}
			}
		}

		fmt.Println("email processed: " + interactionEventId)
	}

	err = s.repositories.UserGmailImportPageTokenRepository.UpdateGmailImportPageTokenForUsername(tenant, username, userMessages.NextPageToken)
	if err != nil {
		return fmt.Errorf("unable to update the gmail page token for username: %v", err)
	}

	return nil
}

func getEmailClassificationByEmail(email string, emailClassifications []*OpenAiEmailClassification) OpenAiEmailClassification {
	for _, emailClassification := range emailClassifications {
		if emailClassification.Email == email {
			return *emailClassification
		}
	}
	return OpenAiEmailClassification{}
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
			emailAddresses = append(emailAddresses, email)
		}

		return emailAddresses
	}

	// If no angle brackets found, assume the input is a single email address
	return []string{input}
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
	return split[len(split)-2] + "." + split[len(split)-1]
}

func (s *emailService) getEmailIdForEmail(ctx context.Context, tenant string, emailClassification OpenAiEmailClassification, now time.Time) (string, error) {
	fromEmailId, err := s.repositories.EmailRepository.GetEmailId(ctx, tenant, emailClassification.Email)
	if err != nil {
		return "", fmt.Errorf("unable to retrieve email id for tenant: %v", err)
	}
	if fromEmailId != "" {
		return fromEmailId, nil
	}

	if emailClassification.IsPersonalEmail {
		emailId, err := s.repositories.EmailRepository.CreateContactWithEmail(ctx, tenant, emailClassification.Email, emailClassification.PersonalFirstName, emailClassification.PersonalLastName, "syng-gmail")
		if err != nil {
			return "", fmt.Errorf("unable to create contact with email: %v", err)
		}
		return emailId, nil
	} else if emailClassification.IsOrganizationEmail {
		domainNode, err := s.repositories.DomainRepository.GetDomain(ctx, extractDomain(emailClassification.Email))
		if err != nil {
			return "", fmt.Errorf("unable to retrieve domain for tenant: %v", err)
		}

		if domainNode != nil {
			organizationNode, err := s.repositories.OrganizationRepository.GetOrganizationWithDomain(ctx, tenant, utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*domainNode), "id"))
			if err != nil {
				return "", fmt.Errorf("unable to retrieve organization for tenant: %v", err)
			}

			if organizationNode == nil {
				organizationNode, err = s.repositories.OrganizationRepository.CreateOrganization(ctx, tenant, emailClassification.OrganizationName, "gmail", "openline", "sync-gmail", now)
				if err != nil {
					return "", fmt.Errorf("unable to create organization for tenant: %v", err)
				}
			}

			emailId, err := s.repositories.EmailRepository.CreateEmailLinkedToOrganization(ctx, tenant, emailClassification.Email, "gmail", "openline", "sync-gmail", utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*organizationNode), "id"), now)
			if err != nil {
				return "", fmt.Errorf("unable to create email linked to organization: %v", err)
			}

			return emailId, nil
		} else {
			domainNode, err := s.repositories.DomainRepository.CreateDomain(ctx, extractDomain(emailClassification.Email), "gmail", "sync-gmail", now)
			if err != nil {
				return "", fmt.Errorf("unable to create domain: %v", err)
			}
			organizationNode, err := s.repositories.OrganizationRepository.CreateOrganization(ctx, tenant, emailClassification.OrganizationName, "gmail", "openline", "sync-gmail", now)
			if err != nil {
				return "", fmt.Errorf("unable to create organization for tenant: %v", err)
			}
			domainId := utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*domainNode), "id")
			organizationId := utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*organizationNode), "id")
			err = s.repositories.OrganizationRepository.LinkDomainToOrganization(ctx, tenant, domainId, organizationId)
			if err != nil {
				return "", fmt.Errorf("unable to link domain to organization: %v", err)
			}

			emailId, err := s.repositories.EmailRepository.CreateEmailLinkedToOrganization(ctx, tenant, emailClassification.Email, "gmail", "openline", "sync-gmail", utils.GetStringPropOrEmpty(utils.GetPropsFromNode(*organizationNode), "id"), now)
			if err != nil {
				return "", fmt.Errorf("unable to create email linked to organization: %v", err)
			}

			return emailId, nil
		}

	} else {
		return "", fmt.Errorf("unable to determine email type: %v", err)
	}
}

func (s *emailService) newGmailService(username string, tenant string) (*gmail.Service, error) {
	tok, err := s.getMailAuthToken(username, tenant)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve mail token for new gmail service: %v", err)
	}
	ctx := context.Background()
	client := tok.Client(ctx)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}

func (s *emailService) getMailAuthToken(identityId, tenant string) (*jwt.Config, error) {
	privateKey, err := s.repositories.ApiKeyRepository.GetApiKeyByTenantService(tenant, repository.GSUITE_SERVICE_PRIVATE_KEY)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve private key for gmail service: %v", err)
	}

	serviceEmail, err := s.repositories.ApiKeyRepository.GetApiKeyByTenantService(tenant, repository.GSUITE_SERVICE_EMAIL_ADDRESS)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve service email for gmail service: %v", err)
	}
	conf := &jwt.Config{
		Email:      serviceEmail,
		PrivateKey: []byte(privateKey),
		TokenURL:   google.JWTTokenURL,
		Scopes:     []string{"https://mail.google.com/"},
		Subject:    identityId,
	}
	return conf, nil
}

func NewEmailService(cfg *config.Config, repositories *repository.Repositories, services *Services) EmailService {
	return &emailService{
		cfg:          cfg,
		repositories: repositories,
		services:     services,
	}
}
