package interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/ai/internal/enum"
)

type LLMClient interface {
	Ask(ctx context.Context, request AskAIRequest) (*string, error)
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
