package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/araddon/dateparse"
	mimemail "github.com/emersion/go-message/mail"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"io"
	"net/url"
	"strings"
	"time"
)

type googleService struct {
	cfg                  *config.GoogleOAuthConfig
	postgresRepositories postgresRepository.Repositories
	services             *Services
}

type GoogleService interface {
	ServiceAccountCredentialsExistsForTenant(ctx context.Context, tenant string) (bool, error)

	GetGmailService(ctx context.Context, username, tenant string) (*gmail.Service, error)

	GetGmailServiceWithServiceAccount(ctx context.Context, username string, tenant string) (*gmail.Service, error)
	GetGCalServiceWithServiceAccount(ctx context.Context, username string, tenant string) (*calendar.Service, error)

	GetGmailServiceWithOauthToken(ctx context.Context, tokenEntity postgresEntity.OAuthTokenEntity) (*gmail.Service, error)
	GetGCalServiceWithOauthToken(ctx context.Context, tokenEntity postgresEntity.OAuthTokenEntity) (*calendar.Service, error)

	ReadEmails(ctx context.Context, batchSize int64, importState *postgresEntity.UserEmailImportState) ([]*postgresEntity.EmailRawData, string, error)

	SendEmail(ctx context.Context, request *postgresEntity.EmailMessage) error
}

func NewGoogleService(cfg *config.GoogleOAuthConfig, postgresRepositories *postgresRepository.Repositories, services *Services) GoogleService {
	return &googleService{
		cfg:                  cfg,
		postgresRepositories: *postgresRepositories,
		services:             services,
	}
}

func (s *googleService) ServiceAccountCredentialsExistsForTenant(ctx context.Context, tenant string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GoogleService.ServiceAccountCredentialsExistsForTenant")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)

	privateKey, err := s.postgresRepositories.GoogleServiceAccountKeyRepository.GetApiKeyByTenantService(ctx, tenant, postgresEntity.GSUITE_SERVICE_PRIVATE_KEY)
	if err != nil {
		return false, nil
	}

	serviceEmail, err := s.postgresRepositories.GoogleServiceAccountKeyRepository.GetApiKeyByTenantService(ctx, tenant, postgresEntity.GSUITE_SERVICE_EMAIL_ADDRESS)
	if err != nil {
		return false, nil
	}

	if privateKey == "" || serviceEmail == "" {
		return false, nil
	}

	return true, nil
}

func (s *googleService) GetGmailService(ctx context.Context, username, tenant string) (*gmail.Service, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GoogleService.GetGmailService")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("username", username))

	serviceAccountExistsForTenant, err := s.ServiceAccountCredentialsExistsForTenant(ctx, tenant)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if serviceAccountExistsForTenant {
		gmailService, err := s.GetGmailServiceWithServiceAccount(ctx, username, tenant)
		if err != nil {
			logrus.Errorf("failed to create gmail service: %v", err)
			return nil, err
		}

		return gmailService, nil
	} else {
		tokenEntity, err := s.postgresRepositories.OAuthTokenRepository.GetByEmail(ctx, tenant, "google", username)
		if err != nil {
			return nil, err
		}
		if tokenEntity != nil && tokenEntity.NeedsManualRefresh {
			return nil, nil
		} else if tokenEntity != nil {
			if tokenEntity.RefreshToken == "" {
				err := s.postgresRepositories.OAuthTokenRepository.MarkForManualRefresh(ctx, tokenEntity.TenantName, tokenEntity.PlayerIdentityId, tokenEntity.Provider)
				return nil, err
			} else {
				gmailService, err := s.GetGmailServiceWithOauthToken(ctx, *tokenEntity)
				if err != nil {
					logrus.Errorf("failed to create gmail service: %v", err)
					return nil, err
				}
				return gmailService, nil
			}
		} else {
			return nil, nil
		}
	}
}

func (s *googleService) GetGmailServiceWithServiceAccount(ctx context.Context, username, tenant string) (*gmail.Service, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GoogleService.GetGmailServiceWithServiceAccount")
	defer span.Finish()
	span.SetTag(tracing.SpanTagTenant, tenant)
	span.LogFields(log.String("username", username))

	tok, err := s.getGmailServiceAccountAuthToken(ctx, username, tenant)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve mail token for new gmail service: %v", err)
	}
	client := tok.Client(ctx)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}

