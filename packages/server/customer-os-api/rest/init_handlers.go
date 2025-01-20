package rest_handlers

import (
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/billing"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/customerbase"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/enrich"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/files"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/integrations"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/public"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/webhooks"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type RestHandlers struct {
	Response             *response.Response
	Billing              *billing.BillingHandler
	Contact              *customerbase.ContactHandler
	Enrich               *enrich.EnrichHandler
	Files                *files.FileHandler
	Integrations         *integrations.IntegrationHandler
	Organization         *customerbase.OrganizationHandler
	Webhooks             *webhooks.WebhookHandler
	WebsiteTrackerEvents *public.WebsiteTrackerEventsHandler
}

func InitRestHandlers(services *cosapi_services.Services) *RestHandlers {
	responseHandler := response.NewRestResponseHandler()

	handlers := RestHandlers{
		Billing:              billing.NewBillingHandler(services, responseHandler),
		Contact:              customerbase.NewContactHandler(services, responseHandler),
		Enrich:               enrich.NewEnrichHandler(services, responseHandler),
		Files:                files.NewFileHandler(services, responseHandler),
		Integrations:         integrations.NewIntegrationHandler(services, responseHandler),
		Organization:         customerbase.NewOrganizationHandler(services, responseHandler),
		Webhooks:             webhooks.NewWebhookHandler(services, responseHandler),
		WebsiteTrackerEvents: public.NewWebsiteTrackerEventsHandler(responseHandler),
	}

	return &handlers
}
