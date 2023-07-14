package entity

import (
	"github.com/openline-ai/openline-customer-os/packages/runner/sync-customer-os-data/utils"
	common_utils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
)

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
		m.CreatedAt = common_utils.TimePtr((*m.CreatedAt).UTC())
	} else {
		m.CreatedAt = common_utils.TimePtr(common_utils.Now())
	}
	if m.UpdatedAt != nil {
		m.UpdatedAt = common_utils.TimePtr((*m.UpdatedAt).UTC())
	} else {
		m.UpdatedAt = common_utils.TimePtr(common_utils.Now())
	}
}

func (m *EmailMessageData) Normalize() {
	m.FormatTimes()
	m.ToEmail = utils.FilterEmpty(m.ToEmail)
	m.CcEmail = utils.FilterEmpty(m.CcEmail)
	m.BccEmail = utils.FilterEmpty(m.BccEmail)
}
