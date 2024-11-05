package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	commonModel "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/entity"
	neo4jmapper "github.com/openline-ai/openline-customer-os/packages/server/customer-os-neo4j-repository/mapper"
	postgresRepository "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"
	"github.com/opentracing/opentracing-go"

	postgresEntity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type azureService struct {
	cfg          *config.AzureOAuthConfig
	repositories *postgresRepository.Repositories
	services     *Services
}

type AzureService interface {
	ReadEmailsFromAzureAd(ctx context.Context, importState *postgresEntity.UserEmailImportState) ([]*postgresEntity.EmailRawData, string, error)

	SendEmail(ctx context.Context, request *postgresEntity.EmailMessage) error
}

func (s *azureService) ReadEmailsFromAzureAd(ctx context.Context, importState *postgresEntity.UserEmailImportState) ([]*postgresEntity.EmailRawData, string, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "AzureService.ReadEmailsFromAzureAd")
	defer span.Finish()

	var results []*postgresEntity.EmailRawData

	var reqUrl string

	oAuthTokenEntity, err := s.repositories.OAuthTokenRepository.GetByEmail(ctx, importState.Tenant, "azure-ad", importState.Username)
	if err != nil {
		return nil, "", fmt.Errorf("unable to get token for user: %v", err)
	}

	if oAuthTokenEntity.ExpiresAt.Before(time.Now().Add(-time.Minute)) {

		refreshTokenResponse, err := s.getNewAccessTokenWithARefreshToken(ctx, oAuthTokenEntity.RefreshToken)
		if err != nil || refreshTokenResponse == nil {
			err := s.repositories.OAuthTokenRepository.MarkForManualRefresh(ctx, oAuthTokenEntity.TenantName, oAuthTokenEntity.PlayerIdentityId, oAuthTokenEntity.Provider)
			if err != nil {
				log.Fatalf("Failed to mark token for manual refresh: %v", err)
				return nil, "", err
			}

			return nil, "", err
		}

		oAuthTokenEntity, err = s.repositories.OAuthTokenRepository.Update(ctx, oAuthTokenEntity.TenantName, oAuthTokenEntity.PlayerIdentityId, oAuthTokenEntity.Provider, refreshTokenResponse.AccessToken, refreshTokenResponse.RefreshToken, time.Now().Add(time.Second*time.Duration(refreshTokenResponse.ExpiresIn)))
		if err != nil {
			log.Fatalf("Failed to update token: %v", err)
			return nil, "", err
		}
	}

	if importState.Cursor != "" {
		reqUrl = importState.Cursor
	} else {
		reqUrl = "https://graph.microsoft.com/v1.0/me/messages"

		queryParams := url.Values{}
		queryParams.Add("$top", fmt.Sprintf("%d", 100))

		reqUrl = reqUrl + "?" + queryParams.Encode()
	}

	client := &http.Client{}

	req, _ := http.NewRequest("GET", reqUrl, nil)
	req.Header.Set("Authorization", "Bearer "+oAuthTokenEntity.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to fetch emails: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusUnauthorized {

		return nil, "", fmt.Errorf("failed to fetch emails: %v", resp.Status)
	} else if resp.StatusCode == http.StatusOK {
		var result MicrosoftRawEmailsResponse
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Fatalf("Failed to parse response: %v", err)
		}

		results = convertToEmailRawData(result)

		if nextLink := result.OdataNextLink; nextLink != "" {
			return results, nextLink, nil
		}

		return results, "", nil
	}

	return nil, "", fmt.Errorf("failed to fetch emails: %v", resp.Status)
}

