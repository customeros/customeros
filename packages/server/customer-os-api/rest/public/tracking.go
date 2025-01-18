package public

import (
	"net/http"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/customeros/mailsherpa/mailvalidate"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

func TrackLinkRequest(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "tracking.trackLinkRequest", c.Request.Header)
		defer span.Finish()

		// Extract the 'c' query parameter
		emailLookupId := c.Query("c")
		if emailLookupId == "" {
			c.String(http.StatusBadRequest, "Missing required parameter")
			return
		}

		// Check email lookup id
		emailLookup, err := s.Repositories.PostgresRepositories.EmailLookupRepository.GetById(ctx, emailLookupId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error retrieving email lookup"))
			c.String(http.StatusInternalServerError, "An error occurred")
			return
		}
		tracing.LogObjectAsJson(span, "emailLookup", emailLookup)

		// if id not found, return 404
		if emailLookup == nil {
			c.String(http.StatusNotFound, "Not found")
			return
		}
		// if email lookup is not of expected type, return 400
		if emailLookup.Type != postgres_entity.EmailLookupTypeLink {
			tracing.TraceErr(span, errors.Wrap(err, "Email lookup is not of expected type"))
			c.String(http.StatusNotFound, "Not found")
			return
		}

		if emailLookup.TrackClicks {
			// Save IP data
			ipAddress, err := saveIP(c, s, emailLookup)

			// Store click data
			_, err = s.Repositories.PostgresRepositories.EmailTrackingRepository.Register(ctx, postgres_entity.EmailTracking{
				Tenant:      emailLookup.Tenant,
				MessageId:   emailLookup.MessageId,
				LinkId:      emailLookup.LinkId,
				EventType:   postgres_entity.EmailTrackingEventTypeLinkClick,
				IP:          ipAddress,
				RecipientId: emailLookup.RecipientId,
				Campaign:    emailLookup.Campaign,
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Email storing click data"))
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

func TrackOpenRequest(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "tracking.trackOpenRequest", c.Request.Header)
		defer span.Finish()

		// Extract the 'c' query parameter
		emailLookupId := c.Query("c")
		span.LogFields(log.String("emailLookupId", emailLookupId))
		if emailLookupId == "" {
			err := errors.New("Missing required parameter")
			tracing.TraceErr(span, err)
			return
		}

		// Check email lookup id
		emailLookup, err := s.Repositories.PostgresRepositories.EmailLookupRepository.GetById(ctx, emailLookupId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error retrieving email lookup"))
			return
		}
		tracing.LogObjectAsJson(span, "emailLookup", emailLookup)

		if emailLookup == nil {
			return
		}
		// if email lookup is not of expected type, return 400
		if emailLookup.Type != postgres_entity.EmailLookupTypeSpyPixel {
			tracing.TraceErr(span, errors.Wrap(err, "Email lookup is not of expected type"))
			return
		}
		if emailLookup.TrackOpens {
			// Get IP address and email address
			ipAddress := c.ClientIP()

			// Store click data
			_, err = s.Repositories.PostgresRepositories.EmailTrackingRepository.Register(ctx, postgres_entity.EmailTracking{
				Tenant:      emailLookup.Tenant,
				MessageId:   emailLookup.MessageId,
				RecipientId: emailLookup.RecipientId,
				Campaign:    emailLookup.Campaign,
				EventType:   postgres_entity.EmailTrackingEventTypeOpen,
				IP:          ipAddress,
			})
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "Error storing email open data"))
			}
		}
	}
}

func TrackUnsubscribeRequest(s *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "tracking.trackUnsubscribeRequest", c.Request.Header)
		defer span.Finish()

		// Extract the 'c' query parameter
		emailLookupId := c.Query("u")
		if emailLookupId == "" {
			c.String(http.StatusBadRequest, "Missing required parameter")
			return
		}

		// Check email lookup id
		emailLookup, err := s.Repositories.PostgresRepositories.EmailLookupRepository.GetById(ctx, emailLookupId)
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error retrieving email lookup"))
			c.String(http.StatusInternalServerError, "An error occurred")
			return
		}
		tracing.LogObjectAsJson(span, "emailLookup", emailLookup)

		// if id not found, return 404
		if emailLookup == nil {
			c.String(http.StatusNotFound, "Not found")
			return
		}
		// if email lookup is not of expected type, return 400
		if emailLookup.Type != postgres_entity.EmailLookupTypeUnsubscribe {
			tracing.TraceErr(span, errors.Wrap(err, "Email lookup is not of expected type"))
			c.String(http.StatusNotFound, "Not found")
			return
		}

		// Get IP address and email address
		ipAddress := c.ClientIP()

		// Store click data
		_, err = s.Repositories.PostgresRepositories.EmailTrackingRepository.Register(ctx, postgres_entity.EmailTracking{
			Tenant:      emailLookup.Tenant,
			MessageId:   emailLookup.MessageId,
			LinkId:      emailLookup.LinkId,
			EventType:   postgres_entity.EmailTrackingEventTypeUnsubscribe,
			IP:          ipAddress,
			RecipientId: emailLookup.RecipientId,
			Campaign:    emailLookup.Campaign,
		})
		if err != nil {
			tracing.TraceErr(span, errors.Wrap(err, "Error storing unsubscribe data"))
		}

		// Redirect to the specified URL
		c.Redirect(http.StatusFound, ensureAbsoluteURL(emailLookup.UnsubscribeUrl))
	}
}

func saveIP(c *gin.Context, s *cosapi_services.Services, emailLookup *postgres_entity.EmailLookup) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "tracking.saveIP")
	defer span.Finish()
	originalIP := c.Request.Header["X-Original-Forwarded-For"][0]
	cloudflareIP := c.Request.Header["Cf-Connecting-Ip"][0]

	var clientIP string
	if cloudflareIP == "" && originalIP == "" {
		return "", nil
	}
	if cloudflareIP != "" {
		clientIP = cloudflareIP
	} else {
		clientIP = originalIP
	}

	emailMessage, err := s.Repositories.PostgresRepositories.EmailMessageRepository.GetByProviderMessageId(ctx, emailLookup.Tenant, emailLookup.MessageId)
	if err != nil {
		tracing.TraceErr(span, err)
		return clientIP, err
	}

	emailVerify := mailvalidate.ValidateEmailSyntax(emailMessage.From)

	details := postgres_entity.EnrichDetailsTracking{
		IP:             clientIP,
		CompanyDomain:  &emailVerify.Domain,
		CompanyWebsite: &emailVerify.Domain,
		SourceEmail:    &emailVerify.CleanEmail,
	}

	err = s.Repositories.PostgresRepositories.EnrichDetailsTrackingRepository.Save(ctx, details)
	if err != nil {
		tracing.TraceErr(span, err)
		return clientIP, err
	}
	return clientIP, nil
}
