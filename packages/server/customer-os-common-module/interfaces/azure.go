package interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type AzureService interface {
	ReadEmailsFromAzureAd(ctx context.Context, importState *entity.UserEmailImportState) ([]*entity.EmailRawData, string, error)
	SendEmail(ctx context.Context, request *entity.EmailMessage) error
}