func (s *azureService) SendEmail(ctx context.Context, request *postgresEntity.EmailMessage) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "AzureService.SendEmail")
	defer span.Finish()

	tenant := common.GetTenantFromContext(ctx)

	var err error

	oAuthTokenEntity, err := s.repositories.OAuthTokenRepository.GetByEmail(ctx, tenant, "azure-ad", request.From)
	if err != nil {
		return fmt.Errorf("unable to get token for user: %v", err)
	}
	if oAuthTokenEntity == nil {
		return fmt.Errorf("unable to get token for user: %v", err)
	}
	if oAuthTokenEntity.NeedsManualRefresh {
		return fmt.Errorf("oauth token needs manual refresh: %v", err)
	}

	// build and send the email
	var message MailRequest
	message.Subject = request.Subject
	message.Body.ContentType = "HTML"
	message.Body.Content = request.Content

	if request.FromName != "" {
		message.From.EmailAddress.Address = fmt.Sprintf("%s <%s>", request.FromName, request.From)
	} else {
		message.From.EmailAddress.Address = request.From
	}

	for _, to := range request.To {
		message.ToRecipients = append(message.ToRecipients, Recipient{
			EmailAddress: struct {
				Address string `json:"address"`
			}{Address: to},
		})
	}

	for _, cc := range request.Cc {
		message.CcRecipients = append(message.CcRecipients, Recipient{
			EmailAddress: struct {
				Address string `json:"address"`
			}{Address: cc},
		})
	}

	for _, bcc := range request.Bcc {
		message.BccRecipients = append(message.BccRecipients, Recipient{
			EmailAddress: struct {
				Address string `json:"address"`
			}{Address: bcc},
		})
	}

	var req *http.Request

	if request.ReplyTo == nil {
		reqBody, err := json.Marshal(message)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %v", err)
		}

		req, err = http.NewRequest("POST", "https://graph.microsoft.com/v1.0/me/messages", bytes.NewBuffer(reqBody))
		if err != nil {
			return fmt.Errorf("failed to create request: %v", err)
		}
	} else {
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

		replyReq := ReplyRequest{
			Message: message,
		}

		reqBody, err := json.Marshal(replyReq)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %v", err)
		}

		req, err = http.NewRequest("POST", fmt.Sprintf("https://graph.microsoft.com/v1.0/me/messages/%s/createReply", emailChannelData.ProviderMessageId), bytes.NewBuffer(reqBody))
		if err != nil {
			return fmt.Errorf("failed to create request: %v", err)
		}
	}

	// Set the headers
	req.Header.Set("Authorization", "Bearer "+oAuthTokenEntity.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status: %s, response: %s", resp.Status, bodyBytes)
	} else {

		var draftResponse DraftResponse
		if err := json.NewDecoder(resp.Body).Decode(&draftResponse); err != nil {
			return fmt.Errorf("failed to decode response body: %v", err)
		}

		sendDraftReq := SendDraftRequest{
			SaveToSentItems: true,
		}

		reqBody, err := json.Marshal(sendDraftReq)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %v", err)
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("https://graph.microsoft.com/v1.0/me/messages/%s/send", draftResponse.Id), bytes.NewBuffer(reqBody))
		if err != nil {
			return fmt.Errorf("failed to create request: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+oAuthTokenEntity.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusAccepted {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			return fmt.Errorf("API request failed with status: %s, response: %s", resp.Status, bodyBytes)
		}

		messageId, threadId, err := getIdAndThreadIdForSentEmail(oAuthTokenEntity.AccessToken, draftResponse.Id)
		if err != nil {
			return fmt.Errorf("failed to get conversation id for sent email: %v", err)
		}

		request.ProviderMessageId = messageId
		request.ProviderThreadId = threadId
		//todo do we need references?
	}

	return nil
}

func getIdAndThreadIdForSentEmail(token, messageId string) (string, string, error) {
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/messages/%s", messageId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return "", "", fmt.Errorf("API request failed with status: %s, response: %s", resp.Status, bodyBytes)
	}

	var emailResponse MicrosoftRawEmailResponse
	if err := json.NewDecoder(resp.Body).Decode(&emailResponse); err != nil {
		return "", "", fmt.Errorf("failed to decode response body: %v", err)
	}

	return emailResponse.Id, emailResponse.ConversationId, nil
}

type MicrosoftRawEmailsResponse struct {
	OdataNextLink string                      `json:"@odata.nextLink"`
	Value         []MicrosoftRawEmailResponse `json:"value"`
}

