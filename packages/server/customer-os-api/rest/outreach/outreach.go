// @openapi 3.0.0
package outreach

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"

	"github.com/customeros/customeros/packages/server/customer-os-api/enum"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

// EmailTrackingRequest represents the request for generating tracking URLs
// @Description Request payload for generating email tracking URLs and pixels
type EmailTrackingRequest struct {
	// Domain to use for tracking URLs (optional, system default used if not provided)
	// required: false
	// format: hostname
	// example: track.example.com
	TrackerDomain string `json:"trackerDomain"`

	// Unique identifier for the email campaign
	// required: true
	// example: camp_123456
	CampaignId string `json:"campaignId"`

	// Unique identifier for the message (optional, generated if not provided)
	// required: false
	// example: msg_123456
	MessageId string `json:"messageId"`

	// Unique identifier for the recipient
	// required: true
	// example: recipient_123456
	RecipientId string `json:"recipientId"`

	// Enable open tracking via pixel
	// required: false
	// default: false
	TrackOpens bool `json:"trackOpens"`

	// Enable click tracking for links
	// required: false
	// default: false
	TrackClicks bool `json:"trackClicks"`

	// Enable unsubscribe link generation
	// required: false
	// default: false
	GenerateUnsubscribeLink bool `json:"generateUnsubscribeLink"`

	// List of URLs to be tracked
	// required: false
	// example: ["https://example.com/page1", "https://example.com/page2"]
	Links []string `json:"links"`

	// URL for unsubscribe page
	// required: false
	// format: uri
	// example: https://example.com/unsubscribe
	UnsubscribeLink string `json:"unsubscribeLink"`
}

// EmailTrackingResponse represents the response containing generated tracking URLs
// @Description Response containing generated tracking URLs and pixel
type EmailTrackingResponse struct {
	// Operation status
	// required: true
	// enum: success
	Status string `json:"status" example:"success"`

	// Generated or provided message ID
	// required: true
	TrackingId string `json:"trackingId" example:"msg_123456"`

	// URL for tracking pixel (only if TrackOpens is true)
	// required: false
	// format: uri
	TrackingPixel string `json:"trackingPixel,omitempty" example:"https://track.example.com/p/abc123"`

	// List of original and tracked URLs (only if TrackClicks is true)
	// required: false
	TrackedLinks []TrackedLink `json:"trackedLinks,omitempty"`

	// Generated unsubscribe link (only if GenerateUnsubscribeLink is true)
	// required: false
	// format: uri
	UnsubscribeLink string `json:"unsubscribeLink,omitempty" example:"https://track.example.com/u/abc123"`
}

// TrackedLink represents an original URL and its tracked version
// @Description Pair of original and tracking-enabled URLs
type TrackedLink struct {
	// Original URL before tracking
	// required: true
	// format: uri
	Original string `json:"original" example:"https://example.com/page1"`

	// Generated tracking URL
	// required: true
	// format: uri
	Tracked string `json:"tracked" example:"https://track.example.com/r/abc123"`
}

// GenerateEmailTrackingUrls generates tracking URLs for email campaigns
// @Summary Generate email tracking URLs
// @Description Generates tracking pixels, click tracking URLs, and unsubscribe links for email campaigns
// @Tags Outreach
// @Accept json
// @Produce json
// @Param body body EmailTrackingRequest true "Tracking URL generation configuration"
// @Success 200 {object} EmailTrackingResponse "Tracking URLs generated successfully"
// @Failure 400 {object} handlers.ErrorResponse "Invalid request body"
// @Failure 401 {object} handlers.ErrorResponse "Unauthorized - Missing or invalid API key"
// @Failure 500 {object} handlers.ErrorResponse "Internal server error"
// @Router /outreach/v1/email/tracking [post]
// @Security ApiKeyAuth
func GenerateEmailTrackingUrls(services *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "GenerateEmailTrackingUrls", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)
		tracing.TagTenant(span, common.GetTenantFromContext(ctx))

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			handlers.SendError(c, span, http.StatusUnauthorized, enum.ErrUnauthorized)
			return
		}

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
			log.Error(ctx, "Invalid request body", err)
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Invalid request body"))
			return
		}
		tracing.LogObjectAsJson(span, "request", request)

		trackerDomain := request.TrackerDomain
		if trackerDomain == "" {
			trackerDomain = services.Cfg.App.TrackingPublicUrl
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
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Error generating email open tracker"))
			return
		}

		// Generate tracked links
		var trackedLinks []map[string]string
		for _, redirectUrl := range request.Links {
			trackedUrl, _, _, err := services.CommonServices.EmailingService.GenerateEmailLinkUrl(ctx, tenant, trackerDomain, redirectUrl, messageId, request.CampaignId, request.RecipientId, request.TrackClicks)
			if err != nil {
				tracing.TraceErr(span, err)
				log.Error(ctx, "Error generating tracked link", err)
				handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Error generating tracking links"))
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
				handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Error generating unsubscribe link"))
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
