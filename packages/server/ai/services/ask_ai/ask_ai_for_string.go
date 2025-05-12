package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
	"github.com/nats-io/nats.go"

	"github.com/customeros/customeros/packages/server/ai/proto/pb"
)

func (s *aiService) AskAIForString(ctx context.Context, msg *nats.Msg) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "AIService.AskAIForString")
	defer span.Finish()

	message, err := s.parseAskAIMessage(ctx, msg)
	if err != nil {
		span.TraceError(err)
		return err
	}

	// Validate the request and set default values
	err = s.validateAIRequest(ctx, message)
	if err != nil {
		span.TraceError(err)
		return err
	}

	// defaults
	message.RequestType = pb.AIRequestType_AI_REQUEST_GENERIC
	message.OutputFormat = pb.AIOutputFormat_AI_OUTPUT_TEXT

	var lastError error
	var answer *string

	for attempt := 0; attempt < int(*message.Retries); attempt++ {
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
					utils.IfNotNilString(message.SystemPrompt),
					*answer,
					lastError.Error())
			} else {
				// We got no response at all
				errorPrompt = fmt.Sprintf(`
I previously asked you to do the following: %s
The request failed with error: %s
Please try again with a valid response.`,
					utils.IfNotNilString(message.SystemPrompt),
					lastError.Error())
			}

			message.SystemPrompt = &errorPrompt
		}

		// Make the API call
		answer, err = s.askAIWithRetry(ctx, message)

		// Success case
		if err == nil && answer != nil {
			trimmedAnswer := strings.TrimSpace(*answer)
			if trimmedAnswer == "" {
				lastError = fmt.Errorf("received empty response")
				continue
			}

			return nil
		}

		// Store the error for the next retry
		lastError = err

		// If error is not retryable, exit immediately
		if err != nil && !s.IsRetryable(err) {
			span.TraceError(err)
			return err
		}
	}

	// We've exhausted all retries
	finalErr := fmt.Errorf("AskAI failed after %d attempts: %w", *message.Retries, lastError)
	span.TraceError(finalErr)
	return finalErr
}
