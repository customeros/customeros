package ai

import (
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/config"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
)

type aiService struct {
	config *config.AnthropicConfig
}

func NewAIService(config *config.AnthropicConfig) interfaces.AIService {
	return &aiService{
		config: config,
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

	// setup client
	client := NewAnthropicClient(s.config, model)
	response, err := client.Invoke(ctx, prompt)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &response, nil
}
