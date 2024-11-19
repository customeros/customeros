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
