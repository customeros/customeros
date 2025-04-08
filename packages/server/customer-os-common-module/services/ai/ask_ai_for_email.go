package ai

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

func (s *aiService) AskAIForEmail(ctx context.Context, request interfaces.AskAIRequest) (*data_fields.EmailResponse, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "AIService.AskAIForEmail")
	defer spans.Finish()
	spans.LogObjectAsJson("requestParams", request)

	err := s.validateAIRequest(ctx, &request)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	// Create output validator
	validator := data_fields.NewEmailSignatureValidator()

	// defaults
	request.RequestType = enum.AIRequestEmail
	request.OutputFormat = enum.AIOutputJson

	// Enhance the system prompt with the expected schema if needed
	if request.SystemPrompt != nil {
		currentPrompt := *request.SystemPrompt
		schema := validator.GetExpectedSchema()
		// Only add the schema example if not already present
		enhancedPrompt := fmt.Sprintf(`%s
IMPORTANT: Your response MUST be in valid JSON format exactly matching this schema:
%s
Do not include any text outside of the JSON object. If you are unable to determine a value for a field, return null for that field.`,
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
				spans.TraceError(err)
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
		description, err := validator.ValidateResponse(*answer)
		if err != nil {
			lastError = err
			continue
		}

		return description, nil
	}

	// If we've exhausted all retries, return the last error
	if lastError != nil {
		err := fmt.Errorf("email parsing failed after %d attempts: %w", *request.Retries, lastError)
		spans.TraceError(err)
		return nil, err
	}

	return nil, fmt.Errorf("email parsing failed after %d attempts", *request.Retries)
}
