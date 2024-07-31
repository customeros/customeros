package route

import (
	"context"
	"github.com/gin-gonic/gin"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/utils"
	postgresentity "github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"github.com/openline-ai/openline-customer-os/packages/server/email-tracking-api/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/email-tracking-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/email-tracking-api/tracing"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func AddEmailTrackRoute(ctx context.Context, route *gin.Engine, log logger.Logger, commonServices *commonservice.Services) {
	route.GET("/v1/l",
		handler.TracingEnhancer(ctx, "/v1/l"),
		trackLinkRequest(commonServices, log))
	route.GET("/v1/s",
		handler.TracingEnhancer(ctx, "/v1/s"),
		trackOpenRequest(commonServices, log))
	route.GET("/v1/u",
		handler.TracingEnhancer(ctx, "/v1/u"),
		trackUnsubscribeRequest(commonServices, log))
	route.POST("/v1/emailTracker",
		handler.TracingEnhancer(ctx, "/v1/emailTracker"),
		generateTrackingUrls(commonServices, log))
}

func trackLinkRequest(commonServices *commonservice.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "trackLinkRequest", c.Request.Header)
		defer span.Finish()

		// Extract the 'c' query parameter
		emailLookupId := c.Query("c")
		if emailLookupId == "" {
			c.String(http.StatusBadRequest, "Missing required parameter")
			return
		}

		// Check email lookup id
		emailLookup, err := commonServices.PostgresRepositories.EmailLookupRepository.GetById(ctx, emailLookupId)
		if err != nil {
			log.Error(ctx, "Error retrieving email lookup", err)
			c.String(http.StatusInternalServerError, "An error occurred")
			return
		}
		// if id not found, return 404
		if emailLookup == nil {
			c.String(http.StatusNotFound, "Not found")
			return
		}
		// if email lookup is not of expected type, return 400
		if emailLookup.Type != postgresentity.EmailLookupTypeLink {
			log.Error(ctx, "Email lookup is not of expected type")
			c.String(http.StatusNotFound, "Not found")
			return
		}

		if emailLookup.TrackClicks {
			// Get IP address and email address
			ipAddress := c.ClientIP()

			// Store click data
			_, err = commonServices.PostgresRepositories.EmailTrackingRepository.Register(ctx, postgresentity.EmailTracking{
				Tenant:      emailLookup.Tenant,
				MessageId:   emailLookup.MessageId,
				LinkId:      emailLookup.LinkId,
				EventType:   postgresentity.EmailTrackingEventTypeLinkClick,
				IP:          ipAddress,
				RecipientId: emailLookup.RecipientId,
				Campaign:    emailLookup.Campaign,
			})
			if err != nil {
				log.Error(ctx, "Error storing click data", err)
			}
		}

		// Redirect to the specified URL
		c.Redirect(http.StatusFound, ensureAbsoluteURL(emailLookup.RedirectUrl))
	}
}

func ensureAbsoluteURL(url string) string {
	// Check if the URL already has a scheme
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}

	// If not, prepend "https://"
	return "https://" + url
}

func trackOpenRequest(commonServices *commonservice.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "trackOpenRequest", c.Request.Header)
		defer span.Finish()

		// Extract the 'c' query parameter
		emailLookupId := c.Query("c")
		if emailLookupId == "" {
			err := errors.New("Missing required parameter")
			tracing.TraceErr(span, err)
			log.Error(ctx, err.Error())
			return
		}

		// Check email lookup id
		emailLookup, err := commonServices.PostgresRepositories.EmailLookupRepository.GetById(ctx, emailLookupId)
		if err != nil {
			log.Error(ctx, "Error retrieving email lookup", err)
			return
		}
		if emailLookup == nil {
			return
		}
		// if email lookup is not of expected type, return 400
		if emailLookup.Type != postgresentity.EmailLookupTypeLink {
			log.Error(ctx, "Email lookup is not of expected type")
			return
		}
		if emailLookup.TrackOpens {
			// Get IP address and email address
			ipAddress := c.ClientIP()

			// Store click data
			_, err = commonServices.PostgresRepositories.EmailTrackingRepository.Register(ctx, postgresentity.EmailTracking{
				Tenant:      emailLookup.Tenant,
				MessageId:   emailLookup.MessageId,
				RecipientId: emailLookup.RecipientId,
				Campaign:    emailLookup.Campaign,
				EventType:   postgresentity.EmailTrackingEventTypeOpen,
				IP:          ipAddress,
			})
			if err != nil {
				log.Error(ctx, "Error storing email open data", err)
			}
		}
	}
}

