package service

type EmailInService interface {
	LoadEmail(rawEmail string) (*EmailMessageData, error)
	ProcessEmailCheck(email *EmailMessageData) HeaderAnalysis
}

type emailInService struct {
	services *Services
}

func NewEmailInService() EmailInService {
	return &emailInService{}
}
