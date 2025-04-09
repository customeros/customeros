package services

import (
	"github.com/customeros/customeros/packages/server/eventstream/internal/config"
	"github.com/customeros/customeros/packages/server/eventstream/internal/repository"
	"github.com/customeros/customeros/packages/server/eventstream/services/event_processor"
	"github.com/customeros/customeros/packages/server/eventstream/services/ipdata"
	"github.com/customeros/customeros/packages/server/eventstream/services/snitcher"
	ssl_cert "github.com/customeros/customeros/packages/server/eventstream/services/ssl_certificate"
	website_registration "github.com/customeros/customeros/packages/server/eventstream/services/website_registation"
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
