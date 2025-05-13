package interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/ai/internal/enum"
)

type AIService interface {
	AskAI(ctx context.Context, message AskAIRequest) (*string, error)
}

type AskAIRequest struct {
	RequestType      enum.AIRequestType
	Model            enum.AIModel
	SystemPrompt     *string
	Prompt           *string
	ModelTemperature *float32
	MaxOutputTokens  *int32
	OutputFormat     enum.AIOutputFormat
	Retries          *int
}
