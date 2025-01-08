package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers/billing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers/customerbase"
	restEnrich "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers/enrich"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers/flows"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers/mailstack"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers/outreach"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers/reveal"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handlers/verify"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
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

func registerBillingRoutes(ctx context.Context, r *gin.Engine, s *service.Services, cache *caches.Cache) {
	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/organizations/:id/invoices", BillingPath),
		handler:   billing.GetInvoicesForOrganization(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})
}

func registerCustomerBaseRoutes(ctx context.Context, r *gin.Engine, s *service.Services, cache *caches.Cache) {
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

func registerEnrichRoutes(ctx context.Context, r *gin.Engine, s *service.Services, cache *caches.Cache) {
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
		path:      fmt.Sprintf("%s/person/organization", EnrichPath),
		handler:   restEnrich.EnrichOrganization(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})
}

func registerFlowRoutes(ctx context.Context, r *gin.Engine, s *service.Services, cache *caches.Cache) {
	// flow builder
	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/flows", FlowsPath),
		handler:   flows.CreateFlow(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/flows", FlowsPath),
		handler:   flows.GetFlows(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	// manage specific flow
	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/flows/:flowId", FlowsPath),
		handler:   flows.GetFlows(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "PUT",
		path:      fmt.Sprintf("%s/flows/:flowId", FlowsPath),
		handler:   flows.UpdateFlow(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "DELETE",
		path:      fmt.Sprintf("%s/flows/:flowId", FlowsPath),
		handler:   flows.DeleteFlow(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/flows/:flowId/on", FlowsPath),
		handler:   flows.TurnOn(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/flows/:flowId/off", FlowsPath),
		handler:   flows.TurnOff(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	// manage flow nodes
	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/flows/:flowId/nodes", FlowsPath),
		handler:   flows.GetFlowNodes(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/flows/:flowId/nodes", FlowsPath),
		handler:   flows.CreateFlowNode(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/flows/:flowId/nodes/:nodeId", FlowsPath),
		handler:   flows.GetFlowNodes(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "PUT",
		path:      fmt.Sprintf("%s/flows/:flowId/nodes/:nodeId", FlowsPath),
		handler:   flows.UpdateFlowNode(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "DELETE",
		path:      fmt.Sprintf("%s/flows/:flowId/nodes/:nodeId", FlowsPath),
		handler:   flows.DeleteFlowNode(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	// manage flow edges
	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/flows/:flowId/edges", FlowsPath),
		handler:   flows.GetFlowEdges(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/flows/:flowId/edges", FlowsPath),
		handler:   flows.CreateFlowEdge(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/flows/:flowId/edges/:edgeId", FlowsPath),
		handler:   flows.GetFlowEdges(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "PUT",
		path:      fmt.Sprintf("%s/flows/:flowId/edges/:edgeId", FlowsPath),
		handler:   flows.UpdateFlowEdge(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "DELETE",
		path:      fmt.Sprintf("%s/flows/:flowId/edges/:edgeId", FlowsPath),
		handler:   flows.DeleteFlowEdge(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	// flow validation
	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/agents", FlowsPath),
		handler:   flows.GetAgents(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/listeners", FlowsPath),
		handler:   flows.GetListeners(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      fmt.Sprintf("%s/transitions", FlowsPath),
		handler:   flows.GetTransitions(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})

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

func registerIDRoutes(ctx context.Context, r *gin.Engine, s *service.Services) {
	registerRoute(ctx, r, RouteConfig{
		method:    "GET",
		path:      "/me",
		handler:   handlers.AuthorizeMe(s),
		routeType: RouteCustomer,
		services:  s,
	})
}

func registerMailstackRoutes(ctx context.Context, r *gin.Engine, s *service.Services, cache *caches.Cache) {
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

func registerOutreachRoutes(ctx context.Context, r *gin.Engine, s *service.Services, cache *caches.Cache) {
	registerRoute(ctx, r, RouteConfig{
		method:    "POST",
		path:      fmt.Sprintf("%s/track/email", OutreachPath),
		handler:   outreach.GenerateEmailTrackingUrls(s),
		routeType: RouteCustomer,
		services:  s,
		cache:     cache,
	})
}

func registerRevealRoutes(ctx context.Context, r *gin.Engine, s *service.Services, cache *caches.Cache) {
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

func registerVerifyRoutes(ctx context.Context, r *gin.Engine, s *service.Services, cache *caches.Cache) {
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
