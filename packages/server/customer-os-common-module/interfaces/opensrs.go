package interfaces

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type OpenSrsService interface {
	SendEmail(ctx context.Context, request *postgres_entity.EmailMessage) error
	SetupDomain(ctx context.Context, tenant, domain string) error
	SetupMailbox(ctx context.Context, tenant, username, password string, forwardingTo []string, webmailEnabled bool) error
	GetMailboxDetails(ctx context.Context, email string) (MailboxDetails, error)
}

type MailboxDetails struct {
	Email             string   `json:"email"`
	ForwardingEnabled bool     `json:"forwardingEnabled"`
	ForwardingTo      []string `json:"forwardingTo"`
	WebmailEnabled    bool     `json:"webmailEnabled"`
}
