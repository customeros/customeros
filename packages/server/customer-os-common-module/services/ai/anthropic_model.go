package ai

// AnthropicApiRequest represents the request structure for Claude API
type AnthropicApiRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

// Message represents a single message in the conversation
type Message struct {
	Role    string `json:"role"` // "user" or "assistant"
	Content any    `json:"content"`
}

// AnthropicApiResponse represents the response structure from Claude API
type AnthropicApiResponse struct {
	Id         string    `json:"id"`
	Type       string    `json:"type"`
	Role       string    `json:"role"`
	Content    []Content `json:"content"`
	Model      string    `json:"model"`
	StopReason string    `json:"stop_reason,omitempty"`
	Usage      Usage     `json:"usage"`
}

// Content represents the content structure in the response
type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Usage represents token usage information
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}
