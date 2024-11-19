package service

import (
	"github.com/google/uuid"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
)

type EmailInService interface {
	LoadEmail(rawEmail *postgresentity.RawEmail) (EmailMessageData, error)
	ProcessEmailCheck(email *EmailMessageData) HeaderAnalysis
	SyncEmail(tenant string, emailId uuid.UUID) (postgresentity.RawState, *string, error)
}

type emailInService struct {
	services *Services
}

func NewEmailInService(services *Services) EmailInService {
	return &emailInService{
		services: services,
	}
}
