package interfaces

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/gmail/v1"
)

type GoogleService interface {
	ServiceAccountCredentialsExistsForTenant(ctx context.Context, tenant string) (bool, error)

	GetGmailService(ctx context.Context, username, tenant string) (*gmail.Service, error)

	GetGmailServiceWithServiceAccount(ctx context.Context, username string, tenant string) (*gmail.Service, error)
	GetGCalServiceWithServiceAccount(ctx context.Context, username string, tenant string) (*calendar.Service, error)

	GetGmailServiceWithOauthToken(ctx context.Context, tokenEntity postgres_entity.OAuthTokenEntity) (*gmail.Service, error)
	GetGCalServiceWithOauthToken(ctx context.Context, tokenEntity postgres_entity.OAuthTokenEntity) (*calendar.Service, error)

	ReadEmails(ctx context.Context, batchSize int64, importState *postgres_entity.UserEmailImportState) ([]*postgres_entity.EmailRawData, string, error)

	SendEmail(ctx context.Context, request *postgres_entity.EmailMessage) error
}
