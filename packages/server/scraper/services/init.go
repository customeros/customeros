package services

import (
	"github.com/customeros/customeros/packages/server/eventstream/services/event_processor"
	"github.com/customeros/customeros/packages/server/eventstream/services/nginx_config_manager"
	ssl_cert "github.com/customeros/customeros/packages/server/eventstream/services/ssl_certificate"
	"github.com/customeros/customeros/packages/server/eventstream/services/websession_producer"
	website_registration "github.com/customeros/customeros/packages/server/eventstream/services/website_registation"
)

type Services struct {
	WebsiteRegistrationService website_registration.WebsiteRegistrationService
	NGINXConfigManager         nginx_config_manager.NGINXConfigManager
	SSLCertificateService      ssl_cert.SSLCertificateService
	EventProcessor             event_processor.EventProcessor
	WebSessionProducer         *websession_producer.WebSessionProducer
}
