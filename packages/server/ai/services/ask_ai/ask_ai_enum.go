package ai

import (
	"fmt"
	"strings"

	"github.com/customeros/customeros/packages/server/ai/interfaces"
)

type EnumValidator interface {
	IsValid(value string) bool
	ValidValues() []string
}

func (s *aiService) enhancePromptWithValidValues(request *interfaces.AskAIRequest, validator EnumValidator) {
	if request.SystemPrompt != nil {
		currentPrompt := *request.SystemPrompt
		validOptions := strings.Join(validator.ValidValues(), ", ")

		if !strings.Contains(currentPrompt, validOptions) {
			enhancedPrompt := fmt.Sprintf(`%s

IMPORTANT: Your response MUST be exactly one of these valid values: %s
Do not include any additional text, comments, or explanation.`,
				currentPrompt, validOptions)
			request.SystemPrompt = &enhancedPrompt
		}
	}
}

func (s *aiService) cleanResponseString(response string) string {
	return strings.TrimSpace(response)
}

func (s *aiService) addRetryInstructions(request *interfaces.AskAIRequest, validator EnumValidator) {
	if request.SystemPrompt != nil {
		retryPrompt := fmt.Sprintf(`%s

Your previous response was invalid. Please ONLY respond with exactly one of these values: %s
No explanation, commentary, or additional text.`,
			*request.SystemPrompt,
			strings.Join(validator.ValidValues(), ", "))
		request.SystemPrompt = &retryPrompt
	}
}
