package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	rest_handlers "github.com/customeros/customeros/packages/server/customer-os-api/rest"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/me"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

const (
	AgentsPath       = "/agents/v1"
	BillingPath      = "/billing/v1"
	CustomerBasePath = "/customerbase/v1"
	EnrichPath       = "/enrich/v1"
	FlowsPath        = "/flows/v1"
	MailstackPath    = "/mailstack/v1"
	OutreachPath     = "/outreach/v1"
	RevealPath       = "/reveal/v1"
	VerifyPath       = "/verify/v1"
	WebhooksPath     = "/webhooks/v1"
)

func registerAgentsRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, h *rest_handlers.RestHandlers) {
	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/registry", AgentsPath),
		handler:   h.Agents.RegisterAgent(),
		routeType: RouteInternal,
		services:  s,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/registry", AgentsPath),
		handler:   h.Agents.AgentRegistry(),
		routeType: RouteCustomer,
		services:  s,
	})
}

func registerBillingRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, h *rest_handlers.RestHandlers) {
	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/organizations/:id/invoices", BillingPath),
		handler:   h.Billing.GetInvoicesForOrganization(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})
}

func registerCustomerBaseRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, h *rest_handlers.RestHandlers) {
	// Organization Routes
	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/organizations", CustomerBasePath),
		handler:   h.Organization.CreateOrganization(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/organizations/:id", CustomerBasePath),
		handler:   h.Organization.GetOrganization(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "PUT",
		path:      fmt.Sprintf("%s/organizations/:id/links/:externalSystem/primary", CustomerBasePath),
		handler:   h.Organization.SetPrimaryExternalSystemId(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	// Contact Routes
	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/contacts", CustomerBasePath),
		handler:   h.Contact.CreateContact(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/contacts/bulk", CustomerBasePath),
		handler:   h.Contact.CreateBulkContacts(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/contacts/import", CustomerBasePath),
		handler:   h.Contact.ImportContacts(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})
}

func registerEnrichRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, h *rest_handlers.RestHandlers) {
	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/person", EnrichPath),
		handler:   h.Enrich.EnrichPerson(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/person/results/:id", EnrichPath),
		handler:   h.Enrich.EnrichPersonCallback(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/organization", EnrichPath),
		handler:   h.Enrich.EnrichOrganization(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})
}

func registerFlowRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, h *rest_handlers.RestHandlers) {
	// webhook admin
	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/hooks", FlowsPath),
		handler:   h.Webhooks.GetActiveWebhooks(CustomerOSAPIURL(), FlowsPath),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/hooks", FlowsPath),
		handler:   h.Webhooks.CreateWebhook(CustomerOSAPIURL(), FlowsPath),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/:tenantId/i/:integrationId/rotate", FlowsPath),
		handler:   h.Webhooks.RotateWebhook(CustomerOSAPIURL(), FlowsPath),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "DELETE",
		path:      fmt.Sprintf("%s/:tenantId/i/:integrationId", FlowsPath),
		handler:   h.Webhooks.DeactivateWebhook(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})
}

func registerIDRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, h *rest_handlers.RestHandlers) {
	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      "/me",
		handler:   me.AuthorizeMe(h),
		routeType: RouteCustomer,
		services:  s,
	})
}

func registerMailstackRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, h *rest_handlers.RestHandlers) {
	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/domains", MailstackPath),
		handler:   h.Mailtstack.RegisterNewDomain(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/domains", MailstackPath),
		handler:   h.Mailtstack.GetDomains(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/domains/recommendations", MailstackPath),
		handler:   h.Mailtstack.RecommendDomain(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/domains/configure", MailstackPath),
		handler:   h.Mailtstack.ConfigureDomain(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/domains/:domain/mailboxes", MailstackPath),
		handler:   h.Mailtstack.RegisterNewMailbox(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/domains/:domain/mailboxes", MailstackPath),
		handler:   h.Mailtstack.GetMailboxes(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/domains/:domain/dns", MailstackPath),
		handler:   h.Mailtstack.DNS(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/domains/:domain/dns", MailstackPath),
		handler:   h.Mailtstack.AddDNSRecord(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "DELETE",
		path:      fmt.Sprintf("%s/domains/:domain/dns", MailstackPath),
		handler:   h.Mailtstack.DeleteDNSRecord(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})
}

func registerOutreachRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, h *rest_handlers.RestHandlers) {
	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/track/email", OutreachPath),
		handler:   h.Outreach.GenerateEmailTrackingUrls(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})
}

func registerRevealRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, h *rest_handlers.RestHandlers) {
	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/trackers", RevealPath),
		handler:   h.WebTracker.ProvisionTracker(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

}

func registerVerifyRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, h *rest_handlers.RestHandlers) {
	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/email", VerifyPath),
		handler:   h.Verify.VerifyEmailAddress(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/email/bulk", VerifyPath),
		handler:   h.Verify.BulkUploadEmailsForVerification(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/email/bulk/results/:requestId", VerifyPath),
		handler:   h.Verify.GetBulkEmailVerificationResults(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/email/bulk/results/:requestId/download", VerifyPath),
		handler:   h.Verify.DownloadBulkEmailVerificationResults(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/ip", VerifyPath),
		handler:   h.Verify.IpIntelligence(),
		routeType: RouteCustomer,
		services:  s,
		cache:     s.Cache,
	})
}
