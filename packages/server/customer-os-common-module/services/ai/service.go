package ai

import (
	"context"
	"errors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/logger"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type aiService struct {
	log             logger.Logger
	anthropicConfig *config.AnthropicConfig
}

func NewAIService(log logger.Logger, config *config.AnthropicConfig) interfaces.AIService {
	return &aiService{
		log:             log,
		anthropicConfig: config,
	}
}

func (s *aiService) AskAI(ctx context.Context, model enum.AIModel, prompt *string) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AIModelService.AskAI")
	defer span.Finish()
	span.LogKV("model", model)
	span.LogKV("prompt", prompt)

	switch model {
	case
		enum.AIModelAnthropicHaiku,
		enum.AIModelAnthropicSonnet:

		return s.askAnthropic(ctx, model, prompt)

	default:
		err := errors.New("Unsupported model")
		return nil, err
	}
}

func (s *aiService) askAnthropic(ctx context.Context, model enum.AIModel, prompt *string) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "AIModelService.AskAnthropic")
	defer span.Finish()

	if s.anthropicConfig.ApiKey == "" || s.anthropicConfig.ApiPath == "" {
		err := errors.New("Anthropic API key or path not set")
		tracing.TraceErr(span, err)
		s.log.Error(err)
		return nil, err
	}

	// setup client
	client := NewAnthropicClient(s.anthropicConfig, model)
	response, err := client.Invoke(ctx, prompt)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &response, nil
}
