package entity

type EmailMessageData struct {
	BaseData
	Html                string   `json:"html,omitempty"`
	Text                string   `json:"text,omitempty"`
	Subject             string   `json:"subject,omitempty"`
	ExternalContactsIds []string `json:"externalContactsIds,omitempty"`
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

func (m *EmailMessageData) Normalize() {
	m.SetTimes()
	m.ToEmail = FilterOutEmpty(m.ToEmail)
	m.CcEmail = FilterOutEmpty(m.CcEmail)
	m.BccEmail = FilterOutEmpty(m.BccEmail)
}

func FilterOutEmpty(emails []string) []string {
	filtered := make([]string, 0)
	for _, email := range emails {
		if email != "" {
			filtered = append(filtered, email)
		}
	}
	return filtered
}
