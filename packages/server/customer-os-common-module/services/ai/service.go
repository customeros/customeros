package ai

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"google.golang.org/api/option"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/dto"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type aiService struct {
	log               logger.Logger
	anthropicConfig   *config.AnthropicConfig
	deepseekConfig    *config.DeepseekConfig
	groqConfig        *config.GroqConfig
	geminiConfig      *config.GeminiConfig
	opensearchService interfaces.OpensearchService
}

func NewAIService(
	log logger.Logger,
	anthropicConfig *config.AnthropicConfig,
	deepseekConfig *config.DeepseekConfig,
	groqConfig *config.GroqConfig,
	geminiConfig *config.GeminiConfig,
	opensearchService interfaces.OpensearchService,
) interfaces.AIService {
	return &aiService{
		log:               log,
		anthropicConfig:   anthropicConfig,
		deepseekConfig:    deepseekConfig,
		groqConfig:        groqConfig,
		geminiConfig:      geminiConfig,
		opensearchService: opensearchService,
	}
}

const MaxAttempts = 3

func (s *aiService) AskAI(ctx context.Context, request interfaces.AskAIRequest) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AIService.AskAI")
	defer span.Finish()
	tracing.LogObjectAsJson(span, "requestParams", request)

	// Validate the request and set default values
	err := s.validateAIRequest(ctx, &request)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	var lastError error
	var answer *string

	for attempt := 0; attempt < *request.Retries; attempt++ {
		// If this isn't the first attempt and we have an error from a previous attempt
		if attempt > 0 && lastError != nil && answer != nil {
			// Create an error prompt that includes feedback from the previous attempt
			errorPrompt := fmt.Sprintf(`
                I previously asked you to do the following: %s
                You gave me an unexpected response of %s
                This resulted in this error: %s
                I'll give you the data again. Please re-evaluate your reply, and ensure your response is valid.`,
				utils.IfNotNilString(request.SystemPrompt),
				*answer,
				lastError.Error())
			request.SystemPrompt = &errorPrompt
		}

		answer, err = s.askAIWithRetry(ctx, request)

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
		if attempt < *request.Retries-1 {
			backoff := utils.BackOffExponentialDelay(attempt)
			time.Sleep(backoff)
		}
	}

	// If we've exhausted all retries, return the last error
	if lastError != nil {
		return nil, fmt.Errorf("askAI failed after %d attempts: %w", *request.Retries, lastError)
	}

	// This handles the case where we didn't get an error but also didn't get a valid answer
	return nil, fmt.Errorf("askAI failed after %d attempts with no specific error and invalid response", *request.Retries)
}

func (s *aiService) validateAIRequest(ctx context.Context, request *interfaces.AskAIRequest) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AIService.validateAIRequest")
	defer span.Finish()
	tracing.LogObjectAsJson(span, "requestParams", request)

	if request.Prompt == nil {
		return errors.New("prompt cannot be empty")
	}

	if request.OutputFormat != enum.AIOutputJson && request.OutputFormat != enum.AIOutputText {
		return errors.New("invalid output format")
	}

	// Set default values for unspecified fields
	if request.ModelTemperature == nil {
		temp := float32(DefaultTemperature)
		request.ModelTemperature = &temp
	}

	if request.MaxOutputTokens == nil {
		maxTokens := int32(MaxTokens)
		request.MaxOutputTokens = &maxTokens
	}

	if request.Retries == nil {
		maxRetries := MaxAttempts
		request.Retries = &maxRetries
	}

	return nil
}

