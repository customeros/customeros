package routes

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	commonUtils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/service"
)

func InitTenantSettingsRoutes(r *gin.Engine, ctx context.Context, services *service.Services) {
	r.POST("/tenant/settings/organizationStage/:id",
		security.TenantUserContextEnhancer(security.USERNAME, services.Repositories.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository, services.Repositories.PostgresRepositories.AppKeyRepository, security.SETTINGS_API),
		func(ginContext *gin.Context) {
			c, cancel := commonUtils.GetContextWithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			tenant, _ := ginContext.Get(security.KEY_TENANT_NAME)
			organizationStageId := ginContext.Param("id")

			ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/tenant/settings/organizationStage/"+organizationStageId, ginContext.Request.Header)
			defer span.Finish()

			var requestData entity.TenantSettingsOpportunityStage
			if err := ginContext.BindJSON(&requestData); err != nil {
				tracing.TraceErr(span, err)
				ginContext.JSON(http.StatusInternalServerError, gin.H{
					"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
				})
				return
			}

			opportunityStage, err := services.Repositories.PostgresRepositories.TenantSettingsOpportunityStageRepository.GetById(ctx, tenant.(string), organizationStageId)
			if err != nil {
				tracing.TraceErr(span, err)
				ginContext.JSON(500, gin.H{"error": err.Error()})
				return
			}

			if opportunityStage == nil {
				tracing.TraceErr(span, err)
				ginContext.JSON(404, gin.H{"error": "Opportunity stage not found"})
				return
			}

			opportunityStage.Label = requestData.Label
			opportunityStage.Order = requestData.Order
			opportunityStage.Visible = requestData.Visible

			opportunityStage, err = services.Repositories.PostgresRepositories.TenantSettingsOpportunityStageRepository.Store(ctx, *opportunityStage)
			if err != nil {
				tracing.TraceErr(span, err)
				ginContext.JSON(500, gin.H{"error": err.Error()})
				return
			}

			ginContext.JSON(200, opportunityStage)
		})
}
