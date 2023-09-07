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
	authEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-auth/repository/postgres/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
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

	GetGmailServiceWithServiceAccount(username string, tenant string) (*gmail.Service, error)
	GetGmailServiceWithOauthToken(tokenEntity authEntity.OAuthTokenEntity) (*gmail.Service, error)

	FindEmailForUser(tenant, userId string) (*entity.EmailEntity, error)
	ReadEmailFromGoogle(gmailService *gmail.Service, userId, messageId string) (*EmailRawData, error)
	ReadNewEmailsForUsername(gmailService *gmail.Service, tenant, username string) error
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

func (s *emailService) ReadEmailFromGoogle(gmailService *gmail.Service, username, providerMessageId string) (*EmailRawData, error) {
	email, err := gmailService.Users.Messages.Get(username, providerMessageId).Format("full").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve email: %v", err)
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

	threadId := email.ThreadId
	references := ""
	inReplyTo := ""

	emailHeaders := make(map[string]string)

	for i := range email.Payload.Headers {
		headerName := strings.ToLower(email.Payload.Headers[i].Name)
		emailHeaders[email.Payload.Headers[i].Name] = email.Payload.Headers[i].Value
		if headerName == "message-id" {
			messageId = email.Payload.Headers[i].Value
		} else if headerName == "subject" {
			emailSubject = email.Payload.Headers[i].Value
			if emailSubject == "" {
				emailSubject = "No Subject"
			} else if strings.HasPrefix(emailSubject, "Re: ") {
				emailSubject = emailSubject[4:]
			}
		} else if headerName == "from" {
			from = email.Payload.Headers[i].Value
		} else if headerName == "to" {
			to = email.Payload.Headers[i].Value
		} else if headerName == "cc" {
			cc = email.Payload.Headers[i].Value
		} else if headerName == "bcc" {
			bcc = email.Payload.Headers[i].Value
		} else if headerName == "references" {
			references = email.Payload.Headers[i].Value
		} else if headerName == "in-reply-to" {
			inReplyTo = email.Payload.Headers[i].Value
		} else if headerName == "date" {
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
		} else if strings.HasPrefix(email.Payload.Parts[i].MimeType, "multipart") {
			for j := range email.Payload.Parts[i].Parts {
				if email.Payload.Parts[i].Parts[j].MimeType == "text/html" {
					emailHtmlBytes, _ := base64.URLEncoding.DecodeString(email.Payload.Parts[i].Parts[j].Body.Data)
					emailHtml = fmt.Sprintf("%s", emailHtmlBytes)
				} else if email.Payload.Parts[i].Parts[j].MimeType == "text/plain" {
					emailTextBytes, _ := base64.URLEncoding.DecodeString(email.Payload.Parts[i].Parts[j].Body.Data)
					emailText = fmt.Sprintf("%s", string(emailTextBytes))
				}
			}
		}
	}

	rawEmailData := &EmailRawData{
		ProviderMessageId: providerMessageId,
		MessageId:         messageId,
		Sent:              emailSentDate,
		Subject:           emailSubject,
		From:              from,
		To:                to,
		Cc:                cc,
		Bcc:               bcc,
		Html:              emailHtml,
		Text:              emailText,
		ThreadId:          threadId,
		InReplyTo:         inReplyTo,
		Reference:         references,
		Headers:           emailHeaders,
	}

	return rawEmailData, nil
}