type MicrosoftRawEmailResponse struct {
	Id                      string    `json:"id"`
	SentDateTime            time.Time `json:"sentDateTime"`
	InternetMessageId       string    `json:"internetMessageId"`
	Subject                 string    `json:"subject"`
	ConversationId          string    `json:"conversationId"`
	ConversationIndex       string    `json:"conversationIndex"`
	InferenceClassification string    `json:"inferenceClassification"`
	Body                    struct {
		ContentType string `json:"contentType"`
		Content     string `json:"content"`
	} `json:"body"`
	Sender struct {
		EmailAddress struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"emailAddress"`
	} `json:"sender"`
	From struct {
		EmailAddress struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"emailAddress"`
	} `json:"from"`
	ToRecipients []struct {
		EmailAddress struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"emailAddress"`
	} `json:"toRecipients"`
	CcRecipients []struct {
		EmailAddress struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"emailAddress"`
	} `json:"ccRecipients"`
	BccRecipients []struct {
		EmailAddress struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"emailAddress"`
	} `json:"bccRecipients"`
	ReplyTo []interface{} `json:"replyTo"`
}

func convertToEmailRawData(microsoftEmails MicrosoftRawEmailsResponse) []*postgresEntity.EmailRawData {
	var emails []*postgresEntity.EmailRawData
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
			Html:              getBodyContent(me.Body, "html"),
			Text:              getBodyContent(me.Body, "text"),
			ThreadId:          me.ConversationId,
			InReplyTo:         "",                      // This field is not provided in the response
			Reference:         "",                      // This field is not provided in the response
			Headers:           make(map[string]string), // Assuming headers are not available in this response
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
}) string {
	var addresses []string
	for _, recipient := range recipients {
		address := fmt.Sprintf("%s <%s>", recipient.EmailAddress.Name, recipient.EmailAddress.Address)
		addresses = append(addresses, address)
	}
	return strings.Join(addresses, ", ")
}

func getBodyContent(body struct {
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}, contentType string) string {
	if body.ContentType == contentType {
		return body.Content
	}
	return ""
}

func (s *azureService) getNewAccessTokenWithARefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "AzureService.getNewAccessTokenWithARefreshToken")
	defer span.Finish()

	data := url.Values{}
	data.Set("client_id", s.cfg.ClientId)
	data.Set("client_secret", s.cfg.ClientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")

	req, err := http.NewRequest("POST", "https://login.microsoftonline.com/common/oauth2/v2.0/token", strings.NewReader(data.Encode()))
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		tracing.TraceErr(span, err)
		return nil, err
	}

	var tokenResponse RefreshTokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &tokenResponse, nil
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

type MailRequest struct {
	Subject string `json:"subject"`
	Body    struct {
		ContentType string `json:"contentType"`
		Content     string `json:"content"`
	} `json:"body"`
	From struct {
		EmailAddress struct {
			Address string `json:"address"`
		} `json:"emailAddress"`
	} `json:"from"`
	ToRecipients  []Recipient `json:"toRecipients"`
	CcRecipients  []Recipient `json:"ccRecipients,omitempty"`
	BccRecipients []Recipient `json:"bccRecipients,omitempty"`
}

type Recipient struct {
	EmailAddress struct {
		Address string `json:"address"`
	} `json:"emailAddress"`
}

type CreateDraftRequest struct {
	Message MailRequest `json:"message"`
}

type DraftResponse struct {
	Id string `json:"id"`
}

type SendDraftRequest struct {
	SaveToSentItems bool `json:"saveToSentItems"`
}

type ReplyRequest struct {
	Message MailRequest `json:"message"`
	Comment string      `json:"comment"`
}

func NewAzureService(cfg *config.AzureOAuthConfig, repositories *postgresRepository.Repositories, services *Services) AzureService {
	return &azureService{
		cfg:          cfg,
		repositories: repositories,
		services:     services,
	}
}
