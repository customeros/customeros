package ai

import (
	"github.com/customeros/customeros/packages/server/ai/interfaces"
	"github.com/customeros/customeros/packages/server/ai/internal/repository"
	"github.com/customeros/customeros/packages/server/ai/services/anthropic"
	"github.com/customeros/customeros/packages/server/ai/services/deepseek"
	"github.com/customeros/customeros/packages/server/ai/services/gemini"
	"github.com/customeros/customeros/packages/server/ai/services/groq"
)

type aiService struct {
	anthropic    *anthropic.AnthropicService
	deepseek     *deepseek.DeepseekService
	groq         *groq.GroqService
	gemini       *gemini.GeminiService
	repositories *repository.Repositories
}

func NewAIService(
	anthropic *anthropic.AnthropicService,
	deepseek *deepseek.DeepseekService,
	groq *groq.GroqService,
	gemini *gemini.GeminiService,
	repositories *repository.Repositories,
) interfaces.AIService {
	return &aiService{
		anthropic:    anthropic,
		deepseek:     deepseek,
		groq:         groq,
		gemini:       gemini,
		repositories: repositories,
	}
}
