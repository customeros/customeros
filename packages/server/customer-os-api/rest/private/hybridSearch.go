package private

import (
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"

	"github.com/customeros/customeros/packages/server/customer-os-api/constants"
	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type HybridSearchHandler struct {
	services        *cosapi_services.Services
	responseHandler *response.Response
}

func NewHybridSearchHandler(services *cosapi_services.Services, responseHandler *response.Response) *HybridSearchHandler {
	return &HybridSearchHandler{
		services:        services,
		responseHandler: responseHandler,
	}
}

type HybridSearchRequest struct {
	Query string `json:"query"`
}

type HybridSearchResponse struct {
	Answer string `json:"answer"`
}

func (h *HybridSearchHandler) Search() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceCustomerOsApi)

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(ctx, "/hybridSearch", c.Request.Header)
		defer span.Finish()
		tracing.SetDefaultServiceSpanTags(ctx, span)

		var request HybridSearchRequest
		err := c.BindJSON(&request)
		if err != nil {
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		answer, err := h.services.CommonServices.SearchService.SearchWebsites(ctx, "nuso.cloud", request.Query)
		if err != nil {
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		h.responseHandler.HandleSuccess(c, HybridSearchResponse{
			Answer: answer,
		})
		return
	}
}
