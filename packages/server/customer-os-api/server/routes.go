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
)

const (
	outreachV1Path     = "/outreach/v1"
	customerBaseV1Path = "/customerbase/v1"
)

func RegisterRestRoutes(ctx context.Context, r *gin.Engine, grpcClients *grpc_client.Clients, services *service.Services, cache *commoncaches.Cache) {
	registerWhoamiRoutes(ctx, r, services, cache)
	registerStreamRoutes(ctx, r, services, cache)
	registerOutreachRoutes(ctx, r, services, cache)
	registerCustomerBaseRoutes(ctx, r, services, grpcClients, cache)
}

func registerCustomerBaseRoutes(ctx context.Context, r *gin.Engine, services *service.Services, grpcClients *grpc_client.Clients, cache *commoncaches.Cache) {
	registerOrganizationRoutes(ctx, r, services, grpcClients, cache)
}

func registerOrganizationRoutes(ctx context.Context, r *gin.Engine, services *service.Services, grpcClients *grpc_client.Clients, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/organizations", customerBaseV1Path), services, cache, rest.CreateOrganization(services, grpcClients))
}

func registerOutreachRoutes(ctx context.Context, r *gin.Engine, services *service.Services, cache *commoncaches.Cache) {
	setupRestRoute(ctx, r, "POST", fmt.Sprintf("%s/track/email", outreachV1Path), services, cache, rest.GenerateEmailTrackingUrls(services))
}

func setupRestRoute(ctx context.Context, r *gin.Engine, method, path string, services *service.Services, cache *commoncaches.Cache, handler gin.HandlerFunc) {
	r.Handle(method, path,
		cosHandler.TracingEnhancer(ctx, method+":"+path),
		security.ApiKeyCheckerHTTP(services.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository, services.Repositories.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_API, security.WithCache(cache)),
		enrichContextMiddleware(services),
		handler)
}

func registerWhoamiRoutes(ctx context.Context, r *gin.Engine, serviceContainer *service.Services, cache *commoncaches.Cache) {
	r.GET("/whoami",
		cosHandler.TracingEnhancer(ctx, "GET:/whoami"),
		security.ApiKeyCheckerHTTP(serviceContainer.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository, serviceContainer.Repositories.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_API, security.WithCache(cache)),
		rest.WhoamiHandler(serviceContainer))
}

func registerStreamRoutes(ctx context.Context, r *gin.Engine, serviceContainer *service.Services, cache *commoncaches.Cache) {
	r.GET("GET /stream/organizations-cache",
		cosHandler.TracingEnhancer(ctx, "GET:/stream/organizations-cache"),
		apiKeyCheckerHTTPMiddleware(serviceContainer.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository, serviceContainer.Repositories.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_API, security.WithCache(cache)),
		tenantUserContextEnhancerMiddleware(security.USERNAME_OR_TENANT, serviceContainer.Repositories.Neo4jRepositories, security.WithCache(cache)),
		rest.OrganizationsCacheHandler(serviceContainer))
	r.GET("GET /stream/organizations-cache-diff",
		cosHandler.TracingEnhancer(ctx, "GET:/stream/organizations-cache-diff"),
		apiKeyCheckerHTTPMiddleware(serviceContainer.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository, serviceContainer.Repositories.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_API, security.WithCache(cache)),
		tenantUserContextEnhancerMiddleware(security.USERNAME_OR_TENANT, serviceContainer.Repositories.Neo4jRepositories, security.WithCache(cache)),
		rest.OrganizationsPatchesCacheHandler(serviceContainer))
}
