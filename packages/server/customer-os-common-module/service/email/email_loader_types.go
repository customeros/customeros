package service

import "time"

// Email Raw Data
type EmailRawData struct {
	ProviderMessageId string            `json:"ProviderMessageId"`
	MessageId         string            `json:"MessageId"`
	Sent              string            `json:"Sent"`
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

// Core message content
type EmailContent struct {
	Html     string
	Text     string
	Subject  string
	SentDate string
}

type EmailParticipant struct {
	Email     string
	FirstName string
	LastName  string
}

// Message routing/delivery information
type EmailParticipants struct {
	From      EmailParticipant
	To        []EmailParticipant
	Cc        []EmailParticipant
	Bcc       []EmailParticipant
	InReplyTo EmailParticipant
}

// identifiers
type EmailIdentifiers struct {
	EmailThreadId       string
	ExternalId          string
	ExternalSystem      string
	ContactsExternalIds []string
	UserExternalId      string
	ProviderMessageId   string
	MessageId           string
	Reference           string
}

// Headers
type EmailHeaders struct {
	AutoSubmitted      bool
	ContentDescription string
	DeliveryStatus     bool
	ListUnsubscribe    bool
	Precedence         string
	ReturnPath         string
	XAutoreply         string
	XAutoresponse      string
	XLoop              bool
	XFailedRecepients  bool
	RawHeaders         map[string]string
}

// Combined message data
type EmailMessageData struct {
	Content      EmailContent
	Headers      EmailHeaders
	Participants EmailParticipants
	Identifiers  EmailIdentifiers
	CreatedAt    time.Time
	Channel      string
	ChannelData  *string
}
