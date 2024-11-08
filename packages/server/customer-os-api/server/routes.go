package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/constants"
	cosHandler "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	restbilling "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/billing"
	restcustomerbase "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/customerbase"
	restenrich "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/enrich"
	restmailstack "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/mailstack"
	restoutreach "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/outreach"
	restverify "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest/verify"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	commoncaches "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

const (
	outreachV1Path     = "/outreach/v1"
	customerBaseV1Path = "/customerbase/v1"
	billingV1Path      = "/billing/v1"
	verifyV1Path       = "/verify/v1"
	enrichV1Path       = "/enrich/v1"
	mailStackV1Path    = "/mailstack/v1"
)

func RegisterRestRoutes(ctx context.Context, r *gin.Engine, grpcClients *grpc_client.Clients, services *service.Services, cache *commoncaches.Cache) {
	registerPublicRoutes(ctx, r, services)
	registerOutreachRoutes(ctx, r, services, cache)
	registerCustomerBaseRoutes(ctx, r, services, grpcClients, cache)
	registerVerifyRoutes(ctx, r, services, cache)
	registerEnrichRoutes(ctx, r, services, cache)
	registerMailStackRoutes(ctx, r, services, cache)
	registerBillingRoutes(ctx, r, services, grpcClients, cache)
}

func registerPublicRoutes(ctx context.Context, r *gin.Engine, services *service.Services) {
	// Redirect to pay invoice link
	r.GET("/invoice/:invoiceId/pay",
		tracing.TracingEnhancer(ctx, "GET:/invoice/:invoiceId/pay"),
		rest.RedirectToPayInvoice(services))
	r.GET("/invoice/:invoiceId/paymentLink",
		tracing.TracingEnhancer(ctx, "GET:/invoice/:invoiceId/paymentLink"),
		rest.GetInvoicePaymentLink(services))
}

func registerEnrichRoutes(ctx context.Context, r *gin.Engine, services *service.Services, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/person", enrichV1Path), services, cache, restenrich.EnrichPerson(services))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/person/results/:id", enrichV1Path), services, cache, restenrich.EnrichPersonCallback(services))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/organizaiton", enrichV1Path), services, cache, restenrich.EnrichOrganization(services))
}

func registerVerifyRoutes(ctx context.Context, r *gin.Engine, services *service.Services, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/email", verifyV1Path), services, cache, restverify.VerifyEmailAddress(services))
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/email/bulk", verifyV1Path), services, cache, restverify.BulkUploadEmailsForVerification(services))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/email/bulk/results/:requestId", verifyV1Path), services, cache, restverify.GetBulkEmailVerificationResults(services))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/email/bulk/results/:requestId/download", verifyV1Path), services, cache, restverify.DownloadBulkEmailVerificationResults(services))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/ip", verifyV1Path), services, cache, restverify.IpIntelligence(services))
}

func registerCustomerBaseRoutes(ctx context.Context, r *gin.Engine, services *service.Services, grpcClients *grpc_client.Clients, cache *commoncaches.Cache) {
	registerOrganizationRoutes(ctx, r, services, grpcClients, cache)
	registerContactRoutes(ctx, r, services, grpcClients, cache)
}

func registerBillingRoutes(ctx context.Context, r *gin.Engine, services *service.Services, grpcClients *grpc_client.Clients, cache *commoncaches.Cache) {
	registerInvoiceRoutes(ctx, r, services, grpcClients, cache)
}

func registerOrganizationRoutes(ctx context.Context, r *gin.Engine, services *service.Services, grpcClients *grpc_client.Clients, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/organizations", customerBaseV1Path), services, cache, restcustomerbase.CreateOrganization(services))
	setupRestRoute(ctx, r, "PUT", fmt.Sprintf("%s/organizations/:id/links/:externalSystem/primary", customerBaseV1Path), services, cache, restcustomerbase.SetPrimaryExternalSystemId(services))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/organizations/:id", customerBaseV1Path), services, cache, restcustomerbase.GetOrganization(services))
}

func registerContactRoutes(ctx context.Context, r *gin.Engine, services *service.Services, grpcClients *grpc_client.Clients, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/contacts", customerBaseV1Path), services, cache, rest.CreateContact(services, grpcClients))
}

func registerInvoiceRoutes(ctx context.Context, r *gin.Engine, services *service.Services, grpcClients *grpc_client.Clients, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/organizations/:id/invoices", billingV1Path), services, cache, restbilling.GetInvoicesForOrganization(services))
}

func registerOutreachRoutes(ctx context.Context, r *gin.Engine, services *service.Services, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/track/email", outreachV1Path), services, cache, restoutreach.GenerateEmailTrackingUrls(services))
}

func registerMailStackRoutes(ctx context.Context, r *gin.Engine, services *service.Services, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/domains", mailStackV1Path), services, cache, restmailstack.RegisterNewDomain(services))
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/domains/configure", mailStackV1Path), services, cache, restmailstack.ConfigureDomain(services))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/domains", mailStackV1Path), services, cache, restmailstack.GetDomains(services))

	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/domains/:domain/mailboxes", mailStackV1Path), services, cache, restmailstack.RegisterNewMailbox(services))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/domains/:domain/mailboxes", mailStackV1Path), services, cache, restmailstack.GetMailboxes(services))
}

func setupRestRoute(ctx context.Context, r *gin.Engine, method, path string, services *service.Services, cache *commoncaches.Cache, handler gin.HandlerFunc) {
	r.Handle(method, path,
		tracing.TracingEnhancer(ctx, method+":"+path),
		security.ApiKeyCheckerHTTP(services.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository, services.Repositories.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_API, security.WithCache(cache)),
		enrichContextMiddleware(constants.AppSourceCustomerOsApiRest),
		cosHandler.StatsSuccessHandler(method+":"+path, services),
		handler)
}
