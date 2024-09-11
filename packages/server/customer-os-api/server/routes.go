package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	cosHandler "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	commoncaches "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

const (
	outreachV1Path     = "/outreach/v1"
	customerBaseV1Path = "/customerbase/v1"
	verifyV1Path       = "/verify/v1"
	enrichV1Path       = "/enrich/v1"
)

func RegisterRestRoutes(ctx context.Context, r *gin.Engine, grpcClients *grpc_client.Clients, services *service.Services, cache *commoncaches.Cache) {
	registerPublicRoutes(ctx, r, services)
	registerStreamRoutes(ctx, r, services, cache)
	registerOutreachRoutes(ctx, r, services, cache)
	registerCustomerBaseRoutes(ctx, r, services, grpcClients, cache)
	registerVerifyRoutes(ctx, r, services, cache)
	registerEnrichRoutes(ctx, r, services, cache)
}

func registerPublicRoutes(ctx context.Context, r *gin.Engine, services *service.Services) {
	// Redirect to pay invoice link
	r.GET("/invoice/:invoiceId/pay",
		tracing.TracingEnhancer(ctx, "GET:/invoice/:invoiceId/pay"),
		rest.RedirectToPayInvoice(services))
}

func registerEnrichRoutes(ctx context.Context, r *gin.Engine, services *service.Services, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/person", enrichV1Path), services, cache, rest.EnrichPerson(services))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/person/results/:id", enrichV1Path), services, cache, rest.EnrichPersonCallback(services))
}

func registerVerifyRoutes(ctx context.Context, r *gin.Engine, services *service.Services, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/email", verifyV1Path), services, cache, rest.VerifyEmailAddress(services))
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/email/bulk", verifyV1Path), services, cache, rest.BulkUploadEmailsForVerification(services))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/email/bulk/results/:requestId", verifyV1Path), services, cache, rest.GetBulkEmailVerificationResults(services))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/email/bulk/results/:requestId/download", verifyV1Path), services, cache, rest.DownloadBulkEmailVerificationResults(services))
	setupRestRoute(ctx, r, "GET", fmt.Sprintf("%s/ip", verifyV1Path), services, cache, rest.IpIntelligence(services))
}

func registerCustomerBaseRoutes(ctx context.Context, r *gin.Engine, services *service.Services, grpcClients *grpc_client.Clients, cache *commoncaches.Cache) {
	registerOrganizationRoutes(ctx, r, services, grpcClients, cache)
}

func registerOrganizationRoutes(ctx context.Context, r *gin.Engine, services *service.Services, grpcClients *grpc_client.Clients, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/organization", customerBaseV1Path), services, cache, rest.CreateOrganization(services, grpcClients))
}

func registerOutreachRoutes(ctx context.Context, r *gin.Engine, services *service.Services, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/track/email", outreachV1Path), services, cache, rest.GenerateEmailTrackingUrls(services))
}

func setupRestRoute(ctx context.Context, r *gin.Engine, method, path string, services *service.Services, cache *commoncaches.Cache, handler gin.HandlerFunc) {
	r.Handle(method, path,
		tracing.TracingEnhancer(ctx, method+":"+path),
		security.ApiKeyCheckerHTTP(services.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository, services.Repositories.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_API, security.WithCache(cache)),
		enrichContextMiddleware(),
		cosHandler.StatsSuccessHandler(method+":"+path, services),
		handler)
}

func registerStreamRoutes(ctx context.Context, r *gin.Engine, serviceContainer *service.Services, cache *commoncaches.Cache) {
	r.GET("/stream/organizations-cache",
		tracing.TracingEnhancer(ctx, "GET:/stream/organizations-cache"),
		apiKeyCheckerHTTPMiddleware(serviceContainer.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository, serviceContainer.Repositories.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_API, security.WithCache(cache)),
		tenantUserContextEnhancerMiddleware(security.USERNAME_OR_TENANT, serviceContainer.Repositories.Neo4jRepositories, security.WithCache(cache)),
		rest.OrganizationsCacheHandler(serviceContainer))
	r.GET("/stream/organizations-cache-diff",
		tracing.TracingEnhancer(ctx, "GET:/stream/organizations-cache-diff"),
		apiKeyCheckerHTTPMiddleware(serviceContainer.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository, serviceContainer.Repositories.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_API, security.WithCache(cache)),
		tenantUserContextEnhancerMiddleware(security.USERNAME_OR_TENANT, serviceContainer.Repositories.Neo4jRepositories, security.WithCache(cache)),
		rest.OrganizationsPatchesCacheHandler(serviceContainer))
}
