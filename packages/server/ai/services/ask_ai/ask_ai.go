package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/customeros/customeros/packages/server/ai/internal/telemetry"
	"github.com/customeros/customeros/packages/server/ai/internal/utils"
	"github.com/customeros/customeros/packages/server/ai/proto/pb"
)

const (
	MAX_ATTEMPTS        = 3
	DEFAULT_TEMPERATURE = 0.1
	DEFAULT_MAX_TOKENS  = 1024
)

func (s *aiService) askAI(ctx context.Context, message *pb.AskAI) (*string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "AIService.AskAI")
	defer spans.Finish()

	// Validate the request and set default values
	err := s.validateAIRequest(ctx, message)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	var lastError error
	var answer *string

	for attempt := 0; attempt < int(*message.Retries); attempt++ {
		// If this isn't the first attempt and we have an error from a previous attempt
		if attempt > 0 && lastError != nil && answer != nil {
			// Create an error prompt that includes feedback from the previous attempt
			errorPrompt := fmt.Sprintf(`
                I previously asked you to do the following: %s
                You gave me an unexpected response of %s
                This resulted in this error: %s
                I'll give you the data again. Please re-evaluate your reply, and ensure your response is valid.`,
				utils.IfNotNilString(message.SystemPrompt),
				*answer,
				lastError.Error())
			message.SystemPrompt = &errorPrompt
		}

		answer, err = s.askAIWithRetry(ctx, message)

		// If successful, return the answer
		if err == nil && answer != nil {
			return answer, nil
		}

		// Store the last error for potential use in the next retry
		lastError = err

		// If the error is not retryable, stop trying
		if err != nil && !s.IsRetryable(err) {
			return nil, err
		}

		// Add a delay before the next retry (except for the last attempt)
		if attempt < int(*message.Retries)-1 {
			backoff := utils.BackOffExponentialDelay(attempt)
			time.Sleep(backoff)
		}
	}

	// If we've exhausted all retries, return the last error
	if lastError != nil {
		return nil, fmt.Errorf("askAI failed after %d attempts: %w", *message.Retries, lastError)
	}

	// This handles the case where we didn't get an error but also didn't get a valid answer
	return nil, fmt.Errorf("askAI failed after %d attempts with no specific error and invalid response", *message.Retries)
}

func (s *aiService) validateAIRequest(ctx context.Context, request *pb.AskAI) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "AIService.validateAIRequest")
	defer spans.Finish()

	spans.LogObjectAsJson("requestParams", request)

	if request.Prompt == nil {
		return errors.New("prompt cannot be empty")
	}

	// Set default values for unspecified fields
	if request.ModelTemperature == nil {
		temp := float32(DEFAULT_TEMPERATURE)
		request.ModelTemperature = &temp
	}

	if request.MaxOutputTokens == nil {
		maxTokens := int32(DEFAULT_MAX_TOKENS)
		request.MaxOutputTokens = &maxTokens
	}

	if request.Retries == nil {
		maxRetries := MAX_ATTEMPTS
		request.Retries = utils.Int32Ptr(maxRetries)
	}

	// If JSON output is requested, enhance the system prompt with JSON validation requirements
	if request.OutputFormat == pb.AIOutputFormat_AI_OUTPUT_JSON && request.SystemPrompt != nil {
		enhancedPrompt := *request.SystemPrompt + `

CRITICAL JSON FORMATTING REQUIREMENTS:
1. Response MUST be a complete, valid JSON object
2. Response MUST start with { and end with }
3. All arrays MUST start with [ and end with ]
4. Every opening bracket MUST have a matching closing bracket
5. Every opening brace MUST have a matching closing brace
6. NO text outside the JSON structure
7. NO truncation of the JSON response
8. All string values MUST be in double quotes
9. All numeric values MUST NOT be in quotes
10. NO trailing commas allowed
11. Do not include reasoning or rationalization in response

VALIDATION CHECKLIST (complete ALL before responding):
1. Verify all opening { have a matching number of closing }
2. Verify all opening [ have a matching number of closing ]
3. Verify all string values are in double quotes
4. Verify all numeric values are unquoted
5. Verify no trailing commas
6. Verify complete JSON structure
7. Test parse the JSON to ensure it's valid
8. Only then return the response

Remember: NEVER return incomplete JSON. If you need more tokens, make the response more concise but ensure it is complete.`

		request.SystemPrompt = &enhancedPrompt
	}

	return nil
}

func (s *aiService) askAIWithRetry(ctx context.Context, message *pb.AskAI) (*string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "AIService.askAIWithRetry")
	defer spans.Finish()

	var result *string
	var err error

	switch message.Model {
	case
		pb.AIModel_AI_MODEL_ANTHROPIC_HAIKU,
		pb.AIModel_AI_MODEL_ANTHROPIC_SONNET:

		result, err = s.anthropic.Ask(ctx, message)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}

	case pb.AIModel_AI_MODEL_DEEPSEEK_CHAT:
		result, err = s.deepseek.Ask(ctx, message)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}

	case
		pb.AIModel_AI_MODEL_DEEPSEEK_QWEN,
		pb.AIModel_AI_MODEL_GEMMA,
		pb.AIModel_AI_MODEL_LLAMA_8B,
		pb.AIModel_AI_MODEL_LLAMA_70B,
		pb.AIModel_AI_MODEL_MIXTRAL,
		pb.AIModel_AI_MODEL_WHISPER:

		result, err = s.groq.Ask(ctx, message)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}

	case
		pb.AIModel_AI_MODEL_GEMINI,
		pb.AIModel_AI_MODEL_GEMINI_LITE:
		result, err = s.gemini.Ask(ctx, message)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}

	default:
		err := s.NewErrorNoRetry("Unsupported model", nil)
		return nil, err
	}

	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	var processedResult *string
	switch message.OutputFormat {
	case pb.AIOutputFormat_AI_OUTPUT_TEXT:
		answer := strings.TrimPrefix(*result, `"""`)
		answer = strings.TrimSuffix(answer, `"""`)
		processedResult = &answer
		spans.LogKV("result.plain", answer)

	case pb.AIOutputFormat_AI_OUTPUT_JSON:
		processedResult = s.extractJsonFromAiResponse(result)
		spans.LogFields(
			log.String("result.plain", utils.IfNotNilString(result)),
			log.String("result.json", utils.IfNotNilString(processedResult)),
		)

	default:
		// For unknown formats, return the original result
		processedResult = result
		spans.LogKV("result.plain", utils.IfNotNilString(*result))
	}

	return processedResult, nil
}

func (s *aiService) extractJsonFromAiResponse(result *string) *string {
	if result == nil {
		return nil
	}

	firstBrace := strings.Index(*result, "{")
	lastBrace := strings.LastIndex(*result, "}")

	if firstBrace >= 0 && lastBrace >= 0 && lastBrace > firstBrace {
		trimmed := (*result)[firstBrace : lastBrace+1]
		return &trimmed
	}

	return result
}

func (s *aiService) parseAskAIMessage(ctx context.Context, msg *nats.Msg) (*pb.AskAI, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "aiService.parseAskAIMessage")
	defer span.Finish()

	message := &pb.AskAI{}
	err := proto.Unmarshal(msg.Data, message)
	if err != nil {
		span.TraceError(err)
		return nil, err
	}
	return message, nil
}
