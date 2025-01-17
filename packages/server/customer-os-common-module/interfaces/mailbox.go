package interfaces

import (
	"context"

	"gorm.io/gorm"
)

type MailboxService interface {
	CreateMailbox(ctx context.Context, tx *gorm.DB, request CreateMailboxRequest) error
	IsDomainAvailable(ctx context.Context, domain string) (ok, available bool)
	RecommendOutboundDomains(ctx context.Context, domainRoot string, count int) []string
	ReputationScore(ctx context.Context, domain, tenant string) (int, error)
}

type CreateMailboxRequest struct {
	Domain          string
	Username        string
	Password        string
	LinkedUserEmail string
	WebmailEnabled  bool
	ForwardingTo    []string

	IgnoreDomainOwnership bool
}
