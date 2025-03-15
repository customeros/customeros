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
	// mailboxes
	RegisterMailbox(ctx context.Context, tenant, domain string, request CreateMailboxRequest) (*RegisterMailboxResponse, error)
	ConfigureMailbox(ctx context.Context, tenant, mailboxId string) error
	// domains
	RegisterNewDomain(ctx context.Context, tenant, domain, website string) (int, string, *DomainRecord, error)
	ConfigureDomain(ctx context.Context, tenant, domain, website string) (int, string, *DomainRecord, error)
	GetDomains(ctx context.Context, tenant string) (int, string, []DomainRecord, error)
	RecommendDomain(ctx context.Context, tenant, baseName string) (int, string, []string, error)
	GetAllMailstackDomains(ctx context.Context) (map[string]string, error)
}
