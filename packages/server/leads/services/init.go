package services

import (
	"github.com/customeros/customeros/packages/server/leads/internal/config"
	"github.com/customeros/customeros/packages/server/leads/internal/repository"
	"github.com/customeros/customeros/packages/server/leads/services/event_processor"
	"github.com/customeros/customeros/packages/server/leads/services/ipdata"
	"github.com/customeros/customeros/packages/server/leads/services/snitcher"
	ssl_cert "github.com/customeros/customeros/packages/server/leads/services/ssl_certificate"
	website_registration "github.com/customeros/customeros/packages/server/leads/services/website_registation"
)

type Services struct {
	WebsiteRegistrationService website_registration.WebsiteRegistrationService
	NGINXConfigManager         p.NGINXConfigManager
	SSLCertificateService      ssl_cert.SSLCertificateService
	EventProcessor             event_processor.EventProcessor
	WebSessionProducer         *websession_producer.WebSessionProducer
	SnitcherService            snitcher.SnitcherService
	IPDataService              ipdata.IPDataService
}

func InitServices(config *config.Config, repositories *repository.Repositories) *Services {
	return nil
}
