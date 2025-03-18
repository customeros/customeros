package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

func (s *aiService) AskAIForIndustryCode(ctx context.Context, request interfaces.AskAIRequest) (*data_fields.IndustryCode, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AIService.AskAIForIndustryCode")
	defer span.Finish()
	tracing.LogObjectAsJson(span, "requestParams", request)

	err := s.validateAIRequest(ctx, &request)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Create output validator
	validator := data_fields.NewIndustryCodeResponseValidator()

	// defaults
	request.RequestType = enum.AIRequestIndustryCode
	request.OutputFormat = enum.AIOutputJson

	// Enhance the system prompt with the expected schema if needed
	if request.SystemPrompt != nil {
		currentPrompt := *request.SystemPrompt
		schema := validator.GetExpectedSchema()
		// Only add the schema example if not already present
		enhancedPrompt := fmt.Sprintf(`%s
IMPORTANT: Your response MUST be in valid JSON format exactly matching this schema:
%s
Example valid response:
{"industry": {"code": "541511", "confidence": 0.95}}`,
			currentPrompt, schema)
		request.SystemPrompt = &enhancedPrompt
	}

	var lastError error
	for attempt := 0; attempt < *request.Retries; attempt++ {
		// Add retry instructions if needed
		if attempt > 0 && lastError != nil {
			retryPrompt := fmt.Sprintf(`%s
Your previous response was not valid. 
Error: %s
Please respond ONLY with valid JSON matching this schema:
%s`,
				*request.SystemPrompt,
				lastError.Error(),
				validator.GetExpectedSchema())
			request.SystemPrompt = &retryPrompt
		}

		// Make the AI call
		answer, err := s.askAIWithRetry(ctx, request)
		// Handle API errors
		if err != nil {
			lastError = err
			if !s.IsRetryable(err) {
				tracing.TraceErr(span, err)
				return nil, err
			}
			continue
		}

		// Handle nil response
		if answer == nil {
			lastError = fmt.Errorf("received nil response from LLM")
			continue
		}

		// Validate the response
		industryCode, err := validator.ValidateResponse(*answer)
		if err != nil {
			lastError = err
			continue
		}

		return industryCode, nil
	}

	// If we've exhausted all retries, return the last error
	if lastError != nil {
		err := fmt.Errorf("industry code identification failed after %d attempts: %w", *request.Retries, lastError)
		tracing.TraceErr(span, err)
		return nil, err
	}

	return nil, fmt.Errorf("industry code identification failed after %d attempts", *request.Retries)
}
