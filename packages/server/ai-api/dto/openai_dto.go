package dto

type OpenAiApiRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`

	//Note that our models may stop before reaching this maximum. This parameter only specifies the absolute maximum number of tokens to generate.
	//Default to: 256
	MaxTokensToSample *int `json:"maxTokensToSample"`

	//Defaults to 1. Ranges from 0 to 1. Use temp closer to 0 for analytical / multiple choice, and closer to 1 for creative and generative tasks.
	Temperature *int `json:"temperature"`
}

type OpenAiApiResponse struct {
	Model string `json:"model"`

	Choices *[]struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`

	Usage *struct {
		CompletionTokens int `json:"completion_tokens"`
		PromptTokens     int `json:"prompt_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`

	Error *OpenAiApiErrorResponse `json:"error"`
}

type OpenAiApiErrorResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
