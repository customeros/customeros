package interfaces

import "context"

type MailstackService interface {
	GetPaymentIntent(ctx context.Context, domains []string, usernames []string, amount int64) (string, error)                                                   // stripe client secret
	RegisterBuyDomainsWithMailboxes(ctx context.Context, test bool, paymentIntentId string, domains []string, usernames []string, redirectWebsite string) error // id, stripe client secret
	GetTenantForMailstackDomain(ctx context.Context, domain string) (string, error)
	// key domain, value tenant
	GetAllMailstackDomains(ctx context.Context) (map[string]string, error)
	ConfigureMailstackDomain(ctx context.Context, domain, redirectWebsite string) error
}
