package ai

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	"github.com/customeros/customeros/packages/server/ai/proto/pb"
)

func (s *aiService) AskAIForCompanyDescription(ctx context.Context, message *pb.AskAI) (*data_fields.CompanyDescription, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "AIService.AskAIForCompanyDescription")
	defer spans.Finish()
	spans.LogObjectAsJson("requestParams", message)

	err := s.validateAIRequest(ctx, message)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	// Create output validator
	validator := data_fields.NewCompanyDescriptionResponseValidator()

	// defaults
	message.RequestType = pb.AIRequestType_AI_REQUEST_COMPANY_DESCRIPTION
	message.OutputFormat = pb.AIOutputFormat_AI_OUTPUT_JSON

	// Enhance the system prompt with the expected schema if needed
	if message.SystemPrompt != nil {
		currentPrompt := *message.SystemPrompt
		schema := validator.GetExpectedSchema()
		// Only add the schema example if not already present
		enhancedPrompt := fmt.Sprintf(`%s
IMPORTANT: Your response MUST be in valid JSON format exactly matching this schema:
%s
Do not include any text outside of the JSON object. The description should be a single paragraph of max 300 characters that clearly explains who the company serves and their revenue model.`,
			currentPrompt, schema)
		message.SystemPrompt = &enhancedPrompt
	}

	var lastError error
	for attempt := 0; attempt < int(*message.Retries); attempt++ {
		// Add retry instructions if needed
		if attempt > 0 && lastError != nil {
			retryPrompt := fmt.Sprintf(`%s
Your previous response was not valid. 
Error: %s
Please respond ONLY with valid JSON matching this schema:
%s`,
				*message.SystemPrompt,
				lastError.Error(),
				validator.GetExpectedSchema())
			message.SystemPrompt = &retryPrompt
		}

		// Make the AI call
		answer, err := s.askAIWithRetry(ctx, message)
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
		err := fmt.Errorf("company description generation failed after %d attempts: %w", *message.Retries, lastError)
		spans.TraceError(err)
		return nil, err
	}

	return nil, fmt.Errorf("company description generation failed after %d attempts", *message.Retries)
}
