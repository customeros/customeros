package services

import (
	"github.com/customeros/customeros/packages/server/customer-os-common-module/caches"

	"github.com/customeros/customeros/packages/server/mailstack/interfaces"
	"github.com/customeros/customeros/packages/server/mailstack/services/imap"
)

type Services struct {
	Cache       *caches.Cache
	IMAPService interfaces.IMAPService
}

func InitServices() *Services {
	services := Services{
		IMAPService: imap.NewIMAPService(),
	}

	return &services
}