func (s *googleService) getGmailServiceAccountAuthToken(ctx context.Context, identityId, tenant string) (*jwt.Config, error) {
	privateKey, err := s.postgresRepositories.GoogleServiceAccountKeyRepository.GetApiKeyByTenantService(ctx, tenant, postgresEntity.GSUITE_SERVICE_PRIVATE_KEY)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve private key for gmail service: %v", err)
	}

	serviceEmail, err := s.postgresRepositories.GoogleServiceAccountKeyRepository.GetApiKeyByTenantService(ctx, tenant, postgresEntity.GSUITE_SERVICE_EMAIL_ADDRESS)
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

func (s *googleService) GetGCalServiceWithServiceAccount(ctx context.Context, username, tenant string) (*calendar.Service, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GoogleService.GetGCalServiceWithServiceAccount")
	defer span.Finish()

	tok, err := s.getGCalServiceAccountAuthToken(ctx, username, tenant)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve mail token for new gmail service: %v", err)
	}
	client := tok.Client(ctx)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	return srv, err
}

func (s *googleService) getGCalServiceAccountAuthToken(ctx context.Context, identityId, tenant string) (*jwt.Config, error) {
	privateKey, err := s.postgresRepositories.GoogleServiceAccountKeyRepository.GetApiKeyByTenantService(ctx, tenant, postgresEntity.GSUITE_SERVICE_PRIVATE_KEY)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve private key for gmail service: %v", err)
	}

	serviceEmail, err := s.postgresRepositories.GoogleServiceAccountKeyRepository.GetApiKeyByTenantService(ctx, tenant, postgresEntity.GSUITE_SERVICE_EMAIL_ADDRESS)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve service email for gmail service: %v", err)
	}
	conf := &jwt.Config{
		Email:      serviceEmail,
		PrivateKey: []byte(privateKey),
		TokenURL:   google.JWTTokenURL,
		Scopes:     []string{"https://calendar.google.com/"},
		Subject:    identityId,
	}
	return conf, nil
}

func (s *googleService) GetGmailServiceWithOauthToken(ctx context.Context, tokenEntity postgresEntity.OAuthTokenEntity) (*gmail.Service, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GoogleService.GetGmailServiceWithOauthToken")
	defer span.Finish()

	oauth2Config := &oauth2.Config{
		ClientID:     s.cfg.ClientId,
		ClientSecret: s.cfg.ClientSecret,
		Endpoint:     google.Endpoint,
	}

	token := oauth2.Token{
		AccessToken:  tokenEntity.AccessToken,
		RefreshToken: tokenEntity.RefreshToken,
		Expiry:       tokenEntity.ExpiresAt,
	}

	tokenSource := oauth2Config.TokenSource(ctx, &token)
	reuseTokenSource := oauth2.ReuseTokenSource(&token, tokenSource)

	if !token.Valid() {
		newToken, err := reuseTokenSource.Token()
		if err != nil && err.(*oauth2.RetrieveError) != nil && err.(*oauth2.RetrieveError).ErrorCode == "invalid_grant" {
			err := s.postgresRepositories.OAuthTokenRepository.MarkForManualRefresh(ctx, tokenEntity.TenantName, tokenEntity.PlayerIdentityId, tokenEntity.Provider)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
			return nil, fmt.Errorf("token is invalid and marked for manual refresh")
		} else if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

		if newToken.AccessToken != tokenEntity.AccessToken {

			_, err := s.postgresRepositories.OAuthTokenRepository.Update(ctx, tokenEntity.TenantName, tokenEntity.PlayerIdentityId, tokenEntity.Provider, newToken.AccessToken, newToken.RefreshToken, newToken.Expiry)
			if err != nil {
				tracing.TraceErr(span, err)
				return nil, err
			}
		}

	}

	gmailService, err := gmail.NewService(ctx, option.WithTokenSource(reuseTokenSource))
	if err != nil && err.(*oauth2.RetrieveError) != nil && err.(*oauth2.RetrieveError).ErrorCode == "invalid_grant" {
		err := s.postgresRepositories.OAuthTokenRepository.MarkForManualRefresh(ctx, tokenEntity.TenantName, tokenEntity.PlayerIdentityId, tokenEntity.Provider)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		return nil, fmt.Errorf("token is invalid and marked for manual refresh")
	} else if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	//Request had invalid authentication credentials. Expected OAuth 2 access token, login cookie or other valid authentication credential.
	//See https://developers.google.com/identity/sign-in/web/devconsole-project.
	_, err2 := gmailService.Users.GetProfile("me").Do()
	if err2 != nil {
		var googleApiErr *googleapi.Error
		var urlErr *url.Error

		switch {
		case errors.As(err2, &googleApiErr) && googleApiErr.Code == 401:
			// Handle 401 Unauthorized error
			err3 := s.postgresRepositories.OAuthTokenRepository.MarkForManualRefresh(ctx, tokenEntity.TenantName, tokenEntity.PlayerIdentityId, tokenEntity.Provider)
			if err3 != nil {
				tracing.TraceErr(span, errors.Wrap(err3, "failed to mark token for manual refresh"))
				return nil, err3
			}
			return nil, fmt.Errorf("token is invalid and marked for manual refresh")

		case errors.As(err2, &urlErr):
			// Handle URL error (e.g., network issues)
			tracing.TraceErr(span, errors.Wrap(urlErr, "network error occurred"))
			return nil, fmt.Errorf("network error occurred: %w", urlErr)

		default:
			// Handle any other errors
			tracing.TraceErr(span, errors.Wrap(err2, "unexpected error occurred"))
			return nil, fmt.Errorf("unexpected error occurred: %w", err2)
		}
	}

	return gmailService, nil
}

