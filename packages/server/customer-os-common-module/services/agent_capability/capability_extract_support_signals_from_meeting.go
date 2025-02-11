package agent_capability

import (
	"context"
	"encoding/json"
	"fmt"

	postgres_entity "github.com/customeros/customeros/packages/server/customer-os-postgres-repository/entity"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/customeros/customeros/packages/server/customer-os-common-module/common"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/coserrors"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/enum"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/interfaces"
	"github.com/customeros/customeros/packages/server/customer-os-common-module/tracing"
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
	return "Extract meeting highlights"
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
	MeetingSummary string   `json:"meetingSummary"`
	ActionItems    []string `json:"actionItems"`
}

func (c *ExtractSupportSignalsFromMeetingCapability) Execute(ctx context.Context, executionContainer interfaces.TypedExecutionContainer[ExtractSupportSignalsFromMeetingInput, postgres_entity.NoConfig]) (bool, ExtractSupportSignalsFromMeetingOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExtractSupportSignalsFromMeetingCapability.Execute")
	defer span.Finish()
	tracing.TagComponentService(span)
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	tracing.LogObjectAsJson(span, "input", executionContainer.InputData)
	tracing.LogObjectAsJson(span, "config", executionContainer.ConfigData)

	result := ExtractSupportSignalsFromMeetingOutput{}

	if err := c.ValidateInput(executionContainer.InputData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid input"))
		return false, result, err
	}
	if err := c.ValidateConfig(executionContainer.ConfigData); err != nil {
		tracing.TraceErr(span, errors.Wrap(err, "invalid config"))
		return false, result, err
	}

	systemPrompt := `Your objective is to capture the main outcome of the meeting for the CEO in two sentences or less.  Then provide clear account of all action points that need to be followed up on and their owners.  Ensure you do not make anything up or leave anything out, and please communicate in simple, clear, and direct language.

Please analyze the meeting and respond in this exact JSON format:
{
    "summary": "This is my clear and direct summary for the CEO.",
    "actionItems": [
        "Name of owner: action item",
        "Name of owner: action item"
    ]
}

Important: Always provide your answer as valid JSON. If there are no clear action items, omit the actionItems parameter and array from your response.`

	answer, err := c.aiService.AskAI(ctx, enum.AIModelAnthropicHaiku, systemPrompt, executionContainer.InputData.MeetingContent)
	if err != nil {
		tracing.TraceErr(span, err)
		return true, result, err
	}
	if answer == nil {
		err := errors.New("No answer from AI")
		tracing.TraceErr(span, err)
		return true, result, err
	}

	result, err = c.parseAnswer(ctx, *answer)
	if err != nil {
		tracing.TraceErr(span, err)
		return true, result, err
	}

	tracing.LogObjectAsJson(span, "result", result)
	return true, result, nil
}

func (c *ExtractSupportSignalsFromMeetingCapability) parseAnswer(ctx context.Context, answer string) (ExtractSupportSignalsFromMeetingOutput, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ExtractSupportSignalsFromMeetingCapability.parseAnswer")
	defer span.Finish()
	tracing.TagComponentService(span)

	var result ExtractSupportSignalsFromMeetingOutput
	err := json.Unmarshal([]byte(answer), &result)
	if err != nil {
		tracing.TraceErr(span, err)
		return ExtractSupportSignalsFromMeetingOutput{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if result.MeetingSummary == "" {
		err := errors.New("Unable to produce meeting summary")
		tracing.TraceErr(span, err)
		return ExtractSupportSignalsFromMeetingOutput{}, err
	}

	return result, nil
}
