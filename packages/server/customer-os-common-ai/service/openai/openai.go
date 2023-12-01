package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/go-querystring/query"
)

const completionsUrl = "https://api.openai.com/v1/chat/completions"

// Error is the error standard response from the API
type Error struct {
	Message string      `json:"message,omitempty"`
	Type    string      `json:"type,omitempty"`
	Param   interface{} `json:"param,omitempty"`
	Code    interface{} `json:"code,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s - %s", e.Code, e.Message)
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type OpenAiClient struct {
	client       HttpClient
	apiKey       string
	organization string
	model        string
}

// NewOpenAiClient creates a new client
func NewOpenAiClient(apiKey, organization, model string) *OpenAiClient {
	return &OpenAiClient{
		apiKey:       apiKey,
		organization: organization,
		client:       &http.Client{},
		model:        model,
	}
}

// Post makes a post request
func (c *OpenAiClient) Post(ctx context.Context, url string, input any) (response []byte, err error) {
	response = make([]byte, 0)

	rJson, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	resp, err := c.Call(ctx, http.MethodPost, url, bytes.NewReader(rJson))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return checkError(response)
}

// Get makes a get request
func (c *OpenAiClient) Get(ctx context.Context, url string, input any) (response []byte, err error) {
	if input != nil {
		vals, _ := query.Values(input)
		query := vals.Encode()

		if query != "" {
			sb := strings.Builder{}
			sb.WriteString(url)
			sb.WriteString("?")
			sb.WriteString(query)
			url = sb.String()
		}
	}

	resp, err := c.Call(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return checkError(response)
}

type ErrorResponse struct {
	Error *Error `json:"error,omitempty"`
}

func checkError(response []byte) ([]byte, error) {
	r := &ErrorResponse{}
	err := json.Unmarshal(response, r)
	if err != nil {
		return nil, err
	}
	if r.Error != nil {
		return nil, r.Error
	}
	return response, nil
}

// Call makes a request
func (c *OpenAiClient) Call(ctx context.Context, method string, url string, body io.Reader) (response *http.Response, err error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	sb := strings.Builder{}
	sb.WriteString("Bearer ")
	sb.WriteString(c.apiKey)
	authHeader := sb.String()

	req.Header.Add("Authorization", authHeader)
	req.Header.Add("Content-Type", "application/json")
	if c.organization != "" {
		req.Header.Add("OpenAI-Organization", c.organization)
	}

	resp, err := c.client.Do(req)
	return resp, err
}

type FunctionCall struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

type Message struct {
	Role         string        `json:"role,omitempty"`
	Content      string        `json:"content,omitempty"`
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
}

type CreateChatCompletionsRequest struct {
	Model            string               `json:"model,omitempty"`
	Messages         []Message            `json:"messages,omitempty"`
	Functions        []CompletionFunciton `json:"functions,omitempty"`
	FunctionCall     *string              `json:"function_call,omitempty"`
	Temperature      float64              `json:"temperature,omitempty"`
	TopP             float64              `json:"top_p,omitempty"`
	N                int                  `json:"n,omitempty"`
	Stream           bool                 `json:"stream,omitempty"`
	Stop             StrArray             `json:"stop,omitempty"`
	MaxTokens        int                  `json:"max_tokens,omitempty"`
	PresencePenalty  float64              `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64              `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]string    `json:"logit_bias,omitempty"`
	User             string               `json:"user,omitempty"`
	ResponseFormat   ResponseFormat       `json:"response_format,omitempty"`
	Seed             int                  `json:"seed,omitempty"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type CompletionFunciton struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Parameters  []byte `json:"parameters,omitempty"`
}

type CreateChatCompletionsResponse struct {
	ID                string                        `json:"id,omitempty"`
	Object            string                        `json:"object,omitempty"`
	Created           int                           `json:"created,omitempty"`
	Model             string                        `json:"model,omitempty"`
	Choices           []CreateChatCompletionsChoice `json:"choices,omitempty"`
	Usage             CreateChatCompletionsUsave    `json:"usage,omitempty"`
	SystemFingerprint string                        `json:"system_fingerprint,omitempty"`
}

type CreateChatCompletionsChoice struct {
	Index        int      `json:"index,omitempty"`
	Message      *Message `json:"message,omitempty"`
	Delta        *Message `json:"delta,omitempty"`
	FinishReason string   `json:"finish_reason,omitempty"`
}

type CreateChatCompletionsUsave struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}

func (c *OpenAiClient) CreateChatCompletionsRaw(ctx context.Context, r *CreateChatCompletionsRequest) ([]byte, error) {
	if r.ResponseFormat.Type == "" {
		r.ResponseFormat.Type = "text"
	}
	return c.Post(ctx, completionsUrl, r)
}

func (c *OpenAiClient) CreateChatCompletions(ctx context.Context, r *CreateChatCompletionsRequest) (response *CreateChatCompletionsResponse, err error) {
	raw, err := c.CreateChatCompletionsRaw(ctx, r)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(raw, &response)
	return response, err
}