func trackUnsubscribeRequest(commonServices *commonservice.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "trackUnsubscribeRequest", c.Request.Header)
		defer span.Finish()

		// Extract the 'c' query parameter
		emailLookupId := c.Query("u")
		if emailLookupId == "" {
			c.String(http.StatusBadRequest, "Missing required parameter")
			return
		}

		// Check email lookup id
		emailLookup, err := commonServices.PostgresRepositories.EmailLookupRepository.GetById(ctx, emailLookupId)
		if err != nil {
			log.Error(ctx, "Error retrieving email lookup", err)
			c.String(http.StatusInternalServerError, "An error occurred")
			return
		}
		// if id not found, return 404
		if emailLookup == nil {
			c.String(http.StatusNotFound, "Not found")
			return
		}
		// if email lookup is not of expected type, return 400
		if emailLookup.Type != postgresentity.EmailLookupTypeUnsubscribe {
			log.Error(ctx, "Email lookup is not of expected type")
			c.String(http.StatusNotFound, "Not found")
			return
		}

		// Get IP address and email address
		ipAddress := c.ClientIP()

		// Store click data
		_, err = commonServices.PostgresRepositories.EmailTrackingRepository.Register(ctx, postgresentity.EmailTracking{
			Tenant:      emailLookup.Tenant,
			MessageId:   emailLookup.MessageId,
			LinkId:      emailLookup.LinkId,
			EventType:   postgresentity.EmailTrackingEventTypeUnsubscribe,
			IP:          ipAddress,
			RecipientId: emailLookup.RecipientId,
			Campaign:    emailLookup.Campaign,
		})
		if err != nil {
			log.Error(ctx, "Error storing unsubscribe data", err)
		}

		// Redirect to the specified URL
		c.Redirect(http.StatusFound, ensureAbsoluteURL(emailLookup.UnsubscribeUrl))
	}
}

func generateTrackingUrls(services *commonservice.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "generateTrackingUrls", c.Request.Header)
		defer span.Finish()

		// Define a struct for the request body
		var request struct {
			Tenant                  string   `json:"tenant"`
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

		// Validate required fields
		if request.Tenant == "" || request.TrackerDomain == "" {
			err := errors.New("Missing required parameters")
			tracing.TraceErr(span, err)
			log.Error(ctx, err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Missing required parameters"})
			return
		}

		messageId := request.MessageId
		if request.MessageId == "" {
			messageId = utils.GenerateRandomString(64)
		}

		// Generate tracking open URL
		trackedOpenUrl, _, err := services.EmailingService.GenerateEmailSpyPixelUrl(ctx, request.Tenant, request.TrackerDomain, messageId, request.CampaignId, request.RecipientId, request.TrackOpens)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Error(ctx, "Error generating spy pixel URL", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Error generating open url"})
			return
		}

		// Generate tracked links
		var trackedLinks []map[string]string
		for _, redirectUrl := range request.Links {
			trackedUrl, _, _, err := services.EmailingService.GenerateEmailLinkUrl(ctx, request.Tenant, request.TrackerDomain, redirectUrl, messageId, request.CampaignId, request.RecipientId, request.TrackClicks)
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
			unsubscribeUrl, _, err := services.EmailingService.GenerateEmailUnsubscribeUrl(ctx, request.Tenant, request.TrackerDomain, request.UnsubscribeLink, messageId, request.CampaignId, request.RecipientId)
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
			"status":    "success",
			"messageId": messageId,
		}
		if trackedOpenUrl != "" {
			response["trackingPixel"] = trackedOpenUrl
		}
		if len(trackedLinks) > 0 {
			response["trackedLinks"] = trackedLinks
		}
		if trackedUnsubscribeLink != "" {
			response["trackedUnsubscribeLink"] = trackedUnsubscribeLink
		}

		c.JSON(http.StatusOK, response)
	}
}
