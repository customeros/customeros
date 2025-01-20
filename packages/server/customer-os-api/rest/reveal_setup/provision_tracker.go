package reveal

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type WebTrackerHandler struct {
	services        *cosapi_services.Services
	responseHandler *response.Response
}

func NewWebTrackerHandler(services *cosapi_services.Services, responseHandler *response.Response) *WebTrackerHandler {
	return &WebTrackerHandler{
		services:        services,
		responseHandler: responseHandler,
	}
}

type TrackerRequest struct {
	Domain string `json:"domain"`
}

type TrackerResponse struct {
	Tracker TrackerRecord `json:"tracker"`
}

type TrackerRecord struct {
	ID        string    `json:"id"`
	Domain    string    `json:"domain"`
	CreatedAt time.Time `json:"createdAt"`
	Code      string    `json:"code,omitempty"`
	Status    string    `json:"status,omitempty"`
}

const TrackerScript = `<!-- CustomerOS Visitor Reveal --> 
<script data-customeros>
(function(c, u, s, t, o, m, e, r, O, S) {
   var customerOS = document.createElement(s);
   customerOS.src = u;
   customerOS.async = true;
   (document.body || document.head).appendChild(customerOS);
})(window, "https://app.customeros.ai/analytics-0.1.js", "script");
</script>`

func (h *WebTrackerHandler) ProvisionTracker() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Reveal.ProvisionTracker", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := common.GetTenantFromContext(ctx)
		if tenant == "" {
			h.responseHandler.HandleError(c, http.StatusNotFound, nil)
			return
		}

		payload, err := h.getTrackerRequestPayload(c)
		if err != nil {
			switch {
			case err.Error() == "No domain":
				message := "domain not provided"
				h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
				return
			case err.Error() == "Invalid domain":
				message := "domain is not valid"
				h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
				return
			default:
				message := "Could not parse request"
				h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
				return
			}
		}

		// whitelist domain
		whitelist := postgres_entity.TrackingAllowedOrigin{
			Tenant:  tenant,
			Origin:  payload.Domain,
			Enabled: true,
		}

		record, err := h.services.Repositories.PostgresRepositories.TrackingAllowedOriginRepository.Create(ctx, whitelist)
		if err != nil || record == nil {
			message := "Could not create tracker"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		// return tracker

		h.responseHandler.HandleSuccess(c, TrackerResponse{
			Tracker: TrackerRecord{
				ID:        record.ID,
				Domain:    record.Origin,
				CreatedAt: record.CreatedAt,
				Code:      TrackerScript,
			},
		})
	}
}

func (h *WebTrackerHandler) getTrackerRequestPayload(c *gin.Context) (TrackerRequest, error) {
	span, _ := tracing.StartTracerSpan(c.Request.Context(), "Reveal.getTrackerRequestPayload")
	defer span.Finish()
	tracing.TagComponentRest(span)

	var req TrackerRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return req, err
	}

	if req.Domain == "" {
		return req, errors.New("No domain")
	}

	req.Domain = cleanDomain(req.Domain)

	// if !isValidDomain(ctx, req.Domain) {
	// 	return req, errors.New("Invalid domain")
	// }

	return req, nil
}

func (h *WebTrackerHandler) isValidDomain(ctx context.Context, domain string) bool {
	span, ctx := tracing.StartTracerSpan(ctx, "Reveal.isValidDomain")
	defer span.Finish()
	tracing.TagComponentRest(span)

	_, err := net.LookupHost(domain)
	if err != nil {
		tracing.TraceErr(span, err)
		return false
	}
	return true
}

func cleanDomain(domain string) string {
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "www.")
	domain = strings.Trim(domain, "/")

	return fmt.Sprintf("https://%s", domain)
}
