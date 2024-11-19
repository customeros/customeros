package service

import "time"

type HeaderAnalysis struct {
	ShouldProcess   bool
	IsBounce        bool
	IsAutoResponder bool
	IsBulkMail      bool
}

// Core message content
type EmailContent struct {
	Html    string
	Text    string
	Subject string
}

// Message routing/delivery information
type EmailParticipants struct {
	FromEmail     string
	FromFirstName string
	FromLastName  string
	ToEmail       []string
	CcEmail       []string
	BccEmail      []string
}

// External system identifiers
type EmailIdentifiers struct {
	EmailThreadId       string
	ExternalId          string
	ExternalSystem      string
	ContactsExternalIds []string
	UserExternalId      string
}

// Headers
type EmailHeaders struct {
	ReturnPath         string
	XAutoreply         string
	XAutoresponse      string
	AutoSubmitted      bool
	XLoop              bool
	Precedence         string
	ListUnsubscribe    bool
	XFailedRecepients  bool
	ContentDescription string
	DeliveryStatus     bool
	RawHeaders         string
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
