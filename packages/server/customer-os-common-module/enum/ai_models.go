package enum

import "fmt"

type AIModel string

const (
	AIModelAnthropicSonnet AIModel = "claude-3-5-sonnet-20241022"
	AIModelAnthropicHaiku  AIModel = "claude-3-5-haiku-20241022"
	AIModelDeepseekChat    AIModel = "deepseek-chat"
)

func (a AIModel) String() string {
	return string(a)
}

func GetAIModel(s string) (AIModel, error) {
	switch AIModel(s) {
	case
		AIModelAnthropicHaiku,
		AIModelAnthropicSonnet,
		AIModelDeepseekChat:
		return AIModel(s), nil

	default:
		return "", fmt.Errorf("invalid AIModel: %s", s)
	}
}
