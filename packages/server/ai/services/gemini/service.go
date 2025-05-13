package gemini

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	common_enum "github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"
	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	postgres_repository "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/repository"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"github.com/customeros/customeros/packages/server/ai/interfaces"
	"github.com/customeros/customeros/packages/server/ai/internal/clients"
	"github.com/customeros/customeros/packages/server/ai/internal/config"
	"github.com/customeros/customeros/packages/server/ai/internal/utils"
	"github.com/customeros/customeros/packages/server/ai/services/common"
)

const (
	DefaultTemperature = 0.7
	MaxTokens          = 1024
	DefaultTimeout     = 45 * time.Second
	DefaultGeminiModel = "gemini-2.0-flash"
)

type GeminiService struct {
	config         *config.GeminiConfig
	warehouseRepos *postgres_repository.WarehouseRepositories
	httpClient     *http.Client // Store your custom HTTP client
}

func NewGeminiService(cfg *config.GeminiConfig, warehouseRepos *postgres_repository.WarehouseRepositories) *GeminiService {
	DEFAULT_TIMEOUT := 45 * time.Second
	return &GeminiService{
		config:         cfg,
		warehouseRepos: warehouseRepos,
		httpClient:     clients.NewLoggingClient(warehouseRepos.APICallLogRepository, common_enum.VendorGemini, &DEFAULT_TIMEOUT),
	}
}

func (s *GeminiService) Ask(ctx context.Context, message interfaces.AskAIRequest) (*string, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "GeminiService.Ask")
	defer span.Finish()
	span.LogObjectAsJson("request", message)

	model := message.Model.String()
	if model == "" {
		model = DefaultGeminiModel
	} else if !strings.HasPrefix(model, "gemini") {
		err := fmt.Errorf("model '%s' is not a Gemini model", model)
		span.TraceError(err)
		return nil, err
	}

	// Validate prompt
	if message.Prompt == nil || *message.Prompt == "" {
		err := errors.New("prompt cannot be empty")
		span.TraceError(err)
		return nil, err
	}

	// Execute with retry
	respText, err := s.executeWithRetry(ctx, message, model)
	if err != nil {
		span.TraceError(err)
		return nil, err
	}

	return &respText, nil
}

func (s *GeminiService) createClient(ctx context.Context) (*genai.Client, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "GeminiService.createClient")
	defer span.Finish()

	if s.config.ApiKey == "" {
		err := errors.New("Gemini API key not configured")
		span.TraceError(err)
		return nil, err
	}

	// Use option.WithHTTPClient to pass your custom client
	client, err := genai.NewClient(ctx, option.WithAPIKey(s.config.ApiKey), option.WithHTTPClient(s.httpClient))
	if err != nil {
		span.TraceError(err)
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return client, nil
}

func applyGenerationConfig(model *genai.GenerativeModel, message interfaces.AskAIRequest) {
	config := &genai.GenerationConfig{}

	// Set temperature
	if message.ModelTemperature != nil {
		config.Temperature = message.ModelTemperature
	} else {
		defaultTemp32 := float32(DefaultTemperature)
		config.Temperature = &defaultTemp32
	}

	// Set max tokens
	if message.MaxOutputTokens != nil {
		maxOut := int32(*message.MaxOutputTokens)
		config.MaxOutputTokens = &maxOut
	} else {
		defaultMaxOut := int32(MaxTokens)
		config.MaxOutputTokens = &defaultMaxOut
	}

	model.GenerationConfig = *config
}

func (s *GeminiService) processResponse(resp *genai.GenerateContentResponse) (string, error) {
	if resp == nil {
		return "", errors.New("empty response from Gemini API")
	}

	var respStr string

	if len(resp.Candidates) > 0 {
		// Assuming only one candidate for simplicity.  Handle multiple candidates if needed.
		for _, part := range resp.Candidates[0].Content.Parts {
			if text, ok := part.(genai.Text); ok {
				respStr += string(text)
			}
		}
	}

	if respStr == "" {
		if len(resp.Candidates) > 0 {
			finishReason := resp.Candidates[0].FinishReason
			if finishReason != genai.FinishReasonStop && finishReason != genai.FinishReasonUnspecified {
				return "", fmt.Errorf("gemini response generation stopped due to: %s", finishReason)
			}
		}
		return "", errors.New("answer from Gemini was empty")
	}
	return respStr, nil
}

