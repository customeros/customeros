package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/config"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/entity"
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-gmail-raw/repository"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"strings"
)

type emailService struct {
	cfg          *config.Config
	repositories *repository.Repositories
	services     *Services
}

type EmailService interface {
	ServiceAccountCredentialsExistsForTenant(tenant string) (bool, error)
	FindEmailForUser(tenant, userId string) (*entity.EmailEntity, error)
	ReadNewEmailsForUsername(tenant, username string) error
}

func (s *emailService) ServiceAccountCredentialsExistsForTenant(tenant string) (bool, error) {
	privateKey, err := s.repositories.ApiKeyRepository.GetApiKeyByTenantService(tenant, repository.GSUITE_SERVICE_PRIVATE_KEY)
	if err != nil {
		return false, nil
	}

	serviceEmail, err := s.repositories.ApiKeyRepository.GetApiKeyByTenantService(tenant, repository.GSUITE_SERVICE_EMAIL_ADDRESS)
	if err != nil {
		return false, nil
	}

	if privateKey == "" || serviceEmail == "" {
		return false, nil
	}

	return true, nil
}

func (s *emailService) FindEmailForUser(tenant, userId string) (*entity.EmailEntity, error) {
	ctx := context.Background()

	email, err := s.repositories.EmailRepository.FindUserByEmail(ctx, tenant, userId)
	if err != nil {
		return nil, fmt.Errorf("unable to find user by email: %v", err)
	}
	if email == nil {
		return nil, nil
	}

	return s.mapDbNodeToEmailEntity(*email), nil
}

func (s *emailService) ReadNewEmailsForUsername(tenant, username string) error {
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
		emailSentDate := ""

		from := ""
		to := ""
		cc := ""
		bcc := ""

		references := ""
		inReplyTo := ""

		emailHeaders := make(map[string]string)

		for i := range email.Payload.Headers {
			emailHeaders[email.Payload.Headers[i].Name] = email.Payload.Headers[i].Value
			if email.Payload.Headers[i].Name == "Message-ID" {
				messageId = email.Payload.Headers[i].Value
			} else if email.Payload.Headers[i].Name == "Subject" {
				emailSubject = email.Payload.Headers[i].Value
				if emailSubject == "" {
					emailSubject = "No Subject"
				} else if strings.HasPrefix(emailSubject, "Re: ") {
					emailSubject = emailSubject[4:]
				}
			} else if email.Payload.Headers[i].Name == "From" {
				from = email.Payload.Headers[i].Value
			} else if email.Payload.Headers[i].Name == "To" {
				to = email.Payload.Headers[i].Value
			} else if email.Payload.Headers[i].Name == "Cc" {
				cc = email.Payload.Headers[i].Value
			} else if email.Payload.Headers[i].Name == "Bcc" {
				bcc = email.Payload.Headers[i].Value
			} else if email.Payload.Headers[i].Name == "References" {
				references = email.Payload.Headers[i].Value
			} else if email.Payload.Headers[i].Name == "In-Reply-To" {
				inReplyTo = email.Payload.Headers[i].Value
			} else if email.Payload.Headers[i].Name == "Date" {
				emailSentDate = email.Payload.Headers[i].Value
			}
		}

		for i := range email.Payload.Parts {
			if email.Payload.Parts[i].MimeType == "text/html" {
				emailHtmlBytes, _ := base64.URLEncoding.DecodeString(email.Payload.Parts[i].Body.Data)
				emailHtml = fmt.Sprintf("%s", emailHtmlBytes)
			} else if email.Payload.Parts[i].MimeType == "text/plain" {
				emailTextBytes, _ := base64.URLEncoding.DecodeString(email.Payload.Parts[i].Body.Data)
				emailText = fmt.Sprintf("%s", string(emailTextBytes))
			}
		}

		rawEmailData := &EmailRawData{
			MessageId: messageId,
			Sent:      emailSentDate,
			Subject:   emailSubject,
			From:      from,
			To:        to,
			Cc:        cc,
			Bcc:       bcc,
			Html:      emailHtml,
			Text:      emailText,
			InReplyTo: inReplyTo,
			Reference: references,
			Headers:   emailHeaders,
		}

		jsonContent, err := JSONMarshal(rawEmailData)
		if err != nil {
			return fmt.Errorf("failed to marshal email content: %v", err)
		}

		err = s.repositories.RawEmailRepository.Store("gmail", tenant, username, messageId, string(jsonContent))
		if err != nil {
			return fmt.Errorf("failed to store email content: %v", err)
		}

	}

	err = s.repositories.UserGmailImportPageTokenRepository.UpdateGmailImportPageTokenForUsername(tenant, username, userMessages.NextPageToken)
	if err != nil {
		return fmt.Errorf("unable to update the gmail page token for username: %v", err)
	}

	return nil
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
	InReplyTo string            `json:"InReplyTo"`
	Reference string            `json:"Reference"`
	Headers   map[string]string `json:"Headers"`
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
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
