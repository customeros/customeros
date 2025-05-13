package services

import (
	"context"

	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/ai/interfaces"
	"github.com/customeros/customeros/packages/server/ai/internal/config"
	nats_internal "github.com/customeros/customeros/packages/server/ai/internal/nats"
	"github.com/customeros/customeros/packages/server/ai/internal/repository"
	"github.com/customeros/customeros/packages/server/ai/internal/telemetry"
	"github.com/customeros/customeros/packages/server/ai/services/anthropic"
	"github.com/customeros/customeros/packages/server/ai/services/deepseek"
	"github.com/customeros/customeros/packages/server/ai/services/gemini"
	"github.com/customeros/customeros/packages/server/ai/services/groq"
	"github.com/customeros/customeros/packages/server/ai/services/webpage_classification"
)

type Services struct {
	WebpageClassification interfaces.NatsService
}

func InitServices(
	config *config.Config,
	repositories *repository.Repositories,
	warehouseRepos *postgres_repository.WarehouseRepositories,
	natsConn *nats_internal.NATSConnections,
) *Services {
	// model providers
	anthropic := anthropic.NewAnthropicService(config.Anthropic, warehouseRepos)
	deepseek := deepseek.NewDeepseekService(config.Deepseek, warehouseRepos)
	gemini := gemini.NewGeminiService(config.Gemini, warehouseRepos)
	groq := groq.NewGroqService(config.Groq, warehouseRepos)

	// AI-enabled services
	services := &Services{
		WebpageClassification: webpage_classification.NewWebpageClassificationService(
			natsConn, anthropic, deepseek, groq, gemini, repositories,
		),
	}

	return services
}

func (s *Services) Start(ctx context.Context) error {
	span, ctx := telemetry.StartServiceSpan(ctx, "Services.Start")
	defer span.Finish()

	err := s.WebpageClassification.Start(ctx)
	if err != nil {
		span.TraceError(err)
		return err
	}
	return nil
}

func (s *Services) Stop(ctx context.Context) {}
