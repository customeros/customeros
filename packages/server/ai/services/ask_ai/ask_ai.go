package ai

import (
	"context"
	"strings"

	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/ai/interfaces"
	"github.com/customeros/customeros/packages/server/ai/internal/enum"
	"github.com/customeros/customeros/packages/server/ai/internal/telemetry"
	"github.com/customeros/customeros/packages/server/ai/internal/utils"
)

const (
	MAX_ATTEMPTS        = 3
	DEFAULT_TEMPERATURE = 0.1
	DEFAULT_MAX_TOKENS  = 1024
)

var ErrEmptyResponse = errors.New("Empty response from LLM")

func (s *aiService) AskAI(ctx context.Context, message *interfaces.AskAIRequest) (*string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "AIService.AskAI")
	defer spans.Finish()

	// Validate the request and set default values
	err := s.validateAIRequest(ctx, message)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}

	answer, err := s.routeAIRequest(ctx, message)
	if err != nil {
		spans.TraceError(err)
		return nil, err
	}
	if answer == nil {
		spans.TraceError(ErrEmptyResponse)
		return nil, ErrEmptyResponse
	}

	return answer, nil
}

func (s *aiService) validateAIRequest(ctx context.Context, request *interfaces.AskAIRequest) error {
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
		request.Retries = &maxRetries
	}

	// If JSON output is requested, enhance the system prompt with JSON validation requirements
	if request.OutputFormat == enum.AIOutputJson && request.SystemPrompt != nil {
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

func (s *aiService) routeAIRequest(ctx context.Context, message *interfaces.AskAIRequest) (*string, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "AIService.routeAIRequest")
	defer spans.Finish()

	var result *string
	var err error

	switch message.Model {
	case
		enum.AIModelAnthropicHaiku,
		enum.AIModelAnthropicSonnet:

		result, err = s.anthropic.Ask(ctx, message)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}

	case enum.AIModelDeepseekChat:
		result, err = s.deepseek.Ask(ctx, message)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}

	case
		enum.AIModelDeepseekQwen,
		enum.AIModelGemma,
		enum.AIModelLlama8B,
		enum.AIModelLlama70B,
		enum.AIModelMixtral,
		enum.AIModelWhisper:

		result, err = s.groq.Ask(ctx, message)
		if err != nil {
			spans.TraceError(err)
			return nil, err
		}

	case
		enum.AIModelGemini,
		enum.AIModelGeminiLite:
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
	case enum.AIOutputText:
		answer := strings.TrimPrefix(*result, `"""`)
		answer = strings.TrimSuffix(answer, `"""`)
		processedResult = &answer
		spans.LogKV("result.plain", answer)

	case enum.AIOutputJson:
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
