package route

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	commonservice "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service"
	"github.com/openline-ai/openline-customer-os/packages/server/email-tracking-api/handler"
	"github.com/openline-ai/openline-customer-os/packages/server/email-tracking-api/logger"
	"github.com/openline-ai/openline-customer-os/packages/server/email-tracking-api/tracing"
	"net/http"
)

type LinkData struct {
	MessageID   string `json:"mid"`
	LinkID      string `json:"lid"`
	RedirectURL string `json:"url"`
	Campaign    string `json:"ca"`
}

func AddEmailLinkRoute(ctx context.Context, route *gin.Engine, log logger.Logger, commonServices *commonservice.Services) {
	route.GET("/l",
		handler.TracingEnhancer(ctx, "/l"),
		handleRequest(ctx, commonServices, log))
}

func handleRequest(ctx context.Context, commonServices *commonservice.Services, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "HandleEmailLink", c.Request.Header)
		defer span.Finish()

		// Extract the 'c' query parameter
		hashedParam := c.Query("c")
		if hashedParam == "" {
			c.String(http.StatusBadRequest, "Missing required parameter")
			return
		}

		// Decode the Base64-encoded JSON
		jsonData, err := base64.URLEncoding.DecodeString(hashedParam)
		if err != nil {
			tracing.TraceErr(span, err)
			log.Error(ctx, "Failed to decode parameter", err)
			c.String(http.StatusBadRequest, "Invalid parameter encoding")
			return
		}

		// Parse the JSON data
		var linkData LinkData
		if err := json.Unmarshal(jsonData, &linkData); err != nil {
			tracing.TraceErr(span, err)
			log.Error(ctx, "Failed to parse link data", err)
			c.String(http.StatusBadRequest, "Invalid parameter format")
			return
		}

		// Implement all below

		//// Check if the message ID and link ID are valid (present in DB)
		//valid, err := commonServices.ValidateLinkIDs(ctx, linkData.MessageID, linkData.LinkID)
		//if err != nil {
		//	log.Error(ctx, "Error validating link IDs", err)
		//	c.String(http.StatusInternalServerError, "An error occurred")
		//	return
		//}
		//if !valid {
		//	c.String(http.StatusGone, "This link is no longer valid")
		//	return
		//}
		//
		//// Get IP address and email address
		//ipAddress := c.ClientIP()
		//emailAddress, err := commonServices.GetEmailFromMessageID(ctx, linkData.MessageID)
		//if err != nil {
		//	log.Error(ctx, "Error retrieving email address", err)
		//	// Continue processing even if email retrieval fails
		//}
		//
		//// Store click data
		//err = commonServices.StoreClickData(ctx, linkData.MessageID, linkData.LinkID, linkData.Campaign, ipAddress, emailAddress)
		//if err != nil {
		//	log.Error(ctx, "Error storing click data", err)
		//	// Continue processing even if storage fails
		//}

		// Redirect to the specified URL
		c.Redirect(http.StatusFound, linkData.RedirectURL)
	}
}
