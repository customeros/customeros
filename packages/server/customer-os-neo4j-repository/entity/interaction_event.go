package entity

import (
	"encoding/json"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/model"
	"time"
)

type InteractionEventEntity struct {
	DataLoaderKey
	Id                           string
	CreatedAt                    time.Time
	UpdatedAt                    time.Time
	Content                      string
	ContentType                  string
	Channel                      string
	ChannelData                  string
	Identifier                   string
	CustomerOSInternalIdentifier string
	EventType                    string
	Hide                         bool
	Source                       DataSource
	SourceOfTruth                DataSource
	AppSource                    string
}

type InteractionEventEntities []InteractionEventEntity

func (InteractionEventEntity) IsTimelineEvent() {
}

func (InteractionEventEntity) TimelineEventLabel() string {
	return model.NodeLabelInteractionEvent
}

func (e *InteractionEventEntity) GetDataloaderKey() string {
	return e.DataloaderKey
}

func (e *InteractionEventEntity) SetDataloaderKey(key string) {
	e.DataloaderKey = key
}

type EmailChannelData struct {
	ProviderMessageId string `json:"providerMessageId"`
	ThreadId          string `json:"threadId"`
	Subject           string `json:"Subject"`
	InReplyTo         string `json:"InReplyTo"`
	Reference         string `json:"Reference"`
}

func BuildEmailChannelData(messageId, threadId, subject, inReplyTo, references string) (*string, error) {
	emailContent := EmailChannelData{
		ProviderMessageId: messageId,
		ThreadId:          threadId,
		Subject:           subject,
		InReplyTo:         inReplyTo,
		Reference:         references,
	}
	jsonContent, err := json.Marshal(emailContent)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal email content: %v", err)
	}
	jsonContentString := string(jsonContent)

	return &jsonContentString, nil
}
