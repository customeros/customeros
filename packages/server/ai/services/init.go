package services

import (
	"context"
	"fmt"

	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"

	"github.com/customeros/customeros/packages/server/ai/interfaces"
	"github.com/customeros/customeros/packages/server/ai/internal/config"
	nats_internal "github.com/customeros/customeros/packages/server/ai/internal/nats"
	"github.com/customeros/customeros/packages/server/ai/internal/repository"
	"github.com/customeros/customeros/packages/server/ai/services/anthropic"
	"github.com/customeros/customeros/packages/server/ai/services/deepseek"
	"github.com/customeros/customeros/packages/server/ai/services/gemini"
	"github.com/customeros/customeros/packages/server/ai/services/groq"
	"github.com/customeros/customeros/packages/server/ai/services/webpage_classification"
	"github.com/customeros/customeros/packages/server/ai/services/webpage_intent_profiler"
)

type Services struct {
	WebpageIntentProfiler interfaces.NatsService
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
			natsConn, anthropic, deepseek, groq, gemini,
		),
		WebpageIntentProfiler: webpage_intent_profiler.NewWebpageIntentProfiler(
			natsConn, anthropic, deepseek, groq, gemini,
		),
	}

	return services
}

func (s *Services) Start(ctx context.Context) error {
	services := []struct {
		name    string
		starter func(context.Context) error
	}{
		{"Webpage Classifier", s.WebpageClassification.Start},
		{"Webpage Intent Profiler", s.WebpageIntentProfiler.Start},
	}

	for _, svc := range services {
		if err := svc.starter(ctx); err != nil {
			return fmt.Errorf("failed to start %s service: %w", svc.name, err)
		}
	}

	return nil
}

func (s *Services) Stop(ctx context.Context) {
	services := []struct {
		name    string
		stopper func(context.Context)
	}{
		{"Webpage Classifier", func(ctx context.Context) { s.WebpageClassification.Stop() }},
		{"Webpage Intent Profiler", func(ctx context.Context) { s.WebpageIntentProfiler.Stop() }},
	}

	for _, service := range services {
		service.stopper(ctx)
	}
}
