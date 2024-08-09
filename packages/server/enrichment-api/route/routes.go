package route

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/service/security"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/service"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"net/http"
)

func RegisterRoutes(ctx context.Context, r *gin.Engine, services *service.Services) {
	r.GET("/enrichPerson",
		tracing.TracingEnhancer(ctx, "/enrichPerson"),
		security.ApiKeyCheckerHTTP(services.CommonServices.PostgresRepositories.TenantWebhookApiKeyRepository, services.CommonServices.PostgresRepositories.AppKeyRepository, security.ENRICHMENT_API),
		enrichPerson(services))
}

func enrichPerson(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(c.Request.Context(), "EnrichPerson", c.Request.Header)
		defer span.Finish()

		var request model.EnrichPersonRequest

		if err := c.BindJSON(&request); err != nil {
			services.Logger.Errorf("Fail reading request: %v", err.Error())
			c.JSON(http.StatusBadRequest, model.EnrichPersonResponse{
				Status:  "error",
				Message: "Invalid request body",
			})
			return
		}

		// for now linkedin url is mandatory
		if request.LinkedinUrl == "" {
			tracing.TraceErr(span, errors.New("Missing linkedin parameter"))
			services.Logger.Errorf("Missing ip parameter")
			c.JSON(http.StatusBadRequest, model.EnrichPersonResponse{
				Status:  "error",
				Message: "Missing ip parameter",
			})
			return
		}
		span.LogFields(log.String("request.linkedin", request.LinkedinUrl))

		var enrichPersonData *model.EnrichedPersonData
		if request.LinkedinUrl != "" {
			response, err := services.PersonScrapeInService.ScrapInPersonProfile(ctx, request.LinkedinUrl)
			if err != nil {
				tracing.TraceErr(span, err)
				c.JSON(http.StatusInternalServerError, model.EnrichPersonResponse{
					Status:  "error",
					Message: "Internal server error",
				})
				return
			}
			enrichPersonData = &model.EnrichedPersonData{
				PersonProfile: response,
			}
		}

		c.JSON(http.StatusOK, model.EnrichPersonResponse{
			Status: "success",
			Data:   enrichPersonData,
		})
	}
}
