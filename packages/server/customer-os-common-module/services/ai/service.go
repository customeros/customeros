package ai

import (
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"

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
}

func NewAIService(log logger.Logger, anthropicConfig *config.AnthropicConfig, deepseekConfig *config.DeepseekConfig) interfaces.AIService {
	return &aiService{
		log:             log,
		anthropicConfig: anthropicConfig,
		deepseekConfig:  deepseekConfig,
	}
}

func (s *aiService) AskAI(ctx context.Context, model enum.AIModel, systemPrompt *string, prompt any) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AIService.AskAI")
	defer span.Finish()
	span.LogKV("model", model)
	span.LogKV("systemPrompt", utils.IfNotNilString(systemPrompt))
	tracing.LogObjectAsJson(span, "prompt", prompt)

	var result *string
	var err error

	switch model {
	case
		enum.AIModelAnthropicHaiku,
		enum.AIModelAnthropicSonnet:

		result, err = s.askAnthropic(ctx, model, systemPrompt, prompt)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}

	case enum.AIModelDeepseekChat:
		result, err = s.askDeepseek(ctx, model, systemPrompt, prompt)
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

func (s *aiService) askDeepseek(ctx context.Context, model enum.AIModel, systemPrompt *string, prompt any) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AIService.askDeepseek")
	defer span.Finish()
	tracing.SetDefaultServiceSpanTags(ctx, span)

	if s.deepseekConfig.ApiKey == "" || s.deepseekConfig.Url == "" {
		err := errors.New("Deepseek API key or path not set")
		tracing.TraceErr(span, err)
		s.log.Error(err)
		return nil, err
	}

	client := NewDeepseekClient(s.deepseekConfig, model)
	response, err := client.AskDeepseek(ctx, systemPrompt, prompt)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return response, nil
}

func (s *aiService) askAnthropic(ctx context.Context, model enum.AIModel, systemPrompt *string, prompt any) (*string, error) {
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
	client := NewAnthropicClient(s.anthropicConfig, model)
	response, err := client.Invoke(ctx, systemPrompt, prompt)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &response, nil
}