func (s *aiService) askAIWithRetry(ctx context.Context, request interfaces.AskAIRequest) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AIService.askAIWithRetry")
	defer span.Finish()

	llmTracker := s.newObservabilityContainer(ctx, span, request)

	var result *string
	var err error

	llmTracker.Temperature = *request.ModelTemperature

	switch request.Model {
	case
		enum.AIModelAnthropicHaiku,
		enum.AIModelAnthropicSonnet:

		result, err = s.askAnthropic(ctx, request)
		if err != nil {
			tracing.TraceErr(span, err)
			s.trackError(ctx, llmTracker, err)
			return nil, err
		}

	case enum.AIModelDeepseekChat:
		result, err = s.askDeepseek(ctx, request)
		if err != nil {
			tracing.TraceErr(span, err)
			s.trackError(ctx, llmTracker, err)
			return nil, err
		}

	case
		enum.AIModelDeepseekQwen,
		enum.AIModelGemma,
		enum.AIModelLlama8B,
		enum.AIModelLlama70B,
		enum.AIModelMixtral,
		enum.AIModelWhisper:

		result, err = s.askGroq(ctx, request)
		if err != nil {
			tracing.TraceErr(span, err)
			s.trackError(ctx, llmTracker, err)
			return nil, err
		}

	case
		enum.AIModelGemini,
		enum.AIModelGeminiLite:
		result, err = s.askGemini(ctx, request)
		if err != nil {
			tracing.TraceErr(span, err)
			s.trackError(ctx, llmTracker, err)
			return nil, err
		}

	default:
		err := s.NewErrorNoRetry("Unsupported model", nil)
		return nil, err
	}

	if err != nil {
		tracing.TraceErr(span, err)
		s.trackError(ctx, llmTracker, err)
		return nil, err
	}

	if request.OutputFormat == enum.AIOutputText {
		answer := strings.TrimPrefix(*result, `"""`)
		answer = strings.TrimSuffix(answer, `"""`)
		s.trackSuccess(ctx, llmTracker, result)
		return &answer, nil
	}

	s.trackSuccess(ctx, llmTracker, result)
	return result, nil
}

func (s *aiService) trackError(ctx context.Context, llmTracker *dto.LLMObservability, aiError error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "AIService.trackError")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "llm", llmTracker)

	if s.opensearchService == nil {
		err := errors.New("Opensearch service is not initialized")
		tracing.TraceErr(span, err)
		return
	}

	if llmTracker == nil || llmTracker.RequestID == "" {
		err := errors.New("llm observability container cannot be nil")
		tracing.TraceErr(span, err)
		return
	}
	if llmTracker.RequestID == "" {
		llmTracker.RequestID = utils.GenerateNanoIdWithPrefix("llm", 16)
	}

	index := fmt.Sprintf("llm-%s", utils.CurrentMonth())
	err := s.opensearchService.LLMObservabilityIndexCheck(ctx, index)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	llmTracker.Success = false
	llmTracker.ErrorMessage = aiError.Error()
	err = s.opensearchService.UpsertDocument(ctx, index, &llmTracker.RequestID, llmTracker)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}
	return
}

func (s *aiService) trackSuccess(ctx context.Context, llmTracker *dto.LLMObservability, result *string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "AIService.trackSuccess")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "llm", llmTracker)

	if s.opensearchService == nil {
		err := errors.New("Opensearch service is not initialized")
		tracing.TraceErr(span, err)
		return
	}

	if llmTracker == nil {
		err := errors.New("llm observability container cannot be nil")
		tracing.TraceErr(span, err)
		return
	}
	if llmTracker.RequestID == "" {
		llmTracker.RequestID = utils.GenerateNanoIdWithPrefix("llm", 16)
	}

	index := fmt.Sprintf("llm-%s", utils.CurrentMonth())
	err := s.opensearchService.LLMObservabilityIndexCheck(ctx, index)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}

	llmTracker.Success = true

	if result != nil {
		llmTracker.Response = *result
	}

	err = s.opensearchService.UpsertDocument(ctx, index, &llmTracker.RequestID, llmTracker)
	if err != nil {
		tracing.TraceErr(span, err)
		return
	}
	return
}

func (s *aiService) newObservabilityContainer(ctx context.Context, span opentracing.Span, request interfaces.AskAIRequest) *dto.LLMObservability {
	var traceID string

	// For Jaeger specifically
	if jaegerSpan, ok := span.(*jaeger.Span); ok {
		spanContext := jaegerSpan.Context().(jaeger.SpanContext)
		traceID = spanContext.TraceID().String()
	} else {
		// Fallback for other tracers - get carrier with all span context info
		carrier := opentracing.TextMapCarrier{}
		err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, carrier)
		if err == nil {
			// Many tracers use these standard field names
			traceID = carrier["uber-trace-id"]
		}
	}

	return &dto.LLMObservability{
		RequestID:    utils.GenerateNanoIdWithPrefix("llm", 16),
		Timestamp:    utils.Now(),
		UserID:       common.GetUserIdFromContext(ctx),
		Tenant:       common.GetTenantFromContext(ctx),
		Model:        request.Model.String(),
		TraceID:      traceID,
		SystemPrompt: utils.IfNotNilString(request.SystemPrompt),
		Prompt:       utils.IfNotNilString(request.Prompt),
	}
}

