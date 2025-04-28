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

		spans, ctx := telemetry.StartRestSpan(ctx, "WebscrapeHandler/Crawl")
		defer spans.Finish()

		// parse request
		request, err := h.parseRequest(c)
		if err != nil {
			message := "Unable to parse request"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		crawled, err := h.services.CommonServices.WebscraperService.Crawl(ctx, request.Url)
		if err != nil {
			spans.TraceError(err)
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
	spans, _ := telemetry.StartRestSpan(c.Request.Context(), "WebscrapeHandler.parseRequest")
	defer spans.Finish()

	var request WebscrapeRequest
	err := c.BindJSON(&request)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	spans.LogKV("request", request)

	return &request, nil
}

type EmbedRequest struct {
	Url string `json:"url"`
}

func (h *WebscrapeHandler) EmbedWebpage() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := common.WithCustomContextFromGinRequest(c, constants.AppSourceCustomerOsApi)

		spans, ctx := telemetry.StartRestSpan(ctx, "WebscrapeHandler.EmbedWebpage")
		defer spans.Finish()

		var request EmbedRequest
		err := c.BindJSON(&request)
		if err != nil {
			spans.TraceError(err)
			h.responseHandler.HandleError(c, http.StatusBadRequest, nil)
			return
		}

		_, err = h.services.CommonServices.WebscraperService.Scrape(ctx, request.Url)
		if err != nil {
			spans.TraceError(err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
		}

		record, err := h.services.CommonServices.PostgresRepositories.ScrapedWebpageRepository.GetWebpage(ctx, request.Url, 180)
		if err != nil {
			spans.TraceError(err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
		}

		embeddingIds, err := h.services.CommonServices.EmbeddingService.EmbedWebpage(ctx, *record)
		if err != nil {
			spans.TraceError(err)
			h.responseHandler.HandleError(c, http.StatusInternalServerError, nil)
		}
		h.responseHandler.HandleSuccess(c, embeddingIds)
	}
}
