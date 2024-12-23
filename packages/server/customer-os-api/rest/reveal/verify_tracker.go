package reveal

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"
	"golang.org/x/net/html"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

func VerifyTracker(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Reveal.VerifyTracker", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, enum.ErrInvalidAPIKey)
			return
		}

		domain, exists := c.GetQuery("domain")
		if !exists || domain == "" {
			rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("must provide domain"))
			return
		}

		domain = cleanDomain(domain)
		// if !isValidDomain(ctx, domain) {
		// 	rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("domain is not valid"))
		// 	return
		// }

		// lookup domain to ensure it exists in whitelist
		query := entity.TrackingAllowedOrigin{
			Origin:  domain,
			Enabled: true,
		}

		record, err := s.Repositories.PostgresRepositories.TrackingAllowedOriginRepository.Find(ctx, query)
		if err != nil {
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}
		if record == nil {
			rest.SendError(c, span, http.StatusNotFound, enum.ErrNotFound.WithMessage("domain not found"))
			return
		}

		isActive, err := verifyTrackerScript(ctx, domain)
		if err != nil {
			tracing.TraceErr(span, err)
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer)
			return
		}

		verified := "unverified"
		if isActive {
			verified = "verified"
		}

		c.JSON(http.StatusOK, TrackerResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Tracker: TrackerRecord{
				ID:        record.ID,
				Domain:    record.Origin,
				CreatedAt: record.CreatedAt,
				Status:    verified,
			},
		})
	}
}

func verifyTrackerScript(ctx context.Context, domain string) (bool, error) {
	span, ctx := tracing.StartTracerSpan(ctx, "Reveal.verifyTrackerScript")
	defer span.Finish()
	tracing.TagComponentRest(span)

	if domain == "" {
		return false, errors.New("domain is required")
	}

	doc, err := fetchAndParsePage(domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return false, err
	}

	return hasCustomerOSScript(doc), nil
}

func fetchAndParsePage(domain string) (*html.Node, error) {
	resp, err := http.Get(domain)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return html.Parse(bytes.NewReader(bodyBytes))
}

func hasCustomerOSScript(n *html.Node) bool {
	if isCustomerOSScript(n) {
		return true
	}

	if hasCustomerOSInGTM(n) {
		return true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if hasCustomerOSScript(c) {
			return true
		}
	}

	return false
}

func isCustomerOSScript(n *html.Node) bool {
	if n.Type != html.ElementNode || n.Data != "script" {
		return false
	}

	isCustomerOS := false
	hasCorrectSrc := false

	for _, attr := range n.Attr {
		if attr.Key == "data-customeros" {
			isCustomerOS = true
		}
		if attr.Key == "src" && attr.Val == "https://app.customeros.ai/analytics-0.1.js" {
			hasCorrectSrc = true
		}
	}

	return isCustomerOS && hasCorrectSrc
}

func hasCustomerOSInGTM(n *html.Node) bool {
	// Check if it's a script node
	if n.Type == html.ElementNode && n.Data == "script" {
		// Check script content
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				content := strings.TrimSpace(c.Data)
				return strings.Contains(content, "app.customeros.ai/analytics-0.1.js") ||
					strings.Contains(content, "customeros") ||
					strings.Contains(content, "gtm") && strings.Contains(content, "customeros")
			}
		}
	}
	return false
}
