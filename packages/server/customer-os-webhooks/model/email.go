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

type PostmarkEmailWebhookData struct {
	FromName      string `json:"FromName"`
	MessageStream string `json:"MessageStream"`
	FromFull      struct {
		Email       string `json:"Email"`
		Name        string `json:"Name"`
		MailboxHash string `json:"MailboxHash"`
	} `json:"FromFull"`
	ToFull []struct {
		Email       string `json:"Email"`
		Name        string `json:"Name"`
		MailboxHash string `json:"MailboxHash"`
	} `json:"ToFull"`
	CcFull []*struct {
		Email       string `json:"Email"`
		Name        string `json:"Name"`
		MailboxHash string `json:"MailboxHash"`
	} `json:"CcFull"`
	BccFull []*struct {
		Email       string `json:"Email"`
		Name        string `json:"Name"`
		MailboxHash string `json:"MailboxHash"`
	} `json:"BccFull"`
	OriginalRecipient string `json:"OriginalRecipient"`
	Subject           string `json:"Subject"`
	MessageID         string `json:"MessageID"`
	ReplyTo           string `json:"ReplyTo"`
	MailboxHash       string `json:"MailboxHash"`
	Date              string `json:"Date"`
	TextBody          string `json:"TextBody"`
	HtmlBody          string `json:"HtmlBody"`
	Tag               string `json:"Tag"`
	Headers           []struct {
		Name  string `json:"Name"`
		Value string `json:"Value"`
	} `json:"Headers"`
}
