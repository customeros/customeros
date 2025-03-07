package enum

import "fmt"

type AIModel string

const (
	AIModelAnthropicHaiku  AIModel = "claude-3-5-haiku-20241022"
	AIModelAnthropicSonnet AIModel = "claude-3-5-sonnet-20241022"
	AIModelDeepseekChat    AIModel = "deepseek-chat"
	AIModelDeepseekQwen    AIModel = "deepseek-r1-distill-qwen-32b"
	AIModelGemini          AIModel = "gemini-2.0-flash"
	AIModelGeminiLite      AIModel = "gemini-2.0-flash-lite"
	AIModelGemma           AIModel = "gemma2-9b-it"
	AIModelJinaEmbeddings  AIModel = "jina-embeddings-v3"
	AIModelLlama8B         AIModel = "llama-3.1-8b-instant"
	AIModelLlama70B        AIModel = "llama-3.3-70b-versatile"
	AIModelMixtral         AIModel = "mixtral-8x7b-32768"
	AIModelWhisper         AIModel = "distil-whisper-large-v3-en"
)

func (a AIModel) String() string {
	return string(a)
}

func GetAIModel(s string) (AIModel, error) {
	switch AIModel(s) {
	case
		AIModelAnthropicHaiku,
		AIModelAnthropicSonnet,
		AIModelDeepseekChat,
		AIModelDeepseekQwen,
		AIModelGemini,
		AIModelGemma,
		AIModelJinaEmbeddings,
		AIModelLlama8B,
		AIModelLlama70B,
		AIModelMixtral,
		AIModelWhisper:
		return AIModel(s), nil

	default:
		return "", fmt.Errorf("invalid AIModel: %s", s)
	}
}

type AIOutputFormat string

const (
	AIOutputText AIOutputFormat = "text"
	AIOutputJson AIOutputFormat = "json"
)

func (a AIOutputFormat) String() string {
	return string(a)
}

type AIRequestType string

const (
	AIRequestGeneric         AIRequestType = "generic"
	AIRequestCompanyName     AIRequestType = "company_name"
	AIRequestContentStage    AIRequestType = "content_journey_stage"
	AIRequestIndustryCode    AIRequestType = "industry_code"
	AIRequestWebpageCategory AIRequestType = "webpage_category"
	AIRequestWebpageTopics   AIRequestType = "webpage_topics"
)

func (a AIRequestType) String() string {
	return string(a)
}
