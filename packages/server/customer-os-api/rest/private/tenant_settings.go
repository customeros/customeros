package private

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/security"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	commonUtils "github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

func CreateOrganizationStage(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := commonUtils.GetContextWithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		tenant, _ := c.Get(security.KEY_TENANT_NAME)
		organizationStageId := c.Param("id")

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c, "/tenant/settings/organizationStage/"+organizationStageId, c.Request.Header)
		defer span.Finish()

		var requestData entity.TenantSettingsOpportunityStage
		if err := c.BindJSON(&requestData); err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
			})
			return
		}

		opportunityStage, err := s.Repositories.PostgresRepositories.TenantSettingsOpportunityStageRepository.GetById(ctx, tenant.(string), organizationStageId)
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if opportunityStage == nil {
			tracing.TraceErr(span, err)
			c.JSON(404, gin.H{"error": "Opportunity stage not found"})
			return
		}

		opportunityStage.Label = requestData.Label
		opportunityStage.Order = requestData.Order
		opportunityStage.Visible = requestData.Visible

		opportunityStage, err = s.Repositories.PostgresRepositories.TenantSettingsOpportunityStageRepository.Store(ctx, *opportunityStage)
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, opportunityStage)
	}
}

func GetAPIKey(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "GetApiKey")
		defer span.Finish()

		tenantValue, _ := c.Get(security.KEY_TENANT_NAME)
		tenant := tenantValue.(string)
		tracing.TagTenant(span, tenant)

		apiKey, err := s.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository.GetFirstApiKeyForTenant(ctx, tenant)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "GetFirstApiKeyForTenant"))
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if apiKey == nil {
			err = s.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository.CreateApiKey(ctx, tenant)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "CreateApiKey"))
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			apiKey, err = s.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository.GetFirstApiKeyForTenant(ctx, tenant)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "GetFirstApiKeyForTenant"))
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(200, apiKey.Key)
	}
}
