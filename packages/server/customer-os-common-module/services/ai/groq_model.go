package ai

type GroqRequest struct {
	Model               string         `json:"model"`
	Messages            []Message      `json:"messages"`
	MaxCompletionTokens int32          `json:"max_completion_tokens,omitempty"`
	Temperature         float32        `json:"temperature,omitempty"`
	ResponseFormat      ResponseFormat `json:"response_format,omitempty"`
}

type ResponseFormat struct {
	Type string `json:"type"` // add value "json_object" to guarantee a valid json response
}

type GroqResponse struct {
	ID                string       `json:"id"`
	Object            string       `json:"object"`
	Created           int64        `json:"created"`
	Model             string       `json:"model"`
	Choices           []GroqChoice `json:"choices"`
	Usage             GroqUsage    `json:"usage"`
	SystemFingerprint string       `json:"system_fingerprint"`
}

// Choice represents a single choice in the response
type GroqChoice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	LogProbs     any     `json:"logprobs"`
	FinishReason string  `json:"finish_reason"`
}

// Usage contains token and time usage statistics
type GroqUsage struct {
	QueueTime        float64 `json:"queue_time"`
	PromptTokens     int     `json:"prompt_tokens"`
	PromptTime       float64 `json:"prompt_time"`
	CompletionTokens int     `json:"completion_tokens"`
	CompletionTime   float64 `json:"completion_time"`
	TotalTokens      int     `json:"total_tokens"`
	TotalTime        float64 `json:"total_time"`
}
