package interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type AIService interface {
	AskAI(ctx context.Context, request AskAIRequest) (*string, error)
}

type AskAIRequest struct {
	Model            enum.AIModel
	SystemPrompt     *string
	Prompt           *string
	ModelTemperature *float32
	MaxOutputTokens  *int32
	OutputFormat     enum.AIOutputFormat
	Retries          *int
}
