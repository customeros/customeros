package private

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	commonEnum "github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/rest"
	cosapi_services "github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/services"
)

type AskAIRequest struct {
	Model   string
	Prompt  string
	AIModel commonEnum.AIModel
}

type AskAIResponse struct {
	enum.BaseResponse
	Model  string `json:"model"`
	Answer string `json:"answer"`
}

// AskAI handles AI requests
func AskAI(services *cosapi_services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "AskAI")
		defer span.Finish()

		// parse request
		request, err := parseRequest(c)
		if err != nil {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("unable to parse request"))
			return
		}
		if request == nil {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("empty request"))
			return
		}
		if request.Model == "" {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("model missing"))
			return
		}
		if request.Prompt == "" {
			handlers.SendError(c, span, http.StatusBadRequest, enum.ErrBadRequest.WithMessage("prompt missing"))
			return
		}

		// call appropriate model
		answer, err := services.CommonServices.AIService.AskAI(ctx, request.AIModel, &request.Prompt)
		if err != nil {
			tracing.TraceErr(span, err)
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Unable to ask AI"))
			return
		}
		if answer == nil {
			handlers.SendError(c, span, http.StatusInternalServerError, enum.ErrInternalServer.WithMessage("Empty response from AI"))
			return
		}

		c.JSON(http.StatusOK, AskAIResponse{
			BaseResponse: enum.BuildBaseResponse(enum.StatusSuccess),
			Model:        request.AIModel.String(),
			Answer:       *answer,
		})

		return
	}
}

func parseRequest(c *gin.Context) (*AskAIRequest, error) {
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
