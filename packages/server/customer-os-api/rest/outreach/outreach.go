package restoutreach

import (
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/common"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	"net/http"
)

func GenerateEmailTrackingUrls(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GenerateEmailTrackingUrls", c.Request.Header)
		defer span.Finish()

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Missing tenant context"})
			return
		}
		span.SetTag(tracing.SpanTagTenant, common.GetTenantFromContext(ctx))

		log := services.Log

		// Define a struct for the request body
		var request struct {
			TrackerDomain           string   `json:"trackerDomain"`
			CampaignId              string   `json:"campaignId"`
			MessageId               string   `json:"messageId"`
			RecipientId             string   `json:"recipientId"`
			TrackOpens              bool     `json:"trackOpens"`
			TrackClicks             bool     `json:"trackClicks"`
			GenerateUnsubscribeLink bool     `json:"generateUnsubscribeLink"`
			Links                   []string `json:"links"`
			UnsubscribeLink         string   `json:"unsubscribeLink"`
		}

		// Bind the JSON request body to the struct
		if err := c.ShouldBindJSON(&request); err != nil {
			tracing.TraceErr(span, err)
			log.Error(ctx, "Invalid request body", err)
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request body"})
			return
		}
		tracing.LogObjectAsJson(span, "request", request)

		trackerDomain := request.TrackerDomain
		if trackerDomain == "" {
			trackerDomain = services.Cfg.AppConfig.TrackingPublicUrl
		}

		messageId := request.MessageId
		if request.MessageId == "" {
			messageId = utils.GenerateRandomString(64)
		}

		// Generate tracking open URL
		trackedOpenUrl, _, err := services.CommonServices.EmailingService.GenerateEmailSpyPixelUrl(ctx, tenant, trackerDomain, messageId, request.CampaignId, request.RecipientId, request.TrackOpens)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Error(ctx, "Error generating spy pixel URL", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error generating open url"})
			return
		}

		// Generate tracked links
		var trackedLinks []map[string]string
		for _, redirectUrl := range request.Links {
			trackedUrl, _, _, err := services.CommonServices.EmailingService.GenerateEmailLinkUrl(ctx, tenant, trackerDomain, redirectUrl, messageId, request.CampaignId, request.RecipientId, request.TrackClicks)
			if err != nil {
				tracing.TraceErr(span, err)
				log.Error(ctx, "Error generating tracked link", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error generating links"})
				return
			}
			trackedLinks = append(trackedLinks, map[string]string{
				"original": redirectUrl,
				"tracked":  trackedUrl,
			})
		}

		// Generate unsubscribe link
		trackedUnsubscribeLink := ""
		if request.GenerateUnsubscribeLink && request.UnsubscribeLink != "" {
			unsubscribeUrl, _, err := services.CommonServices.EmailingService.GenerateEmailUnsubscribeUrl(ctx, tenant, trackerDomain, request.UnsubscribeLink, messageId, request.CampaignId, request.RecipientId)
			if err != nil {
				tracing.TraceErr(span, err)
				log.Error(ctx, "Error generating unsubscribe URL", err)
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error generating unsubscribe link"})
				return
			}
			trackedUnsubscribeLink = unsubscribeUrl
		}

		// Prepare and send the response
		response := gin.H{
			"status":     "success",
			"trackingId": messageId,
		}
		if trackedOpenUrl != "" {
			response["trackingPixel"] = trackedOpenUrl
		}
		if len(trackedLinks) > 0 {
			response["trackedLinks"] = trackedLinks
		}
		if trackedUnsubscribeLink != "" {
			response["unsubscribeLink"] = trackedUnsubscribeLink
		}

		c.JSON(http.StatusOK, response)
	}
}
