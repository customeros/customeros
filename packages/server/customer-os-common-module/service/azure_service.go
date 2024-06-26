package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/DusanKasan/parsemail"
	mimemail "github.com/emersion/go-message/mail"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/dto"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
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
	repositories *postgresRepository.Repositories
	services     *Services
}

type AzureService interface {
	ReadEmailsFromAzureAd(ctx context.Context, importState *postgresEntity.UserEmailImportState) ([]*postgresEntity.EmailRawData, string, error)

	SendEmail(ctx context.Context, tenant string, request dto.MailRequest) (*parsemail.Email, error)
}

func (s *azureService) ReadEmailsFromAzureAd(ctx context.Context, importState *postgresEntity.UserEmailImportState) ([]*postgresEntity.EmailRawData, string, error) {
	var results []*postgresEntity.EmailRawData

	var reqUrl string

	oAuthTokenEntity, err := s.repositories.OAuthTokenRepository.GetByEmail(ctx, importState.Tenant, "azure-ad", importState.Username)
	if err != nil {
		return nil, "", fmt.Errorf("unable to get token for user: %v", err)
	}

	if importState.Cursor != "" {
		reqUrl = importState.Cursor
	} else {
		reqUrl = "https://graph.microsoft.com/v1.0/me/messages"

		queryParams := url.Values{}
		queryParams.Add("$top", fmt.Sprintf("%d", 1))

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
		err := s.repositories.OAuthTokenRepository.MarkForManualRefresh(ctx, oAuthTokenEntity.TenantName, oAuthTokenEntity.PlayerIdentityId, oAuthTokenEntity.Provider)
		if err != nil {
			log.Fatalf("Failed to mark token for manual refresh: %v", err)
		}

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

func (s *azureService) SendEmail(ctx context.Context, tenant string, request dto.MailRequest) (*parsemail.Email, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "AzureService.SendEmail")
	defer span.Finish()

	var err error

	oAuthTokenEntity, err := s.repositories.OAuthTokenRepository.GetByEmail(ctx, tenant, "azure-ad", request.From)
	if err != nil {
		return nil, fmt.Errorf("unable to get token for user: %v", err)
	}
	if oAuthTokenEntity == nil {
		return nil, fmt.Errorf("unable to get token for user: %v", err)
	}
	if oAuthTokenEntity.NeedsManualRefresh {
		return nil, fmt.Errorf("oauth token needs manual refresh: %v", err)
	}

	//build internal object for transfer
	retMail := parsemail.Email{}
	retMail.HTMLBody = request.Content
	retMail.Subject = *request.Subject

	sentByName := "" //TODO

	fromAddress := []*mimemail.Address{{sentByName, request.From}}
	retMail.From = fromAddress
	var toAddress []*mimemail.Address
	var ccAddress []*mimemail.Address
	var bccAddress []*mimemail.Address
	for _, to := range request.To {
		toAddress = append(toAddress, &mimemail.Address{Address: to})
		retMail.To = toAddress
	}
	if request.Cc != nil {
		for _, cc := range request.Cc {
			ccAddress = append(ccAddress, &mimemail.Address{Address: cc})
			retMail.Cc = ccAddress
		}
	}
	if request.Bcc != nil {
		for _, bcc := range request.Bcc {
			bccAddress = append(bccAddress, &mimemail.Address{Address: bcc})
			retMail.Bcc = bccAddress
		}
	}

	// build and send the email
	var message MailRequest
	message.Subject = *request.Subject
	message.Body.ContentType = "HTML"
	message.Body.Content = request.Content
	message.From.EmailAddress.Address = request.From

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
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}

		req, err = http.NewRequest("POST", "https://graph.microsoft.com/v1.0/me/messages", bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}
	} else {
		interactionEventNode, err := s.services.Neo4jRepositories.InteractionEventReadRepository.GetInteractionEvent(ctx, tenant, *request.ReplyTo)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		interactionEvent := neo4jmapper.MapDbNodeToInteractionEventEntity(interactionEventNode)

		emailChannelData := dto.EmailChannelData{}
		err = json.Unmarshal([]byte(interactionEvent.ChannelData), &emailChannelData)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, fmt.Errorf("unable to parse email channel data for %s", *request.ReplyTo)
		}

		replyReq := ReplyRequest{
			Message: message,
		}

		reqBody, err := json.Marshal(replyReq)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}

		req, err = http.NewRequest("POST", fmt.Sprintf("https://graph.microsoft.com/v1.0/me/messages/%s/createReply", emailChannelData.ProviderMessageId), bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}
	}

	// Set the headers
	req.Header.Set("Authorization", "Bearer "+oAuthTokenEntity.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status: %s, response: %s", resp.Status, bodyBytes)
	} else {

		var draftResponse DraftResponse
		if err := json.NewDecoder(resp.Body).Decode(&draftResponse); err != nil {
			return nil, fmt.Errorf("failed to decode response body: %v", err)
		}

		sendDraftReq := SendDraftRequest{
			SaveToSentItems: true,
		}

		reqBody, err := json.Marshal(sendDraftReq)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %v", err)
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("https://graph.microsoft.com/v1.0/me/messages/%s/send", draftResponse.Id), bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+oAuthTokenEntity.AccessToken)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusAccepted {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			return nil, fmt.Errorf("API request failed with status: %s, response: %s", resp.Status, bodyBytes)
		}

		messageId, threadId, err := getIdAndThreadIdForSentEmail(oAuthTokenEntity.AccessToken, draftResponse.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to get conversation id for sent email: %v", err)
		}

		//this is used to store the thread id
		retMail.Header = map[string][]string{
			"Message-Id": {messageId},
			"Thread-Id":  {threadId},
		}
	}

	return &retMail, nil
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

func NewAzureService(repositories *postgresRepository.Repositories, services *Services) AzureService {
	return &azureService{
		repositories: repositories,
		services:     services,
	}
}
