package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/service"
)

func enrichPerson(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "enrichPerson")
		defer span.Finish()

		var request model.EnrichPersonRequest

		if err := c.BindJSON(&request); err != nil {
			tracing.TraceErr(span, err)
			services.Logger.Errorf("Fail reading request: %v", err.Error())
			c.JSON(http.StatusBadRequest, model.EnrichPersonScrapinResponse{
				Status:      "error",
				Message:     "Invalid request body",
				PersonFound: false,
			})
			return
		}
		request.Normalize()

		tracing.LogObjectAsJson(span, "request", request)

		// validate mandatory parameters
		if request.LinkedinUrl == "" && request.Email == "" && request.Domain == "" {
			tracing.TraceErr(span, errors.New("Missing linkedin, email and domain parameters"))
			services.Logger.Errorf("Missing linkedin, email and domain parameters")
			c.JSON(http.StatusBadRequest, model.EnrichPersonScrapinResponse{
				Status:      "error",
				Message:     "Missing linkedin, email and domain parameters",
				PersonFound: false,
			})
			return
		}

		var scrapinRecordId uint64
		var enrichPersonData *model.EnrichedPersonData

		// Step 1 - Scrapin by linked in url
		if request.LinkedinUrl != "" {
			recordId, response, err := services.ScrapeInService.ScrapInPersonProfile(ctx, request.LinkedinUrl)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "ScrapInPersonProfile"))
				c.JSON(http.StatusInternalServerError, model.EnrichPersonScrapinResponse{
					Status:      "error",
					Message:     "Internal server error",
					PersonFound: false,
				})
				return
			}
			enrichPersonData = &model.EnrichedPersonData{
				PersonProfile: response,
			}
			scrapinRecordId = recordId
		}

		foundByLinkedInUrl := enrichPersonData != nil && enrichPersonData.PersonProfile != nil && enrichPersonData.PersonProfile.Person != nil

		// Step 2 - Search by email, domain and company name
		if !foundByLinkedInUrl && (request.Email != "" || request.Domain != "" || request.CompanyName != "") {
			recordId, response, err := services.ScrapeInService.ScrapInSearchPerson(ctx, request.Email, request.FirstName, request.LastName, request.Domain, request.CompanyName)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "ScrapInSearchPerson"))
				c.JSON(http.StatusInternalServerError, model.EnrichPersonScrapinResponse{
					Status:      "error",
					Message:     "Internal server error",
					PersonFound: false,
				})
				return
			}
			enrichPersonData = &model.EnrichedPersonData{
				PersonProfile: response,
			}
			scrapinRecordId = recordId
		}

		c.JSON(http.StatusOK, model.EnrichPersonScrapinResponse{
			Status:      "success",
			RecordId:    scrapinRecordId,
			PersonFound: enrichPersonData.PersonProfile != nil && enrichPersonData.PersonProfile.Person != nil,
			Data:        enrichPersonData,
		})
	}
}
