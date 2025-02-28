package private

import (
	"net/http"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/constants"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type WebscrapeHandler struct {
	services        *cosapi_services.Services
	responseHandler *response.Response
}

func NewWebscrapeHandler(services *cosapi_services.Services, responseHandler *response.Response) *WebscrapeHandler {
	return &WebscrapeHandler{
		services:        services,
		responseHandler: responseHandler,
	}
}

type WebscrapeRequest struct {
	Url string `json:"url"`
}

type WebscrapeResponse struct {
	ScrapedUrls []string `json:"scrapedUrls"`
}

func (h *WebscrapeHandler) Crawl() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceCustomerOsApi)

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(ctx, "/crawl", c.Request.Header)
		defer span.Finish()
		tracing.SetDefaultServiceSpanTags(ctx, span)

		// parse request
		request, err := h.parseRequest(c)
		if err != nil {
			message := "Unable to parse request"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		crawled, err := h.services.CommonServices.WebscraperService.Crawl(ctx, request.Url)
		if err != nil {
			tracing.TraceErr(span, err)
			message := "Unable to crawl website"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		h.responseHandler.HandleSuccess(c, WebscrapeResponse{
			ScrapedUrls: crawled,
		})

		return
	}
}

func (h *WebscrapeHandler) parseRequest(c *gin.Context) (*WebscrapeRequest, error) {
	span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "parseRequest")
	defer span.Finish()
	tracing.TagComponentRest(span)

	var request WebscrapeRequest
	err := c.BindJSON(&request)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	tracing.LogObjectAsJson(span, "request", request)

	return &request, nil
}

type EmbedRequest struct {
	Url string `json:"url"`
}

func (h *WebscrapeHandler) EmbedWebpage() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceCustomerOsApi)

		ctx, span := tracing.StartHttpServerTracerSpanWithHeader(ctx, "/embedWebpage", c.Request.Header)
		defer span.Finish()
		tracing.SetDefaultServiceSpanTags(ctx, span)

		var request EmbedRequest
		err := c.BindJSON(&request)
		if err != nil {
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		_, err = h.services.CommonServices.WebscraperService.Scrape(ctx, request.Url)
		if err != nil {
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
		}

		record, err := h.services.CommonServices.PostgresRepositories.ScrapedWebpageRepository.GetWebpage(ctx, request.Url, 180)
		if err != nil {
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
		}

		embeddingIds, err := h.services.CommonServices.EmbeddingService.EmbedWebpage(ctx, *record)
		if err != nil {
			tracing.TraceErr(span, err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
		}
		h.responseHandler.HandleSuccess(c, embeddingIds)
	}
}
