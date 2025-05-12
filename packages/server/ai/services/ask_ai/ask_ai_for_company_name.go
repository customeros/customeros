package ai

import (
	"context"
	"fmt"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	"github.com/customeros/customeros/packages/server/ai/proto/pb"
)

func (s *aiService) AskAIForCompanyName(ctx context.Context, message *pb.AskAI) (*data_fields.CompanyIdentification, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "AIService.AskAIForCompanyName")
	defer spans.Finish()
	spans.LogObjectAsJson("requestParams", message)

	err := s.validateAIRequest(ctx, message)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	// Create output validator
	validator := data_fields.NewCompanyResponseValidator()

	// defaults
	message.RequestType = pb.AIRequestType_AI_REQUEST_COMPANY_NAME
	message.OutputFormat = pb.AIOutputFormat_AI_OUTPUT_JSON

	// Enhance the system prompt with the expected schema if needed
	if message.SystemPrompt != nil {
		currentPrompt := *message.SystemPrompt
		schema := validator.GetExpectedSchema()
		// Only add the schema example if not already present
		enhancedPrompt := fmt.Sprintf(`%s
IMPORTANT: Your response MUST be in valid JSON format exactly matching this schema:
%s
Do not include any text outside of the JSON object. Provide the company name with high confidence if clearly identifiable, or lower confidence if uncertain.`,
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
		company, err := validator.ValidateResponse(*answer)
		if err != nil {
			lastError = err
			continue
		}

		return company, nil
	}

	// If we've exhausted all retries, return the last error
	if lastError != nil {
		err := fmt.Errorf("company identification failed after %d attempts: %w", *message.Retries, lastError)
		spans.TraceError(err)
		return nil, err
	}

	return nil, fmt.Errorf("company identification failed after %d attempts", *message.Retries)
}
