package private

import (
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
	api_tenant_settings "github.com/customeros/customeros/packages/server/customer-os-api/services/tenant_settings"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgresentity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

type PrivateIntegrationHandler struct {
	services        *cosapi_services.Services
	responseHandler *response.Response
}

func NewPrivateIntegrationHandler(services *cosapi_services.Services, responseHandler *response.Response) *PrivateIntegrationHandler {
	return &PrivateIntegrationHandler{
		services:        services,
		responseHandler: responseHandler,
	}
}

func (h *PrivateIntegrationHandler) GetIntegrations() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "PrivateIntegrationHandler.GetIntegrations")
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenantName := c.Keys["TenantName"].(string)
		tenantIntegrationSettings, activeServices, err := h.services.TenantSettingsService.GetForTenant(tenantName)
		if err != nil {
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
			return
		}

		h.responseHandler.HandleSuccess(c, h.mapTenantSettingsEntityToDTO(tenantIntegrationSettings, activeServices))
	}
}

func (h *PrivateIntegrationHandler) CreateIntegration() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "PrivateIntegrationHandler.CreateIntegration")
		defer span.Finish()
		tracing.TagComponentRest(span)

		var request map[string]interface{}

		if err := c.BindJSON(&request); err != nil {
			tracing.TraceErr(span, err)
			c.AbortWithStatus(500) // todo
			return
		}

		tenantName := c.Keys["TenantName"].(string)

		tenantIntegrationSettings, activeServices, err := h.services.TenantSettingsService.SaveIntegrationData(tenantName, request)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		h.responseHandler.HandleSuccess(c, h.mapTenantSettingsEntityToDTO(tenantIntegrationSettings, activeServices))
	}
}

func (h *PrivateIntegrationHandler) DeleteIntegrations() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "PrivateIntegrationHandler.DeleteIntegration")
		defer span.Finish()

		tracing.TagComponentRest(span)
		identifier := c.Param("identifier")
		if identifier == "" {
			c.JSON(500, gin.H{"error": "integration identifier is empty"})
			return
		}
		tenantName := c.Keys["TenantName"].(string)

		data, activeServices, err := h.services.TenantSettingsService.ClearIntegrationData(tenantName, identifier)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		h.responseHandler.HandleSuccess(c, h.mapTenantSettingsEntityToDTO(data, activeServices))
	}
}

