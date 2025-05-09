package deepseek

type Message struct {
	Role    string `json:"role"` // "user" or "assistant"
	Content any    `json:"content"`
}

type DeepseekRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type DeepseekResponse struct {
	ID                string        `json:"id"`
	Object            string        `json:"object"`
	Created           int64         `json:"created"`
	Model             string        `json:"model"`
	Choices           []Choice      `json:"choices"`
	Usage             DeepseekUsage `json:"usage"`
	SystemFingerprint string        `json:"system_fingerprint"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	LogProbs     any     `json:"logprobs"`
	FinishReason string  `json:"finish_reason"`
}

type DeepseekUsage struct {
	PromptTokens          int                 `json:"prompt_tokens"`
	CompletionTokens      int                 `json:"completion_tokens"`
	TotalTokens           int                 `json:"total_tokens"`
	PromptTokensDetails   PromptTokensDetails `json:"prompt_tokens_details"`
	PromptCacheHitTokens  int                 `json:"prompt_cache_hit_tokens"`
	PromptCacheMissTokens int                 `json:"prompt_cache_miss_tokens"`
}

type PromptTokensDetails struct {
	CachedTokens int `json:"cached_tokens"`
}
