package interfaces

import "context"

type CloudflareService interface {
	SetupDomainForMailStack(ctx context.Context, tenant, domain, destinationUrl string) ([]string, error)
	AddDNSRecord(ctx context.Context, zoneID, recordType, name, content string, ttl int, proxied bool, priority *int) error
	GetDNSRecords(ctx context.Context, domain string) (*[]DNSRecord, error)
	DeleteDNSRecord(ctx context.Context, zoneID string, recordID string) error
	CheckDomainExists(ctx context.Context, domain string) (bool, string, error)
}

type DNSRecord struct {
	ID      string `json:"id"`
	ZoneID  string `json:"zone_id"`
	Name    string `json:"zone_name"`
	Type    string `json:"type"`
	Content string `json:"content"`
}
