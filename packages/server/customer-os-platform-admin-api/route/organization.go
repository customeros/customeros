package route

import (
	"context"
	commontracing "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	commoncaches "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-platform-admin-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-platform-admin-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-platform-admin-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-platform-admin-api/tracing"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
)

func AddOrganizationRoutes(ctx context.Context, route *gin.Engine, services *service.Services, log logger.Logger, cache *commoncaches.Cache) {
	route.POST("/organization/refreshLastTouchpoint",
		commontracing.TracingEnhancer(ctx, "/organization/refreshLastTouchpoint"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.PLATFORM_ADMIN_API, security.WithCache(cache)),
		security.TenantUserContextEnhancer(security.USERNAME_OR_TENANT, services.CommonServices.Neo4jRepositories, security.WithCache(cache)),
		refreshLastTouchpointHandler(services, log))
}

func refreshLastTouchpointHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RefreshLastTouchpoint", c.Request.Header)
		defer span.Finish()

		userId := ""
		if user, ok := c.Get(security.KEY_USER_ID); ok {
			userId = user.(string)
		}

		tenants, err := services.Repositories.TenantRepository.GetTenants(ctx)
		if err != nil {
			log.Error(ctx, err)
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}

		var wg sync.WaitGroup
		wg.Add(len(tenants))

		for _, tenantName := range tenants {

			go func(tenantName string) {
				defer wg.Done()

				orgsInTenantToRefresh, err := services.Repositories.OrganizationRepository.CountOrganizationsForLastTouchpointRefresh(ctx, tenantName)
				if err != nil {
					log.Error(ctx, err)
					return
				}

				var wgTenant sync.WaitGroup

				limit := 100
				for skip := 0; skip < int(orgsInTenantToRefresh); skip += limit {

					wgTenant.Add(1)

					go func(skip, limit int) {
						defer wgTenant.Done()

						orgs, err := services.Repositories.OrganizationRepository.GetOrganizationsForLastTouchpointRefresh(ctx, tenantName, skip, limit)
						if err != nil {
							log.Error(ctx, err)
							return
						}

						for _, orgId := range orgs {
							_, err := services.GrpcClients.OrganizationClient.RefreshLastTouchpoint(ctx, &organizationpb.OrganizationIdGrpcRequest{
								Tenant:         tenantName,
								OrganizationId: orgId,
								AppSource:      constants.AppSourceCustomerOsPlatformAdminApi,
								LoggedInUserId: userId,
							})
							if err != nil {
								log.Error(ctx, err)
								return
							}
						}
					}(skip, limit)

					time.Sleep(10 * time.Second)
				}
			}(tenantName)
		}

		c.JSON(http.StatusOK, "OK")
	}
}
