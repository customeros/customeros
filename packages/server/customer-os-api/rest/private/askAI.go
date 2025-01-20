package private

import (
	"errors"
	"net/http"

	commonEnum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-api/rest/response"
	cosapi_services "github.com/customeros/customeros/packages/server/customer-os-api/services"
)

type AskAIHandler struct {
	services        *cosapi_services.Services
	responseHandler *response.Response
}

func NewAskAIHandler(services *cosapi_services.Services, responseHandler *response.Response) *AskAIHandler {
	return &AskAIHandler{
		services:        services,
		responseHandler: responseHandler,
	}
}

type AskAIRequest struct {
	Model   string
	Prompt  string
	AIModel commonEnum.AIModel
}

type AskAIResponse struct {
	Model  string `json:"model"`
	Answer string `json:"answer"`
}

// AskAI handles AI requests
func (h *AskAIHandler) AskAI() gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "AskAI")
		defer span.Finish()
		tracing.TagComponentRest(span)

		// parse request
		request, err := h.parseRequest(c)
		if err != nil {
			message := "Unable to parse request"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}
		if request == nil {
			message := "Empty request"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}
		if request.Model == "" {
			message := "model missing"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}
		if request.Prompt == "" {
			message := "prompt missing"
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		// call appropriate model
		answer, err := h.services.CommonServices.AIService.AskAI(ctx, request.AIModel, &request.Prompt)
		if err != nil {
			tracing.TraceErr(span, err)
			message := "Unable to ask AI"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}
		if answer == nil {
			message := "Empty response from AI"
			h.responseHandler.HandleError(c, http.StatusInternalServerError, &message)
			return
		}

		h.responseHandler.HandleSuccess(c, AskAIResponse{
			Model:  request.AIModel.String(),
			Answer: *answer,
		})

		return
	}
}

func (h *AskAIHandler) parseRequest(c *gin.Context) (*AskAIRequest, error) {
	span, _ := opentracing.StartSpanFromContext(c.Request.Context(), "parseRequest")
	defer span.Finish()

	var request AskAIRequest
	err := c.BindJSON(&request)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	aiModel, err := commonEnum.GetAIModel(request.Model)
	if err != nil {
		err = errors.New("Invalid AI model: " + request.Model)
		tracing.TraceErr(span, err)
		return nil, err
	}

	request.AIModel = aiModel

	return &request, nil
}