func (s *googleService) GetGCalServiceWithOauthToken(ctx context.Context, tokenEntity postgresEntity.OAuthTokenEntity) (*calendar.Service, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "GoogleService.GetGCalServiceWithOauthToken")
	defer span.Finish()

	oauth2Config := &oauth2.Config{
		ClientID:     s.cfg.ClientId,
		ClientSecret: s.cfg.ClientSecret,
		Endpoint:     google.Endpoint,
	}

	token := oauth2.Token{
		AccessToken:  tokenEntity.AccessToken,
		RefreshToken: tokenEntity.RefreshToken,
		Expiry:       tokenEntity.ExpiresAt,
	}

	tokenSource := oauth2Config.TokenSource(ctx, &token)
	reuseTokenSource := oauth2.ReuseTokenSource(&token, tokenSource)

	if !token.Valid() {
		newToken, err := reuseTokenSource.Token()
		if err != nil && err.(*oauth2.RetrieveError) != nil && (err.(*oauth2.RetrieveError).ErrorCode == "invalid_grant" || err.(*oauth2.RetrieveError).ErrorCode == "unauthorized_client") {
			err := s.postgresRepositories.OAuthTokenRepository.MarkForManualRefresh(ctx, tokenEntity.TenantName, tokenEntity.PlayerIdentityId, tokenEntity.Provider)
			if err != nil {
				logrus.Errorf("failed to mark token for manual refresh: %v", err)
				return nil, err
			}
			return nil, fmt.Errorf("token is invalid and marked for manual refresh")
		} else if err != nil {
			logrus.Errorf("failed to get new token: %v", err)
			return nil, err
		}

		if newToken.AccessToken != tokenEntity.AccessToken {

			_, err := s.postgresRepositories.OAuthTokenRepository.Update(ctx, tokenEntity.TenantName, tokenEntity.PlayerIdentityId, tokenEntity.Provider, newToken.AccessToken, newToken.RefreshToken, newToken.Expiry)
			if err != nil {
				logrus.Errorf("failed to update token: %v", err)
				return nil, err
			}
		}

	}

	gCalService, err := calendar.NewService(ctx, option.WithTokenSource(reuseTokenSource))
	if err != nil && err.(*oauth2.RetrieveError) != nil && err.(*oauth2.RetrieveError).ErrorCode == "invalid_grant" {
		err := s.postgresRepositories.OAuthTokenRepository.MarkForManualRefresh(ctx, tokenEntity.TenantName, tokenEntity.PlayerIdentityId, tokenEntity.Provider)
		if err != nil {
			logrus.Errorf("failed to mark token for manual refresh: %v", err)
			return nil, err
		}
		return nil, fmt.Errorf("token is invalid and marked for manual refresh")
	} else if err != nil {
		logrus.Errorf("failed to create gmail service for token: %v", err)
		return nil, err
	}

	//Request had invalid authentication credentials. Expected OAuth 2 access token, login cookie or other valid authentication credential.
	//See https://developers.google.com/identity/sign-in/web/devconsole-project.
	_, err2 := gCalService.CalendarList.Get("primary").Do()
	if err2 != nil && err2.(*googleapi.Error) != nil && err2.(*googleapi.Error).Code == 401 {
		err3 := s.postgresRepositories.OAuthTokenRepository.MarkForManualRefresh(ctx, tokenEntity.TenantName, tokenEntity.PlayerIdentityId, tokenEntity.Provider)
		if err3 != nil {
			logrus.Errorf("failed to mark token for manual refresh: %v", err)
			return nil, err3
		}
		return nil, fmt.Errorf("token is invalid and marked for manual refresh")
	} else if err2 != nil {
		logrus.Errorf("failed to get new token: %v", err)
		return nil, err2
	}

	return gCalService, nil
}

