package routes

import (
	"context"
	"fmt"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	commonUtils "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/settings-api/service"
)

func InitTenantSettingsRoutes(ctx context.Context, r *gin.Engine, services *service.Services) {
	r.POST("/tenant/settings/organizationStage/:id",
		tracing.TracingEnhancer(ctx, "/enrichPerson"),
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API, security.WithCache(services.CommonServices.Cache)),
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

			opportunityStage, err := services.CommonServices.PostgresRepositories.TenantSettingsOpportunityStageRepository.GetById(ctx, tenant.(string), organizationStageId)
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

			opportunityStage, err = services.CommonServices.PostgresRepositories.TenantSettingsOpportunityStageRepository.Store(ctx, *opportunityStage)
			if err != nil {
				tracing.TraceErr(span, err)
				ginContext.JSON(500, gin.H{"error": err.Error()})
				return
			}

			ginContext.JSON(200, opportunityStage)
		})
	r.GET("/tenant/settings/apiKey",
		tracing.TracingEnhancer(ctx, "GET /tenant/settings/apiKey"),
		security.TenantUserContextEnhancer(security.USERNAME, services.CommonServices.Neo4jRepositories),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.SETTINGS_API, security.WithCache(services.CommonServices.Cache)),
		func(c *gin.Context) {
			span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "GetApiKey")
			defer span.Finish()

			tenantValue, _ := c.Get(security.KEY_TENANT_NAME)
			tenant := tenantValue.(string)
			tracing.TagTenant(span, tenant)

			apiKey, err := services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository.GetFirstApiKeyForTenant(ctx, tenant)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "GetFirstApiKeyForTenant"))
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			if apiKey == nil {
				err = services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository.CreateApiKey(ctx, tenant)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "CreateApiKey"))
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				apiKey, err = services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository.GetFirstApiKeyForTenant(ctx, tenant)
				if err != nil {
					tracing.TraceErr(span, errors.Wrap(err, "GetFirstApiKeyForTenant"))
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
			}

			c.JSON(200, apiKey.Key)
		})
}
