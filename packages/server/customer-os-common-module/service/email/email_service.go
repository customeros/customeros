package service

type EmailService interface {
	LoadEmail(rawEmail string) (*EmailMessageData, error)
	AnalyzeEmail(email *EmailMessageData) HeaderAnalysis
}

type emailService struct{}

func NewEmailService() EmailService {
	return &emailService{}
}
