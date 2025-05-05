package private

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/services/security"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	commonUtils "github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

func CreateOrganizationStage(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := commonUtils.GetContextWithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		tenant, _ := c.Get(security.KEY_TENANT_NAME)
		organizationStageId := c.Param("id")

		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "CreateOrganizationStage")
		defer spans.Finish()

		var requestData postgres_entity.TenantSettingsOpportunityStage
		if err := c.BindJSON(&requestData); err != nil {
			spans.TraceError(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": fmt.Sprintf("unable to parse json: %v", err.Error()),
			})
			return
		}

		opportunityStage, err := s.Repositories.PostgresRepositories.TenantSettingsOpportunityStageRepository.GetById(ctx, tenant.(string), organizationStageId)
		if err != nil {
			spans.TraceError(err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if opportunityStage == nil {
			spans.TraceError(err)
			c.JSON(404, gin.H{"error": "Opportunity stage not found"})
			return
		}

		opportunityStage.Label = requestData.Label
		opportunityStage.Order = requestData.Order
		opportunityStage.Visible = requestData.Visible

		opportunityStage, err = s.Repositories.PostgresRepositories.TenantSettingsOpportunityStageRepository.Store(ctx, *opportunityStage)
		if err != nil {
			spans.TraceError(err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, opportunityStage)
	}
}

func GetAPIKey(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		spans, ctx := telemetry.StartRestSpan(c.Request.Context(), "TenantSettings.GetAPIKey")
		defer spans.Finish()

		tenantValue, _ := c.Get(security.KEY_TENANT_NAME)
		tenant := tenantValue.(string)
		spans.TagTenant(tenant)

		apiKey, err := GetApiKeyForTenant(ctx, s, tenant)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()}); return
		}

		c.JSON(200, apiKey)
	}
}

func GetApiKeyForTenant(ctx context.Context, s *cosapi_services.Services, tenant string) (string, error) {
	spans, ctx := telemetry.StartRestSpan(ctx, "GetApiKeyForTenant")
	defer spans.Finish()

	apiKey, err := s.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository.GetFirstApiKeyForTenant(ctx, tenant)
	if err != nil {
		spans.TraceError(errors.Wrap(err, "GetFirstApiKeyForTenant"))
		return "", err
	}
	if apiKey == nil {
		err = s.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository.CreateApiKey(ctx, tenant)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "CreateApiKey"))
			return "", err
		}
		apiKey, err = s.Repositories.PostgresRepositories.TenantWebhookApiKeyRepository.GetFirstApiKeyForTenant(ctx, tenant)
		if err != nil {
			spans.TraceError(errors.Wrap(err, "GetFirstApiKeyForTenant"))
			return "", err
		}
	}
	if apiKey == nil {
		spans.TraceError(errors.New("api key is nil"))
		return "", errors.New("api key not found")
	}
	return apiKey.Key, nil
}