func (s *emailService) ReadNewEmailsForUsername(gmailService *gmail.Service, tenant, username string) error {
	forUsername, err := s.repositories.UserGmailImportPageTokenRepository.GetGmailImportPageTokenForUsername(tenant, username)
	if err != nil {
		return fmt.Errorf("unable to retrieve history id for username: %v", err)
	}

	if forUsername == nil {
		emptyString := ""
		forUsername = &emptyString
	}

	countEmailsExists := int64(0)

	userMessages, err := gmailService.Users.Messages.List(username).MaxResults(s.cfg.SyncData.BatchSize).PageToken(*forUsername).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve emails for user: %v", err)
	}

	if userMessages != nil && len(userMessages.Messages) > 0 {
		for _, message := range userMessages.Messages {

			emailRawData, err := s.ReadEmailFromGoogle(gmailService, username, message.Id)
			if err != nil {
				return fmt.Errorf("unable to read email from google: %v", err)
			}

			if emailRawData.MessageId == "" {
				continue
			}

			emailExists, err := s.repositories.RawEmailRepository.EmailExistsByMessageId("gmail", tenant, username, emailRawData.MessageId)
			if err != nil {
				return fmt.Errorf("unable to check if email exists: %v", err)
			}

			//counting emails that are already imported based on the batch size
			//if the job is stopped in the middle of execution and we haven't saved the latest token
			//we are going to loose the history
			if emailExists {
				countEmailsExists = countEmailsExists + 1

				if countEmailsExists >= s.cfg.SyncData.BatchSize {
					err = s.repositories.UserGmailImportPageTokenRepository.UpdateGmailImportPageTokenForUsername(tenant, username, "")
					if err != nil {
						return fmt.Errorf("unable to update the gmail page token for username: %v", err)
					}
				}

				continue
			}

			jsonContent, err := JSONMarshal(emailRawData)
			if err != nil {
				return fmt.Errorf("failed to marshal email content: %v", err)
			}

			err = s.repositories.RawEmailRepository.Store("gmail", tenant, username, emailRawData.ProviderMessageId, emailRawData.MessageId, string(jsonContent))
			if err != nil {
				return fmt.Errorf("failed to store email content: %v", err)
			}
		}
	}

	err = s.repositories.UserGmailImportPageTokenRepository.UpdateGmailImportPageTokenForUsername(tenant, username, userMessages.NextPageToken)
	if err != nil {
		return fmt.Errorf("unable to update the gmail page token for username: %v", err)
	}

	return nil
}

func (s *emailService) GetGmailServiceWithServiceAccount(username string, tenant string) (*gmail.Service, error) {
	tok, err := s.getServiceAccountAuthToken(username, tenant)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve mail token for new gmail service: %v", err)
	}
	ctx := context.Background()
	client := tok.Client(ctx)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}

func (s *emailService) getServiceAccountAuthToken(identityId, tenant string) (*jwt.Config, error) {
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

func (s *emailService) GetGmailServiceWithOauthToken(tokenEntity authEntity.OAuthTokenEntity) (*gmail.Service, error) {
	oauth2Config := &oauth2.Config{
		ClientID:     s.cfg.GoogleOAuth.ClientId,
		ClientSecret: s.cfg.GoogleOAuth.ClientSecret,
		Endpoint:     google.Endpoint,
	}

	token := oauth2.Token{
		AccessToken:  tokenEntity.AccessToken,
		RefreshToken: tokenEntity.RefreshToken,
		Expiry:       tokenEntity.ExpiresAt,
	}

	tokenSource := oauth2Config.TokenSource(context.TODO(), &token)
	reuseTokenSource := oauth2.ReuseTokenSource(&token, tokenSource)

	if !token.Valid() {
		newToken, err := reuseTokenSource.Token()
		if err != nil {
			logrus.Errorf("failed to get new token: %v", err)
			return nil, err
		}

		if newToken.AccessToken != tokenEntity.AccessToken {

			_, err := s.repositories.OAuthRepositories.OAuthTokenRepository.Update(tokenEntity.PlayerIdentityId, tokenEntity.Provider, newToken.AccessToken, newToken.RefreshToken, newToken.Expiry)
			if err != nil {
				logrus.Errorf("failed to update token: %v", err)
				return nil, err
			}
		}

	}

	gmailService, err := gmail.NewService(context.TODO(), option.WithTokenSource(reuseTokenSource))
	if err != nil {
		logrus.Errorf("failed to create gmail service for token: %v", err)
		return nil, err
	}

	return gmailService, nil
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

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
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
