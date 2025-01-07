package service

import (
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/clients/anthropic_client"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/config"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/tracing"
)

type AIService interface {
	AskAI(ctx context.Context, model enum.AIModel, prompt *string) (*string, error)
}

type aiService struct {
	config   *config.GlobalConfig
	services *Services
}

func NewAIService(config *config.GlobalConfig, services *Services) AIService {
	return &aiService{
		config:   config,
		services: services,
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
	client := anthropic_client.NewAnthropicClient(s.config, model)
	response, err := client.Invoke(ctx, prompt)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}

	return &response, nil
}
