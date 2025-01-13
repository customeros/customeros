package interfaces

import "context"

type NamecheapService interface {
	CheckDomainAvailability(ctx context.Context, domain string) (bool, bool, error)
	PurchaseDomain(ctx context.Context, tenant, domain string) error
	GetDomainPrice(ctx context.Context, domain string) (float64, error)
	GetDomainInfo(ctx context.Context, tenant, domain string) (NamecheapDomainInfo, error)
	UpdateNameservers(ctx context.Context, tenant, domain string, nameservers []string) error
}

type NamecheapDomainInfo struct {
	DomainName  string   `json:"domainName"`
	CreatedDate string   `json:"createdDate"`
	ExpiredDate string   `json:"expiredDate"`
	Nameservers []string `json:"nameservers"`
	WhoisGuard  bool     `json:"whoisGuard"`
}
