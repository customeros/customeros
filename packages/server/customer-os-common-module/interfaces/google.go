package interfaces

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"google.golang.org/api/gmail/v1"
)

type GoogleService interface {
	GetGmailService(ctx context.Context, username, tenant string) (*gmail.Service, error)

	GetGmailServiceWithOauthToken(ctx context.Context, tokenEntity postgres_entity.OAuthTokenEntity) (*gmail.Service, error)

	ReadEmails(ctx context.Context, batchSize int64, importState *postgres_entity.IngestEmailImportState) ([]*postgres_entity.EmailRawData, string, error)

	SendEmail(ctx context.Context, request *postgres_entity.EmailMessage) error

	GetAccessToken(ctx context.Context, tenant, email string) (string, error)

	GetRefreshToken(ctx context.Context, tenant, email string) (string, error)
}
