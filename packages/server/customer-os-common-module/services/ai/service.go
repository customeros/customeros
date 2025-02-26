package ai

import (
	"context"
	"errors"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/api/option"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/utils"
)

type aiService struct {
	log             logger.Logger
	anthropicConfig *config.AnthropicConfig
	deepseekConfig  *config.DeepseekConfig
	groqConfig      *config.GroqConfig
	geminiConfig    *config.GeminiConfig
}

func NewAIService(
	log logger.Logger,
	anthropicConfig *config.AnthropicConfig,
	deepseekConfig *config.DeepseekConfig,
	groqConfig *config.GroqConfig,
	geminiConfig *config.GeminiConfig,
) interfaces.AIService {
	return &aiService{
		log:             log,
		anthropicConfig: anthropicConfig,
		deepseekConfig:  deepseekConfig,
		groqConfig:      groqConfig,
		geminiConfig:    geminiConfig,
	}
}

func (s *aiService) AskAI(ctx context.Context, request interfaces.AskAIRequest) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AIService.AskAI")
	defer span.Finish()
	tracing.LogObjectAsJson(span, "requestParams", request)

	var result *string
	var err error

	if request.Prompt == nil {
		err := errors.New("Prompt cannot be empty")
		tracing.TraceErr(span, err)
		return nil, err
	}

	if request.ModelTemperature == nil {
		temp := float32(DefaultTemperature)
		request.ModelTemperature = &temp
	}
	if request.MaxOutputTokens == nil {
		maxTokens := int32(MaxTokens)
		request.MaxOutputTokens = &maxTokens
	}

	switch request.Model {
	case
		enum.AIModelAnthropicHaiku,
		enum.AIModelAnthropicSonnet:

		result, err = s.askAnthropic(ctx, request)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

	case enum.AIModelDeepseekChat:
		result, err = s.askDeepseek(ctx, request)
		if err != nil {
			tracing.TraceErr(span, err)
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
			return nil, err
		}

	case
		enum.AIModelGemini,
		enum.AIModelGeminiLite:
		result, err = s.askGemini(ctx, request)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

	default:
		err := errors.New("Unsupported model")
		return nil, err
	}

	span.LogKV("result", utils.IfNotNilString(result))
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	return result, nil
}

func (s *aiService) askGemini(ctx context.Context, request interfaces.AskAIRequest) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AIService.askGemini")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if s.geminiConfig.ApiKey == "" {
		err := errors.New("Gemini API key not set")
		tracing.TraceErr(span, err)
		s.log.Error(err)
		return nil, err
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(s.geminiConfig.ApiKey))
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	defer client.Close()

	model := client.GenerativeModel(request.Model.String())
	model.SetTemperature(*request.ModelTemperature)
	model.SetMaxOutputTokens(*request.MaxOutputTokens)

	switch request.OutputFormat {
	case enum.AIOutputText:
		model.ResponseMIMEType = "text/plain"
	case enum.AIOutputJson:
		model.ResponseMIMEType = "application/json"
	default:
		return nil, errors.New("unsupported output type")
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
	return &respStr, nil
}

func (s *aiService) askDeepseek(ctx context.Context, request interfaces.AskAIRequest) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AIService.askDeepseek")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if s.deepseekConfig.ApiKey == "" {
		err := errors.New("Deepseek API key not set")
		tracing.TraceErr(span, err)
		s.log.Error(err)
		return nil, err
	}

	if s.deepseekConfig.Url == "" {
		err := errors.New("Deepseek Url not set")
		tracing.TraceErr(span, err)
		s.log.Error(err)
		return nil, err
	}

	client := NewDeepseekClient(s.deepseekConfig)
	response, err := client.AskDeepseek(ctx, request)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return response, nil
}

func (s *aiService) askAnthropic(ctx context.Context, request interfaces.AskAIRequest) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AIService.askAnthropic")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if s.anthropicConfig.ApiKey == "" || s.anthropicConfig.ApiPath == "" {
		err := errors.New("Anthropic API key or path not set")
		tracing.TraceErr(span, err)
		s.log.Error(err)
		return nil, err
	}

	// setup client
	client := NewAnthropicClient(s.anthropicConfig)
	response, err := client.Invoke(ctx, request)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &response, nil
}

func (s *aiService) askGroq(ctx context.Context, request interfaces.AskAIRequest) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AIService.askGroq")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if s.groqConfig.ApiKey == "" || s.groqConfig.Url == "" {
		err := errors.New("Groq API key or URL not set")
		tracing.TraceErr(span, err)
		s.log.Error(err)
		return nil, err
	}

	// setup client
	client := NewGroqClient(s.groqConfig)
	response, err := client.Invoke(ctx, request)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &response, nil
}