func (s *googleService) ReadEmails(ctx context.Context, batchSize int64, importState *postgresEntity.UserEmailImportState) ([]*postgresEntity.EmailRawData, string, error) {
	var results []*postgresEntity.EmailRawData

	gmailService, err := s.GetGmailService(ctx, importState.Username, importState.Tenant)
	if err != nil {
		logrus.Errorf("failed to create gmail service: %v", err)
		return nil, "", fmt.Errorf("failed to create gmail service: %v", err)
	}

	userMessages, err := gmailService.Users.Messages.List(importState.Username).Q("in:anywhere -label:draft").MaxResults(batchSize).PageToken(importState.Cursor).Do()
	if err != nil {
		return nil, "", fmt.Errorf("unable to retrieve emails for user: %v", err)
	}

	if userMessages != nil && len(userMessages.Messages) > 0 {
		for _, message := range userMessages.Messages {

			emailRawData, err := s.ReadEmailFromGoogle(gmailService, importState.Username, message.Id)
			if err != nil {
				return nil, "", fmt.Errorf("unable to read email from google: %v", err)
			}

			if emailRawData.MessageId == "" {
				continue
			}

			results = append(results, emailRawData)
		}
	}

	return results, userMessages.NextPageToken, nil
}

func (s *googleService) ReadEmailFromGoogle(gmailService *gmail.Service, username, providerMessageId string) (*postgresEntity.EmailRawData, error) {
	email, err := gmailService.Users.Messages.Get(username, providerMessageId).Format("full").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve email: %v", err)
	}

	messageId := ""
	emailSubject := ""
	emailHtml := ""
	emailText := ""
	emailSentDate := time.Time{}

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
			emailSentDate, err = convertToUTC(email.Payload.Headers[i].Value)
		}
	}

	if email.Payload.Parts != nil && len(email.Payload.Parts) > 0 {
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
	} else if email.Payload.Body != nil && email.Payload.Body.Data != "" {
		n, err := base64.URLEncoding.DecodeString(email.Payload.Body.Data)
		if err != nil {
			return nil, fmt.Errorf("unable to decode email body: %v", err)
		}
		emailText = fmt.Sprintf("%s", n)
	}

	rawEmailData := &postgresEntity.EmailRawData{
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

func (s *googleService) SendEmail(ctx context.Context, request *postgresEntity.EmailMessage) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "GoogleService.SendEmail")
	defer span.Finish()

	tenant := common.GetTenantFromContext(ctx)

	gSrv, err := s.GetGmailService(ctx, request.From, tenant)
	if err != nil {
		tracing.TraceErr(span, err)
		return fmt.Errorf("unable to retrieve mail token for new gmail service: %v", err)
	}

	if gSrv == nil {
		tracing.TraceErr(span, errors.New("unable to build a gmail service with service account or auth token"))
		return err
	}

	fromAddress := []*mimemail.Address{{request.FromName, request.From}}
	var toAddress []*mimemail.Address
	var ccAddress []*mimemail.Address
	var bccAddress []*mimemail.Address
	for _, to := range request.To {
		toAddress = append(toAddress, &mimemail.Address{Address: to})
	}
	if request.Cc != nil {
		for _, cc := range request.Cc {
			ccAddress = append(ccAddress, &mimemail.Address{Address: cc})
		}
	}
	if request.Bcc != nil {
		for _, bcc := range request.Bcc {
			bccAddress = append(bccAddress, &mimemail.Address{Address: bcc})
		}
	}

	var b bytes.Buffer

	var h mimemail.Header

	h.SetDate(time.Now())
	h.SetAddressList("From", fromAddress)
	h.SetAddressList("To", toAddress)
	h.SetAddressList("Cc", ccAddress)
	h.SetAddressList("Bcc", bccAddress)

	if request.Subject != "" {
		h.SetSubject(request.Subject)
	}

	threadId := ""

	if request.ReplyTo != nil {
		span.LogFields(log.String("replyTo", *request.ReplyTo))
		interactionEventNode, err := s.services.Neo4jRepositories.CommonReadRepository.GetById(ctx, tenant, *request.ReplyTo, commonModel.NodeLabelInteractionEvent)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
		interactionEvent := neo4jmapper.MapDbNodeToInteractionEventEntity(interactionEventNode)

		emailChannelData := entity.EmailChannelData{}
		err = json.Unmarshal([]byte(interactionEvent.ChannelData), &emailChannelData)
		if err != nil {
			tracing.TraceErr(span, err)
			return fmt.Errorf("unable to parse email channel data for %s", *request.ReplyTo)
		}

		h.Set("References", emailChannelData.Reference+" "+interactionEvent.Identifier)
		h.Set("In-Reply-To", interactionEvent.Identifier)

		interactionSessionNode, err := s.services.Neo4jRepositories.InteractionSessionReadRepository.GetForInteractionEvent(ctx, tenant, *request.ReplyTo)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}

		if interactionSessionNode == nil {
			tracing.TraceErr(span, errors.New("interaction session not found"))
			return errors.New("interaction session not found")
		}

		interactionSession := neo4jmapper.MapDbNodeToInteractionSessionEntity(interactionSessionNode)

		if interactionSession != nil && interactionSession.Identifier != "" {
			span.LogFields(log.String("threadId", interactionSession.Identifier))
			threadId = interactionSession.Identifier
		}
	}

	// Create a new mail writer
	mw, err := mimemail.CreateWriter(&b, h)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	// Create a text part
	tw, err := mw.CreateInline()
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	var th mimemail.InlineHeader
	th.Set("Content-Type", "text/html")
	w, err := tw.CreatePart(th)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	io.WriteString(w, request.Content)
	w.Close()
	tw.Close()

	mw.Close()

	raw := base64.StdEncoding.EncodeToString(b.Bytes())
	msgToSend := &gmail.Message{
		Raw: raw,
	}

	if threadId != "" {
		msgToSend.ThreadId = threadId
	}

	result, err := gSrv.Users.Messages.Send("me", msgToSend).Do()
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	request.ProviderThreadId = result.ThreadId

	generatedMessage, err := gSrv.Users.Messages.Get("me", result.Id).Do()
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	for _, header := range generatedMessage.Payload.Headers {
		if strings.EqualFold(header.Name, "Message-ID") {
			request.ProviderMessageId = header.Value
			break
		}
		if strings.EqualFold(header.Name, "References") {
			request.ProviderReferences = header.Value
			break
		}
	}

	return nil
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
		parsedTime, err = dateparse.ParseAny(datetimeStr)
		if err != nil {
			return time.Time{}, fmt.Errorf("unable to parse datetime string: %s", datetimeStr)
		}
	}

	return parsedTime.UTC(), nil
}
