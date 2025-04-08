package ai

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

func (s *aiService) AskAIForContentStage(ctx context.Context, request interfaces.AskAIRequest) (enum.CustomerJourneyStage, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "AIService.AskAIForContentStage")
	defer spans.Finish()
	spans.LogObjectAsJson("requestParams", request)

	// validate input
	err := s.validateAIRequest(ctx, &request)
	if err != nil {
		spans.TraceError(err)
		return "", err
	}

	// defaults
	request.RequestType = enum.AIRequestContentStage
	request.OutputFormat = enum.AIOutputText

	// Get the output validator
	validator := enum.GetCustomerJourneyStageValidator()

	// Add valid values to the prompt
	s.enhancePromptWithValidValues(&request, validator)

	var lastError error
	var emptyValue enum.CustomerJourneyStage

	// Retry loop
	for attempt := 0; attempt < *request.Retries; attempt++ {
		// Add more explicit instructions on retry
		if attempt > 0 && lastError != nil {
			s.addRetryInstructions(&request, validator)
		}

		answer, err := s.askAIWithRetry(ctx, request)
		if err != nil {
			lastError = err
			if !s.IsRetryable(err) {
				spans.TraceError(err)
				return emptyValue, err
			}
			continue
		}

		if answer == nil {
			lastError = fmt.Errorf("received nil response from LLM")
			continue
		}

		// Clean and validate the response
		cleanResponse := s.cleanResponseString(*answer)
		if validator.IsValid(cleanResponse) {
			// Parse the response to the correct type
			stage, err := validator.ParseJourneyStage(cleanResponse)
			if err != nil {
				lastError = err
				continue
			}
			return stage, nil
		}

		// Invalid response
		lastError = fmt.Errorf("invalid customer journey stage value: %s", cleanResponse)
	}

	// If we get here, all retries failed
	err = fmt.Errorf("customer journey stage validation failed after %d attempts: %w", *request.Retries, lastError)
	spans.TraceError(err)
	return emptyValue, err
}
