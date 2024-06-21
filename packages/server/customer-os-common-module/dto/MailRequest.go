package dto

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

type MailRequest struct {
	From         string   `json:"from"`
	FromProvider string   `json:"fromProvider"`
	To           []string `json:"to"`
	Cc           []string `json:"cc"`
	Bcc          []string `json:"bcc"`
	Subject      *string  `json:"subject"`
	Content      string   `json:"content"`

	ReplyTo *string `json:"replyTo,omitempty"`

	UniqueInternalIdentifier *string
}

type EmailChannelData struct {
	ProviderMessageId string   `json:"providerMessageId"`
	ThreadId          string   `json:"threadId"`
	Subject           string   `json:"Subject"`
	InReplyTo         []string `json:"InReplyTo"`
	Reference         []string `json:"Reference"`
}

func BuildEmailChannelData(messageId, threadId, subject string, inReplyTo, references []string) (*string, error) {
	emailContent := EmailChannelData{
		ProviderMessageId: messageId,
		ThreadId:          threadId,
		Subject:           subject,
		InReplyTo:         utils.EnsureEmailRfcIds(inReplyTo),
		Reference:         utils.EnsureEmailRfcIds(references),
	}
	jsonContent, err := json.Marshal(emailContent)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal email content: %v", err)
	}
	jsonContentString := string(jsonContent)

	return &jsonContentString, nil
}
