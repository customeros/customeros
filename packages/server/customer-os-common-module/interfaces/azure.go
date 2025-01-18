package interfaces

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type AzureService interface {
	ReadEmailsFromAzureAd(ctx context.Context, importState *postgres_entity.UserEmailImportState) ([]*postgres_entity.EmailRawData, string, error)
	SendEmail(ctx context.Context, request *postgres_entity.EmailMessage) error
}
