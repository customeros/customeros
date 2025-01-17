package interfaces

import (
	"context"

	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-common-module/enum"
)

type AIService interface {
	AskAI(ctx context.Context, model enum.AIModel, prompt *string) (*string, error)
}
