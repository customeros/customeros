package interfaces

import "context"

type CreateMailboxRequest struct {
	IgnoreDomainOwnership bool
	Domain                string
	Username              string
	Password              string
	WebmailEnabled        bool
	ForwardingTo          []string
	LinkedUserEmail       string
}

type MailboxRecord struct {
	ID                string   `json:"id"`
	Email             string   `json:"email"`
	Password          string   `json:"password"`
	ForwardingTo      []string `json:"forwardingTo"`
	ForwardingEnabled bool     `json:"forwardingEnabled"`
	WebmailEnabled    bool     `json:"webmailEnabled"`
}

type RegisterMailboxResponse struct {
	Mailbox    *MailboxRecord
	StatusCode int
	ErrorMsg   string
}

type DomainRecord struct {
	Domain      string   `json:"domain"`
	CreatedDate string   `json:"createdDate"`
	ExpiredDate string   `json:"expiredDate"`
	Nameservers []string `json:"nameservers"`
}

type MailstackService interface {
	GetPaymentIntent(ctx context.Context, domains []string, usernames []string, amount int64) (string, error)
	RegisterBuyDomainsWithMailboxes(ctx context.Context, test bool, paymentIntentId string, domains []string, usernames []string, redirectWebsite string) error
	GetTenantForMailstackDomain(ctx context.Context, domain string) (string, error)
	GetAllMailstackDomains(ctx context.Context) (map[string]string, error)
	RegisterMailbox(ctx context.Context, tenant string, domain string, request CreateMailboxRequest) (*RegisterMailboxResponse, error)
	ConfigureMailbox(ctx context.Context, tenant string, mailboxId string) error
	RegisterNewDomain(ctx context.Context, tenant string, domain string, website string) (int, string, *DomainRecord, error)
}
