package dto

type AnthropicApiRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model"`

	//Note that our models may stop before reaching this maximum. This parameter only specifies the absolute maximum number of tokens to generate.
	//Default to: 256
	MaxTokensToSample *int `json:"maxTokensToSample"`

	//Our models stop on "\n\nHuman:", and may include additional built-in stop sequences in the future.
	//By providing the stop_sequences parameter, you may include additional strings that will cause the model to stop generating.
	StopSequences *[]string `json:"stopSequences"`

	//Defaults to 1. Ranges from 0 to 1. Use temp closer to 0 for analytical / multiple choice, and closer to 1 for creative and generative tasks.
	Temperature *int `json:"temperature"`
}

type AnthropicApiResponse struct {
	Completion *string `json:"completion"`
	StopReason *string `json:"stopReason"`
	Model      *string `json:"model"`
	Error      *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}
