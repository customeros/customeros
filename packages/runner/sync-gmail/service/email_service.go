package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail/repository"
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
	repositories *repository.Repositories
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

	userMessages, err := googleServer.Users.Messages.List(username).MaxResults(1).PageToken(*forUsername).Do()
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

		//vasi@openline.ai ( alex si edi si are si user )
		//match (t:Tenant)--(e:Workspace{email:{openline.ai})
		//check if email address is a user ( email has the Domain of the tenant with the domain name )
		//if it's a user, use the email id
		//if it's not a user, we need to create an email and a contact for the email address
		//grpc call to create an email with a contact ( will create 2 events in event store

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

		sessionId, err := s.repositories.InteractionEventRepository.MergeInteractionSession(ctx, tenant, time.Now().UTC(), emailForCustomerOS)
		if err != nil {
			logrus.Errorf("failed merge interaction session with external reference %v for tenant %v :%v", message.Id, tenant, err)
			return err
		}

		interactionEventId, err := s.repositories.InteractionEventRepository.MergeEmailInteractionEvent(ctx, tenant, time.Now().UTC(), emailForCustomerOS)
		if err != nil {
			logrus.Errorf("failed merge interaction event with external reference %v for tenant %v :%v", message.Id, tenant, err)
			return err
		}

		err = s.repositories.InteractionEventRepository.LinkInteractionEventToSession(ctx, tenant, interactionEventId, sessionId)
		if err != nil {
			logrus.Errorf("failed to associate interaction event to session %v for tenant %v :%v", message.Id, tenant, err)
			return err
		}

		//from
		//check if domain exists for tenant by email. if so, link the email to the user otherwise create a contact and link the email to the contact
		fromEmailId, err := s.getEmailIdForEmail(ctx, tenant, from)
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
			toEmailId, err := s.getEmailIdForEmail(ctx, tenant, toEmail)
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
			ccEmailId, err := s.getEmailIdForEmail(ctx, tenant, ccEmail)
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
			bccEmailId, err := s.getEmailIdForEmail(ctx, tenant, bccEmail)
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

	err = s.repositories.UserGmailImportPageTokenRepository.UpdateGmailImportPageTokenForUsername(tenant, username, userMessages.NextPageToken)
	if err != nil {
		return fmt.Errorf("unable to update the gmail page token for username: %v", err)
	}

	return nil
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
	return parts[1]
}

func (s *emailService) getEmailIdForEmail(ctx context.Context, tenant, email string) (string, error) {
	fromEmailId := ""
	fromDomain, err := s.repositories.WorkspaceRepository.GetWorkspaceForTenantByName(ctx, tenant, extractDomain(email))
	if err != nil {
		return "", fmt.Errorf("unable to retrieve workspace for tenant: %v", err)
	}

	if fromDomain == nil {
		fromEmailId, err = s.repositories.EmailRepository.GetEmailIdOrCreateContactByEmail(ctx, tenant, email, "dummy", "name", "syng-gmail")
		if err != nil {
			return "", fmt.Errorf("unable to retrieve email id for tenant: %v", err)
		}
	} else {
		fromEmailId, err = s.repositories.EmailRepository.GetEmailId(ctx, tenant, email)
		if err != nil {
			return "", fmt.Errorf("unable to retrieve email id for tenant: %v", err)
		}

		if fromEmailId == "" {
			//todo create user
		}
	}
	return fromEmailId, nil
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

func NewEmailService(repositories *repository.Repositories) EmailService {
	return &emailService{
		repositories: repositories,
	}
}
