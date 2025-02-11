package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	neoEntity "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/entity"
	"github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/mapper"
	neo4j_repository "github.com/customeros/customeros/packages/server/customer-os-neo4j-repository/repository"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/model"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

const (
	graphAPIBaseURL  = "https://graph.microsoft.com/v1.0/me"
	tokenEndpoint    = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	maxEmailsPerPage = 100
	htmlContentType  = "html"
	textContentType  = "text"
)

type azureService struct {
	cfg        *config.AzureOAuthConfig
	postgres   *postgres_repository.Repositories
	neo4j      *neo4j_repository.Repositories
	httpClient *http.Client
}

type azureApiCall struct {
	method         string
	url            string
	token          string
	body           []byte
	expectedStatus int
	response       interface{}
}

func NewAzureService(cfg *config.AzureOAuthConfig, postgres *postgres_repository.Repositories, neo4j *neo4j_repository.Repositories) interfaces.AzureService {
	return &azureService{
		cfg:        cfg,
		postgres:   postgres,
		neo4j:      neo4j,
		httpClient: &http.Client{},
	}
}

func (s *azureService) ReadEmailsFromAzureAd(ctx context.Context, importState *postgres_entity.UserEmailImportState) ([]*postgres_entity.EmailRawData, string, error) {
	span, ctx := s.initializeTracing(ctx, "AzureService.ReadEmailsFromAzureAd")
	defer span.Finish()

	token, err := s.getValidToken(ctx, span, importState.Tenant, importState.Username)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "token error"))
		return nil, "", fmt.Errorf("token error: %w", err)
	}

	reqURL := s.buildEmailsRequestURL(importState.Cursor)
	emails, nextLink, err := s.fetchEmails(ctx, span, reqURL, token)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "fetch emails error"))
		return nil, "", fmt.Errorf("fetch error: %w", err)
	}

	emailsWithHeaders, err := s.getEmailHeaders(ctx, span, emails, token)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "fetch email headers error"))
		return nil, "", fmt.Errorf("fetch error: %w", err)
	}

	return emailsWithHeaders, nextLink, nil
}

func (s *azureService) SendEmail(ctx context.Context, request *postgres_entity.EmailMessage) error {
	span, ctx := s.initializeTracing(ctx, "AzureService.SendEmail")
	defer span.Finish()

	token, err := s.getValidToken(ctx, span, common.GetTenantFromContext(ctx), request.From)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "token error"))
		return fmt.Errorf("token error: %w", err)
	}

	message := s.buildMailRequest(request)

	var draftID string
	if request.ReplyTo == nil {
		draftID, err = s.createNewEmail(ctx, span, message, token)
	} else {
		draftID, err = s.createReplyEmail(ctx, span, message, *request.ReplyTo, token)
	}
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error creating email"))
		return err
	}

	if err := s.sendDraft(ctx, span, draftID, token); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error sending draft"))
		return err
	}

	messageID, threadID, err := s.getEmailMetadata(token, span, draftID)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting email metadata"))
		return err
	}

	request.ProviderMessageId = messageID
	request.ProviderThreadId = threadID

	return nil
}

func (s *azureService) getEmailHeaders(
	ctx context.Context,
	span opentracing.Span,
	rawEmails []*postgres_entity.EmailRawData,
	token string,
) ([]*postgres_entity.EmailRawData, error) {
	for _, rawEmail := range rawEmails {
		url := fmt.Sprintf("%s/messages/%s/?$select=internetMessageHeaders",
			graphAPIBaseURL, rawEmail.ProviderMessageId)

		var response MicrosoftEmailHeaderResponse
		apiErr := s.makeRequest(ctx, span, azureApiCall{
			method:         "GET",
			url:            url,
			token:          token,
			expectedStatus: http.StatusOK,
			response:       &response,
		})
		if apiErr != nil {
			return nil, apiErr
		}

		rawEmail.Headers = response.InternetMessageHeaders
	}

	return rawEmails, nil
}

// makeRequest is a generic function to handle Azure HTTP requests
func (s *azureService) makeRequest(ctx context.Context, span opentracing.Span, opts azureApiCall) error {
	var reqBody io.Reader
	if opts.body != nil {
		reqBody = bytes.NewBuffer(opts.body)
	}

	req, err := http.NewRequestWithContext(ctx, opts.method, opts.url, reqBody)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed building new request"))
		return fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+opts.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed making request"))
		return fmt.Errorf("do request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != opts.expectedStatus {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("API error: %s, response: %s", resp.Status, body)
		tracing.TraceErr(span, err)
		return err
	}

	if opts.response != nil {
		if err := json.NewDecoder(resp.Body).Decode(opts.response); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "failed to decode json"))
			return fmt.Errorf("decode error: %w", err)
		}
	}

	return nil
}