func (h *PrivateIntegrationHandler) mapTenantSettingsEntityToDTO(tenantSettings *postgresentity.TenantSettings, activeServices map[string]bool) *map[string]interface{} {

	responseMap := make(map[string]interface{})

	for service, isActive := range activeServices {
		if isActive {
			responseMap[service] = make(map[string]interface{})
			responseMap[service].(map[string]interface{})["state"] = "ACTIVE"
		}
	}

	if tenantSettings == nil {
		return &responseMap
	}

	if tenantSettings != nil && tenantSettings.SmartSheetId != nil && tenantSettings.SmartSheetAccessToken != nil {
		responseMap[api_tenant_settings.SERVICE_SMARTSHEET] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_SMARTSHEET].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.JiraAPIToken != nil && tenantSettings.JiraDomain != nil && tenantSettings.JiraEmail != nil {
		responseMap[api_tenant_settings.SERVICE_JIRA] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_JIRA].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.TrelloAPIToken != nil && tenantSettings.TrelloAPIKey != nil {
		responseMap[api_tenant_settings.SERVICE_TRELLO] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_TRELLO].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.AhaAPIUrl != nil && tenantSettings.AhaAPIKey != nil {
		responseMap[api_tenant_settings.SERVICE_AHA] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_AHA].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.AirtablePersonalAccessToken != nil {
		responseMap[api_tenant_settings.SERVICE_AIRTABLE] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_AIRTABLE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.AmplitudeSecretKey != nil && tenantSettings.AmplitudeAPIKey != nil {
		responseMap[api_tenant_settings.SERVICE_AMPLITUDE] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_AMPLITUDE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.AsanaAccessToken != nil {
		responseMap[api_tenant_settings.SERVICE_ASANA] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_ASANA].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.BatonAPIKey != nil {
		responseMap[api_tenant_settings.SERVICE_BATON] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_BATON].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.BabelforceRegionEnvironment != nil && tenantSettings.BabelforceAccessKeyId != nil && tenantSettings.BabelforceAccessToken != nil {
		responseMap[api_tenant_settings.SERVICE_BABELFORCE] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_BABELFORCE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.BigQueryServiceAccountKey != nil {
		responseMap[api_tenant_settings.SERVICE_BIGQUERY] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_BIGQUERY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.BraintreeEnvironment != nil && tenantSettings.BraintreeMerchantId != nil && tenantSettings.BraintreePublicKey != nil && tenantSettings.BraintreePrivateKey != nil {
		responseMap[api_tenant_settings.SERVICE_BRAINTREE] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_BRAINTREE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.CallRailAccount != nil && tenantSettings.CallRailApiToken != nil {
		responseMap[api_tenant_settings.SERVICE_CALLRAIL] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_CALLRAIL].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ChargebeeProductCatalog != nil && tenantSettings.ChargebeeApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_CHARGEBEE] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_CHARGEBEE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ChargifyApiKey != nil && tenantSettings.ChargifyDomain != nil {
		responseMap[api_tenant_settings.SERVICE_CHARGIFY] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_CHARGIFY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ClickUpApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_CLICKUP] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_CLICKUP].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.CloseComApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_CLOSECOM] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_CLOSECOM].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.CodaAuthToken != nil && tenantSettings.CodaDocumentId != nil {
		responseMap[api_tenant_settings.SERVICE_CODA] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_CODA].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ConfluenceApiToken != nil && tenantSettings.ConfluenceDomain != nil && tenantSettings.ConfluenceLoginEmail != nil {
		responseMap[api_tenant_settings.SERVICE_CONFLUENCE] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_CONFLUENCE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.CourierApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_COURIER] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_COURIER].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.CustomerIoApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_CUSTOMERIO] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_CUSTOMERIO].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.DatadogApiKey != nil && tenantSettings.DatadogApplicationKey != nil {
		responseMap[api_tenant_settings.SERVICE_DATADOG] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_DATADOG].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.DelightedApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_DELIGHTED] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_DELIGHTED].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.DixaApiToken != nil {
		responseMap[api_tenant_settings.SERVICE_DIXA] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_DIXA].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.DriftApiToken != nil {
		responseMap[api_tenant_settings.SERVICE_DRIFT] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_DRIFT].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.EmailOctopusApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_EMAILOCTOPUS] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_EMAILOCTOPUS].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FacebookMarketingAccessToken != nil {
		responseMap[api_tenant_settings.SERVICE_FACEBOOK_MARKETING] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_FACEBOOK_MARKETING].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FastbillApiKey != nil && tenantSettings.FastbillProjectId != nil {
		responseMap[api_tenant_settings.SERVICE_FASTBILL] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_FASTBILL].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FlexportApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_FLEXPORT] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_FLEXPORT].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FreshcallerApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_FRESHCALLER] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_FRESHCALLER].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FreshdeskApiKey != nil && tenantSettings.FreshdeskDomain != nil {
		responseMap[api_tenant_settings.SERVICE_FRESHDESK] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_FRESHDESK].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FreshsalesApiKey != nil && tenantSettings.FreshsalesDomain != nil {
		responseMap[api_tenant_settings.SERVICE_FRESHSALES] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_FRESHSALES].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.FreshserviceApiKey != nil && tenantSettings.FreshserviceDomain != nil {
		responseMap[api_tenant_settings.SERVICE_FRESHSERVICE] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_FRESHSERVICE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.GenesysRegion != nil && tenantSettings.GenesysClientId != nil && tenantSettings.GenesysClientSecret != nil {
		responseMap[api_tenant_settings.SERVICE_GENESYS] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_GENESYS].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.GitHubAccessToken != nil {
		responseMap[api_tenant_settings.SERVICE_GITHUB] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_GITHUB].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.GitLabAccessToken != nil {
		responseMap[api_tenant_settings.SERVICE_GITLAB] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_GITLAB].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.GoCardlessAccessToken != nil && tenantSettings.GoCardlessEnvironment != nil && tenantSettings.GoCardlessVersion != nil {
		responseMap[api_tenant_settings.SERVICE_GOCARDLESS] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_GOCARDLESS].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.GongApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_GONG] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_GONG].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.HarvestAccountId != nil && tenantSettings.HarvestAccessToken != nil {
		responseMap[api_tenant_settings.SERVICE_HARVEST] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_HARVEST].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.InsightlyApiToken != nil {
		responseMap[api_tenant_settings.SERVICE_INSIGHTLY] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_INSIGHTLY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.InstagramAccessToken != nil {
		responseMap[api_tenant_settings.SERVICE_INSTAGRAM] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_INSTAGRAM].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.InstatusApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_INSTATUS] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_INSTATUS].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.IntercomAccessToken != nil {
		responseMap[api_tenant_settings.SERVICE_INTERCOM] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_INTERCOM].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.KlaviyoApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_KLAVIYO] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_KLAVIYO].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.KustomerApiToken != nil {
		responseMap[api_tenant_settings.SERVICE_KUSTOMER] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_KUSTOMER].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.LookerClientId != nil && tenantSettings.LookerClientSecret != nil && tenantSettings.LookerDomain != nil {
		responseMap[api_tenant_settings.SERVICE_LOOKER] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_LOOKER].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.MailchimpApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_MAILCHIMP] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_MAILCHIMP].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.MailjetEmailApiKey != nil && tenantSettings.MailjetEmailApiSecret != nil {
		responseMap[api_tenant_settings.SERVICE_MAILJETEMAIL] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_MAILJETEMAIL].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.MarketoClientId != nil && tenantSettings.MarketoClientSecret != nil && tenantSettings.MarketoDomainUrl != nil {
		responseMap[api_tenant_settings.SERVICE_MARKETO] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_MARKETO].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.MicrosoftTeamsTenantId != nil && tenantSettings.MicrosoftTeamsClientId != nil && tenantSettings.MicrosoftTeamsClientSecret != nil {
		responseMap[api_tenant_settings.SERVICE_MICROSOFT_TEAMS] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_MICROSOFT_TEAMS].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.MondayApiToken != nil {
		responseMap[api_tenant_settings.SERVICE_MONDAY] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_MONDAY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.NotionInternalAccessToken != nil || (tenantSettings.NotionPublicClientId != nil && tenantSettings.NotionPublicClientSecret != nil && tenantSettings.NotionPublicAccessToken != nil) {
		responseMap[api_tenant_settings.SERVICE_NOTION] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_NOTION].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.OracleNetsuiteAccountId != nil && tenantSettings.OracleNetsuiteConsumerKey != nil && tenantSettings.OracleNetsuiteConsumerSecret != nil && tenantSettings.OracleNetsuiteTokenId != nil && tenantSettings.OracleNetsuiteTokenSecret != nil {
		responseMap[api_tenant_settings.SERVICE_ORACLE_NETSUITE] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_ORACLE_NETSUITE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.OrbApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_ORB] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_ORB].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.OrbitApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_ORBIT] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_ORBIT].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PagerDutyApikey != nil {
		responseMap[api_tenant_settings.SERVICE_PAGERDUTY] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_PAGERDUTY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PaypalTransactionClientId != nil && tenantSettings.PaypalTransactionSecret != nil {
		responseMap[api_tenant_settings.SERVICE_PAYPAL_TRANSACTION] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_PAYPAL_TRANSACTION].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PaystackSecretKey != nil {
		responseMap[api_tenant_settings.SERVICE_PAYSTACK] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_PAYSTACK].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PendoApiToken != nil {
		responseMap[api_tenant_settings.SERVICE_PENDO] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_PENDO].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PipedriveApiToken != nil {
		responseMap[api_tenant_settings.SERVICE_PIPEDRIVE] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_PIPEDRIVE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PlaidAccessToken != nil {
		responseMap[api_tenant_settings.SERVICE_PLAID] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_PLAID].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PlausibleSiteId != nil && tenantSettings.PlausibleApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_PLAUSIBLE] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_PLAUSIBLE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.PostHogApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_POSTHOG] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_POSTHOG].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.QualarooApiKey != nil && tenantSettings.QualarooApiToken != nil {
		responseMap[api_tenant_settings.SERVICE_QUALAROO] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_QUALAROO].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.QuickBooksClientId != nil && tenantSettings.QuickBooksClientSecret != nil && tenantSettings.QuickBooksRealmId != nil && tenantSettings.QuickBooksRefreshToken != nil {
		responseMap[api_tenant_settings.SERVICE_QUICKBOOKS] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_QUICKBOOKS].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.RechargeApiToken != nil {
		responseMap[api_tenant_settings.SERVICE_RECHARGE] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_RECHARGE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.RecruiteeCompanyId != nil && tenantSettings.RecruiteeApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_RECRUITEE] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_RECRUITEE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.RecurlyApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_RECURLY] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_RECURLY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.RetentlyApiToken != nil {
		responseMap[api_tenant_settings.SERVICE_RETENTLY] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_RETENTLY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.SalesloftApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_SALESLOFT] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_SALESLOFT].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.SendgridApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_SENDGRID] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_SENDGRID].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.SentryProject != nil && tenantSettings.SentryOrganization != nil && tenantSettings.SentryAuthenticationToken != nil {
		responseMap[api_tenant_settings.SERVICE_SENTRY] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_SENTRY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.SlackApiToken != nil && tenantSettings.SlackChannelFilter != nil {
		responseMap[api_tenant_settings.SERVICE_SLACK] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_SLACK].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.StripeAccountId != nil && tenantSettings.StripeSecretKey != nil {
		responseMap[api_tenant_settings.SERVICE_STRIPE] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_STRIPE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.SurveySparrowAccessToken != nil {
		responseMap[api_tenant_settings.SERVICE_SURVEYSPARROW] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_SURVEYSPARROW].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.SurveyMonkeyAccessToken != nil {
		responseMap[api_tenant_settings.SERVICE_SURVEYMONKEY] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_SURVEYMONKEY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.TalkdeskApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_TALKDESK] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_TALKDESK].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.TikTokAccessToken != nil {
		responseMap[api_tenant_settings.SERVICE_TIKTOK] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_TIKTOK].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.TodoistApiToken != nil {
		responseMap[api_tenant_settings.SERVICE_TODOIST] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_TODOIST].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.TypeformApiToken != nil {
		responseMap[api_tenant_settings.SERVICE_TYPEFORM] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_TYPEFORM].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.VittallyApiKey != nil {
		responseMap[api_tenant_settings.SERVICE_VITTALLY] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_VITTALLY].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.WrikeAccessToken != nil && tenantSettings.WrikeHostUrl != nil {
		responseMap[api_tenant_settings.SERVICE_WRIKE] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_WRIKE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.XeroClientId != nil && tenantSettings.XeroClientSecret != nil && tenantSettings.XeroTenantId != nil && tenantSettings.XeroScopes != nil {
		responseMap[api_tenant_settings.SERVICE_XERO] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_XERO].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ZendeskAPIKey != nil && tenantSettings.ZendeskSubdomain != nil && tenantSettings.ZendeskAdminEmail != nil {
		responseMap[api_tenant_settings.SERVICE_ZENDESK_SUPPORT] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_ZENDESK_SUPPORT].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ZendeskChatSubdomain != nil && tenantSettings.ZendeskChatAccessKey != nil {
		responseMap[api_tenant_settings.SERVICE_ZENDESK_CHAT] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_ZENDESK_CHAT].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ZendeskTalkSubdomain != nil && tenantSettings.ZendeskTalkAccessKey != nil {
		responseMap[api_tenant_settings.SERVICE_ZENDESK_TALK] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_ZENDESK_TALK].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ZendeskSellApiToken != nil {
		responseMap[api_tenant_settings.SERVICE_ZENDESK_SELL] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_ZENDESK_SELL].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ZendeskSunshineSubdomain != nil && tenantSettings.ZendeskSunshineApiToken != nil && tenantSettings.ZendeskSunshineEmail != nil {
		responseMap[api_tenant_settings.SERVICE_ZENDESK_SUNSHINE] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_ZENDESK_SUNSHINE].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.ZenefitsToken != nil {
		responseMap[api_tenant_settings.SERVICE_ZENEFITS] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_ZENEFITS].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && h.hasMixpanelKeys(tenantSettings) {
		responseMap[api_tenant_settings.SERVICE_MIXPANEL] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_MIXPANEL].(map[string]interface{})["state"] = "ACTIVE"
	}

	if tenantSettings != nil && tenantSettings.LinkedInCredential != nil {
		responseMap[api_tenant_settings.SERVICE_LINKEDIN] = make(map[string]interface{})
		responseMap[api_tenant_settings.SERVICE_LINKEDIN].(map[string]interface{})["state"] = "ACTIVE"
	}

	return &responseMap
}

func (h *PrivateIntegrationHandler) hasMixpanelKeys(tenantSettings *postgresentity.TenantSettings) bool {
	return tenantSettings.MixpanelUsername != nil || tenantSettings.MixpanelSecret != nil || tenantSettings.MixpanelProjectId != nil || tenantSettings.MixpanelProjectSecret != nil || tenantSettings.MixpanelProjectTimezone != nil || tenantSettings.MixpanelRegion != nil
}
