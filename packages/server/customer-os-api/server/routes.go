package server

import (
	"context"
	"github.com/gin-gonic/gin"
	cosHandler "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	commoncaches "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/grpc_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
)

func RegisterRoutes(ctx context.Context, r *gin.Engine, grpcClients *grpc_client.Clients, serviceContainer *service.Services, cache *commoncaches.Cache) {
	registerWhoamiRoutes(ctx, r, serviceContainer, cache)
	registerStreamRoutes(ctx, r, serviceContainer, cache)
}

func registerWhoamiRoutes(ctx context.Context, r *gin.Engine, serviceContainer *service.Services, cache *commoncaches.Cache) {
	r.GET("/whoami",
		cosHandler.TracingEnhancer(ctx, "/whoami"),
		security.ApiKeyCheckerHTTP(serviceContainer.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository, serviceContainer.Repositories.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_API, security.WithCache(cache)),
		rest.WhoamiHandler(serviceContainer))
}

func registerStreamRoutes(ctx context.Context, r *gin.Engine, serviceContainer *service.Services, cache *commoncaches.Cache) {
	r.GET("/stream/organizations-cache",
		cosHandler.TracingEnhancer(ctx, "/stream/organizations-cache"),
		apiKeyCheckerHTTPMiddleware(serviceContainer.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository, serviceContainer.Repositories.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_API, security.WithCache(cache)),
		tenantUserContextEnhancerMiddleware(security.USERNAME_OR_TENANT, serviceContainer.Repositories.Neo4jRepositories, security.WithCache(cache)),
		rest.OrganizationsCacheHandler(serviceContainer))
	r.GET("/stream/organizations-cache-diff",
		cosHandler.TracingEnhancer(ctx, "/stream/organizations-cache-diff"),
		apiKeyCheckerHTTPMiddleware(serviceContainer.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository, serviceContainer.Repositories.PostgresRepositories.AppKeyRepository, security.CUSTOMER_OS_API, security.WithCache(cache)),
		tenantUserContextEnhancerMiddleware(security.USERNAME_OR_TENANT, serviceContainer.Repositories.Neo4jRepositories, security.WithCache(cache)),
		rest.OrganizationsPatchesCacheHandler(serviceContainer))
}
