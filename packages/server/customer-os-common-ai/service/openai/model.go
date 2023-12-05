package openai

import "context"

func NewModel(apiKey, organization, model string) *OpenAiModel {
	// TODO: use openai config object plus add logger to openai client
	return &OpenAiModel{
		Client: NewOpenAiClient(apiKey, organization, model),
	}
}

type OpenAiModel struct {
	Client *OpenAiClient
}

func limitTokens(s string, n int) string {
	if len(s)/4 > n {
		return s[:n]
	}
	return s
}

func (m *OpenAiModel) Inference(ctx context.Context, input string) (string, error) {
	input = limitTokens(input, 16385)
	request := &CreateChatCompletionsRequest{
		Model: m.Client.model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant designed to output JSON.",
			},
			{
				Role:    "user",
				Content: input,
			},
		},
		Temperature:    0.7,
		Seed:           42,                                  // always use same seed to increase likelihood of consistent results
		ResponseFormat: ResponseFormat{Type: "json_object"}, // https://platform.openai.com/docs/guides/text-generation/json-mode
	}

	response, err := m.Client.CreateChatCompletions(ctx, request)
	if err != nil {
		return "", err
	}
	return response.Choices[0].Message.Content, nil
}
