package route

import (
	"context"
	"github.com/gin-gonic/gin"
	commoncaches "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/caches"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-platform-admin-api/constants"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-platform-admin-api/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-platform-admin-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-platform-admin-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-platform-admin-api/tracing"
	organizationpb "github.com/openline-ai/openline-customer-os/packages/server/events-processing-proto/gen/proto/go/api/grpc/v1/organization"
	"net/http"
	"sync"
)

func AddOrganizationRoutes(ctx context.Context, route *gin.Engine, services *service.Services, log logger.Logger, cache *commoncaches.Cache) {
	route.POST("/organization/refreshLastTouchpoint",
		handler.TracingEnhancer(ctx, "/organization/refreshLastTouchpoint"),
		commonservice.ApiKeyCheckerHTTP(services.CommonServices.CommonRepositories.AppKeyRepository, commonservice.PLATFORM_ADMIN_API, commonservice.WithCache(cache)),
		commonservice.TenantUserContextEnhancer(commonservice.USERNAME, services.CommonServices.CommonRepositories, commonservice.WithCache(cache)),
		refreshLastTouchpointHandler(services, log))
}

func refreshLastTouchpointHandler(services *service.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "RefreshLastTouchpoint", c.Request.Header)
		defer span.Finish()

		userId := c.Keys[commonservice.KEY_USER_ID].(string)

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
				}
			}(tenantName)
		}

		c.JSON(http.StatusOK, "OK")
	}
}