func (s *GeminiService) executeWithRetry(ctx context.Context, message interfaces.AskAIRequest, modelName string) (string, error) {
	span, ctx := telemetry.StartServiceSpan(ctx, "GeminiService.executeWithRetry")
	defer span.Finish()

	var lastErr error
	var apiResponse *genai.GenerateContentResponse // Capture the raw API response

	for attempt := 1; attempt <= common.MAX_RETRIES; attempt++ {
		client, err := s.createClient(ctx)
		if err != nil {
			span.TraceError(err)
			return "", fmt.Errorf("failed to create gemini client: %w", err)
		}
		defer client.Close()

		// Get the model handle
		model := client.GenerativeModel(modelName)

		// Apply generation parameters
		applyGenerationConfig(model, message)

		var contents []*genai.Content

		if message.SystemPrompt != nil && *message.SystemPrompt != "" {
			contents = append(contents, &genai.Content{
				Parts: []genai.Part{genai.Text(*message.SystemPrompt)},
				Role:  "user",
			})
		}

		// Add user prompt
		contents = append(contents, &genai.Content{
			Parts: []genai.Part{genai.Text(*message.Prompt)},
			Role:  "user",
		})

		// Flatten contents into a slice of parts
		var parts []genai.Part
		for _, content := range contents {
			parts = append(parts, content.Parts...)
		}

		apiResponse, err = model.GenerateContent(ctx, parts...)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: error calling gemini api: %w", attempt, err)
			if attempt == common.MAX_RETRIES {
				span.TraceError(lastErr)
				return "", lastErr
			}
			common.HandleRetry(attempt, err)
			continue
		}

		response, processErr := s.processResponse(apiResponse)
		if processErr != nil {
			lastErr = fmt.Errorf("attempt %d: error processing gemini response: %w", attempt, processErr)
			if !common.ShouldRetry(attempt) {
				span.TraceError(lastErr)
				return "", lastErr
			}
			common.HandleRetry(attempt, processErr)
			continue
		}

		span.LogKV("result_length", len(response))
		// Log response body here if needed

		// Now the function also logs request body and response
		go s.logGeminiCall(message, apiResponse, ctx, modelName)
		return response, nil
	}

	return "", lastErr
}

func (s *GeminiService) logGeminiCall(message interfaces.AskAIRequest, resp *genai.GenerateContentResponse, ctx context.Context, modelName string) {
	span, ctx := telemetry.StartServiceSpan(ctx, "GeminiService.logGeminiCall")
	defer span.Finish()

	requestID := utils.GenerateAPIRequestID()

	// Log request start time
	requestStartTime := time.Now()

	// Prepare request body for logging
	requestDataForLog := map[string]any{
		"model":       modelName,
		"prompt":      *message.Prompt,
		"temperature": message.ModelTemperature,
		"max_tokens":  message.MaxOutputTokens,
		// ... Add more relevant request params
	}

	requestBytes, marshalErr := json.Marshal(requestDataForLog)

	if marshalErr != nil {
		span.TraceError(fmt.Errorf("Failed to marshal Gemini request for logging: %w", marshalErr))
		return
	}

	// Prepare the response body for logging
	responseDataForLog := map[string]any{
		"candidates": resp.Candidates, // Log the full response for debugging
	}

	responseBytes, marshalResponseErr := json.Marshal(responseDataForLog)

	if marshalResponseErr != nil {
		span.TraceError(fmt.Errorf("Failed to marshal Gemini response for logging: %w", marshalResponseErr))
		return
	}

	// Calculate duration
	duration := time.Since(requestStartTime)
	durationMs := int(duration.Milliseconds())

	// Create log entry
	logEntry := &postgres_entity.APICallLog{
		ID:          requestID,
		Vendor:      common_enum.VendorGemini, // Assuming a default vendor
		Method:      "POST",
		URL:         fmt.Sprintf("gemini.api/%s/generateContent", modelName),
		RequestBody: requestBytes, // Can be nil if marshaling failed
		Timestamp:   requestStartTime,
		Duration:    durationMs,
	}

	// Handle response or error
	if resp != nil {
		statusCode := 200 // Assuming success (HTTP 200 OK)
		logEntry.StatusCode = &statusCode
		if len(responseBytes) > 0 {
			logEntry.ResponseBody = &responseBytes
		}

	} else { // If response is nil, then most likely an error occurred
		errMsg := "No response from Gemini API" // Custom error message
		logEntry.ErrorMessage = &errMsg
	}

	// Log asynchronously to not block the request
	go func() {
		bgCtx, cancelLog := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelLog()

		if logErr := s.warehouseRepos.APICallLogRepository.Create(bgCtx, logEntry); logErr != nil {
			span.TraceError(fmt.Errorf("Failed to log API call to DB (RequestID: %s): %w", requestID, logErr))
			return
		}
	}()
}
