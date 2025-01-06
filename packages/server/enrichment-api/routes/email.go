package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/model"
	"github.com/openline-ai/openline-customer-os/packages/server/enrichment-api/service"
)

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
