package interfaces

import (
	"context"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
)

type OpenSrsService interface {
	SendEmail(ctx context.Context, request *postgres_entity.EmailMessage) error
	SetMailstackService(mailstack MailstackService)
}
