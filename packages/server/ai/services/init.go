package services

import (
	"context"

	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/ai/interfaces"
	"github.com/customeros/customeros/packages/server/ai/internal/config"
	"github.com/customeros/customeros/packages/server/ai/internal/repository"
	"github.com/customeros/customeros/packages/server/ai/internal/telemetry"
	nats_internal "github.com/customeros/customeros/packages/server/ai/nats"
	"github.com/customeros/customeros/packages/server/ai/services/anthropic"
	ai "github.com/customeros/customeros/packages/server/ai/services/ask_ai"
	"github.com/customeros/customeros/packages/server/ai/services/deepseek"
	"github.com/customeros/customeros/packages/server/ai/services/gemini"
	"github.com/customeros/customeros/packages/server/ai/services/groq"
)

type Services struct {
	AskAI     interfaces.NatsService
	Anthropic *anthropic.AnthropicService
	Deepseek  *deepseek.DeepseekService
	Gemini    *gemini.GeminiService
	Groq      *groq.GroqService
}

func InitServices(
	config *config.Config,
	repositories *repository.Repositories,
	warehouseRepos *postgres_repository.WarehouseRepositories,
	natsConn *nats_internal.NATSConnections,
) *Services {
	services := &Services{
		Anthropic: anthropic.NewAnthropicService(config.Anthropic, warehouseRepos),
		Deepseek:  deepseek.NewDeepseekService(config.Deepseek, warehouseRepos),
		Gemini:    gemini.NewGeminiService(config.Gemini, warehouseRepos),
		Groq:      groq.NewGroqService(config.Groq, warehouseRepos),
	}

	services.AskAI = ai.NewAIService(
		natsConn,
		services.Anthropic,
		services.Deepseek,
		services.Groq,
		services.Gemini,
		repositories,
	)
	return services
}

func (s *Services) Start(ctx context.Context) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "Services.Start")
	defer span.Finish()

	err := s.AskAI.Start(ctx)
	if err != nil {
		span.TraceError(err)
		return err
	}
	return nil
}

func (s *Services) Stop(ctx context.Context) {}
