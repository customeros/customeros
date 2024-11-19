package service

import "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/repository"

type EmailInService interface {
	LoadEmail(rawEmail string) (*EmailMessageData, error)
	ProcessEmailCheck(email *EmailMessageData) HeaderAnalysis
}

type emailInService struct {
	repositories *repository.Repositories
	services     *Services
}

func NewEmailInService() EmailInService {
	return &emailInService{}
}
