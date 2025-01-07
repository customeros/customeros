package private

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-api/service"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type AskAIRequest struct {
	Model   string
	Prompt  string
	AIModel enum.AIModel
}

// AskAI handles AI requests
func AskAI(services *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		span, ctx := opentracing.StartSpanFromContext(c.Request.Context(), "AskAI")
		defer span.Finish()

		// parse request
		request, err := parseRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "unable to parse request",
			})
			return
		}
		if request == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "empty request",
			})
			return
		}
		if request.Model == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "model missing",
			})
			return
		}
		if request.Prompt == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "prompt missing",
			})
			return
		}

		// call appropriate model
		answer, err := services.AIService.AskAI(ctx, request.AIModel, &request.Prompt)
		if err != nil {
			tracing.TraceErr(span, err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "unable to ask AI",
			})
			return
		}
		if answer == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "empty response from AI",
			})
			return
		}

		c.JSON(200, gin.H{
			"status":   "success",
			"model":    request.AIModel.String(),
			"response": answer,
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

	aiModel, err := enum.GetAIModel(request.Model)
	if err != nil {
		err = errors.New("Invalid AI model: " + request.Model)
		tracing.TraceErr(span, err)
		return nil, err
	}

	request.AIModel = aiModel

	return &request, nil
}
