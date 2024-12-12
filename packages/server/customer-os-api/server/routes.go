package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	commoncaches "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	cosHandler "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	restbilling "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/billing"
	restcustomerbase "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/customerbase"
	restenrich "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/enrich"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/flows"
	restmailstack "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/mailstack"
	restoutreach "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/outreach"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/tracking"
	restverify "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/verify"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

const (
	outreachV1Path     = "/outreach/v1"
	customerBaseV1Path = "/customerbase/v1"
	billingV1Path      = "/billing/v1"
	verifyV1Path       = "/verify/v1"
	enrichV1Path       = "/enrich/v1"
	mailStackV1Path    = "/mailstack/v1"
	flowsV1Path        = "/flows/v1"
	webhooksV1Path     = "/webhooks/v1"
)

func RegisterRestRoutes(ctx context.Context, r *gin.Engine, grpcClients *grpc_client.Clients, s *service.Services, cache *commoncaches.Cache) {
	registerPublicRoutes(ctx, r, s)
	registerHealthRoutes(ctx, r, s, cache)
	registerBillingRoutes(ctx, r, s, grpcClients, cache)
	registerCustomerBaseRoutes(ctx, r, s, grpcClients, cache)
	registerEnrichRoutes(ctx, r, s, cache)
	registerFlowRoutes(ctx, r, s, cache)
	registerMailStackRoutes(ctx, r, s, cache)
	registerOutreachRoutes(ctx, r, s, cache)
	registerVerifyRoutes(ctx, r, s, cache)
}

func registerPublicRoutes(ctx context.Context, r *gin.Engine, s *service.Services) {
	// Redirect to pay invoice link
	setupPublicRoute(ctx, r, "GET", "/invoice/:invoiceId/pay", rest.RedirectToPayInvoice(s))
	setupPublicRoute(ctx, r, "GET", "/invoice/:invoiceId/paymentLink", rest.GetInvoicePaymentLink(s))

	// Internal Webhooks
	setupPublicRoute(ctx, r, "POST", fmt.Sprintf("%s/dmarc", webhooksV1Path), flows.PostmarkDMARCMonitor(s))

	// Flow Wehbooks
	setupPublicRoute(ctx, r, "POST", fmt.Sprintf("%s/:tenantId/i/:integrationId", flowsV1Path), flows.HandleWebhook(s, flowsV1Path))

	//tracking
	setupPublicRoute(ctx, r, "GET", "/v1/l", tracking.TrackLinkRequest(s))
	setupPublicRoute(ctx, r, "GET", "/v1/s", tracking.TrackOpenRequest(s))
	setupPublicRoute(ctx, r, "GET", "/v1/u", tracking.TrackUnsubscribeRequest(s))
}

func registerHealthRoutes(ctx context.Context, r *gin.Engine, s *service.Services, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "GET", "/me", s, cache, rest.AuthorizeMe(s))
}

func registerEnrichRoutes(ctx context.Context, r *gin.Engine, s *service.Services, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/person", enrichV1Path), s, cache, restenrich.EnrichPerson(s))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/person/results/:id", enrichV1Path), s, cache, restenrich.EnrichPersonCallback(s))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/organization", enrichV1Path), s, cache, restenrich.EnrichOrganization(s))
}

func registerVerifyRoutes(ctx context.Context, r *gin.Engine, s *service.Services, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/email", verifyV1Path), s, cache, restverify.VerifyEmailAddress(s))
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/email/bulk", verifyV1Path), s, cache, restverify.BulkUploadEmailsForVerification(s))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/email/bulk/results/:requestId", verifyV1Path), s, cache, restverify.GetBulkEmailVerificationResults(s))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/email/bulk/results/:requestId/download", verifyV1Path), s, cache, restverify.DownloadBulkEmailVerificationResults(s))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/ip", verifyV1Path), s, cache, restverify.IpIntelligence(s))
}

func registerCustomerBaseRoutes(ctx context.Context, r *gin.Engine, s *service.Services, grpcClients *grpc_client.Clients, cache *commoncaches.Cache) {
	registerOrganizationRoutes(ctx, r, s, grpcClients, cache)
	registerContactRoutes(ctx, r, s, grpcClients, cache)
}

