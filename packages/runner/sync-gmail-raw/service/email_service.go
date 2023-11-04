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
	"google.golang.org/api/gmail/v1"
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
	ReadEmailFromGoogle(gmailService *gmail.Service, userId, providerMessageId string) (*EmailRawData, error)
	ReadEmailsForUsername(gmailService *gmail.Service, gmailImportState *entity.UserGmailImportState) (*entity.UserGmailImportState, error)
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
			if err != nil {
				return nil, fmt.Errorf("unable to parse email sent date: %v", err)
			}
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

func (s *emailService) ReadEmailsForUsername(gmailService *gmail.Service, gmailImportState *entity.UserGmailImportState) (*entity.UserGmailImportState, error) {
	countEmailsExists := int64(0)

	userInboxMessages, err := gmailService.Users.Messages.List(gmailImportState.Username).Q("in:inbox").MaxResults(s.cfg.SyncData.BatchSize).PageToken(gmailImportState.Cursor).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve inbox emails for user: %v", err)
	}
	userSentMessages, err := gmailService.Users.Messages.List(gmailImportState.Username).Q("in:sent").MaxResults(s.cfg.SyncData.BatchSize).PageToken(gmailImportState.Cursor).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve sent emails for user: %v", err)
	}
	//create a slice that contains inbox and sent emails for user
	userMessages := append(userInboxMessages.Messages, userSentMessages.Messages...)

	if userMessages != nil && len(userMessages) > 0 {
		for _, message := range userMessages {
			emailRawData, err := s.ReadEmailFromGoogle(gmailService, gmailImportState.Username, message.Id)
			if err != nil {
				return nil, fmt.Errorf("unable to read email from google: %v", err)
			}

			if emailRawData.MessageId == "" {
				continue
			}

			emailExists, err := s.repositories.RawEmailRepository.EmailExistsByMessageId("gmail", gmailImportState.Tenant, gmailImportState.Username, emailRawData.MessageId)
			if err != nil {
				return nil, fmt.Errorf("unable to check if email exists: %v", err)
			}

			//counting emails that are already imported based on the batch size
			//if the job is stopped in the middle of execution and we haven't saved the latest token
			//we are going to loose the history
			if emailExists {

				if gmailImportState.State == entity.REAL_TIME {
					countEmailsExists = countEmailsExists + 1

					if countEmailsExists >= s.cfg.SyncData.BatchSize {
						gmailImportState, err = s.repositories.UserGmailImportPageTokenRepository.UpdateGmailImportState(gmailImportState.Tenant, gmailImportState.Username, gmailImportState.State, "")
						if err != nil {
							return nil, fmt.Errorf("unable to update the gmail page token for username: %v", err)
						}
						return gmailImportState, nil
					}
				}

				continue
			} else {

				if gmailImportState.StopDate != nil && emailRawData.Sent.Before(*gmailImportState.StopDate) {
					gmailImportState, err = s.repositories.UserGmailImportPageTokenRepository.UpdateGmailImportState(gmailImportState.Tenant, gmailImportState.Username, gmailImportState.State, "")
					if err != nil {
						return nil, fmt.Errorf("unable to update the gmail page token for username: %v", err)
					}
					return gmailImportState, nil
				}

			}

			jsonContent, err := JSONMarshal(emailRawData)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal email content: %v", err)
			}

			err = s.repositories.RawEmailRepository.Store("gmail", gmailImportState.Tenant, gmailImportState.Username, emailRawData.ProviderMessageId, emailRawData.MessageId, string(jsonContent), emailRawData.Sent, gmailImportState.State)
			if err != nil {
				return nil, fmt.Errorf("failed to store email content: %v", err)
			}
		}
	}
	//next page token value assigned from sent or inbox
	var nextPageToken string
	if userInboxMessages.NextPageToken == "" {
		nextPageToken = userSentMessages.NextPageToken
	} else {
		nextPageToken = userInboxMessages.NextPageToken
	}
	gmailImportState, err = s.repositories.UserGmailImportPageTokenRepository.UpdateGmailImportState(gmailImportState.Tenant, gmailImportState.Username, gmailImportState.State, nextPageToken)
	if err != nil {
		return nil, fmt.Errorf("unable to update the gmail page token for username: %v", err)
	}

	return gmailImportState, nil
}

type EmailRawData struct {
	ProviderMessageId string            `json:"ProviderMessageId"`
	MessageId         string            `json:"MessageId"`
	Sent              time.Time         `json:"Sent"`
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
		return time.Time{}, fmt.Errorf("unable to parse datetime string: %s", datetimeStr)
	}

	return parsedTime.UTC(), nil
}

func NewEmailService(cfg *config.Config, repositories *repository.Repositories, services *Services) EmailService {
	return &emailService{
		cfg:          cfg,
		repositories: repositories,
		services:     services,
	}
}
