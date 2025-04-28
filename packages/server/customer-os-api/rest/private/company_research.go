package private

import (
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type CompanyResearchHandler struct {
	services        *cosapi_services.Services
	responseHandler *response.Response
}

func NewCompanyResearchHandler(services *cosapi_services.Services, responseHandler *response.Response) *CompanyResearchHandler {
	return &CompanyResearchHandler{
		services:        services,
		responseHandler: responseHandler,
	}
}

type CompanyResearchRequest struct {
	Domain string `json:"domain"`
}

func (h *CompanyResearchHandler) GenerateCompanyBrief() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceCustomerOsApi)

		spans, ctx := telemetry.StartRestSpan(ctx, "CompanyResearchHandler.GenerateCompanyBrief")
		defer spans.Finish()

		// parse request
		_, err := h.parseRequest(c)
		if err != nil {
			message := "Unable to parse request"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		h.responseHandler.HandleAccepted(c)

		_, err = h.services.CommonServices.CompanyResearch.GenerateCompanyBriefForTenant(ctx)
		if err != nil {
			spans.TraceError(err)
			message := "Unable to crawl website"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}
	}
}

func (h *CompanyResearchHandler) parseRequest(c *gin.Context) (*WebscrapeRequest, error) {
	spans, _ := telemetry.StartRestSpan(c.Request.Context(), "CompanyResearchHandler.parseRequest")
	defer spans.Finish()

	var request WebscrapeRequest
	err := c.BindJSON(&request)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogObjectAsJson("request", request)

	return &request, nil
}