func registerFlowRoutes(ctx context.Context, r *gin.Engine, s *service.Services, cache *commoncaches.Cache) {
    // flow builder
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/flows", flowsV1Path), s, cache, flows.CreateFlow(s))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/flows", flowsV1Path), s, cache, ...)

    // manage specific flow
    setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/flows/:flowId", flowsV1Path), s, cache, ...)
    setupRestRoute(ctx, r, "PUT", fmt.Sprintf("%s/flows/:flowId", flowsV1Path), s, cache, ...)
    setupRestRoute(ctx, r, "DELETE", fmt.Sprintf("%s/flows/:flowId", flowsV1Path), s, cache, ...)

    // manage flow nodes
    setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/flows/:flowId/nodes", flowsV1Path), s, cache, ...)
    setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/flows/:flowId/nodes", flowsV1Path), s, cache, ...)
    setupRestRoute(ctx, r, "PUT", fmt.Sprintf("%s/flows/:flowId/nodes/:nodeId", flowsV1Path), s, cache, ...)
    setupRestRoute(ctx, r, "DELETE", fmt.Sprintf("%s/flows/:flowId/nodes/:nodeId", flowsV1Path), s, cache, ...)

    // manage flow edges
    setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/flows/:flowId/edges", flowsV1Path), s, cache, ...)
    setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/flows/:flowId/edges", flowsV1Path), s, cache, ...)
    setupRestRoute(ctx, r, "DELETE", fmt.Sprintf("%s/flows/:flowId/edges/:edgeId", flowsV1Path), s, cache, ...)

    // flow validation
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/actions", flowsV1Path), s, cache, flows.GetActions(s))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/listeners", flowsV1Path), s, cache, flows.GetListeners(s))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/transitions", flowsV1Path), s, cache, flows.GetTransitions(s))

    // webhook admin
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/hooks", flowsV1Path), s, cache, flows.GetActiveWebhooks(s, CustomerOSAPIURL(), flowsV1Path))
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/hooks", flowsV1Path), s, cache, flows.CreateWebhook(s, CustomerOSAPIURL(), flowsV1Path))
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/:tenantId/i/:integrationId/rotate", flowsV1Path), s, cache, flows.RotateWebhook(s, CustomerOSAPIURL(), flowsV1Path))
	setupRestRoute(ctx, r, "DELETE", fmt.Sprintf("%s/:tenantId/i/:integrationId", flowsV1Path), s, cache, flows.DeactivateWebhook(s))
}

func registerBillingRoutes(ctx context.Context, r *gin.Engine, s *service.Services, grpcClients *grpc_client.Clients, cache *commoncaches.Cache) {
	registerInvoiceRoutes(ctx, r, s, grpcClients, cache)
}

func registerOrganizationRoutes(ctx context.Context, r *gin.Engine, s *service.Services, grpcClients *grpc_client.Clients, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/organizations", customerBaseV1Path), s, cache, restcustomerbase.CreateOrganization(s))
	// setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/organizations/bulk", customerBaseV1Path), s, cache, restcustomerbase.CreateBulkOrganizations(s))
	// setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/organizations/import", customerBaseV1Path), s, cache, restcustomerbase.ImportOrganizations(s))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/organizations/:id", customerBaseV1Path), s, cache, restcustomerbase.GetOrganization(s))
	setupRestRoute(ctx, r, "PUT", fmt.Sprintf("%s/organizations/:id/links/:externalSystem/primary", customerBaseV1Path),
		s, cache, restcustomerbase.SetPrimaryExternalSystemId(s))
}

func registerContactRoutes(ctx context.Context, r *gin.Engine, s *service.Services, grpcClients *grpc_client.Clients, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/contacts", customerBaseV1Path), s, cache, restcustomerbase.CreateContact(s))
	// setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/contacts/:id", customerBaseV1Path), services, cache, restcustomerbase.GetContact(services))
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/contacts/bulk", customerBaseV1Path), s, cache, restcustomerbase.CreateBulkContacts(s))
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/contacts/import", customerBaseV1Path), s, cache, restcustomerbase.ImportContacts(s))
}

func registerInvoiceRoutes(ctx context.Context, r *gin.Engine, s *service.Services, grpcClients *grpc_client.Clients, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/organizations/:id/invoices", billingV1Path), s, cache, restbilling.GetInvoicesForOrganization(s))
}

func registerOutreachRoutes(ctx context.Context, r *gin.Engine, s *service.Services, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/track/email", outreachV1Path), s, cache, restoutreach.GenerateEmailTrackingUrls(s))
}

func registerMailStackRoutes(ctx context.Context, r *gin.Engine, s *service.Services, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/domains", mailStackV1Path), s, cache, restmailstack.RegisterNewDomain(s))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/domains", mailStackV1Path), s, cache, restmailstack.GetDomains(s))
	// setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/domains/recommendations", mailStackV1Path), s, cache, restmailstack.GetDomainRecommendations(s))
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/domains/configure", mailStackV1Path), s, cache, restmailstack.ConfigureDomain(s))

	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/domains/:domain/mailboxes", mailStackV1Path), s, cache, restmailstack.RegisterNewMailbox(s))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/domains/:domain/mailboxes", mailStackV1Path), s, cache, restmailstack.GetMailboxes(s))
}

func setupRestRoute(ctx context.Context, r *gin.Engine, method, path string, s *service.Services, cache *commoncaches.Cache, handler gin.HandlerFunc) {
	r.Handle(method, path,
		tracing.TracingEnhancer(ctx, method+":"+path),
		security.ApiKeyCheckerHTTP(s.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository,
			s.Repositories.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_API, security.WithCache(cache)),
		enrichContextMiddleware(constants.AppSourceCustomerOsApiRest),
		cosHandler.StatsSuccessHandler(method+":"+path, s),
		handler)
}

func setupPublicRoute(ctx context.Context, r *gin.Engine, method, path string, handler gin.HandlerFunc) {
	r.Handle(method, path, tracing.TracingEnhancer(ctx, method+":"+path), handler)
}
