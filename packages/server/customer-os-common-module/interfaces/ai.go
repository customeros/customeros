package interfaces

import (
	"context"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/data_fields"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
)

type AIService interface {
	AskAI(ctx context.Context, request AskAIRequest) (*string, error)
	AskAIForCompanyDescription(ctx context.Context, request AskAIRequest) (*data_fields.CompanyDescription, error)
	AskAIForCompanyName(ctx context.Context, request AskAIRequest) (*data_fields.CompanyIdentification, error)
	AskAIForEmail(ctx context.Context, request AskAIRequest) (*data_fields.EmailResponse, error)
	AskAIForIndustryCode(ctx context.Context, request AskAIRequest) (*data_fields.IndustryCode, error)
	AskAIForString(ctx context.Context, request AskAIRequest) (*string, error)
	AskAIForWebpageCategory(ctx context.Context, request AskAIRequest) (enum.WebpageCategory, error)
	AskAIForWebpageTopics(ctx context.Context, request AskAIRequest) ([]data_fields.Topic, error)
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
