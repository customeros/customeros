package interfaces

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type MailService interface {
	SetContactService(ContactService)
	SetEmailService(EmailService)
	SetInteractionEventService(InteractionEventService)
	SetOrganizationService(OrganizationService)
	IsInitialized() bool

	ExtractEmails(s string) []string
	GetEmailsForProcessingForUser(ctx context.Context, tenant, userEmailAddress string)
	LoadEmail(ctx context.Context, rawEmail *entity.RawEmail) (EmailMessageData, error)
	ProcessEmailCheck(ctx context.Context, tenant string, email *EmailMessageData) HeaderAnalysis
	ProcessEmail(ctx context.Context, tenant string, emailId uuid.UUID) entity.UpdateRawEmailTable
	ProcessEmailByMessageId(ctx context.Context, tenant, usernameSource, messageId string) entity.UpdateRawEmailTable
	SendMail(ctx context.Context, emailMessage *entity.EmailMessage) error
	ProcessSentEmail(ctx context.Context, tx *neo4j.ManagedTransaction, emailMessage *entity.EmailMessage) (*string, error)
	GetEmailIdForEmail(ctx context.Context, txWithPostCommit *utils.TxWithPostCommit, tenant, email string, source string) (string, error)
}

type HeaderAnalysis struct {
	ProcessEmail    bool
	IsBounce        bool
	IsAutoResponder bool
	IsBulkMail      bool
	SkipReason      string
	BouncedEmails   []string
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
	ReplyTo   []EmailParticipant
	AllEmails []string
}

func (ep EmailParticipants) GetToEmailAddresses() []string {
	output := make([]string, 0)
	for _, email := range ep.To {
		if email.Email != "" {
			output = append(output, email.Email)
		}
	}
	return utils.RemoveDuplicates(output)
}

func (ep EmailParticipants) GetCcEmailAddresses() []string {
	output := make([]string, 0)
	for _, email := range ep.Cc {
		if email.Email != "" {
			output = append(output, email.Email)
		}
	}
	return utils.RemoveDuplicates(output)
}

func (ep EmailParticipants) GetBccEmailAddresses() []string {
	output := make([]string, 0)
	for _, email := range ep.Bcc {
		if email.Email != "" {
			output = append(output, email.Email)
		}
	}
	return utils.RemoveDuplicates(output)
}

func (ep EmailParticipants) GetReplyToEmailAddresses() []string {
	output := make([]string, 0)
	for _, email := range ep.ReplyTo {
		if email.Email != "" {
			output = append(output, email.Email)
		}
	}
	return utils.RemoveDuplicates(output)
}

// identifiers
type EmailIdentifiers struct {
	EmailThreadId       string
	ExternalSystem      string
	ContactsExternalIds []string
	UserExternalId      string
	ProviderMessageId   string
	MessageId           string
	References          []string
}

// Headers
type EmailHeaders struct {
	AutoSubmitted      bool
	ContentDescription string
	DeliveryStatus     bool
	ListUnsubscribe    bool
	Precedence         string
	ReturnPath         string
	ReturnPathExists   bool
	XAutoreply         string
	XAutoresponse      string
	XLoop              bool
	XFailedRecepients  []string
	ReplyTo            string
	ReplyToExists      bool
	Sender             string
	ForwardedFor       string
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
