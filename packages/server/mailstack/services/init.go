package services

import (
	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/services/imap"
)

type Services struct {
	IMAPService interfaces.IMAPService
}

func InitServices() *Services {
	services := Services{
		IMAPService: imap.NewIMAPService(),
	}

	return &services
}
