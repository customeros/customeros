package interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
)

type GoogleService interface {
	ServiceAccountCredentialsExistsForTenant(ctx context.Context, tenant string) (bool, error)

	GetGmailService(ctx context.Context, username, tenant string) (*gmail.Service, error)

	GetGmailServiceWithServiceAccount(ctx context.Context, username string, tenant string) (*gmail.Service, error)
	GetGCalServiceWithServiceAccount(ctx context.Context, username string, tenant string) (*calendar.Service, error)

	GetGmailServiceWithOauthToken(ctx context.Context, tokenEntity entity.OAuthTokenEntity) (*gmail.Service, error)
	GetGCalServiceWithOauthToken(ctx context.Context, tokenEntity entity.OAuthTokenEntity) (*calendar.Service, error)

	ReadEmails(ctx context.Context, batchSize int64, importState *entity.UserEmailImportState) ([]*entity.EmailRawData, string, error)

	SendEmail(ctx context.Context, request *entity.EmailMessage) error
}
