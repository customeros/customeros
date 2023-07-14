package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

/*
{
  "html": "<html><body>Hello World!</body></html>",
  "text": "Hello World!",
  "subject": "Test Email",
  "contactsExternalIds": [
    "1234",
    "5678"
  ],
  "externalUserId": "user123",
  "messageId": "abc123",
  "threadId": "xyz456",
  "fromEmail": "john@email.com",
  "toEmail": [
    "jane@email.com",
    "bob@email.com"
  ],
  "ccEmail": [
    "anna@email.com"
  ],
  "bccEmail": [
    "mark@email.com"
  ],
  "direction": "inbound",
  "firstName": "John",
  "lastName": "Doe",

  "skip": false,
  "skipReason": "draft data",
  "id": "1234",
  "externalId": "abcd1234",
  "externalSystem": "HubSpot",
  "createdAt": "2022-02-28T19:52:05Z",
  "updatedAt": "2022-03-01T11:23:45Z",
  "syncId": "sync_1234"
}
*/

type EmailMessageData struct {
	BaseData
	Html                string   `json:"html,omitempty"`
	Text                string   `json:"text,omitempty"`
	Subject             string   `json:"subject,omitempty"`
	ContactsExternalIds []string `json:"contactsExternalIds,omitempty"`
	ExternalUserId      string   `json:"externalUserId,omitempty"`
	EmailMessageId      string   `json:"messageId,omitempty"`
	EmailThreadId       string   `json:"threadId,omitempty"`
	FromEmail           string   `json:"fromEmail,omitempty"`
	ToEmail             []string `json:"toEmail,omitempty"`
	CcEmail             []string `json:"ccEmail,omitempty"`
	BccEmail            []string `json:"bccEmail,omitempty"`
	Direction           string   `json:"direction,omitempty"`
	FromFirstName       string   `json:"firstName,omitempty"`
	FromLastName        string   `json:"lastName,omitempty"`
}

func (m *EmailMessageData) FormatTimes() {
	if m.CreatedAt != nil {
		m.CreatedAt = utils.TimePtr((*m.CreatedAt).UTC())
	} else {
		m.CreatedAt = utils.TimePtr(utils.Now())
	}
	if m.UpdatedAt != nil {
		m.UpdatedAt = utils.TimePtr((*m.UpdatedAt).UTC())
	} else {
		m.UpdatedAt = utils.TimePtr(utils.Now())
	}
}

func (m *EmailMessageData) Normalize() {
	m.FormatTimes()
	m.ToEmail = utils.FilterEmpty(m.ToEmail)
	m.CcEmail = utils.FilterEmpty(m.CcEmail)
	m.BccEmail = utils.FilterEmpty(m.BccEmail)
}
