package interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type AIService interface {
	AskAI(ctx context.Context, model enum.AIModel, systemPrompt *string, prompt any) (*string, error)
}
