package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/billing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/customerbase"
	restEnrich "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/enrich"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/flows"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/mailstack"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/outreach"
	reveal "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/reveal_setup"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/verify"
	cosapi_services "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services"
)

const (
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

func registerBillingRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, cache *caches.Cache) {
	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/organizations/:id/invoices", BillingPath),
		handler:   billing.GetInvoicesForOrganization(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})
}

func registerCustomerBaseRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, cache *caches.Cache) {
	// Organization Routes
	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/organizations", CustomerBasePath),
		handler:   customerbase.CreateOrganization(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/organizations/:id", CustomerBasePath),
		handler:   customerbase.GetOrganization(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "PUT",
		path:      fmt.Sprintf("%s/organizations/:id/links/:externalSystem/primary", CustomerBasePath),
		handler:   customerbase.SetPrimaryExternalSystemId(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	// Contact Routes
	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/contacts", CustomerBasePath),
		handler:   customerbase.CreateContact(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/contacts/bulk", CustomerBasePath),
		handler:   customerbase.CreateBulkContacts(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/contacts/import", CustomerBasePath),
		handler:   customerbase.ImportContacts(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})
}

func registerEnrichRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, cache *caches.Cache) {
	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/person", EnrichPath),
		handler:   restEnrich.EnrichPerson(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/person/results/:id", EnrichPath),
		handler:   restEnrich.EnrichPersonCallback(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/organization", EnrichPath),
		handler:   restEnrich.EnrichOrganization(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})
}

func registerFlowRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, cache *caches.Cache) {

	// webhook admin
	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/hooks", FlowsPath),
		handler:   flows.GetActiveWebhooks(s, CustomerOSAPIURL(), FlowsPath),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/hooks", FlowsPath),
		handler:   flows.CreateWebhook(s, CustomerOSAPIURL(), FlowsPath),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/:tenantId/i/:integrationId/rotate", FlowsPath),
		handler:   flows.RotateWebhook(s, CustomerOSAPIURL(), FlowsPath),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "DELETE",
		path:      fmt.Sprintf("%s/:tenantId/i/:integrationId", FlowsPath),
		handler:   flows.DeactivateWebhook(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})
}

func registerIDRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services) {
	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      "/me",
		handler:   handlers.AuthorizeMe(),
		routeType: RouteCustomer,
		services:  s,
	})
}

func registerMailstackRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, cache *caches.Cache) {
	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/domains", MailstackPath),
		handler:   mailstack.RegisterNewDomain(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/domains", MailstackPath),
		handler:   mailstack.GetDomains(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/domains/recommendations", MailstackPath),
		handler:   mailstack.RecommendDomain(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/domains/configure", MailstackPath),
		handler:   mailstack.ConfigureDomain(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/domains/:domain/mailboxes", MailstackPath),
		handler:   mailstack.RegisterNewMailbox(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/domains/:domain/mailboxes", MailstackPath),
		handler:   mailstack.GetMailboxes(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/domains/:domain/dns", MailstackPath),
		handler:   mailstack.DNS(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/domains/:domain/dns", MailstackPath),
		handler:   mailstack.AddDNSRecord(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "DELETE",
		path:      fmt.Sprintf("%s/domains/:domain/dns", MailstackPath),
		handler:   mailstack.DeleteDNSRecord(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})
}

func registerOutreachRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, cache *caches.Cache) {
	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/track/email", OutreachPath),
		handler:   outreach.GenerateEmailTrackingUrls(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})
}

func registerRevealRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, cache *caches.Cache) {
	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/trackers", RevealPath),
		handler:   reveal.ProvisionTracker(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/verify", RevealPath),
		handler:   reveal.VerifyTracker(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})
}

func registerVerifyRoutes(ctx context.Context, r *gin.Engine, s *cosapi_services.Services, cache *caches.Cache) {
	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/email", VerifyPath),
		handler:   verify.VerifyEmailAddress(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/email/bulk", VerifyPath),
		handler:   verify.BulkUploadEmailsForVerification(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/email/bulk/results/:requestId", VerifyPath),
		handler:   verify.GetBulkEmailVerificationResults(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/email/bulk/results/:requestId/download", VerifyPath),
		handler:   verify.DownloadBulkEmailVerificationResults(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/ip", VerifyPath),
		handler:   verify.IpIntelligence(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})
}
