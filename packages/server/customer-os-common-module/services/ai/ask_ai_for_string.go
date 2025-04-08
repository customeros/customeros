package ai

import (
	"context"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

func (s *aiService) AskAIForString(ctx context.Context, request interfaces.AskAIRequest) (*string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "AIService.AskAIForString")
	defer spans.Finish()
	spans.LogObjectAsJson("requestParams", request)

	// Validate the request and set default values
	err := s.validateAIRequest(ctx, &request)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	// defaults
	request.RequestType = enum.AIRequestGeneric
	request.OutputFormat = enum.AIOutputText

	var lastError error
	var answer *string

	for attempt := 0; attempt < *request.Retries; attempt++ {
		// If this isn't the first attempt and we have a previous error
		if attempt > 0 && lastError != nil {
			var errorPrompt string

			// If we got a response but it had issues
			if answer != nil {
				errorPrompt = fmt.Sprintf(`
I previously asked you to do the following: %s
You gave me this response: %s
This resulted in this error: %s
Please re-evaluate your reply, and ensure your response is valid.`,
					utils.IfNotNilString(request.SystemPrompt),
					*answer,
					lastError.Error())
			} else {
				// We got no response at all
				errorPrompt = fmt.Sprintf(`
I previously asked you to do the following: %s
The request failed with error: %s
Please try again with a valid response.`,
					utils.IfNotNilString(request.SystemPrompt),
					lastError.Error())
			}

			request.SystemPrompt = &errorPrompt
		}

		// Make the API call
		answer, err = s.askAIWithRetry(ctx, request)

		// Success case
		if err == nil && answer != nil {
			trimmedAnswer := strings.TrimSpace(*answer)
			if trimmedAnswer == "" {
				lastError = fmt.Errorf("received empty response")
				continue
			}

			return answer, nil
		}

		// Store the error for the next retry
		lastError = err

		// If error is not retryable, exit immediately
		if err != nil && !s.IsRetryable(err) {
			spans.TraceError(err)
			return nil, err
		}
	}

	// We've exhausted all retries
	finalErr := fmt.Errorf("AskAI failed after %d attempts: %w", *request.Retries, lastError)
	spans.TraceError(finalErr)
	return nil, finalErr
}