func (s *azureService) fetchEmails(ctx context.Context, span opentracing.Span, reqURL string, token string) ([]*postgres_entity.EmailRawData, string, error) {
	var result MicrosoftRawEmailsResponse
	err := s.makeRequest(ctx, span, azureApiCall{
		method:         "GET",
		url:            reqURL,
		token:          token,
		expectedStatus: http.StatusOK,
		response:       &result,
	})
	if err != nil {
		return nil, "", err
	}
	return convertToEmailRawData(result), result.OdataNextLink, nil
}

func (s *azureService) getEmailMetadata(token string, span opentracing.Span, messageID string) (string, string, error) {
	var result MicrosoftRawEmailResponse
	url := fmt.Sprintf("%s/messages/%s", graphAPIBaseURL, messageID)

	err := s.makeRequest(context.Background(), span, azureApiCall{
		method:         "GET",
		url:            url,
		token:          token,
		expectedStatus: http.StatusOK,
		response:       &result,
	})
	if err != nil {
		return "", "", err
	}

	return result.Id, result.ConversationId, nil
}

func (s *azureService) createEmailRequest(ctx context.Context, span opentracing.Span, url string, body []byte, token string) (string, error) {
	var result DraftResponse
	err := s.makeRequest(ctx, span, azureApiCall{
		method:         "POST",
		url:            url,
		token:          token,
		body:           body,
		expectedStatus: http.StatusCreated,
		response:       &result,
	})
	if err != nil {
		return "", err
	}
	return result.Id, nil
}

func (s *azureService) sendDraft(ctx context.Context, span opentracing.Span, draftID, token string) error {
	sendReq := SendDraftRequest{SaveToSentItems: true}
	reqBody, err := json.Marshal(sendReq)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal json for sendDraft"))
		return fmt.Errorf("marshal error: %w", err)
	}

	return s.makeRequest(ctx, span, azureApiCall{
		method:         "POST",
		url:            fmt.Sprintf("%s/messages/%s/send", graphAPIBaseURL, draftID),
		token:          token,
		body:           reqBody,
		expectedStatus: http.StatusAccepted,
	})
}

func (s *azureService) getValidToken(ctx context.Context, span opentracing.Span, tenant, email string) (string, error) {
	token, err := s.postgres.OAuthTokenRepository.GetByEmail(ctx, tenant, enum.WorkspaceProviderAzure.String(), email)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error getting oauth token"))
		return "", fmt.Errorf("get token error: %w", err)
	}
	if token == nil || token.NeedsManualRefresh {
		tracing.TraceErr(span, errors.Wrap(err, "invalid token state"))
		return "", fmt.Errorf("invalid token state")
	}

	if token.ExpiresAt.Before(time.Now().Add(-time.Minute)) {
		refreshed, err := s.refreshToken(ctx, span, token)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error refreshing token"))
			return "", fmt.Errorf("refresh error: %w", err)
		}
		token = refreshed
	}

	return token.AccessToken, nil
}

func (s *azureService) refreshToken(ctx context.Context, span opentracing.Span, token *postgres_entity.OAuthTokenEntity) (*postgres_entity.OAuthTokenEntity, error) {
	data := url.Values{
		"client_id":     {s.cfg.ClientId},
		"client_secret": {s.cfg.ClientSecret},
		"refresh_token": {token.RefreshToken},
		"grant_type":    {"refresh_token"},
	}

	resp, err := s.httpClient.PostForm(tokenEndpoint, data)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "token refresh request error"))
		return nil, fmt.Errorf("refresh request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if err := s.postgres.OAuthTokenRepository.MarkForManualRefresh(ctx, token.TenantName, token.PlayerIdentityId, token.Provider); err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "error marking for manual refresh"))
			return nil, fmt.Errorf("mark for refresh error: %w", err)
		}
		tracing.TraceErr(span, fmt.Errorf("token refresh failed: %s", resp.Status))
		return nil, fmt.Errorf("refresh failed: %s", resp.Status)
	}

	var refreshResp RefreshTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&refreshResp); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "error decoding token refresh response"))
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return s.postgres.OAuthTokenRepository.Update(
		ctx,
		token.TenantName,
		token.PlayerIdentityId,
		token.Provider,
		refreshResp.AccessToken,
		refreshResp.RefreshToken,
		time.Now().Add(time.Second*time.Duration(refreshResp.ExpiresIn)),
	)
}

func (s *azureService) buildEmailsRequestURL(cursor string) string {
	if cursor != "" {
		return cursor
	}

	params := url.Values{}
	params.Add("$top", fmt.Sprintf("%d", maxEmailsPerPage))
	return fmt.Sprintf("%s/messages?%s", graphAPIBaseURL, params.Encode())
}

