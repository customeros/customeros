package private

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

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
	Model        string             `json:"model"`
	SystemPrompt *string            `json:"systemPrompt,omitempty"`
	Prompt       any                `json:"prompt"`
	AIModel      commonEnum.AIModel `json:"-"`
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

		// validate request
		err = h.validateAIRequest(request)
		if err != nil {
			tracing.TraceErr(span, err)
			message := err.Error()
			h.responseHandler.HandleError(c, http.StatusBadRequest, &message)
			return
		}

		// call appropriate model
		answer, err := h.services.CommonServices.AIService.AskAI(ctx, request.AIModel, request.SystemPrompt, request.Prompt)
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
	tracing.TagComponentRest(span)

	var request AskAIRequest
	err := c.BindJSON(&request)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	tracing.LogObjectAsJson(span, "request", request)

	aiModel, err := commonEnum.GetAIModel(request.Model)
	if err != nil {
		err = errors.New("Invalid AI model: " + request.Model)
		tracing.TraceErr(span, err)
		return nil, err
	}

	request.AIModel = aiModel

	return &request, nil
}

func (h *AskAIHandler) validateAIRequest(req *AskAIRequest) error {
	if req.Model == "" {
		return fmt.Errorf("model is required")
	}

	switch v := req.Prompt.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return fmt.Errorf("prompt cannot be empty")
		}
	case map[string]interface{}, []interface{}:
		if len(fmt.Sprintf("%v", v)) == 0 {
			return fmt.Errorf("prompt content cannot be empty")
		}
	default:
		return fmt.Errorf("prompt must be string or structured data")
	}

	return nil
}
