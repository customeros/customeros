package service

import postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"

type EmailInService interface {
	LoadEmail(rawEmail *postgresentity.RawEmail) (EmailMessageData, error)
	ProcessEmailCheck(email *EmailMessageData) HeaderAnalysis
}

type emailInService struct {
	services *Services
}

func NewEmailInService() EmailInService {
	return &emailInService{}
}
