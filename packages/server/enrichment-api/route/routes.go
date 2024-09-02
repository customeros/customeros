package route

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/service"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"net/http"
)

func RegisterRoutes(ctx context.Context, r *gin.Engine, services *service.Services) {
	r.GET("/enrichPerson",
		tracing.TracingEnhancer(ctx, "GET /enrichPerson"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.ENRICHMENT_API),
		enrichPerson(services))
	r.GET("/findWorkEmail",
		tracing.TracingEnhancer(ctx, "GET /findWorkEmail"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.ENRICHMENT_API),
		findWorkEmail(services))
}

func enrichPerson(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "enrichPerson")
		defer span.Finish()

		var request model.EnrichPersonRequest

		if err := c.BindJSON(&request); err != nil {
			tracing.TraceErr(span, err)
			services.Logger.Errorf("Fail reading request: %v", err.Error())
			c.JSON(http.StatusBadRequest, model.EnrichPersonResponse{
				Status:      "error",
				Message:     "Invalid request body",
				PersonFound: false,
			})
			return
		}
		request.Normalize()

		tracing.LogObjectAsJson(span, "request", request)

		// validate mandatory parameters
		if request.LinkedinUrl == "" && request.Email == "" {
			tracing.TraceErr(span, errors.New("Missing linkedin and email parameters"))
			services.Logger.Errorf("Missing linkedin and email parameters")
			c.JSON(http.StatusBadRequest, model.EnrichPersonResponse{
				Status:      "error",
				Message:     "Missing linkedin and email parameters",
				PersonFound: false,
			})
			return
		}

		var scrapinRecordId uint64
		var enrichPersonData *model.EnrichedPersonData

		// Step 1 - Scrapin by linked in url
		if request.LinkedinUrl != "" {
			recordId, response, err := services.PersonScrapeInService.ScrapInPersonProfile(ctx, request.LinkedinUrl)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "ScrapInPersonProfile"))
				c.JSON(http.StatusInternalServerError, model.EnrichPersonResponse{
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

		// Step 2 - Scrapin by email
		if !foundByLinkedInUrl && request.Email != "" {
			recordId, response, err := services.PersonScrapeInService.ScrapInSearchPerson(ctx, request.Email, request.FirstName, request.LastName, request.Domain)
			if err != nil {
				tracing.TraceErr(span, errors.Wrap(err, "ScrapInSearchPerson"))
				c.JSON(http.StatusInternalServerError, model.EnrichPersonResponse{
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

		c.JSON(http.StatusOK, model.EnrichPersonResponse{
			Status:      "success",
			RecordId:    scrapinRecordId,
			PersonFound: enrichPersonData.PersonProfile != nil && enrichPersonData.PersonProfile.Person != nil,
			Data:        enrichPersonData,
		})
	}
}

func findWorkEmail(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "findWorkEmail")
		defer span.Finish()

		var request model.FindWorkEmailRequest

		if err := c.BindJSON(&request); err != nil {
			services.Logger.Errorf("Fail reading request: %v", err.Error())
			c.JSON(http.StatusBadRequest, model.FindWorkEmailResponse{
				Status:  "error",
				Message: "Invalid request body",
			})
			return
		}
		tracing.LogObjectAsJson(span, "request", request)

		recordId, requestId, response, err := services.BettercontactService.FindWorkEmail(ctx, request.LinkedinUrl, request.FirstName, request.LastName, request.CompanyName, request.CompanyDomain, request.EnrichPhoneNumber)
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, model.FindWorkEmailResponse{
				Status:  "error",
				Message: "Internal server error",
			})
			return
		}

		c.JSON(http.StatusOK, model.FindWorkEmailResponse{
			Status:                 "success",
			RecordId:               recordId,
			BetterContactRequestId: requestId,
			Data:                   response,
		})
	}
}