func (s *aiService) askGemini(ctx context.Context, request interfaces.AskAIRequest) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AIService.askGemini")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "requestParams", request)

	if s.geminiConfig.ApiKey == "" {
		err := s.NewErrorNoRetry("Gemini API key not set", nil)
		tracing.TraceErr(span, err)
		return nil, err
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(s.geminiConfig.ApiKey))
	if err != nil {
		err := s.NewErrorNoRetry("cannot initiate Gemini client", err)
		tracing.TraceErr(span, err)
		return nil, err
	}
	defer client.Close()

	model := client.GenerativeModel(request.Model.String())
	if request.ModelTemperature == nil {
		model.SetTemperature(DefaultTemperature)
	} else {
		model.SetTemperature(*request.ModelTemperature)
	}

	if request.MaxOutputTokens == nil {
		model.SetMaxOutputTokens(MaxTokens)
	} else {
		model.SetMaxOutputTokens(*request.MaxOutputTokens)
	}

	switch request.OutputFormat {
	case enum.AIOutputText:
		model.ResponseMIMEType = "text/plain"
	case enum.AIOutputJson:
		model.ResponseMIMEType = "application/json"
	default:
		err := s.NewErrorNoRetry("Unsupported model output type", nil)
		return nil, err
	}

	if request.SystemPrompt != nil {
		model.SystemInstruction = &genai.Content{
			Parts: []genai.Part{genai.Text(*request.SystemPrompt)},
		}
	}

	session := model.StartChat()
	session.History = []*genai.Content{}

	resp, err := session.SendMessage(ctx, genai.Text(*request.Prompt))
	if err != nil {
		err = s.NewRetryableError("Could not ask Gemini", err)
		tracing.TraceErr(span, err)
		return nil, err
	}

	// Build response string from all parts
	var responseBuilder strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		// Check if part is a text part
		if textPart, ok := part.(genai.Text); ok {
			responseBuilder.WriteString(string(textPart))
		}
	}

	respStr := responseBuilder.String()
	if respStr == "" {
		err := s.NewRetryableError("Answer from Gemini was empty", nil)
		return nil, err
	}

	return &respStr, nil
}

func (s *aiService) askDeepseek(ctx context.Context, request interfaces.AskAIRequest) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AIService.askDeepseek")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "requestParams", request)

	if s.deepseekConfig.ApiKey == "" {
		err := s.NewErrorNoRetry("Deepseek API key not set", nil)
		tracing.TraceErr(span, err)
		s.log.Error(err)
		return nil, err
	}

	if s.deepseekConfig.Url == "" {
		err := s.NewErrorNoRetry("Deepseek URL not set", nil)
		tracing.TraceErr(span, err)
		s.log.Error(err)
		return nil, err
	}

	client := NewDeepseekClient(s.deepseekConfig)
	response, err := client.AskDeepseek(ctx, request)
	if err != nil {
		err = s.NewErrorNoRetry("Unable to initiate Deepseek client", err)
		tracing.TraceErr(span, err)
		return nil, err
	}

	if response == nil {
		err = s.NewRetryableError("Empty response from Deepseek", nil)
		return nil, err
	}

	return response, nil
}

func (s *aiService) askAnthropic(ctx context.Context, request interfaces.AskAIRequest) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AIService.askAnthropic")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "requestParams", request)

	if s.anthropicConfig.ApiKey == "" || s.anthropicConfig.ApiPath == "" {
		err := s.NewErrorNoRetry("Anthropic API key or path not set", nil)
		tracing.TraceErr(span, err)
		return nil, err
	}

	// setup client
	client := NewAnthropicClient(s.anthropicConfig)
	response, err := client.Invoke(ctx, request)
	if err != nil {
		err = s.NewErrorNoRetry("Unable to initiate Anthropic client", err)
		tracing.TraceErr(span, err)
		return nil, err
	}

	if response == "" {
		err = s.NewRetryableError("Empty response from Anthropic", nil)
		return nil, err
	}

	return &response, nil
}

func (s *aiService) askGroq(ctx context.Context, request interfaces.AskAIRequest) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AIService.askGroq")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)
	tracing.LogObjectAsJson(span, "requestParams", request)

	if s.groqConfig.ApiKey == "" || s.groqConfig.Url == "" {
		err := s.NewErrorNoRetry("Groq API Key or URL not set", nil)
		tracing.TraceErr(span, err)
		return nil, err
	}

	// setup client
	client := NewGroqClient(s.groqConfig)
	response, err := client.Invoke(ctx, request)
	if err != nil {
		err = s.NewErrorNoRetry("Unable to initiate Groq client", err)
		tracing.TraceErr(span, err)
		return nil, err
	}

	if response == "" {
		err = s.NewRetryableError("Empty response from Groq", nil)
		return nil, err
	}

	return &response, nil
}
