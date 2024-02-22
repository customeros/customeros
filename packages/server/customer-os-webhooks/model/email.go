package model

type EmailData struct {
	BaseData
	Content     string `json:"content,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Subject     string `json:"subject,omitempty"`
	SentBy      string `json:"sentBy,omitempty"`
	SentTo      string `json:"sentTo,omitempty"`
	Bcc         string `json:"bcc,omitempty"`
	Cc          string `json:"cc,omitempty"`
	InReplyTo   string `json:"inReplyTo,omitempty"`
	ChannelData string `json:"channelData,omitempty"`
	Hide        bool   `json:"hide,omitempty"`
	ThreadId    string `json:"threadId,omitempty"`
}
