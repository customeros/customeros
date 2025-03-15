package interfaces

import (
	"context"
)

type MailboxService interface {
	ReputationScore(ctx context.Context, domain, tenant string) (int, error)
}
