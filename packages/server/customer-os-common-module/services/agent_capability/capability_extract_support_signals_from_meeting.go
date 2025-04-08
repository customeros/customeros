package agent_capability

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/telemetry"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"

	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
)

type ExtractSupportSignalsFromMeetingCapability struct {
	aiService interfaces.AIService
}

func NewExtractSupportSignalsFromMeetingCapability(aiService interfaces.AIService) *ExtractSupportSignalsFromMeetingCapability {
	return &ExtractSupportSignalsFromMeetingCapability{
		aiService: aiService,
	}
}

// Compile-time interface checks
var (
	_ interfaces.AgentCapability[ExtractSupportSignalsFromMeetingInput, ExtractSupportSignalsFromMeetingOutput, postgres_entity.NoConfig] = (*ExtractSupportSignalsFromMeetingCapability)(nil)
)

func (c *ExtractSupportSignalsFromMeetingCapability) Type() enum.AgentCapability {
	return enum.CapabilityExtractSupportSignalsFromMeeting
}

func (c *ExtractSupportSignalsFromMeetingCapability) Name() string {
	return "Extract support signals from meeting"
}

func (c *ExtractSupportSignalsFromMeetingCapability) NewInput() ExtractSupportSignalsFromMeetingInput {
	return ExtractSupportSignalsFromMeetingInput{}
}

func (c *ExtractSupportSignalsFromMeetingCapability) NewConfig() postgres_entity.NoConfig {
	return postgres_entity.NoConfig{}
}

func (c *ExtractSupportSignalsFromMeetingCapability) DefaultConfig() any {
	config := c.NewConfig()
	return &config
}

func (c *ExtractSupportSignalsFromMeetingCapability) DefaultActive() bool {
	return true
}

func (c *ExtractSupportSignalsFromMeetingCapability) ValidateInput(input ExtractSupportSignalsFromMeetingInput) error {
	if input.MeetingContent == "" {
		return coserrors.ErrMeetingContentMissing
	}
	return nil
}

func (c *ExtractSupportSignalsFromMeetingCapability) ValidateConfig(config postgres_entity.NoConfig) error {
	return nil
}

type ExtractSupportSignalsFromMeetingInput struct {
	MeetingContent string `json:"meetingContent"`
}

type ExtractSupportSignalsFromMeetingOutput struct {
	HelpNeeded []string `json:"helpNeeded"`
}

func (c *ExtractSupportSignalsFromMeetingCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[ExtractSupportSignalsFromMeetingInput, postgres_entity.NoConfig]) (enum.CapabilityExecutionStatus, ExtractSupportSignalsFromMeetingOutput, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ExtractSupportSignalsFromMeetingCapability.Execute")
	defer spans.Finish()

	spans.LogObjectAsJson("input", executionContainer.InputData)
	spans.LogObjectAsJson("config", executionContainer.ConfigData)

	result := ExtractSupportSignalsFromMeetingOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid input"))
		return enum.CapabilityExecutionError, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		spans.TraceError(errors.Wrap(err, "invalid config"))
		return enum.CapabilityExecutionError, result, err
	}

	systemPrompt := `Your objective is to identify everywhere the user or prospect needs help so our team can follow up on all points. Ensure you do not make anything up or leave anything out, and please communicate in simple, clear, and direct language.

Please analyze the meeting and respond in this exact JSON format:
{
    "helpNeeded": [
        "Name of requester: short description of help required",
        "Name of requester: short description of help required"
    ]
}

Important: Always provide your answer as valid JSON. If there is no help required, simply return help_needed with an empty array.`

	answer, err := c.aiService.AskAI(ctx, interfaces.AskAIRequest{
		Model:        enum.AIModelGemini,
		SystemPrompt: &systemPrompt,
		Prompt:       &executionContainer.InputData.MeetingContent,
	})
	if err != nil {
		spans.TraceError(err)
		return enum.CapabilityExecutionError, result, err
	}
	if answer == nil {
		err := errors.New("No answer from AI")
		spans.TraceError(err)
		return enum.CapabilityExecutionError, result, err
	}

	result, err = c.parseAnswer(ctx, *answer)
	if err != nil {
		spans.TraceError(err)
		return enum.CapabilityExecutionError, result, err
	}

	spans.LogObjectAsJson("result", result)
	return enum.CapabilityExecutionCompleted, result, nil
}

func (c *ExtractSupportSignalsFromMeetingCapability) parseAnswer(ctx context.Context, answer string) (ExtractSupportSignalsFromMeetingOutput, error) {
	spans, ctx := telemetry.StartServiceSpan(ctx, "ExtractSupportSignalsFromMeetingCapability.parseAnswer")
	defer spans.Finish()

	var result ExtractSupportSignalsFromMeetingOutput
	err := json.Unmarshal([]byte(answer), &result)
	if err != nil {
		spans.TraceError(err)
		return ExtractSupportSignalsFromMeetingOutput{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return result, nil
}
