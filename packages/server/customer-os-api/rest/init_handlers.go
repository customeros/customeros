package rest_handlers

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/billing"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/customerbase"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/enrich"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/files"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/integrations"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/mailstack"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/outreach"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/private"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/public"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	reveal "github.com/customeros/customeros/packages/server/customer-os-api/rest/reveal_setup"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/verify"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/webhooks"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type RestHandlers struct {
	Response *response.Response

	AskAI                *private.AskAIHandler
	Billing              *billing.BillingHandler
	Contact              *customerbase.ContactHandler
	Enrich               *enrich.EnrichHandler
	Files                *files.FileHandler
	Integrations         *integrations.IntegrationHandler
	Mail                 *private.MailHandler
	Mailtstack           *mailstack.MailstackHandler
	Outreach             *outreach.OutreachHandler
	Organization         *customerbase.OrganizationHandler
	PrivateIntegrations  *private.PrivateIntegrationHandler
	Verify               *verify.VerifyHandler
	Webhooks             *webhooks.WebhookHandler
	WebTracker           *reveal.WebTrackerHandler
	WebsiteTrackerEvents *public.WebsiteTrackerEventsHandler
}

func InitRestHandlers(services *cosapi_services.Services) *RestHandlers {
	responseHandler := response.NewRestResponseHandler()
	integrationsHandler := integrations.NewIntegrationHandler(services, responseHandler)

	handlers := RestHandlers{
		AskAI:                private.NewAskAIHandler(services, responseHandler),
		Billing:              billing.NewBillingHandler(services, responseHandler),
		Contact:              customerbase.NewContactHandler(services, responseHandler),
		Enrich:               enrich.NewEnrichHandler(services, responseHandler),
		Files:                files.NewFileHandler(services, responseHandler),
		Integrations:         integrationsHandler,
		Mail:                 private.NewMailHandler(services, responseHandler),
		Mailtstack:           mailstack.NewMailstackHandler(services, responseHandler),
		Outreach:             outreach.NewOutreackHandler(services, responseHandler),
		Organization:         customerbase.NewOrganizationHandler(services, responseHandler),
		PrivateIntegrations:  private.NewPrivateIntegrationHandler(services, responseHandler),
		Verify:               verify.NewVerifyHandler(services, responseHandler),
		Webhooks:             webhooks.NewWebhookHandler(services, responseHandler, integrationsHandler),
		WebTracker:           reveal.NewWebTrackerHandler(services, responseHandler),
		WebsiteTrackerEvents: public.NewWebsiteTrackerEventsHandler(services, responseHandler),
	}

	return &handlers
}
