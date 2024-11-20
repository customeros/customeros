package service

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

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

const (
	graphAPIBaseURL  = "https://graph.microsoft.com/v1.0/me"
	tokenEndpoint    = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	maxEmailsPerPage = 100
	htmlContentType  = "html"
	textContentType  = "text"
)

type AzureService interface {
	ReadEmailsFromAzureAd(ctx context.Context, importState *postgresEntity.UserEmailImportState) ([]*postgresEntity.EmailRawData, string, error)
	SendEmail(ctx context.Context, request *postgresEntity.EmailMessage) error
}

type azureService struct {
	cfg          *config.AzureOAuthConfig
	repositories *postgresRepository.Repositories
	services     *Services
	httpClient   *http.Client
}

func NewAzureService(cfg *config.AzureOAuthConfig, repositories *postgresRepository.Repositories, services *Services) AzureService {
	return &azureService{
		cfg:          cfg,
		repositories: repositories,
		services:     services,
		httpClient:   &http.Client{},
	}
}

func (s *azureService) ReadEmailsFromAzureAd(ctx context.Context, importState *postgresEntity.UserEmailImportState) ([]*postgresEntity.EmailRawData, string, error) {
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
		tracing.TraceErr(span, errors.Wrap(err, "fetch error"))
		return nil, "", fmt.Errorf("fetch error: %w", err)
	}

	return emails, nextLink, nil
}

func (s *azureService) SendEmail(ctx context.Context, request *postgresEntity.EmailMessage) error {
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

func (s *azureService) getValidToken(ctx context.Context, span opentracing.Span, tenant, email string) (string, error) {
	token, err := s.repositories.OAuthTokenRepository.GetByEmail(ctx, tenant, "azure-ad", email)
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

func (s *azureService) refreshToken(ctx context.Context, span opentracing.Span, token *postgresEntity.OAuthTokenEntity) (*postgresEntity.OAuthTokenEntity, error) {
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
		if err := s.repositories.OAuthTokenRepository.MarkForManualRefresh(ctx, token.TenantName, token.PlayerIdentityId, token.Provider); err != nil {
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

	return s.repositories.OAuthTokenRepository.Update(
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

func (s *azureService) fetchEmails(ctx context.Context, span opentracing.Span, reqURL string, token string) ([]*postgresEntity.EmailRawData, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed building new fetchEmails request"))
		return nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed making fetchEmails request"))
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed reading fetchEmails response"))
		return nil, "", err
	}

	if resp.StatusCode != http.StatusOK {
		tracing.TraceErr(span, fmt.Errorf("fetchEmails call failed: %s", resp.Status))
		return nil, "", fmt.Errorf("API error: %s", resp.Status)
	}

	var result MicrosoftRawEmailsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to unmarshal fetchEmail response"))
		return nil, "", err
	}

	return convertToEmailRawData(result), result.OdataNextLink, nil
}

func (s *azureService) buildMailRequest(request *postgresEntity.EmailMessage) MailRequest {
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

func (s *azureService) createEmailRequest(ctx context.Context, span opentracing.Span, url string, body []byte, token string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to createEmailRequest"))
		return "", fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to make createEmailRequest call"))
		return "", fmt.Errorf("do request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		tracing.TraceErr(span, fmt.Errorf("API error: %s, response: %s", resp.Status, body))
		return "", fmt.Errorf("API error: %s, response: %s", resp.Status, body)
	}

	var draftResponse DraftResponse
	if err := json.NewDecoder(resp.Body).Decode(&draftResponse); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode json response"))
		return "", fmt.Errorf("decode error: %w", err)
	}

	return draftResponse.Id, nil
}

func (s *azureService) sendDraft(ctx context.Context, span opentracing.Span, draftID, token string) error {
	sendReq := SendDraftRequest{SaveToSentItems: true}
	reqBody, err := json.Marshal(sendReq)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to marshal json for sendDraft"))
		return fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/messages/%s/send", graphAPIBaseURL, draftID), bytes.NewBuffer(reqBody))
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request for sendDraft"))
		return fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to send request for sendDraft"))
		return fmt.Errorf("do request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		body, _ := io.ReadAll(resp.Body)
		r := fmt.Errorf("API error: %s, response: %s", resp.Status, body)
		tracing.TraceErr(span, r)
		return r
	}

	return nil
}

func (s *azureService) getEmailMetadata(token string, span opentracing.Span, messageID string) (string, string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/messages/%s", graphAPIBaseURL, messageID), nil)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to create request for getEmailMetadata"))
		return "", "", fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to make request for getEmailMetadata"))
		return "", "", fmt.Errorf("do request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		r := fmt.Errorf("API error: %s, response: %s", resp.Status, body)
		tracing.TraceErr(span, r)
		return "", "", r
	}

	var emailResponse MicrosoftRawEmailResponse
	if err := json.NewDecoder(resp.Body).Decode(&emailResponse); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to decode json"))
		return "", "", fmt.Errorf("decode error: %w", err)
	}

	return emailResponse.Id, emailResponse.ConversationId, nil
}

func (s *azureService) getReplyEmailData(ctx context.Context, span opentracing.Span, replyToID string) (*entity.EmailChannelData, error) {
	tenant := common.GetTenantFromContext(ctx)
	interactionEventNode, err := s.services.Neo4jRepositories.CommonReadRepository.GetById(ctx, tenant, replyToID, commonModel.NodeLabelInteractionEvent)
	if err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "failed to get interaction event"))
		return nil, fmt.Errorf("get interaction event error: %w", err)
	}

	interactionEvent := neo4jmapper.MapDbNodeToInteractionEventEntity(interactionEventNode)
	var emailData entity.EmailChannelData
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

func convertToEmailRawData(microsoftEmails MicrosoftRawEmailsResponse) []*postgresEntity.EmailRawData {
	emails := make([]*postgresEntity.EmailRawData, 0, len(microsoftEmails.Value))
	for _, me := range microsoftEmails.Value {
		email := &postgresEntity.EmailRawData{
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