func (s *azureService) buildMailRequest(request *postgres_entity.EmailMessage) MailRequest {
	message := MailRequest{
		Subject: request.Subject,
		Body: struct {
			ContentType string `json:"contentType"`
			Content     string `json:"content"`
		}{
			ContentType: "HTML",
			Content:     request.Content,
		},
	}

	if request.FromName != "" {
		message.From.EmailAddress.Address = fmt.Sprintf("%s <%s>", request.FromName, request.From)
	} else {
		message.From.EmailAddress.Address = request.From
	}

	message.ToRecipients = buildRecipients(request.To)
	message.CcRecipients = buildRecipients(request.Cc)
	message.BccRecipients = buildRecipients(request.Bcc)

	return message
}

func (s *azureService) createNewEmail(ctx context.Context, span opentracing.Span, message MailRequest, token string) (string, error) {
	reqBody, err := json.Marshal(message)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal json for createNewEmail"))
		return "", fmt.Errorf("marshal error: %w", err)
	}

	return s.createEmailRequest(ctx, span, fmt.Sprintf("%s/messages", graphAPIBaseURL), reqBody, token)
}

func (s *azureService) createReplyEmail(ctx context.Context, span opentracing.Span, message MailRequest, replyToID string, token string) (string, error) {
	emailData, err := s.getReplyEmailData(ctx, span, replyToID)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get email reply data"))
		return "", err
	}

	replyReq := ReplyRequest{
		Message: message,
	}

	reqBody, err := json.Marshal(replyReq)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal json for createReplyEmail"))
		return "", fmt.Errorf("marshal error: %w", err)
	}

	return s.createEmailRequest(ctx, span, fmt.Sprintf("%s/messages/%s/createReply", graphAPIBaseURL, emailData.ProviderMessageId), reqBody, token)
}

func (s *azureService) getReplyEmailData(ctx context.Context, span opentracing.Span, replyToID string) (*neoEntity.EmailChannelData, error) {
	tenant := common.GetTenantFromContext(ctx)
	interactionEventNode, err := s.neo4j.CommonReadRepository.GetById(ctx, tenant, replyToID, model.NodeLabelInteractionEvent)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get interaction event"))
		return nil, fmt.Errorf("get interaction event error: %w", err)
	}

	interactionEvent := mapper.MapDbNodeToInteractionEventEntity(interactionEventNode)
	var emailData neoEntity.EmailChannelData
	if err := json.Unmarshal([]byte(interactionEvent.ChannelData), &emailData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to parse channel data"))
		return nil, fmt.Errorf("parse channel data error: %w", err)
	}

	return &emailData, nil
}

func buildRecipients(addresses []string) []Recipient {
	recipients := make([]Recipient, len(addresses))
	for i, addr := range addresses {
		recipients[i] = Recipient{
			EmailAddress: struct {
				Address string `json:"address"`
			}{
				Address: addr,
			},
		}
	}
	return recipients
}

func convertToEmailRawData(microsoftEmails MicrosoftRawEmailsResponse) []*postgres_entity.EmailRawData {
	emails := make([]*postgres_entity.EmailRawData, 0, len(microsoftEmails.Value))
	for _, me := range microsoftEmails.Value {
		email := &postgres_entity.EmailRawData{
			ProviderMessageId: me.Id,
			MessageId:         me.InternetMessageId,
			Sent:              me.SentDateTime,
			Subject:           me.Subject,
			From:              me.From.EmailAddress.Address,
			To:                concatenateEmailAddresses(me.ToRecipients),
			Cc:                concatenateEmailAddresses(me.CcRecipients),
			Bcc:               concatenateEmailAddresses(me.BccRecipients),
			Html:              getBodyContent(me.Body, htmlContentType),
			Text:              getBodyContent(me.Body, textContentType),
			ThreadId:          me.ConversationId,
			InReplyTo:         "",
			Reference:         "",
			Headers:           make(map[string]string),
		}
		emails = append(emails, email)
	}
	return emails
}

func concatenateEmailAddresses(recipients []struct {
	EmailAddress struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	} `json:"emailAddress"`
},
) string {
	addresses := make([]string, 0, len(recipients))
	for _, recipient := range recipients {
		address := fmt.Sprintf("%s <%s>", recipient.EmailAddress.Name, recipient.EmailAddress.Address)
		addresses = append(addresses, address)
	}
	return strings.Join(addresses, ", ")
}

func getBodyContent(body struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}, contentType string,
) string {
	if strings.EqualFold(body.ContentType, contentType) {
		return body.Content
	}
	return ""
}

func (s *azureService) initializeTracing(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
	tracing.SetDefaultServiceSpanTags(ctx, span)
	return span, ctx
}
