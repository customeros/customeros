package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type mailService struct {
	services *Services
}

type MailService interface {
	LoadEmail(rawEmail *postgresentity.RawEmail) (EmailMessageData, error)
	ProcessEmailCheck(email *EmailMessageData) HeaderAnalysis
	SendMail(ctx context.Context, emailMessage *entity.EmailMessage) error
	SyncEmail(tenant string, emailId uuid.UUID) (postgresentity.RawState, *string, error)
}

func NewMailService(services *Services) MailService {
	return &mailService{
		services: services,
	}
}
