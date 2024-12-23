package reveal

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-postgres-repository/entity"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
)

type TrackerRequest struct {
	Domain string `json:"domain"`
}

type TrackerResponse struct {
	enum.BaseResponse
	Tracker TrackerRecord `json:"tracker"`
}

type TrackerRecord struct {
	ID        string `json:"id"`
	Domain    string `json:"domain"`
	CreatedAt string `json:"createdAt"`
	Code      string `json:"code"`
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

func ProvisionTracker(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "Reveal.ProvisionTracker", c.Request.Header)
		defer span.Finish()
		tracing.TagComponentRest(span)

		tenant := rest.ValidateTenant(c, ctx, span)
		if tenant == "" {
			rest.SendError(c, span, http.StatusUnauthorized, enum.ErrInvalidAPIKey)
			return
		}

		payload, err := getTrackerRequestPayload(c)
		if err != nil {
			switch {
			case err.Error() == "No domain":
				rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("domain not provided"))
				return
			default:
				rest.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("Could not parse request"))
				return
			}
		}

		// whitelist domain
		whitelist := entity.TrackingAllowedOrigin{
			Tenant:  tenant,
			Origin:  payload.Domain,
			Enabled: true,
		}

		record, err := s.Repositories.PostgresRepositories.TrackingAllowedOriginRepository.Create(ctx, whitelist)
		if err != nil || record == nil {
			rest.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("could not create tracker"))
			return
		}

		// return tracker

		c.JSON(http.StatusOK, TrackerResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Tracker: TrackerRecord{
				ID:        record.ID,
				Domain:    record.Domain,
				CreatedAt: record.CreatedAt,
				Code:      TrackerScript,
			},
		})
	}
}

func getTrackerRequestPayload(c *gin.Context) (TrackerRequest, error) {
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

	domain := strings.TrimPrefix(req.Domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "www.")
	domain = strings.Trim(domain, "/")

	req.Domain = fmt.Sprintf("https://%s", domain)

	return req, nil
}
