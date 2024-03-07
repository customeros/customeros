package model

type EmailData struct {
	BaseData
	Content     string         `json:"content,omitempty"`
	ContentType string         `json:"contentType,omitempty"`
	Subject     string         `json:"subject,omitempty"`
	SentBy      EmailAddress   `json:"sentBy,omitempty"`
	SentTo      []EmailAddress `json:"sentTo,omitempty"`
	Bcc         []EmailAddress `json:"bcc,omitempty"`
	Cc          []EmailAddress `json:"cc,omitempty"`
	InReplyTo   string         `json:"inReplyTo,omitempty"`
	ChannelData string         `json:"channelData,omitempty"`
	Hide        bool           `json:"hide,omitempty"`
	ThreadId    string         `json:"threadId,omitempty"`
}

type EmailAddress struct {
	Name    string `json:"name"`
	Address string `json:"email"`
}
